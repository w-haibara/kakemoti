package db

import (
	"bytes"
	"encoding/gob"

	"github.com/glebarez/sqlite"
	"github.com/w-haibara/kakemoti/compiler"
	"gorm.io/gorm"
)

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

func RemoveWorkflow(name string) error {
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

func ListWorkflow(name string) ([]string, error) {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var w []Workflows
	if err := db.Find(&w).Error; err != nil {
		return nil, err
	}

	res := []string{}
	for _, v := range w {
		res = append(res, v.Name)
	}

	return res, nil
}
