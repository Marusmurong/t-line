package academic

import (
	"time"

	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
)

type Coach struct {
	ID           int64           `gorm:"primaryKey" json:"id"`
	UserID       *int64          `gorm:"uniqueIndex" json:"user_id"`
	Title        string          `gorm:"size:64" json:"title"`
	Specialties  datatypes.JSON  `gorm:"type:jsonb;default:'[]'" json:"specialties"`
	Bio          string          `gorm:"type:text" json:"bio"`
	HourlyRate   decimal.Decimal `gorm:"type:decimal(10,2);default:0" json:"hourly_rate"`
	Rating       decimal.Decimal `gorm:"type:decimal(3,2);default:5.0" json:"rating"`
	StudentCount int             `gorm:"default:0" json:"student_count"`
	Status       int             `gorm:"default:1" json:"status"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

func (Coach) TableName() string { return "coaches" }

type CourseSchedule struct {
	ID                 int64      `gorm:"primaryKey" json:"id"`
	CoachID            int64      `gorm:"not null" json:"coach_id"`
	VenueID            int64      `gorm:"not null" json:"venue_id"`
	StudentID          *int64     `json:"student_id"`
	ProductID          *int64     `json:"product_id"`
	Date               string     `gorm:"type:date;not null" json:"date"`
	StartTime          string     `gorm:"type:time;not null" json:"start_time"`
	EndTime            string     `gorm:"type:time;not null" json:"end_time"`
	Type               string     `gorm:"size:20;default:'private'" json:"type"`
	Status             string     `gorm:"size:20;default:'scheduled'" json:"status"`
	RecurrenceRule     string     `gorm:"size:128" json:"recurrence_rule"`
	RecurrenceGroupID  *uuid.UUID `gorm:"type:uuid" json:"recurrence_group_id"`
	SubstituteCoachID  *int64     `json:"substitute_coach_id"`
	CreatedAt          time.Time  `json:"created_at"`
	UpdatedAt          time.Time  `json:"updated_at"`
}

func (CourseSchedule) TableName() string { return "course_schedules" }

type CoachLeave struct {
	ID        int64     `gorm:"primaryKey" json:"id"`
	CoachID   int64     `gorm:"not null" json:"coach_id"`
	StartDate string    `gorm:"type:date;not null" json:"start_date"`
	EndDate   string    `gorm:"type:date;not null" json:"end_date"`
	Reason    string    `gorm:"size:256" json:"reason"`
	Status    string    `gorm:"size:20;default:'pending'" json:"status"`
	CreatedAt time.Time `json:"created_at"`
}

func (CoachLeave) TableName() string { return "coach_leaves" }

type StudentCourseRecord struct {
	ID            int64     `gorm:"primaryKey" json:"id"`
	ScheduleID    int64     `gorm:"not null" json:"schedule_id"`
	StudentID     int64     `gorm:"not null" json:"student_id"`
	CoachID       int64     `gorm:"not null" json:"coach_id"`
	Attendance    string    `gorm:"size:20;default:'present'" json:"attendance"`
	CoachFeedback string    `gorm:"type:text" json:"coach_feedback"`
	Rating        *int      `json:"rating"`
	RatingComment string    `gorm:"type:text" json:"rating_comment"`
	CreatedAt     time.Time `json:"created_at"`
}

func (StudentCourseRecord) TableName() string { return "student_course_records" }
