package service

import (
	"fmt"
	"log/slog"

	"current-account-service/internal/models"
)

type Notifier interface {
	Send(to, subject, body string) error
}

type NotificationService struct {
	email    Notifier
	telegram Notifier
	events   chan models.ApplicationEvent
}

func NewNotificationService(email, telegram Notifier) *NotificationService {
	ns := &NotificationService{
		email:    email,
		telegram: telegram,
		events:   make(chan models.ApplicationEvent, 128),
	}
	go ns.worker()
	return ns
}

func (ns *NotificationService) Send(ev models.ApplicationEvent) {
	select {
	case ns.events <- ev:
	default:
		slog.Warn("notification queue full, dropping event",
			"app_id", ev.ApplicationID)
	}
}

func (ns *NotificationService) worker() {
	for ev := range ns.events {
		subject, body := ns.buildMessage(ev)
		if ns.email != nil {
			if err := ns.email.Send("", subject, body); err != nil {
				slog.Error("email notification failed",
					"err", err, "app_id", ev.ApplicationID)
			}
		}
		if ns.telegram != nil {
			_ = ns.telegram.Send("", subject, body)
		}
	}
}

func (ns *NotificationService) buildMessage(ev models.ApplicationEvent) (string, string) {
	comment := ""
	if ev.Comment != nil {
		comment = *ev.Comment
	}
	subject := fmt.Sprintf("[Account Service] Заявка #%d → %s",
		ev.ApplicationID, ev.ToState)
	body := fmt.Sprintf(
		"Заявка #%d изменила статус\nС: %s\nНа: %s\nСобытие: %s\nКомментарий: %s",
		ev.ApplicationID, ev.FromState, ev.ToState, ev.Event, comment,
	)
	return subject, body
}
