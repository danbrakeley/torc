package eff

import (
	"testing"
)

func TestEff_New(t *testing.T) {
	e := New()
	if e == nil {
		t.Errorf("Expected New() to not return nil")
	}
}

func TestEff_NewHasStack(t *testing.T) {
	e := New()
	if len(e.Stack()) == 0 {
		t.Errorf("Expected New() to fill in a stack trace")
	}
}

func TestEff_NewHasNoFields(t *testing.T) {
	e := New()
	for k, _ := range e.Fields() {
		if k != Field_Stack {
			t.Errorf("Expected New() to start with only known fields")
		}
	}
}

func TestEff_NewWithFieldAddsField(t *testing.T) {
	e := New()
	startingCount := len(e.Fields())
	e = e.WithField("test", "value")
	if len(e.Fields()) != startingCount+1 {
		t.Errorf("Expected WithField to add 1 field")
	}
	if e.Fields()["test"].(string) != "value" {
		t.Errorf("Expected \"test\" field to be \"value\", not \"%v\"", e.Fields()["test"])
	}
}

func TestEff_NewWithFieldDoesNotAlterOriginal(t *testing.T) {
	e := New()
	f := e.WithField("test", "value")
	if len(e.Fields()) != len(f.Fields())-1 {
		t.Errorf("Expected WithField to not alter original Eff")
	}
}

func TestEff_NewWithErrorIncorporatesEffFields(t *testing.T) {
	e := New().WithField("test", "value")
	f := NewErr(e)
	if f.Fields()["test"].(string) != "value" {
		t.Errorf("Expected \"test\" field to be \"value\", not \"%v\"", e.Fields()["test"])
	}
}
