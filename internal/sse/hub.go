package sse

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

type Message struct {
	ApplicationID int64     `json:"application_id"`
	State         string    `json:"state"`
	Event         string    `json:"event"`
	ActorID       int64     `json:"actor_id"`
	OccuredAt     time.Time `json:"occured_at"`
}

type client struct {
	ch     chan Message
	userID int64
}

type Hub struct {
	mu      sync.RWMutex
	clients map[string]*client
	Publish chan Message
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[string]*client),
		Publish: make(chan Message, 256),
	}
}

func (h *Hub) Run() {
	ticker := time.NewTicker(25 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case msg := <-h.Publish:
			h.broadcast(msg)
		case <-ticker.C:
			h.broadcast(Message{})
		}
	}
}

func (h *Hub) broadcast(msg Message) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, c := range h.clients {
		select {
		case c.ch <- msg:
		default:
		}
	}
}

func (h *Hub) Subscribe(clientID string, userID int64) chan Message {
	ch := make(chan Message, 32)
	h.mu.Lock()
	h.clients[clientID] = &client{
		ch:     ch,
		userID: userID,
	}
	h.mu.Unlock()
	return ch
}

func (h *Hub) Unsubscribe(clientID string) {
	h.mu.Lock()
	if c, ok := h.clients[clientID]; ok {
		close(c.ch)
		delete(h.clients, clientID)
	}
	h.mu.Unlock()
}

func (h *Hub) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "Streaming unsupported", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")

	userID := r.Context().Value("userID").(int64)
	clientID := r.Header.Get("X-Client-ID")

	if clientID == "" {
		clientID = fmt.Sprintf("%d-%d", userID, time.Now().UnixNano())
	}

	ch := h.Subscribe(clientID, userID)
	defer h.Unsubscribe(clientID)

	enc := json.NewEncoder(w)

	for {
		select {
		case msg, ok := <-ch:
			if !ok {
				return
			}
			if msg.ApplicationID == 0 {
				// heartbeat
				fmt.Fprintf(w, ": ping\n\n")
			} else {
				fmt.Fprint(w, "data: ")
				enc.Encode(msg)
				fmt.Fprint(w, "\n")
			}
			flusher.Flush()
		case <-r.Context().Done():
			return
		}
	}
}
