package db

import (
	"bytes"
	"encoding/base64"
	"encoding/gob"
	"fmt"
	"time"

	"github.com/glebarez/sqlite"
	"github.com/ohler55/ojg/jp"
	"github.com/w-haibara/kakemoti/controller/compiler"
	"gorm.io/gorm"
)

func init() {
	registeerTypesForGob()
}

type Workflows struct {
	Name      string `gorm:"primaryKey"`
	ASL       string
	Workflow  []byte
	CreatedAt time.Time
}

func (w *Workflows) EncodeAndSetASL(asl []byte) {
	w.ASL = base64.StdEncoding.EncodeToString(asl)
}

func (w Workflows) DecodeASL() (string, error) {
	b, err := base64.StdEncoding.DecodeString(w.ASL)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func (w *Workflows) EncodeAndSetsWorkflow(workflow compiler.Workflow) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)

	if err := enc.Encode(&workflow); err != nil {
		return err
	}

	w.Workflow = b.Bytes()

	return nil
}

func (w Workflows) DecodeWorkflow() (compiler.Workflow, error) {
	wb := bytes.NewBuffer(w.Workflow)
	dec := gob.NewDecoder(wb)

	var workflow compiler.Workflow
	if err := dec.Decode(&workflow); err != nil {
		return compiler.Workflow{}, err
	}

	return workflow, nil
}

func MustMigrateWorkflows(db *gorm.DB) {
	if err := db.AutoMigrate(&Workflows{}); err != nil {
		panic(err.Error())
	}
}

var ErrWorkflowNameAlreadyExists = func(name string) error {
	return fmt.Errorf("the workflow name already exists: %s", name)
}

func RegisterWorkflow(name string, w compiler.Workflow, asl []byte) error {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return err
	}

	exists, err := isExistsWorkflowName(db, name)
	if err != nil {
		return err
	}
	if exists {
		return ErrWorkflowNameAlreadyExists(name)
	}

	wf := &Workflows{
		Name:      name,
		CreatedAt: time.Now(),
	}
	wf.EncodeAndSetASL(asl)
	if err := wf.EncodeAndSetsWorkflow(w); err != nil {
		return err
	}

	db.Create(wf)

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

func GetWorkflow(name string) (Workflows, error) {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return Workflows{}, err
	}

	var w Workflows
	db.First(&w, "name = ?", name)

	return w, nil
}

func ListWorkflow() ([]Workflows, error) {
	db, err := gorm.Open(sqlite.Open(dbFileName), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	var w []Workflows
	if err := db.Find(&w).Error; err != nil {
		return nil, err
	}

	return w, nil
}

func registeerTypesForGob() {
	gob.Register(map[string]interface{}{})
	gob.Register([]interface{}{})

	gob.Register(compiler.ChoiceState{})
	gob.Register(compiler.CommonState5{})
	gob.Register(compiler.FailState{})
	gob.Register(compiler.MapState{})
	gob.Register(compiler.ParallelState{})
	gob.Register(compiler.PassState{})
	gob.Register(compiler.SucceedState{})
	gob.Register(compiler.TaskState{})
	gob.Register(compiler.WaitState{})

	gob.Register(compiler.AndRule{})
	gob.Register(compiler.OrRule{})
	gob.Register(compiler.NotRule{})
	gob.Register(compiler.StringEqualsRule{})
	gob.Register(compiler.StringEqualsPathRule{})
	gob.Register(compiler.StringLessThanRule{})
	gob.Register(compiler.StringLessThanPathRule{})
	gob.Register(compiler.StringLessThanEqualsRule{})
	gob.Register(compiler.StringLessThanEqualsPathRule{})
	gob.Register(compiler.StringGreaterThanRule{})
	gob.Register(compiler.StringGreaterThanPathRule{})
	gob.Register(compiler.StringGreaterThanEqualsRule{})
	gob.Register(compiler.StringGreaterThanEqualsPathRule{})
	gob.Register(compiler.StringMatchesRule{})
	gob.Register(compiler.NumericEqualsRule{})
	gob.Register(compiler.NumericEqualsPathRule{})
	gob.Register(compiler.NumericLessThanRule{})
	gob.Register(compiler.NumericLessThanPathRule{})
	gob.Register(compiler.NumericLessThanEqualsRule{})
	gob.Register(compiler.NumericLessThanEqualsPathRule{})
	gob.Register(compiler.NumericGreaterThanRule{})
	gob.Register(compiler.NumericGreaterThanPathRule{})
	gob.Register(compiler.NumericGreaterThanEqualsRule{})
	gob.Register(compiler.NumericGreaterThanEqualsPathRule{})
	gob.Register(compiler.BooleanEqualsRule{})
	gob.Register(compiler.BooleanEqualsPathRule{})
	gob.Register(compiler.TimestampEqualsRule{})
	gob.Register(compiler.TimestampEqualsPathRule{})
	gob.Register(compiler.TimestampLessThanRule{})
	gob.Register(compiler.TimestampLessThanPathRule{})
	gob.Register(compiler.TimestampLessThanEqualsRule{})
	gob.Register(compiler.TimestampLessThanEqualsPathRule{})
	gob.Register(compiler.TimestampGreaterThanRule{})
	gob.Register(compiler.TimestampGreaterThanPathRule{})
	gob.Register(compiler.TimestampGreaterThanEqualsRule{})
	gob.Register(compiler.TimestampGreaterThanEqualsPathRule{})
	gob.Register(compiler.IsNullRule{})
	gob.Register(compiler.IsPresentRule{})
	gob.Register(compiler.IsNumericRule{})
	gob.Register(compiler.IsStringRule{})
	gob.Register(compiler.IsBooleanRule{})
	gob.Register(compiler.IsTimestampRule{})

	gob.Register(jp.Root(0))
	gob.Register(jp.Child(""))
}
