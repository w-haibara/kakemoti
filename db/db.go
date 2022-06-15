package db

import (
	"fmt"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

var dbFileName = "/tmp/kakemoti/workflows.db"

func init() {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		panic(err.Error())
	}

	MustMigrateWorkflows(db)

	fmt.Println("===========================================")
	w, err := ListWorkflow()
	if err != nil {
		panic(err.Error())
	}
	for i, v := range w {
		fmt.Println(i, v.Name)
	}
	fmt.Println("===========================================")
}
