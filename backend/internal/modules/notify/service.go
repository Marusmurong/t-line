package notify

import (
	"context"
	"time"

	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"gorm.io/datatypes"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// CreateNotification creates a notification for a user.
// This is intended to be called by other modules.
func (s *Service) CreateNotification(ctx context.Context, userID int64, notifType, title, content string, extraData datatypes.JSON) error {
	n := &Notification{
		UserID:    userID,
		Type:      notifType,
		Title:     title,
		Content:   content,
		IsRead:    false,
		ExtraData: extraData,
		CreatedAt: time.Now(),
	}

	if err := s.repo.Create(ctx, n); err != nil {
		return apperrors.ErrInternal
	}
	return nil
}

// ListNotifications returns paginated notifications for a user.
func (s *Service) ListNotifications(ctx context.Context, userID int64, offset, limit int) ([]NotificationResp, int64, error) {
	notifications, total, err := s.repo.ListByUser(ctx, userID, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]NotificationResp, 0, len(notifications))
	for i := range notifications {
		result = append(result, ToNotificationResp(&notifications[i]))
	}
	return result, total, nil
}

// MarkRead marks a single notification as read.
func (s *Service) MarkRead(ctx context.Context, id, userID int64) error {
	if err := s.repo.MarkRead(ctx, id, userID); err != nil {
		return apperrors.ErrInternal
	}
	return nil
}

// MarkAllRead marks all unread notifications as read for a user.
func (s *Service) MarkAllRead(ctx context.Context, userID int64) error {
	if err := s.repo.MarkAllRead(ctx, userID); err != nil {
		return apperrors.ErrInternal
	}
	return nil
}

// UnreadCount returns the number of unread notifications.
func (s *Service) UnreadCount(ctx context.Context, userID int64) (int64, error) {
	count, err := s.repo.CountUnread(ctx, userID)
	if err != nil {
		return 0, apperrors.ErrInternal
	}
	return count, nil
}

// --- Response DTO ---

type NotificationResp struct {
	ID        int64          `json:"id"`
	UserID    int64          `json:"user_id"`
	Type      string         `json:"type"`
	Title     string         `json:"title"`
	Content   string         `json:"content"`
	IsRead    bool           `json:"is_read"`
	ExtraData datatypes.JSON `json:"extra_data"`
	CreatedAt string         `json:"created_at"`
}

func ToNotificationResp(n *Notification) NotificationResp {
	return NotificationResp{
		ID:        n.ID,
		UserID:    n.UserID,
		Type:      n.Type,
		Title:     n.Title,
		Content:   n.Content,
		IsRead:    n.IsRead,
		ExtraData: n.ExtraData,
		CreatedAt: n.CreatedAt.Format("2006-01-02T15:04:05Z07:00"),
	}
}
