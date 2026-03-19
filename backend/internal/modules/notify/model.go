package notify

import (
	"time"

	"gorm.io/datatypes"
)

type Notification struct {
	ID        int64          `gorm:"primaryKey" json:"id"`
	UserID    int64          `json:"user_id"`
	Type      string         `gorm:"size:40;not null" json:"type"`
	Title     string         `gorm:"size:128;not null" json:"title"`
	Content   string         `gorm:"type:text" json:"content"`
	IsRead    bool           `gorm:"default:false" json:"is_read"`
	ExtraData datatypes.JSON `gorm:"type:jsonb;default:'{}'" json:"extra_data"`
	CreatedAt time.Time      `json:"created_at"`
}

func (Notification) TableName() string { return "notifications" }
