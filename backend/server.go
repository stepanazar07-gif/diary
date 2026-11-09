
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

		
		var count int64

		err = db.Model(&User{}).
			Where("username = ?", user.Username).
			Or("email = ?", user.Email).
			Count(&count).Error

		if err != nil {
			http.Error(w, "Ошибка проверки пользователя", http.StatusInternalServerError)
			return
		}

		if count > 0 {
			http.Error(w, "Пользователь уже существует", http.StatusConflict)
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
			http.Error(w, "Не удалось создать пользователя", http.StatusInternalServerError)
			return
		}

		
		user.Password = ""

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusCreated)

		json.NewEncoder(w).Encode(user)
	}
}




func loginUser(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var request struct {
			Username string `json:"username"`
			Password string `json:"password"`
		}

		err := json.NewDecoder(r.Body).Decode(&request)

		if err != nil {
			http.Error(w, "Неправильные данные", http.StatusBadRequest)
			return
		}

		if request.Username == "" || request.Password == "" {
			http.Error(w, "Введите логин и пароль", http.StatusBadRequest)
			return
		}

		var user User

		err = db.Where("username = ?", request.Username).First(&user).Error

		if err != nil {
			http.Error(w, "Пользователь не найден", http.StatusNotFound)
			return
		}

	
		err = bcrypt.CompareHashAndPassword(
			[]byte(user.Password),
			[]byte(request.Password),
		)

		if err != nil {
			http.Error(w, "Неверный пароль", http.StatusUnauthorized)
			return
		}

	
		token, err := generateToken(user)

		if err != nil {
			http.Error(w, "Ошибка создания токена", http.StatusInternalServerError)
			return
		}

	
		user.Password = ""

		
		response := map[string]interface{}{
			"token": token,
			"user":  user,
		}

		w.Header().Set("Content-Type", "application/json")

		json.NewEncoder(w).Encode(response)
	}
}



func getUsers(db *gorm.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		var users []User

		err := db.Find(&users).Error
		if err != nil {
			http.Error(w, "Ошибка получения пользователей", http.StatusInternalServerError)
			return
		}

		// Скрываем пароли
		for i := range users {
			users[i].Password = ""
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(users)
	}
}



func StartServer(db *gorm.DB) {

	mux := http.NewServeMux()


	mux.HandleFunc("GET /", handleHome)


	mux.HandleFunc("POST /users", createUser(db))
	mux.Handle("GET /users", authMiddleware(getUsers(db)))
	mux.HandleFunc("POST /login", loginUser(db))

		mux.Handle("GET /classes",
		authMiddleware(getClasses(db)),
	)

	
	mux.Handle("POST /classes",
		authMiddleware(teacherOnly(createClass(db))),
	)


	mux.Handle("PUT /classes/{id}",
		authMiddleware(teacherOnly(updateClass(db))),
	)


	mux.Handle("DELETE /classes/{id}",
		authMiddleware(teacherOnly(deleteClass(db))),
	)

	fmt.Println("Сервер запущен: http://localhost:8080")

	err := http.ListenAndServe(":8080", mux)

	if err != nil {
		fmt.Println("Ошибка запуска сервера:", err)
	}
}

