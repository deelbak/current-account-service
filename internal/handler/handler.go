package handler

import (
	"current-account-service/internal/service"
)

type Handler struct {
	auth   *service.AuthService
	app    *service.ApplicationService
	client *service.ClientService
	task   *service.TaskService
}

func New(
	auth *service.AuthService,
	app *service.ApplicationService,
	client *service.ClientService,
	task *service.TaskService,
) *Handler {
	return &Handler{auth: auth, app: app, client: client, task: task}
}
