package events

import (
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"time"

	"github.com/nats-io/nats.go"
	"github.com/novelhive/notification-service/internal/domain"
	"github.com/novelhive/notification-service/internal/push"
	libraryv1 "github.com/novelhive/proto/library/v1"
	"go.uber.org/zap"
)

type ChapterPublishedEvent struct {
	NovelID    string `json:"novel_id"`
	NovelTitle string `json:"novel_title"`
	NovelSlug  string `json:"novel_slug"`
	ChapterID  string `json:"chapter_id"`
	Number     int    `json:"number"`
	Title      string `json:"title"`
}

type NATSSubscriber struct {
	nc             *nats.Conn
	js             nats.JetStreamContext
	repo           domain.NotificationRepository
	libraryClient  libraryv1.LibraryServiceClient
	pusher         *push.FirebasePusher
	logger         *zap.Logger
}

func NewNATSSubscriber(
	nc *nats.Conn,
	repo domain.NotificationRepository,
	libraryClient libraryv1.LibraryServiceClient,
	pusher *push.FirebasePusher,
	logger *zap.Logger,
) (*NATSSubscriber, error) {
	js, err := nc.JetStream()
	if err != nil {
		return nil, err
	}

	return &NATSSubscriber{
		nc:            nc,
		js:            js,
		repo:          repo,
		libraryClient: libraryClient,
		pusher:        pusher,
		logger:        logger,
	}, nil
}

func (s *NATSSubscriber) Start() error {
	_, err := s.js.Subscribe("novel.chapter.published", func(m *nats.Msg) {
		var event ChapterPublishedEvent
		if err := json.Unmarshal(m.Data, &event); err != nil {
			s.logger.Error("failed to unmarshal chapter event", zap.Error(err))
			m.Nak()
			return
		}

		s.logger.Info("received chapter published event",
			zap.String("novel_id", event.NovelID),
			zap.Int("chapter", event.Number),
		)

		if err := s.processChapterPublished(event); err != nil {
			s.logger.Error("failed to process chapter published", zap.Error(err))
			m.Nak()
			return
		}

		m.Ack()
	}, nats.Durable("notification-service-chapter"), nats.ManualAck())

	return err
}

func (s *NATSSubscriber) processChapterPublished(event ChapterPublishedEvent) error {
	// 1. Get users who have this novel in library
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := s.libraryClient.GetUsersByNovel(ctx, &libraryv1.GetUsersByNovelRequest{
		NovelId: event.NovelID,
	})
	if err != nil {
		return fmt.Errorf("failed to get users from library service: %w", err)
	}

	s.logger.Info("sending notifications to users",
		zap.Int("user_count", len(resp.UserIds)),
		zap.String("novel_id", event.NovelID),
	)

	// 2. Create notification for each user
	now := time.Now()
	for _, userID := range resp.UserIds {
		notif := &domain.Notification{
			ID:     genID(),
			UserID: userID,
			Type:   "new_chapter",
			Title:  fmt.Sprintf("Chapter Baru: %s", event.NovelTitle),
			Body:   fmt.Sprintf("Chapter %d - %s sudah tersedia!", event.Number, event.Title),
			Payload: map[string]interface{}{
				"novel_id":       event.NovelID,
				"chapter_id":     event.ChapterID,
				"chapter_number": event.Number,
			},
			IsRead:    false,
			CreatedAt: now,
		}

		if err := s.repo.Create(notif); err != nil {
			s.logger.Error("failed to create notification in db",
				zap.String("user_id", userID),
				zap.Error(err),
			)
			continue
		}

		// 3. Send Push Notification via Firebase
		tokens, err := s.repo.GetFCMTokens(userID)
		if err == nil && len(tokens) > 0 {
			pushData := map[string]string{
				"novel_id":       event.NovelID,
				"chapter_id":     event.ChapterID,
				"chapter_number": fmt.Sprintf("%d", event.Number),
			}
			err = s.pusher.SendPushNotification(tokens, notif.Title, notif.Body, pushData)
			if err != nil {
				s.logger.Error("failed to send push notification",
					zap.String("user_id", userID),
					zap.Error(err),
				)
			}
		}
	}

	return nil
}

func genID() string {
	b := make([]byte, 16)
	rand.Read(b)
	return hex.EncodeToString(b)
}
