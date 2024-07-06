package main

import (
	"log"

	"github.com/jinzhu/gorm"
)

var db *gorm.DB
var err error

func InitDB() {
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Fatal(err)
	}
	db.AutoMigrate(&VM{})
}

type VM struct {
	gorm.Model
	Name  string
	Image string
	State string
}
