package backend

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"gorm.io/gorm"
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

func createClass(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var class Class

		err := json.NewDecoder(r.Body).Decode(&class)

		if err != nil {
			http.Error(w, "Неправильные данные", http.StatusBadRequest)
			return
		}

		if class.Title == "" || class.Topic == "" ||
			class.Date.IsZero() || class.Time == "" {

			http.Error(w, "Заполните обязательные поля", http.StatusBadRequest)
			return
		}

		claims := r.Context().Value("user").(*Claims)

		class.TeacherID = claims.UserID

		class.ID = 0

		err = db.Create(&class).Error

		if err != nil {
			http.Error(w, "Ошибка создания занятия", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(class)
	}
}

func getClasses(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		var classes []Class

		err := db.Find(&classes).Error

		if err != nil {
			http.Error(w, "Ошибка получения занятий", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(classes)
	}
}

func updateClass(db *gorm.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		id, err := strconv.Atoi(r.PathValue("id"))

		if err != nil || id <= 0 {
			http.Error(w, "Неправильный ID", http.StatusBadRequest)
			return
		}

		var class Class

		err = db.First(&class, id).Error

		if err != nil {
			http.Error(w, "Занятие не найдено", http.StatusNotFound)
			return
		}

		var updatedClass Class

		err = json.NewDecoder(r.Body).Decode(&updatedClass)

		if err != nil {
			http.Error(w, "Неправильные данные", http.StatusBadRequest)
			return
		}

		if updatedClass.Title == "" || updatedClass.Topic == "" ||
			updatedClass.Date.IsZero() || updatedClass.Time == "" ||
			updatedClass.TeacherID == 0 {

			http.Error(w, "Заполните обязательные поля", http.StatusBadRequest)
			return
		}

		claims := r.Context().Value("user").(*Claims)

		if class.TeacherID != claims.UserID {
			http.Error(w, "Нельзя изменять чужое занятие", http.StatusForbidden)
			return
		}

		class.Title = updatedClass.Title
		class.Topic = updatedClass.Topic
		class.Date = updatedClass.Date
		class.Time = updatedClass.Time
		class.Homework = updatedClass.Homework
		

		err = db.Save(&class).Error

		if err != nil {
			http.Error(w, "Ошибка изменения занятия", http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(class)
	}
}

func deleteClass(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

	
		id := r.PathValue("id")

	
		var class Class

		err := db.First(&class, id).Error

		if err != nil {
			http.Error(w, "Занятие не найдено", http.StatusNotFound)
			return
		}

		
		claims := r.Context().Value("user").(*Claims)

		if class.TeacherID != claims.UserID {
			http.Error(w, "Нельзя удалить чужое занятие", http.StatusForbidden)
			return
		}

	
		err = db.Delete(&class).Error

		if err != nil {
			http.Error(w, "Ошибка удаления занятия", http.StatusInternalServerError)
			return
		}

		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Занятие успешно удалено"))
	}
}