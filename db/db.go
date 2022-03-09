package db

import (
	"bytes"
	"encoding/gob"
	"os"
	"path/filepath"

	"github.com/glebarez/sqlite"
	"github.com/k0kubun/pp"
	"github.com/w-haibara/kakemoti/compiler"
	"github.com/w-haibara/kakemoti/config"
	"gorm.io/gorm"
)

var dbFileName = ""

func init() {
	dbFileName = filepath.Join(config.ConfigDir(), "workflow.db")
	_, _ = pp.Println(dbFileName)

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

func FetchWorkflow(name string) (compiler.Workflow, error) {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return compiler.Workflow{}, err
	}

	var w Workflows
	db.First(&w, "name = ?", name)
	_, _ = pp.Println("db read", w)

	var wb bytes.Buffer
	dec := gob.NewDecoder(&wb)

	var workflow compiler.Workflow
	if err := dec.Decode(&workflow); err != nil {
		return compiler.Workflow{}, err
	}

	return workflow, nil
}
