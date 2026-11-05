package main

import (
	"log"

	"diary/backend"
)

func main() {
	
	db, err := backend.ConnectDatabase()
	if err != nil {
		log.Fatal("Ошибка подключения к базе:", err)
	}

	log.Println("Дневник занятий запущен!")

	
	backend.StartServer(db)
}