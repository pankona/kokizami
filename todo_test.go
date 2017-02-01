package todo

import "testing"

func TestAdd(t *testing.T) {
	Initialize()
	todo, err := Add("test")
	if err != nil {
		t.Error("Add returned error")
	}
	if todo == nil {
		t.Error("Add returned nil")
	}
	if todo.desc != "test" {
		t.Error("Add returned unexpected value")
	}
}
