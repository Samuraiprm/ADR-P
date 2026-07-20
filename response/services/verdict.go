package services

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/adr-p/response/metrics"
	"github.com/adr-p/response/redis"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	VERDICT_STREAM = "abuse:verdicts"
	GROUP_NAME     = "response"
	CONSUMER_NAME  = "responder-1"
)

type VerdictService struct {
	redis   *redis.Client
	db      *sql.DB
	running bool
}

func NewVerdictService(redis *redis.Client, db *sql.DB) *VerdictService {
	return &VerdictService{
		redis: redis,
		db:    db,
	}
}

func (s *VerdictService) StartConsuming(ctx context.Context, responseService *ResponseService) {
	s.running = true

	_ = s.redis.XGroupCreate(ctx, VERDICT_STREAM, GROUP_NAME, "0")

	for s.running {
		select {
		case <-ctx.Done():
			return
		default:
			s.processMessages(ctx, responseService)
			if s.running {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}
}

func (s *VerdictService) processMessages(ctx context.Context, responseService *ResponseService) {
	streams, err := s.redis.XReadGroup(ctx, GROUP_NAME, CONSUMER_NAME, map[string]string{VERDICT_STREAM: ">"}, 10, 5*time.Second)
	if err != nil {
		return
	}

	for _, stream := range streams {
		for _, msg := range stream.Messages {
			metrics.VerdictsConsumed.Inc()
			s.handleVerdict(ctx, msg.Values, responseService)
			_ = s.redis.XAck(ctx, VERDICT_STREAM, GROUP_NAME, msg.ID)
		}
	}
}

func (s *VerdictService) handleVerdict(ctx context.Context, values map[string]interface{}, responseService *ResponseService) {
	eventID, _ := values["event_id"].(string)
	verdict, _ := values["verdict"].(string)
	score, _ := values["score"].(string)

	if eventID == "" || verdict == "" {
		log.Printf("Skipping verdict with missing fields: event_id=%q verdict=%q", eventID, verdict)
		metrics.VerdictsFailed.Inc()
		return
	}

	log.Printf("Processing verdict: event=%s verdict=%s score=%s", eventID, verdict, score)

	switch verdict {
	case "BLOCK":
		responseService.BlockUser(ctx, eventID)
		metrics.VerdictsProcessed.With(prometheus.Labels{"action": "BLOCK"}).Inc()
	case "WARN":
		responseService.SendWarning(ctx, eventID, score)
		metrics.VerdictsProcessed.With(prometheus.Labels{"action": "WARN"}).Inc()
	}
}

func (s *VerdictService) Stop() {
	s.running = false
}
