package activity

import (
	"context"
	"time"

	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
)

const timeLayout = time.RFC3339

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// --- User-facing ---

func (s *Service) ListActivities(ctx context.Context, actType, status *string, offset, limit int) ([]ActivityResp, int64, error) {
	activities, total, err := s.repo.List(ctx, actType, status, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]ActivityResp, 0, len(activities))
	for i := range activities {
		result = append(result, ToActivityResp(&activities[i]))
	}
	return result, total, nil
}

func (s *Service) GetActivity(ctx context.Context, id int64) (*ActivityResp, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrActivityNotFound
	}
	resp := ToActivityResp(a)
	return &resp, nil
}

func (s *Service) Register(ctx context.Context, activityID, userID int64) (*RegistrationResp, error) {
	a, err := s.repo.GetByID(ctx, activityID)
	if err != nil {
		return nil, apperrors.ErrActivityNotFound
	}

	// Check registration is open
	if a.Status != "registration" && a.Status != "published" {
		return nil, apperrors.ErrActivityClosed
	}
	if time.Now().After(a.RegistrationDeadline) {
		return nil, apperrors.ErrActivityClosed
	}

	// Check capacity
	if a.MaxParticipants > 0 && a.CurrentParticipants >= a.MaxParticipants {
		return nil, apperrors.ErrActivityFull
	}

	// Check duplicate registration
	existing, _ := s.repo.GetRegistration(ctx, activityID, userID)
	if existing != nil {
		return nil, apperrors.ErrAlreadyRegistered
	}

	reg := &ActivityRegistration{
		ActivityID: activityID,
		UserID:     userID,
		Status:     "registered",
	}

	if err := s.repo.CreateRegistration(ctx, reg); err != nil {
		return nil, apperrors.ErrInternal
	}

	if err := s.repo.IncrementParticipants(ctx, activityID, 1); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToRegistrationResp(reg)
	return &resp, nil
}

func (s *Service) CancelRegistration(ctx context.Context, activityID, userID int64) error {
	if _, err := s.repo.GetByID(ctx, activityID); err != nil {
		return apperrors.ErrActivityNotFound
	}

	reg, err := s.repo.GetRegistration(ctx, activityID, userID)
	if err != nil {
		return apperrors.ErrRecordNotFound
	}

	reg.Status = "cancelled"
	if err := s.repo.UpdateRegistration(ctx, reg); err != nil {
		return apperrors.ErrInternal
	}

	if err := s.repo.IncrementParticipants(ctx, activityID, -1); err != nil {
		return apperrors.ErrInternal
	}

	return nil
}

// --- Admin ---

func (s *Service) AdminListActivities(ctx context.Context, actType, status *string, offset, limit int) ([]ActivityResp, int64, error) {
	return s.ListActivities(ctx, actType, status, offset, limit)
}

func (s *Service) AdminCreateActivity(ctx context.Context, req CreateActivityReq) (*ActivityResp, error) {
	price, err := decimal.NewFromString(req.Price)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	startAt, err := time.Parse(timeLayout, req.StartAt)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}
	endAt, err := time.Parse(timeLayout, req.EndAt)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}
	regDeadline, err := time.Parse(timeLayout, req.RegistrationDeadline)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	levelReq := "all"
	if req.LevelRequirement != "" {
		levelReq = req.LevelRequirement
	}

	a := &Activity{
		Title:                req.Title,
		Type:                 req.Type,
		Description:          req.Description,
		CoverImage:           req.CoverImage,
		VenueID:              req.VenueID,
		StartAt:              startAt,
		EndAt:                endAt,
		RegistrationDeadline: regDeadline,
		MinParticipants:      req.MinParticipants,
		MaxParticipants:      req.MaxParticipants,
		Price:                price,
		LevelRequirement:     levelReq,
		Status:               "draft",
	}

	if err := s.repo.Create(ctx, a); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToActivityResp(a)
	return &resp, nil
}

func (s *Service) AdminUpdateActivity(ctx context.Context, id int64, req UpdateActivityReq) (*ActivityResp, error) {
	a, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrActivityNotFound
	}

	if req.Title != nil {
		a.Title = *req.Title
	}
	if req.Type != nil {
		a.Type = *req.Type
	}
	if req.Description != nil {
		a.Description = *req.Description
	}
	if req.CoverImage != nil {
		a.CoverImage = *req.CoverImage
	}
	if req.VenueID != nil {
		a.VenueID = req.VenueID
	}
	if req.StartAt != nil {
		t, pErr := time.Parse(timeLayout, *req.StartAt)
		if pErr != nil {
			return nil, apperrors.ErrInvalidParams
		}
		a.StartAt = t
	}
	if req.EndAt != nil {
		t, pErr := time.Parse(timeLayout, *req.EndAt)
		if pErr != nil {
			return nil, apperrors.ErrInvalidParams
		}
		a.EndAt = t
	}
	if req.RegistrationDeadline != nil {
		t, pErr := time.Parse(timeLayout, *req.RegistrationDeadline)
		if pErr != nil {
			return nil, apperrors.ErrInvalidParams
		}
		a.RegistrationDeadline = t
	}
	if req.MinParticipants != nil {
		a.MinParticipants = *req.MinParticipants
	}
	if req.MaxParticipants != nil {
		a.MaxParticipants = *req.MaxParticipants
	}
	if req.Price != nil {
		p, pErr := decimal.NewFromString(*req.Price)
		if pErr != nil {
			return nil, apperrors.ErrInvalidParams
		}
		a.Price = p
	}
	if req.LevelRequirement != nil {
		a.LevelRequirement = *req.LevelRequirement
	}
	if req.Status != nil {
		a.Status = *req.Status
	}

	if err := s.repo.Update(ctx, a); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToActivityResp(a)
	return &resp, nil
}

func (s *Service) AdminDeleteActivity(ctx context.Context, id int64) error {
	if _, err := s.repo.GetByID(ctx, id); err != nil {
		return apperrors.ErrActivityNotFound
	}
	return s.repo.Delete(ctx, id)
}
