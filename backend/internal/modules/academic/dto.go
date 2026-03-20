package academic

import (
	"encoding/json"
	"time"
)

// --- Coach DTOs ---

type CreateCoachReq struct {
	UserID      *int64 `json:"user_id"`
	Title       string `json:"title" binding:"max=64"`
	Specialties []string `json:"specialties"`
	Bio         string `json:"bio"`
	HourlyRate  string `json:"hourly_rate" binding:"required"`
}

type UpdateCoachReq struct {
	Title       *string  `json:"title" binding:"omitempty,max=64"`
	Specialties []string `json:"specialties"`
	Bio         *string  `json:"bio"`
	HourlyRate  *string  `json:"hourly_rate"`
	Status      *int     `json:"status"`
}

type CoachResp struct {
	ID           int64    `json:"id"`
	UserID       *int64   `json:"user_id"`
	Title        string   `json:"title"`
	Specialties  []string `json:"specialties"`
	Bio          string   `json:"bio"`
	HourlyRate   string   `json:"hourly_rate"`
	Rating       string   `json:"rating"`
	StudentCount int      `json:"student_count"`
	Status       int      `json:"status"`
	CreatedAt    string   `json:"created_at"`
	UpdatedAt    string   `json:"updated_at"`
}

type CoachPerformanceResp struct {
	CoachID      int64  `json:"coach_id"`
	TotalCourses int64  `json:"total_courses"`
	AvgRating    string `json:"avg_rating"`
	StudentCount int    `json:"student_count"`
	TotalIncome  string `json:"total_income"`
}

// --- Schedule DTOs ---

type CreateScheduleReq struct {
	CoachID        int64  `json:"coach_id" binding:"required"`
	VenueID        int64  `json:"venue_id" binding:"required"`
	StudentID      *int64 `json:"student_id"`
	ProductID      *int64 `json:"product_id"`
	Date           string `json:"date" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	EndTime        string `json:"end_time" binding:"required"`
	Type           string `json:"type" binding:"omitempty,oneof=private group"`
	RecurrenceRule string `json:"recurrence_rule"`
}

type UpdateScheduleReq struct {
	VenueID   *int64  `json:"venue_id"`
	StudentID *int64  `json:"student_id"`
	Date      *string `json:"date"`
	StartTime *string `json:"start_time"`
	EndTime   *string `json:"end_time"`
	Type      *string `json:"type" binding:"omitempty,oneof=private group"`
	Status    *string `json:"status" binding:"omitempty,oneof=scheduled completed cancelled"`
}

type BatchScheduleReq struct {
	CoachID        int64  `json:"coach_id" binding:"required"`
	VenueID        int64  `json:"venue_id" binding:"required"`
	StudentID      *int64 `json:"student_id"`
	ProductID      *int64 `json:"product_id"`
	Date           string `json:"date" binding:"required"`
	StartTime      string `json:"start_time" binding:"required"`
	EndTime        string `json:"end_time" binding:"required"`
	Type           string `json:"type" binding:"omitempty,oneof=private group"`
	RecurrenceRule string `json:"recurrence_rule" binding:"required"`
	Weeks          int    `json:"weeks" binding:"required,min=1,max=52"`
}

type ConflictCheckReq struct {
	CoachID   int64  `json:"coach_id" binding:"required"`
	VenueID   int64  `json:"venue_id" binding:"required"`
	Date      string `json:"date" binding:"required"`
	StartTime string `json:"start_time" binding:"required"`
	EndTime   string `json:"end_time" binding:"required"`
}

type SubstituteReq struct {
	SubstituteCoachID int64 `json:"substitute_coach_id" binding:"required"`
}

type ScheduleResp struct {
	ID                int64  `json:"id"`
	CoachID           int64  `json:"coach_id"`
	VenueID           int64  `json:"venue_id"`
	StudentID         *int64 `json:"student_id"`
	ProductID         *int64 `json:"product_id"`
	Date              string `json:"date"`
	StartTime         string `json:"start_time"`
	EndTime           string `json:"end_time"`
	Type              string `json:"type"`
	Status            string `json:"status"`
	RecurrenceRule    string `json:"recurrence_rule"`
	RecurrenceGroupID string `json:"recurrence_group_id"`
	SubstituteCoachID *int64 `json:"substitute_coach_id"`
	CreatedAt         string `json:"created_at"`
	UpdatedAt         string `json:"updated_at"`
}

// --- Leave DTOs ---

type CreateLeaveReq struct {
	CoachID   int64  `json:"coach_id" binding:"required"`
	StartDate string `json:"start_date" binding:"required"`
	EndDate   string `json:"end_date" binding:"required"`
	Reason    string `json:"reason" binding:"max=256"`
}

type UpdateLeaveReq struct {
	StartDate *string `json:"start_date"`
	EndDate   *string `json:"end_date"`
	Reason    *string `json:"reason" binding:"omitempty,max=256"`
	Status    *string `json:"status" binding:"omitempty,oneof=pending approved rejected"`
}

type LeaveResp struct {
	ID        int64  `json:"id"`
	CoachID   int64  `json:"coach_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
	Reason    string `json:"reason"`
	Status    string `json:"status"`
	CreatedAt string `json:"created_at"`
}

// --- Student DTOs ---

type StudentListItem struct {
	UserID       int64  `json:"user_id"`
	TotalCourses int64  `json:"total_courses"`
	LastCourseAt string `json:"last_course_at"`
}

type RecordResp struct {
	ID            int64  `json:"id"`
	ScheduleID    int64  `json:"schedule_id"`
	StudentID     int64  `json:"student_id"`
	CoachID       int64  `json:"coach_id"`
	Attendance    string `json:"attendance"`
	CoachFeedback string `json:"coach_feedback"`
	Rating        *int   `json:"rating"`
	RatingComment string `json:"rating_comment"`
	CreatedAt     string `json:"created_at"`
}

type FeedbackReq struct {
	CoachFeedback string `json:"coach_feedback" binding:"required"`
	Attendance    string `json:"attendance" binding:"omitempty,oneof=present absent late"`
}

type RatingReq struct {
	Rating        int    `json:"rating" binding:"required,min=1,max=5"`
	RatingComment string `json:"rating_comment"`
}

// --- Converters ---

func ToCoachResp(c *Coach) CoachResp {
	var specialties []string
	if c.Specialties != nil {
		_ = json.Unmarshal(c.Specialties, &specialties)
	}
	if specialties == nil {
		specialties = []string{}
	}

	return CoachResp{
		ID:           c.ID,
		UserID:       c.UserID,
		Title:        c.Title,
		Specialties:  specialties,
		Bio:          c.Bio,
		HourlyRate:   c.HourlyRate.StringFixed(2),
		Rating:       c.Rating.StringFixed(2),
		StudentCount: c.StudentCount,
		Status:       c.Status,
		CreatedAt:    c.CreatedAt.Format(time.RFC3339),
		UpdatedAt:    c.UpdatedAt.Format(time.RFC3339),
	}
}

func ToScheduleResp(s *CourseSchedule) ScheduleResp {
	groupID := ""
	if s.RecurrenceGroupID != nil {
		groupID = s.RecurrenceGroupID.String()
	}

	return ScheduleResp{
		ID:                s.ID,
		CoachID:           s.CoachID,
		VenueID:           s.VenueID,
		StudentID:         s.StudentID,
		ProductID:         s.ProductID,
		Date:              s.Date,
		StartTime:         s.StartTime,
		EndTime:           s.EndTime,
		Type:              s.Type,
		Status:            s.Status,
		RecurrenceRule:    s.RecurrenceRule,
		RecurrenceGroupID: groupID,
		SubstituteCoachID: s.SubstituteCoachID,
		CreatedAt:         s.CreatedAt.Format(time.RFC3339),
		UpdatedAt:         s.UpdatedAt.Format(time.RFC3339),
	}
}

func ToLeaveResp(l *CoachLeave) LeaveResp {
	return LeaveResp{
		ID:        l.ID,
		CoachID:   l.CoachID,
		StartDate: l.StartDate,
		EndDate:   l.EndDate,
		Reason:    l.Reason,
		Status:    l.Status,
		CreatedAt: l.CreatedAt.Format(time.RFC3339),
	}
}

func ToRecordResp(r *StudentCourseRecord) RecordResp {
	return RecordResp{
		ID:            r.ID,
		ScheduleID:    r.ScheduleID,
		StudentID:     r.StudentID,
		CoachID:       r.CoachID,
		Attendance:    r.Attendance,
		CoachFeedback: r.CoachFeedback,
		Rating:        r.Rating,
		RatingComment: r.RatingComment,
		CreatedAt:     r.CreatedAt.Format(time.RFC3339),
	}
}
