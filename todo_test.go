package todo

import "testing"

func TestNewToDo(t *testing.T) {
	todo := NewToDo()
	if todo == nil {
		t.Error("NewToDo returned nil. failed")
	}
}
