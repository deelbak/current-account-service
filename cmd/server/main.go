package main

import (
	"database/sql"
	"log/slog"
	"net/http"
	"os"

	_ "github.com/lib/pq"

	"current-account-service/internal/handler"
	"current-account-service/internal/middleware"
	"current-account-service/internal/repository/postgres"
	"current-account-service/internal/service"
	"current-account-service/internal/sse"
	"current-account-service/pkg/statemachine"
)

func main() {
	db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		slog.Error("db connect", "err", err)
		os.Exit(1)
	}
	defer db.Close()

	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	// Repos
	userRepo := postgres.NewUserRepository(db)
	clientRepo := postgres.NewClientRepository(db)
	appRepo := postgres.NewApplicationRepository(db)
	docRepo := postgres.NewDocumentRepository(db)
	taskRepo := postgres.NewTaskRepository(db)
	accountRepo := postgres.NewAccountRepository(db)

	// SSE hub
	hub := sse.NewHub()
	go hub.Run()

	// Services
	sm := statemachine.New()
	notify := service.NewNotificationService(nil, nil) // подключи email/tg
	authSvc := service.NewAuthService(userRepo)
	clientSvc := service.NewClientService(clientRepo)
	taskSvc := service.NewTaskService(taskRepo)
	appSvc := service.NewApplicationService(appRepo, docRepo, accountRepo, taskRepo, sm, hub, notify)

	// Handlers
	h := handler.New(authSvc, appSvc, clientSvc, taskSvc)

	mux := http.NewServeMux()

	// Public
	mux.HandleFunc("POST /login", h.Login)

	// Manager + UL Head
	mux.HandleFunc("GET /clients/{iin}", middleware.Auth(h.GetClient))
	mux.HandleFunc("GET /applications/{id}", middleware.Auth(h.GetApplication))
	mux.HandleFunc("GET /tasks/my", middleware.Auth(h.MyTasks))
	mux.HandleFunc("GET /events", middleware.Auth(hub.ServeHTTP))

	// Manager only
	mux.HandleFunc("POST /applications", middleware.Auth(h.CreateApplication))
	mux.HandleFunc("POST /applications/{id}/documents", middleware.Auth(h.UploadDocument))
	mux.HandleFunc("POST /applications/{id}/submit", middleware.Auth(h.Submit))

	// UL Head only
	mux.HandleFunc("POST /applications/{id}/approve", middleware.RequireRole("ul_head", h.Approve))
	mux.HandleFunc("POST /applications/{id}/revision", middleware.RequireRole("ul_head", h.RequestRevision))
	mux.HandleFunc("POST /applications/{id}/reject", middleware.RequireRole("ul_head", h.Reject))

	addr := ":8080"
	slog.Info("server started", "addr", addr)
	if err := http.ListenAndServe(addr, mux); err != nil {
		slog.Error("server error", "err", err)
		os.Exit(1)
	}
}
