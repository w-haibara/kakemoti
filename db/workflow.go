package db

import (
	"bytes"
	"encoding/gob"
	"fmt"

	"github.com/glebarez/sqlite"
	"github.com/w-haibara/kakemoti/compiler"
	"gorm.io/gorm"
)

type Workflows struct {
	Name     string `gorm:"primaryKey"`
	Workflow []byte
}

func MustMigrateWorkflows(db *gorm.DB) {
	if err := db.AutoMigrate(&Workflows{}); err != nil {
		panic(err.Error())
	}
}

var ErrWorkflowNameAlreadyExists = func(name string) error {
	return fmt.Errorf("the workflow name already exists: %s", name)
}

func RegisterWorkflow(name string, w compiler.Workflow, force bool) error {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return err
	}

	exists, err := isExistsWorkflowName(db, name)
	if err != nil {
		return err
	}
	if !force && exists {
		return ErrWorkflowNameAlreadyExists(name)
	}

	var wb bytes.Buffer
	enc := gob.NewEncoder(&wb)

	if err := enc.Encode(w); err != nil {
		return err
	}

	if force && exists {
		db.Save(&Workflows{
			Name:     name,
			Workflow: wb.Bytes(),
		})
	} else {
		db.Create(&Workflows{
			Name:     name,
			Workflow: wb.Bytes(),
		})
	}

	return nil
}

func isExistsWorkflowName(db *gorm.DB, name string) (bool, error) {
	var w Workflows
	res := db.Find(&w, "name = ?", name)
	if err := res.Error; err != nil {
		return false, err
	}
	if res.RowsAffected == 0 {
		return false, nil
	}
	return true, nil
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

func DropWorkflow() error {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return err
	}

	if err := db.Migrator().DropTable(&Workflows{}); err != nil {
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

func ListWorkflow() ([]string, error) {
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
