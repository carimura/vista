package main

import "testing"

func TestSomething(t *testing.T) {
	if false {
		t.Error("Expected no errors, but got one.")
	}
}
