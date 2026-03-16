package main

import (
	"fmt"
	"log"
	"net/http"

	"current-account-service/internal/config"
	"current-account-service/internal/handler"
	"current-account-service/internal/repository"
	"current-account-service/internal/service"

	"github.com/go-chi/chi/v5"
)

func main() {

	cfg := config.LoadConfig()

	db, err := repository.NewPostgres(cfg)
	if err != nil {
		log.Fatal(err)
	}

	err = db.Ping()
	if err != nil {
		log.Fatal("database not reachable")
	}

	fmt.Println("database connected")

	// repositories
	clientRepo := repository.NewClientRepository(db)

	// services
	clientService := service.NewClientService(clientRepo)

	// handlers
	clientHandler := handler.NewClientHandler(clientService)

	r := chi.NewRouter()
	r.Get("/clients/{iin}", clientHandler.GetClient)

	fmt.Println("server started on :8080")

	log.Fatal(http.ListenAndServe(":8080", r))
}
