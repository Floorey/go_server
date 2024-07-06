package main

import (
	"log"
	"sync"

	"github.com/jinzhu/gorm"
	_ "github.com/mattn/go-sqlite3"
)

var dbMutex sync.Mutex
var dbOps = make(chan func(), 1)
var db *gorm.DB

func InitDB() {
	dbMutex.Lock()
	defer dbMutex.Unlock()

	if db == nil {
		var err error
		db, err = gorm.Open("sqlite3", "test.db")
		if err != nil {
			log.Fatal(err)
		}
		db.AutoMigrate(&VM{})

		// Start a Goroutine to handle database operations
		go func() {
			for op := range dbOps {
				op()
			}
		}()
	}
}

type VM struct {
	gorm.Model
	Name  string
	Image string
	State string
}
