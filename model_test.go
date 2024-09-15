package main

import (
	"errors"
	"testing"
)

func TestModelViewInit(t *testing.T) {
	m := initialModel("", "", "")
	view := m.View()

	if view != "Initializing..." {
		t.Errorf("View should return initializing message for initial instance")
	}
}

func TestModelViewError(t *testing.T) {
	m := initialModel("", "", "")
	m.err = errors.New("ERROR")
	view := m.View()

	if view != "Something went wrong: ERROR" {
		t.Errorf("View should return error message when error is ocured")
	}
}
