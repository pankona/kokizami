package todo

import (
	"errors"
	"log"
)

// ToDo represents a struct of ToDo item
type ToDo struct {
	desc string
}

// NewToDo allocates a ToDo item
func NewToDo(desc string) *ToDo {
	return &ToDo{desc: desc}
}

// Add adds a specified ToDo item to DB
func Add(t *ToDo) error {
	log.Println("Add. t =", t.desc)
	// TODO: store desc to DB
	return errors.New("Add failed")
}
