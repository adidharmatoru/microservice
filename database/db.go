package database

import (
    "microservice/models"
    "log"

    "github.com/jinzhu/gorm"
    _ "github.com/jinzhu/gorm/dialects/postgres"
    _ "github.com/jinzhu/gorm/dialects/sqlite"
)

func ConnectDatabase() {
    var err error
    models.DB, err = gorm.Open("sqlite3", "local.db")
    if err != nil {
        log.Fatal("Failed to connect to database!", err)
    }

    models.DB.AutoMigrate(&models.User{})
}
