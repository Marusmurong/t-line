package academic

import (
	"context"

	"gorm.io/gorm"
)

type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

// --- Coach ---

func (r *Repository) CreateCoach(ctx context.Context, c *Coach) error {
	return r.db.WithContext(ctx).Create(c).Error
}

func (r *Repository) GetCoachByID(ctx context.Context, id int64) (*Coach, error) {
	var c Coach
	if err := r.db.WithContext(ctx).First(&c, id).Error; err != nil {
		return nil, err
	}
	return &c, nil
}

func (r *Repository) ListCoaches(ctx context.Context, status *int, offset, limit int) ([]Coach, int64, error) {
	var coaches []Coach
	var total int64

	query := r.db.WithContext(ctx).Model(&Coach{})
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&coaches).Error
	return coaches, total, err
}

func (r *Repository) UpdateCoach(ctx context.Context, c *Coach) error {
	return r.db.WithContext(ctx).Save(c).Error
}

func (r *Repository) DeleteCoach(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&Coach{}, id).Error
}

// --- Schedule ---

func (r *Repository) CreateSchedule(ctx context.Context, s *CourseSchedule) error {
	return r.db.WithContext(ctx).Create(s).Error
}

func (r *Repository) GetScheduleByID(ctx context.Context, id int64) (*CourseSchedule, error) {
	var s CourseSchedule
	if err := r.db.WithContext(ctx).First(&s, id).Error; err != nil {
		return nil, err
	}
	return &s, nil
}

func (r *Repository) ListSchedules(ctx context.Context, coachID, venueID *int64, date, status *string, offset, limit int) ([]CourseSchedule, int64, error) {
	var schedules []CourseSchedule
	var total int64

	query := r.db.WithContext(ctx).Model(&CourseSchedule{})
	if coachID != nil {
		query = query.Where("coach_id = ?", *coachID)
	}
	if venueID != nil {
		query = query.Where("venue_id = ?", *venueID)
	}
	if date != nil {
		query = query.Where("date = ?", *date)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("date DESC, start_time ASC").Offset(offset).Limit(limit).Find(&schedules).Error
	return schedules, total, err
}

func (r *Repository) UpdateSchedule(ctx context.Context, s *CourseSchedule) error {
	return r.db.WithContext(ctx).Save(s).Error
}

func (r *Repository) DeleteSchedule(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&CourseSchedule{}, id).Error
}

// --- Conflict detection queries ---

func (r *Repository) FindCoachConflict(ctx context.Context, coachID int64, date, startTime, endTime string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&CourseSchedule{}).
		Where("coach_id = ? AND date = ? AND status != 'cancelled' AND id != ?", coachID, date, excludeID).
		Where("start_time < ? AND end_time > ?", endTime, startTime).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) FindVenueConflict(ctx context.Context, venueID int64, date, startTime, endTime string, excludeID int64) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&CourseSchedule{}).
		Where("venue_id = ? AND date = ? AND status != 'cancelled' AND id != ?", venueID, date, excludeID).
		Where("start_time < ? AND end_time > ?", endTime, startTime).
		Count(&count).Error
	return count > 0, err
}

func (r *Repository) FindCoachLeaveConflict(ctx context.Context, coachID int64, date string) (bool, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&CoachLeave{}).
		Where("coach_id = ? AND status = 'approved' AND start_date <= ? AND end_date >= ?", coachID, date, date).
		Count(&count).Error
	return count > 0, err
}

// --- Leave ---

func (r *Repository) CreateLeave(ctx context.Context, l *CoachLeave) error {
	return r.db.WithContext(ctx).Create(l).Error
}

func (r *Repository) GetLeaveByID(ctx context.Context, id int64) (*CoachLeave, error) {
	var l CoachLeave
	if err := r.db.WithContext(ctx).First(&l, id).Error; err != nil {
		return nil, err
	}
	return &l, nil
}

func (r *Repository) ListLeaves(ctx context.Context, coachID *int64, status *string, offset, limit int) ([]CoachLeave, int64, error) {
	var leaves []CoachLeave
	var total int64

	query := r.db.WithContext(ctx).Model(&CoachLeave{})
	if coachID != nil {
		query = query.Where("coach_id = ?", *coachID)
	}
	if status != nil {
		query = query.Where("status = ?", *status)
	}

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&leaves).Error
	return leaves, total, err
}

func (r *Repository) UpdateLeave(ctx context.Context, l *CoachLeave) error {
	return r.db.WithContext(ctx).Save(l).Error
}

func (r *Repository) DeleteLeave(ctx context.Context, id int64) error {
	return r.db.WithContext(ctx).Delete(&CoachLeave{}, id).Error
}

// --- Schedule conflicts during leave ---

func (r *Repository) ListSchedulesDuringLeave(ctx context.Context, coachID int64, startDate, endDate string) ([]CourseSchedule, error) {
	var schedules []CourseSchedule
	err := r.db.WithContext(ctx).
		Where("coach_id = ? AND date >= ? AND date <= ? AND status != 'cancelled'", coachID, startDate, endDate).
		Find(&schedules).Error
	return schedules, err
}

// --- Student records ---

func (r *Repository) CreateRecord(ctx context.Context, rec *StudentCourseRecord) error {
	return r.db.WithContext(ctx).Create(rec).Error
}

func (r *Repository) GetRecordByID(ctx context.Context, id int64) (*StudentCourseRecord, error) {
	var rec StudentCourseRecord
	if err := r.db.WithContext(ctx).First(&rec, id).Error; err != nil {
		return nil, err
	}
	return &rec, nil
}

func (r *Repository) ListRecordsByStudent(ctx context.Context, studentID int64, offset, limit int) ([]StudentCourseRecord, int64, error) {
	var records []StudentCourseRecord
	var total int64

	query := r.db.WithContext(ctx).Model(&StudentCourseRecord{}).Where("student_id = ?", studentID)

	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	err := query.Order("created_at DESC").Offset(offset).Limit(limit).Find(&records).Error
	return records, total, err
}

func (r *Repository) UpdateRecord(ctx context.Context, rec *StudentCourseRecord) error {
	return r.db.WithContext(ctx).Save(rec).Error
}

// --- Student list (distinct students with course count) ---

func (r *Repository) ListStudents(ctx context.Context, offset, limit int) ([]StudentListItem, int64, error) {
	var items []StudentListItem
	var total int64

	countQuery := r.db.WithContext(ctx).
		Model(&StudentCourseRecord{}).
		Select("COUNT(DISTINCT student_id)")

	if err := countQuery.Scan(&total).Error; err != nil {
		return nil, 0, err
	}

	err := r.db.WithContext(ctx).
		Model(&StudentCourseRecord{}).
		Select("student_id, COUNT(*) as total_courses, MAX(created_at) as last_course_at").
		Group("student_id").
		Order("total_courses DESC").
		Offset(offset).Limit(limit).
		Scan(&items).Error
	return items, total, err
}

// --- Performance stats ---

func (r *Repository) CountCoachSchedules(ctx context.Context, coachID int64) (int64, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&CourseSchedule{}).
		Where("coach_id = ? AND status = 'completed'", coachID).
		Count(&count).Error
	return count, err
}

func (r *Repository) AvgCoachRating(ctx context.Context, coachID int64) (float64, error) {
	var avg *float64
	err := r.db.WithContext(ctx).
		Model(&StudentCourseRecord{}).
		Select("AVG(rating)").
		Where("coach_id = ? AND rating IS NOT NULL", coachID).
		Scan(&avg).Error
	if avg == nil {
		return 0, err
	}
	return *avg, err
}

func (r *Repository) CountCoachStudents(ctx context.Context, coachID int64) (int, error) {
	var count int64
	err := r.db.WithContext(ctx).
		Model(&StudentCourseRecord{}).
		Where("coach_id = ?", coachID).
		Distinct("student_id").
		Count(&count).Error
	return int(count), err
}

func (r *Repository) Transaction(fn func(tx *gorm.DB) error) error {
	return r.db.Transaction(fn)
}
