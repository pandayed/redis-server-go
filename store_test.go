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

func TestStore_Get_NonExistent(t *testing.T) {
	store := newStore()

	_, exists := store.Get("nonexistent")
	if exists {
		t.Error("Expected nonexistent key to not exist")
	}
}

func TestStore_Exists(t *testing.T) {
	store := newStore()

	if store.Exists("key1") {
		t.Error("Expected key1 to not exist initially")
	}

	store.Set("key1", "value1")

	if !store.Exists("key1") {
		t.Error("Expected key1 to exist after setting")
	}
}

func TestStore_Delete(t *testing.T) {
	store := newStore()

	deleted := store.Delete("nonexistent")
	if deleted {
		t.Error("Expected delete of nonexistent key to return false")
	}

	store.Set("key1", "value1")
	deleted = store.Delete("key1")
	if !deleted {
		t.Error("Expected delete of existing key to return true")
	}

	if store.Exists("key1") {
		t.Error("Expected key1 to not exist after deletion")
	}
}

func TestStore_Incr(t *testing.T) {
	store := newStore()

	result, err := store.Incr("counter")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 1 {
		t.Errorf("Expected 1, got %d", result)
	}

	result, err = store.Incr("counter")
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
	if result != 2 {
		t.Errorf("Expected 2, got %d", result)
	}
}

func TestStore_Incr_InvalidValue(t *testing.T) {
	store := newStore()

	store.Set("invalid", "not_a_number")

	_, err := store.Incr("invalid")
	if err == nil {
		t.Error("Expected error when incrementing non-numeric value")
	}
}
