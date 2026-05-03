package grpcserver

import (
	"context"
	"encoding/json"

	"github.com/novelhive/notification-service/internal/domain"
	notificationv1 "github.com/novelhive/proto/notification/v1"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type NotificationServiceServer struct {
	notificationv1.UnimplementedNotificationServiceServer
	repo   domain.NotificationRepository
	logger *zap.Logger
}

func NewNotificationServiceServer(repo domain.NotificationRepository, logger *zap.Logger) *NotificationServiceServer {
	return &NotificationServiceServer{
		repo:   repo,
		logger: logger,
	}
}

func (s *NotificationServiceServer) GetNotifications(ctx context.Context, req *notificationv1.GetNotificationsRequest) (*notificationv1.GetNotificationsResponse, error) {
	page := int(req.Page)
	if page < 1 {
		page = 1
	}
	pageSize := int(req.PageSize)
	if pageSize < 1 {
		pageSize = 20
	}

	notifs, total, err := s.repo.ListByUser(req.UserId, page, pageSize)
	if err != nil {
		s.logger.Error("failed to list notifications", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to list notifications")
	}

	var pbNotifs []*notificationv1.Notification
	for _, n := range notifs {
		payload, _ := json.Marshal(n.Payload)
		pbNotifs = append(pbNotifs, &notificationv1.Notification{
			Id:        n.ID,
			UserId:    n.UserID,
			Type:      n.Type,
			Title:     n.Title,
			Body:      n.Body,
			Payload:   string(payload),
			IsRead:    n.IsRead,
			CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z"),
		})
	}

	return &notificationv1.GetNotificationsResponse{
		Notifications: pbNotifs,
		Total:         int32(total),
	}, nil
}

func (s *NotificationServiceServer) MarkAsRead(ctx context.Context, req *notificationv1.MarkAsReadRequest) (*notificationv1.MarkAsReadResponse, error) {
	if err := s.repo.MarkAsRead(req.Id, req.UserId); err != nil {
		s.logger.Error("failed to mark notification as read", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to mark as read")
	}
	return &notificationv1.MarkAsReadResponse{Success: true}, nil
}

func (s *NotificationServiceServer) RegisterFCMToken(ctx context.Context, req *notificationv1.RegisterFCMTokenRequest) (*notificationv1.RegisterFCMTokenResponse, error) {
	if err := s.repo.SaveFCMToken(req.UserId, req.Token); err != nil {
		s.logger.Error("failed to register FCM token", zap.Error(err))
		return nil, status.Error(codes.Internal, "failed to register token")
	}
	return &notificationv1.RegisterFCMTokenResponse{Success: true}, nil
}
