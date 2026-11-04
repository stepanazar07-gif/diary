package backend

import (
	"fmt"
	"net/http"

	
)

func handleUsers(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w,"Hello World")
}
func StartServer() {
	http.HandleFunc("/", handleUsers)
	
	fmt.Println("Сервер запущен на:http://localhost:8080")

	if err := http.ListenAndServe(":8080", nil);err != nil {
		fmt.Println("Error starting server:", err)
	}
}