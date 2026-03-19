package activity

import "time"

// --- Request DTOs ---

type CreateActivityReq struct {
	Title                string `json:"title" binding:"required,max=128"`
	Type                 string `json:"type" binding:"required,oneof=open_play group_class tournament themed"`
	Description          string `json:"description"`
	CoverImage           string `json:"cover_image" binding:"omitempty,max=512"`
	VenueID              *int64 `json:"venue_id"`
	StartAt              string `json:"start_at" binding:"required"`
	EndAt                string `json:"end_at" binding:"required"`
	RegistrationDeadline string `json:"registration_deadline" binding:"required"`
	MinParticipants      int    `json:"min_participants" binding:"min=0"`
	MaxParticipants      int    `json:"max_participants" binding:"required,min=1"`
	Price                string `json:"price" binding:"required"`
	LevelRequirement     string `json:"level_requirement" binding:"omitempty,oneof=beginner intermediate advanced all"`
}

type UpdateActivityReq struct {
	Title                *string `json:"title" binding:"omitempty,max=128"`
	Type                 *string `json:"type" binding:"omitempty,oneof=open_play group_class tournament themed"`
	Description          *string `json:"description"`
	CoverImage           *string `json:"cover_image" binding:"omitempty,max=512"`
	VenueID              *int64  `json:"venue_id"`
	StartAt              *string `json:"start_at"`
	EndAt                *string `json:"end_at"`
	RegistrationDeadline *string `json:"registration_deadline"`
	MinParticipants      *int    `json:"min_participants"`
	MaxParticipants      *int    `json:"max_participants"`
	Price                *string `json:"price"`
	LevelRequirement     *string `json:"level_requirement" binding:"omitempty,oneof=beginner intermediate advanced all"`
	Status               *string `json:"status" binding:"omitempty,oneof=draft published registration confirmed ongoing completed cancelled"`
}

// --- Response DTOs ---

type ActivityResp struct {
	ID                   int64  `json:"id"`
	Title                string `json:"title"`
	Type                 string `json:"type"`
	Description          string `json:"description"`
	CoverImage           string `json:"cover_image"`
	VenueID              *int64 `json:"venue_id"`
	StartAt              string `json:"start_at"`
	EndAt                string `json:"end_at"`
	RegistrationDeadline string `json:"registration_deadline"`
	MinParticipants      int    `json:"min_participants"`
	MaxParticipants      int    `json:"max_participants"`
	CurrentParticipants  int    `json:"current_participants"`
	Price                string `json:"price"`
	LevelRequirement     string `json:"level_requirement"`
	Status               string `json:"status"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
}

type RegistrationResp struct {
	ID         int64  `json:"id"`
	ActivityID int64  `json:"activity_id"`
	UserID     int64  `json:"user_id"`
	Status     string `json:"status"`
	CreatedAt  string `json:"created_at"`
}

func ToActivityResp(a *Activity) ActivityResp {
	return ActivityResp{
		ID:                   a.ID,
		Title:                a.Title,
		Type:                 a.Type,
		Description:          a.Description,
		CoverImage:           a.CoverImage,
		VenueID:              a.VenueID,
		StartAt:              a.StartAt.Format(time.RFC3339),
		EndAt:                a.EndAt.Format(time.RFC3339),
		RegistrationDeadline: a.RegistrationDeadline.Format(time.RFC3339),
		MinParticipants:      a.MinParticipants,
		MaxParticipants:      a.MaxParticipants,
		CurrentParticipants:  a.CurrentParticipants,
		Price:                a.Price.StringFixed(2),
		LevelRequirement:     a.LevelRequirement,
		Status:               a.Status,
		CreatedAt:            a.CreatedAt.Format(time.RFC3339),
		UpdatedAt:            a.UpdatedAt.Format(time.RFC3339),
	}
}

func ToRegistrationResp(r *ActivityRegistration) RegistrationResp {
	return RegistrationResp{
		ID:         r.ID,
		ActivityID: r.ActivityID,
		UserID:     r.UserID,
		Status:     r.Status,
		CreatedAt:  r.CreatedAt.Format(time.RFC3339),
	}
}
