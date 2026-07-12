package services

import (
	"bytes"
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/adr-p/response/db"
)

type ResponseService struct {
	db        *sql.DB
	tgToken   string
	tgChatID  string
	httpClient *http.Client
}

func NewResponseService(db *sql.DB) *ResponseService {
	return &ResponseService{
		db:       db,
		tgToken:  os.Getenv("TELEGRAM_BOT_TOKEN"),
		tgChatID: os.Getenv("TELEGRAM_CHAT_ID"),
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (s *ResponseService) BlockUser(ctx context.Context, eventID string) {
	log.Printf("BLOCK action for event %s", eventID)

	err := db.UpdateEventVerdict(s.db, eventID, "BLOCKED")
	if err != nil {
		log.Printf("Failed to update verdict: %v", err)
		return
	}

	if s.tgToken != "" && s.tgChatID != "" {
		text := fmt.Sprintf("🚫 BLOCKED event %s\nUser has been blocked.", eventID)
		s.sendTelegramMessage(text)
	}

	log.Printf("User blocked for event %s", eventID)
}

func (s *ResponseService) SendWarning(ctx context.Context, eventID string, score string) {
	log.Printf("WARN action for event %s (score: %s)", eventID, score)

	err := db.UpdateEventVerdict(s.db, eventID, "WARNED")
	if err != nil {
		log.Printf("Failed to update verdict: %v", err)
		return
	}

	if s.tgToken != "" && s.tgChatID != "" {
		text := fmt.Sprintf("⚠️ WARN event %s (score: %s)\nPlease review.", eventID, score)
		keyboard := map[string]interface{}{
			"inline_keyboard": [][]map[string]string{
				{
					{"text": "Confirm", "callback_data": fmt.Sprintf("confirm:%s", eventID)},
					{"text": "Reject", "callback_data": fmt.Sprintf("reject:%s", eventID)},
				},
			},
		}
		s.sendTelegramMessageWithKeyboard(text, keyboard)
	}

	log.Printf("Warning sent for event %s", eventID)
}

func (s *ResponseService) ConfirmAction(ctx context.Context, eventID string) error {
	return db.UpdateEventVerdict(s.db, eventID, "CONFIRMED")
}

func (s *ResponseService) RejectAction(ctx context.Context, eventID string) error {
	return db.UpdateEventVerdict(s.db, eventID, "REJECTED")
}

func (s *ResponseService) sendTelegramMessage(text string) {
	if s.tgToken == "" || s.tgChatID == "" {
		return
	}

	payload := map[string]interface{}{
		"chat_id": s.tgChatID,
		"text":    text,
	}
	s.postTelegram("sendMessage", payload)
}

func (s *ResponseService) sendTelegramMessageWithKeyboard(text string, keyboard map[string]interface{}) {
	if s.tgToken == "" || s.tgChatID == "" {
		return
	}

	payload := map[string]interface{}{
		"chat_id":      s.tgChatID,
		"text":         text,
		"reply_markup": keyboard,
	}
	s.postTelegram("sendMessage", payload)
}

func (s *ResponseService) postTelegram(method string, payload map[string]interface{}) {
	url := fmt.Sprintf("https://api.telegram.org/bot%s/%s", s.tgToken, method)

	body, err := json.Marshal(payload)
	if err != nil {
		log.Printf("Failed to marshal telegram payload: %v", err)
		return
	}

	resp, err := s.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		log.Printf("Telegram API error: %v", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		respBody, _ := io.ReadAll(resp.Body)
		log.Printf("Telegram API returned %d: %s", resp.StatusCode, string(respBody))
	}
}
