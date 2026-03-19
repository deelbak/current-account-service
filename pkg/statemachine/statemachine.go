package statemachine

import (
	"fmt"
	"time"
)

type State string
type Event string

const (
	StateDraft            State = "draft"
	StateDocsUploading    State = "docs_uploading"
	StateSubmitted        State = "submitted"
	StateUnderReview      State = "under_review"
	StateRevisionRequired State = "revision_required"
	StateApproved         State = "approved"
	StateRejected         State = "rejected"
	StateAccountOpened    State = "account_opened"
)

const (
	EventCreate      Event = "create"
	EventUploadDooc  Event = "upload_doc"
	EventSubmit      Event = "submit"
	EventReview      Event = "review"
	EventApprove     Event = "approve"
	EventRevision    Event = "revision"
	EventReject      Event = "reject"
	EventOpenAccount Event = "open_account"
)

type Transition struct {
	From  State
	Event Event
	To    State
}

var transitions = []Transition{
	{StateDraft, EventCreate, StateDocsUploading},
	{StateDocsUploading, EventSubmit, StateSubmitted},
	{StateSubmitted, EventReview, StateUnderReview},
	{StateUnderReview, EventApprove, StateApproved},
	{StateUnderReview, EventRevision, StateRevisionRequired},
	{StateUnderReview, EventReject, StateRejected},
	{StateRevisionRequired, EventSubmit, StateSubmitted},
	{StateApproved, EventOpenAccount, StateAccountOpened},
}

type Machine struct {
	transitions map[State]map[Event]State
}

func New() *Machine {
	m := &Machine{transitions: make(map[State]map[Event]State)}

	for _, t := range transitions {
		if m.transitions[t.From] == nil {
			m.transitions[t.From] = make(map[Event]State)
		}
		m.transitions[t.From][t.Event] = t.To
	}
	return m
}

func (m *Machine) Transition(current State, event Event) (State, error) {
	events, ok := m.transitions[current]
	if !ok {
		return current, fmt.Errorf("invalid current state: %s", current)
	}

	next, ok := events[event]
	if !ok {
		return current, fmt.Errorf("invalid event: %s for state: %s", event, current)
	}
	return next, nil
}

type ApplicationEvent struct {
	ApplicationID int64
	FromState     State
	ToState       State
	Event         Event
	ActorID       int64
	ActorRole     string
	Comment       string
	OccuredAt     time.Time
}
