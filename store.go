package main

import (
	"strconv"
	"sync"
)

// Grouped declaration; Common convention for package-level variables
var (
	storeInstance DataStore
	once          sync.Once
)

func GetStore() DataStore {
	once.Do(func() {
		storeInstance = newStore()
	})
	return storeInstance
}

type store struct {
	data map[string]string
	mu   sync.RWMutex
}

type DataStore interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Exists(key string) bool
	Delete(key string) bool
	Incr(key string) (int, error)
	Decr(key string) (int, error)
	Save() error
}

func newStore() DataStore {
	return &store{
		data: make(map[string]string),
	}
}

func (s *store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.data[key] = value
}

func (s *store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.data[key]
	return value, exists
}

func (s *store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	_, exists := s.data[key]
	return exists
}

func (s *store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	_, exists := s.data[key]
	if exists {
		delete(s.data, key)
	}
	return exists
}

func (s *store) Incr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.data[key]
	if !exists {
		s.data[key] = "1"
		return 1, nil
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	num++
	s.data[key] = strconv.Itoa(num)
	return num, nil
}

func (s *store) Decr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.data[key]
	if !exists {
		s.data[key] = "-1"
		return -1, nil
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	num--
	s.data[key] = strconv.Itoa(num)
	return num, nil
}

func (s *store) Save() error {
	// TODO
	return nil
}
