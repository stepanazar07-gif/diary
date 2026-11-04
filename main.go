package main

import (
	"fmt"
	"diary/backend"
)

func main() {
	fmt.Println("Запуск сервера")
	backend.StartServer()
}