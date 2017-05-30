package service

import (
	"../types"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"log"
)

type Service struct {
	db gorm.DB
}

func (s *Service) SetupDatabase() { //todo arguments dialect database
	db, err := gorm.Open("sqlite3", "test.db")
	if err != nil {
		log.Panic("Failed to open database!")
	}

	db.AutoMigrate(&types.User{})
	s.db = *db
}
