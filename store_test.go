package main

import (
	"testing"
)

func TestStore_Set_Get(t *testing.T) {
	store := newStore()

	store.Set("key1", "value1")

	value, exists := store.Get("key1")
	if !exists {
		t.Error("Expected key1 to exist")
	}
	if value != "value1" {
		t.Errorf("Expected value1, got %s", value)
	}
}
