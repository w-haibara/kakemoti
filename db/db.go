package db

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/w-haibara/kakemoti/compiler"
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

type Workflows struct {
	gorm.Model
	Name     string
	Workflow []byte
}

func RegisterWorkflow(name string, w compiler.Workflow) error {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return err
	}

	var wb bytes.Buffer
	enc := gob.NewEncoder(&wb)

	if err := enc.Encode(w); err != nil {
		return err
	}

	db.Create(&Workflows{
		Name:     name,
		Workflow: wb.Bytes(),
	})

	return nil
}

func RmWorkflow(name string) error {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.Delete(&Workflows{}, "name = ?", name).Error; err != nil {
		return err
	}

	return nil
}

func FetchWorkflow(name string) (compiler.Workflow, error) {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return compiler.Workflow{}, err
	}

	var w Workflows
	db.First(&w, "name = ?", name)

	wb := bytes.NewBuffer(w.Workflow)
	dec := gob.NewDecoder(wb)

	var workflow compiler.Workflow
	if err := dec.Decode(&workflow); err != nil {
		return compiler.Workflow{}, err
	}

	return workflow, nil
}
