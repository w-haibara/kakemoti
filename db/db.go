package db

import (
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/w-haibara/kakemoti/config"
	"gorm.io/gorm"
)

var dbFileName = ""

func init() {
	dbFileName = filepath.Join(config.ConfigDir(), "workflow.db")

	if _, err := os.Stat(dbFileName); err == nil {
		return
	}

	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	if err := db.AutoMigrate(&Workflows{}); err != nil {
		panic(err.Error())
	}
}
