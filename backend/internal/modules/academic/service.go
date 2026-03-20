package academic

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	apperrors "github.com/t-line/backend/internal/pkg/errors"
	"gorm.io/datatypes"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

// ============================================================
// Coach management
// ============================================================

func (s *Service) ListCoaches(ctx context.Context, status *int, offset, limit int) ([]CoachResp, int64, error) {
	coaches, total, err := s.repo.ListCoaches(ctx, status, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]CoachResp, 0, len(coaches))
	for i := range coaches {
		result = append(result, ToCoachResp(&coaches[i]))
	}
	return result, total, nil
}

func (s *Service) GetCoach(ctx context.Context, id int64) (*CoachResp, error) {
	c, err := s.repo.GetCoachByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}
	resp := ToCoachResp(c)
	return &resp, nil
}

func (s *Service) CreateCoach(ctx context.Context, req CreateCoachReq) (*CoachResp, error) {
	hourlyRate, err := decimal.NewFromString(req.HourlyRate)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	specialtiesJSON, err := json.Marshal(req.Specialties)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	c := &Coach{
		UserID:      req.UserID,
		Title:       req.Title,
		Specialties: datatypes.JSON(specialtiesJSON),
		Bio:         req.Bio,
		HourlyRate:  hourlyRate,
		Rating:      decimal.NewFromFloat(5.0),
		Status:      1,
	}

	if err := s.repo.CreateCoach(ctx, c); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToCoachResp(c)
	return &resp, nil
}

func (s *Service) UpdateCoach(ctx context.Context, id int64, req UpdateCoachReq) (*CoachResp, error) {
	c, err := s.repo.GetCoachByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	if req.Title != nil {
		c.Title = *req.Title
	}
	if req.Specialties != nil {
		specJSON, jsonErr := json.Marshal(req.Specialties)
		if jsonErr != nil {
			return nil, apperrors.ErrInvalidParams
		}
		c.Specialties = datatypes.JSON(specJSON)
	}
	if req.Bio != nil {
		c.Bio = *req.Bio
	}
	if req.HourlyRate != nil {
		rate, pErr := decimal.NewFromString(*req.HourlyRate)
		if pErr != nil {
			return nil, apperrors.ErrInvalidParams
		}
		c.HourlyRate = rate
	}
	if req.Status != nil {
		c.Status = *req.Status
	}

	if err := s.repo.UpdateCoach(ctx, c); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToCoachResp(c)
	return &resp, nil
}

func (s *Service) DeleteCoach(ctx context.Context, id int64) error {
	if _, err := s.repo.GetCoachByID(ctx, id); err != nil {
		return apperrors.ErrRecordNotFound
	}
	return s.repo.DeleteCoach(ctx, id)
}

func (s *Service) GetCoachPerformance(ctx context.Context, coachID int64) (*CoachPerformanceResp, error) {
	c, err := s.repo.GetCoachByID(ctx, coachID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	totalCourses, err := s.repo.CountCoachSchedules(ctx, coachID)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	avgRating, err := s.repo.AvgCoachRating(ctx, coachID)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	studentCount, err := s.repo.CountCoachStudents(ctx, coachID)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	totalIncome := c.HourlyRate.Mul(decimal.NewFromInt(totalCourses))

	return &CoachPerformanceResp{
		CoachID:      coachID,
		TotalCourses: totalCourses,
		AvgRating:    fmt.Sprintf("%.2f", avgRating),
		StudentCount: studentCount,
		TotalIncome:  totalIncome.StringFixed(2),
	}, nil
}

// ============================================================
// Schedule management
// ============================================================

func (s *Service) CreateSchedule(ctx context.Context, req CreateScheduleReq) (*ScheduleResp, error) {
	if err := s.CheckConflict(ctx, req.CoachID, req.VenueID, req.Date, req.StartTime, req.EndTime, 0); err != nil {
		return nil, err
	}

	scheduleType := "private"
	if req.Type != "" {
		scheduleType = req.Type
	}

	schedule := &CourseSchedule{
		CoachID:        req.CoachID,
		VenueID:        req.VenueID,
		StudentID:      req.StudentID,
		ProductID:      req.ProductID,
		Date:           req.Date,
		StartTime:      req.StartTime,
		EndTime:        req.EndTime,
		Type:           scheduleType,
		Status:         "scheduled",
		RecurrenceRule: req.RecurrenceRule,
	}

	if err := s.repo.CreateSchedule(ctx, schedule); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToScheduleResp(schedule)
	return &resp, nil
}

func (s *Service) ListSchedules(ctx context.Context, coachID, venueID *int64, date, status *string, offset, limit int) ([]ScheduleResp, int64, error) {
	schedules, total, err := s.repo.ListSchedules(ctx, coachID, venueID, date, status, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]ScheduleResp, 0, len(schedules))
	for i := range schedules {
		result = append(result, ToScheduleResp(&schedules[i]))
	}
	return result, total, nil
}

func (s *Service) GetSchedule(ctx context.Context, id int64) (*ScheduleResp, error) {
	schedule, err := s.repo.GetScheduleByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrCourseNotFound
	}
	resp := ToScheduleResp(schedule)
	return &resp, nil
}

func (s *Service) UpdateSchedule(ctx context.Context, id int64, req UpdateScheduleReq) (*ScheduleResp, error) {
	schedule, err := s.repo.GetScheduleByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrCourseNotFound
	}

	if req.VenueID != nil {
		schedule.VenueID = *req.VenueID
	}
	if req.StudentID != nil {
		schedule.StudentID = req.StudentID
	}
	if req.Date != nil {
		schedule.Date = *req.Date
	}
	if req.StartTime != nil {
		schedule.StartTime = *req.StartTime
	}
	if req.EndTime != nil {
		schedule.EndTime = *req.EndTime
	}
	if req.Type != nil {
		schedule.Type = *req.Type
	}
	if req.Status != nil {
		schedule.Status = *req.Status
	}

	// Re-check conflict if date/time/venue changed
	if req.Date != nil || req.StartTime != nil || req.EndTime != nil || req.VenueID != nil {
		if err := s.CheckConflict(ctx, schedule.CoachID, schedule.VenueID, schedule.Date, schedule.StartTime, schedule.EndTime, schedule.ID); err != nil {
			return nil, err
		}
	}

	if err := s.repo.UpdateSchedule(ctx, schedule); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToScheduleResp(schedule)
	return &resp, nil
}

func (s *Service) CancelSchedule(ctx context.Context, id int64) error {
	schedule, err := s.repo.GetScheduleByID(ctx, id)
	if err != nil {
		return apperrors.ErrCourseNotFound
	}

	schedule.Status = "cancelled"
	return s.repo.UpdateSchedule(ctx, schedule)
}

func (s *Service) SubstituteCoach(ctx context.Context, scheduleID int64, substituteCoachID int64) (*ScheduleResp, error) {
	schedule, err := s.repo.GetScheduleByID(ctx, scheduleID)
	if err != nil {
		return nil, apperrors.ErrCourseNotFound
	}

	// Check substitute coach exists
	if _, err := s.repo.GetCoachByID(ctx, substituteCoachID); err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	// Check substitute coach conflict
	if err := s.CheckConflict(ctx, substituteCoachID, schedule.VenueID, schedule.Date, schedule.StartTime, schedule.EndTime, 0); err != nil {
		return nil, err
	}

	schedule.SubstituteCoachID = &substituteCoachID
	if err := s.repo.UpdateSchedule(ctx, schedule); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToScheduleResp(schedule)
	return &resp, nil
}

// BatchCreateSchedules creates recurring schedules for N weeks.
func (s *Service) BatchCreateSchedules(ctx context.Context, req BatchScheduleReq) ([]ScheduleResp, error) {
	baseDate, err := time.Parse("2006-01-02", req.Date)
	if err != nil {
		return nil, apperrors.ErrInvalidParams
	}

	scheduleType := "private"
	if req.Type != "" {
		scheduleType = req.Type
	}

	groupID := uuid.New()
	results := make([]ScheduleResp, 0, req.Weeks)
	var conflicts []string

	for i := 0; i < req.Weeks; i++ {
		date := baseDate.AddDate(0, 0, i*7)
		dateStr := date.Format("2006-01-02")

		if conflictErr := s.CheckConflict(ctx, req.CoachID, req.VenueID, dateStr, req.StartTime, req.EndTime, 0); conflictErr != nil {
			conflicts = append(conflicts, fmt.Sprintf("第%d周(%s): %s", i+1, dateStr, conflictErr.Error()))
			continue
		}

		schedule := &CourseSchedule{
			CoachID:           req.CoachID,
			VenueID:           req.VenueID,
			StudentID:         req.StudentID,
			ProductID:         req.ProductID,
			Date:              dateStr,
			StartTime:         req.StartTime,
			EndTime:           req.EndTime,
			Type:              scheduleType,
			Status:            "scheduled",
			RecurrenceRule:    req.RecurrenceRule,
			RecurrenceGroupID: &groupID,
		}

		if createErr := s.repo.CreateSchedule(ctx, schedule); createErr != nil {
			conflicts = append(conflicts, fmt.Sprintf("第%d周(%s): 创建失败", i+1, dateStr))
			continue
		}

		results = append(results, ToScheduleResp(schedule))
	}

	if len(results) == 0 && len(conflicts) > 0 {
		return nil, apperrors.New(apperrors.ErrScheduleConflict.Code, fmt.Sprintf("所有排课均冲突: %v", conflicts))
	}

	return results, nil
}

// ============================================================
// Coach leave management
// ============================================================

func (s *Service) CreateLeave(ctx context.Context, req CreateLeaveReq) (*LeaveResp, error) {
	if _, err := s.repo.GetCoachByID(ctx, req.CoachID); err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	leave := &CoachLeave{
		CoachID:   req.CoachID,
		StartDate: req.StartDate,
		EndDate:   req.EndDate,
		Reason:    req.Reason,
		Status:    "pending",
	}

	if err := s.repo.CreateLeave(ctx, leave); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToLeaveResp(leave)
	return &resp, nil
}

func (s *Service) ListLeaves(ctx context.Context, coachID *int64, status *string, offset, limit int) ([]LeaveResp, int64, error) {
	leaves, total, err := s.repo.ListLeaves(ctx, coachID, status, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]LeaveResp, 0, len(leaves))
	for i := range leaves {
		result = append(result, ToLeaveResp(&leaves[i]))
	}
	return result, total, nil
}

func (s *Service) GetLeave(ctx context.Context, id int64) (*LeaveResp, error) {
	leave, err := s.repo.GetLeaveByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}
	resp := ToLeaveResp(leave)
	return &resp, nil
}

func (s *Service) UpdateLeave(ctx context.Context, id int64, req UpdateLeaveReq) (*LeaveResp, error) {
	leave, err := s.repo.GetLeaveByID(ctx, id)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	if req.StartDate != nil {
		leave.StartDate = *req.StartDate
	}
	if req.EndDate != nil {
		leave.EndDate = *req.EndDate
	}
	if req.Reason != nil {
		leave.Reason = *req.Reason
	}
	if req.Status != nil {
		leave.Status = *req.Status
	}

	if err := s.repo.UpdateLeave(ctx, leave); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToLeaveResp(leave)
	return &resp, nil
}

func (s *Service) DeleteLeave(ctx context.Context, id int64) error {
	if _, err := s.repo.GetLeaveByID(ctx, id); err != nil {
		return apperrors.ErrRecordNotFound
	}
	return s.repo.DeleteLeave(ctx, id)
}

// GetLeaveConflicts returns schedules that overlap with a leave period.
func (s *Service) GetLeaveConflicts(ctx context.Context, coachID int64, startDate, endDate string) ([]ScheduleResp, error) {
	schedules, err := s.repo.ListSchedulesDuringLeave(ctx, coachID, startDate, endDate)
	if err != nil {
		return nil, apperrors.ErrInternal
	}

	result := make([]ScheduleResp, 0, len(schedules))
	for i := range schedules {
		result = append(result, ToScheduleResp(&schedules[i]))
	}
	return result, nil
}

// ============================================================
// Student & record management
// ============================================================

func (s *Service) ListStudents(ctx context.Context, offset, limit int) ([]StudentListItem, int64, error) {
	items, total, err := s.repo.ListStudents(ctx, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}
	return items, total, nil
}

func (s *Service) ListStudentRecords(ctx context.Context, studentID int64, offset, limit int) ([]RecordResp, int64, error) {
	records, total, err := s.repo.ListRecordsByStudent(ctx, studentID, offset, limit)
	if err != nil {
		return nil, 0, apperrors.ErrInternal
	}

	result := make([]RecordResp, 0, len(records))
	for i := range records {
		result = append(result, ToRecordResp(&records[i]))
	}
	return result, total, nil
}

func (s *Service) AddFeedback(ctx context.Context, recordID int64, req FeedbackReq) (*RecordResp, error) {
	rec, err := s.repo.GetRecordByID(ctx, recordID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	rec.CoachFeedback = req.CoachFeedback
	if req.Attendance != "" {
		rec.Attendance = req.Attendance
	}

	if err := s.repo.UpdateRecord(ctx, rec); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToRecordResp(rec)
	return &resp, nil
}

func (s *Service) AddRating(ctx context.Context, recordID, userID int64, req RatingReq) (*RecordResp, error) {
	rec, err := s.repo.GetRecordByID(ctx, recordID)
	if err != nil {
		return nil, apperrors.ErrRecordNotFound
	}

	// Ensure the student owns this record
	if rec.StudentID != userID {
		return nil, apperrors.ErrForbidden
	}

	rec.Rating = &req.Rating
	rec.RatingComment = req.RatingComment

	if err := s.repo.UpdateRecord(ctx, rec); err != nil {
		return nil, apperrors.ErrInternal
	}

	resp := ToRecordResp(rec)
	return &resp, nil
}

// ListMyRecords returns course records for the authenticated user.
func (s *Service) ListMyRecords(ctx context.Context, userID int64, offset, limit int) ([]RecordResp, int64, error) {
	return s.ListStudentRecords(ctx, userID, offset, limit)
}
