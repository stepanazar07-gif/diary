package backend

import (
	"encoding/json"
	"fmt"
	"net/http"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func handleHome(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Дневник занятий работает!")
}
func createUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var user User
		var request struct {
			Username string `json:"username"`
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)
		if err != nil {
			http.Error(w, "Неправильные данные", http.StatusBadRequest)
			return
		}

		user.Username = request.Username
		user.Email = request.Email
		user.Role = request.Role

		if user.Role == "" {
			user.Role = "student"
		}

	
		err = validateUser(user.Username, user.Email, request.Password, user.Role)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

	
		hash, err := bcrypt.GenerateFromPassword(
			[]byte(request.Password),
			bcrypt.DefaultCost,
		)

		if err != nil {
			http.Error(w, "Ошибка обработки пароля", http.StatusInternalServerError)
			return
		}

	
		user.Password = string(hash)

		err = db.Create(&user).Error
		if err != nil {
			http.Error(w, "Не удалось создать пользователя", http.StatusConflict)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)
		json.NewEncoder(w).Encode(user)
	}
}
	func StartServer(db *gorm.DB) {

	http.HandleFunc("/", handleHome)
	http.HandleFunc("POST /users", createUser(db))

	fmt.Println("Сервер запущен: http://localhost:8080")

	err := http.ListenAndServe(":8080", nil)

	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}
