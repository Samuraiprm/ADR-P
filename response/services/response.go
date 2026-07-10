package services

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"

	"github.com/adr-p/response/db"
)

type ResponseService struct {
	db *sql.DB
}

func NewResponseService(db *sql.DB) *ResponseService {
	return &ResponseService{db: db}
}

func (s *ResponseService) BlockUser(ctx context.Context, eventID string) {
	log.Printf("BLOCK action for event %s", eventID)

	err := db.UpdateEventVerdict(s.db, eventID, "BLOCKED")
	if err != nil {
		log.Printf("Failed to update verdict: %v", err)
		return
	}

	// TODO: Call real TG/VK ban API
	// For now, just log the action
	log.Printf("User blocked for event %s", eventID)
}

func (s *ResponseService) SendWarning(ctx context.Context, eventID string, score string) {
	log.Printf("WARN action for event %s (score: %s)", eventID, score)

	// Create inline keyboard for Telegram
	keyboard := map[string]interface{}{
		"inline_keyboard": [][]map[string]string{
			{
				{"text": "Подтвердить", "callback_data": fmt.Sprintf("confirm:%s", eventID)},
				{"text": "Отклонить", "callback_data": fmt.Sprintf("reject:%s", eventID)},
			},
		},
	}

	keyboardJSON, _ := json.Marshal(keyboard)

	// TODO: Send actual Telegram message with keyboard
	log.Printf("Would send warning with keyboard: %s", string(keyboardJSON))
}

func (s *ResponseService) ConfirmAction(ctx context.Context, eventID string) error {
	return db.UpdateEventVerdict(s.db, eventID, "CONFIRMED")
}

func (s *ResponseService) RejectAction(ctx context.Context, eventID string) error {
	return db.UpdateEventVerdict(s.db, eventID, "REJECTED")
}
