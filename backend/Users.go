package backend

import (
	"fmt"
	"net/mail"
	"strings"
	"time"
)

type User struct {
	ID        uint      `gorm:"primaryKey" json:"id"`
	Username  string    `gorm:"unique;not null" json:"username"`
	Email     string    `gorm:"unique;not null" json:"email"`
	CreatedAt time.Time `json:"created_at"`
	Role      string    `gorm:"default:student" json:"role"`
	Password  string    `gorm:"not null" json:"-"`
}

func validateUser(username, email, password, role string) error {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if len([]rune(username)) < 3 {
		return fmt.Errorf("Логин должен сожержать минимум 3 символа!")

	}
	if _, err := mail.ParseAddress(email); err != nil {
		return fmt.Errorf("Неверный формат email!")
	}
	if len(password) < 6 {
		return fmt.Errorf("Пароль должен содержать минимум 6 символов!")
	}
	if role != "student" && role != "teacher" {
		return fmt.Errorf("Роль должна быть либо 'student', либо 'teacher'!")
	}

	return nil

}
