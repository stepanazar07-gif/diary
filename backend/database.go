package backend

import (
    "fmt"

    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
)

func ConnectDatabase() (*gorm.DB, error) {

    db, err := gorm.Open(sqlite.Open("diary.db"), &gorm.Config{})
    if err != nil {
        return nil, err
    }

    fmt.Println("Подключение к SQLite успешно!")

    err = db.AutoMigrate(&User{}, &Class{})
    if err != nil {
        return nil, err
    }

    fmt.Println("Таблицы успешно созданы!")

    return db, nil
}