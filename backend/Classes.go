package backend

import (
	"time"
)

type Class struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
    Title     string    `gorm:"not null" json:"title"`
    Topic     string    `gorm:"not null" json:"topic"`
    Date      time.Time `gorm:"not null" json:"date"`
    Time      string    `gorm:"not null" json:"time"`
    Homework  string    `json:"homework"`
    TeacherID uint      `gorm:"not null" json:"teacher_id"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}





