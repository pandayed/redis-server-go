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
	strings map[string]string
	lists   map[string][]string
	sets    map[string]map[string]struct{}
	hashes  map[string]map[string]string
	mu      sync.RWMutex
}

type DataStore interface {
	Set(key, value string)
	Get(key string) (string, bool)
	Exists(key string) bool
	Delete(key string) bool
	Incr(key string) (int, error)
	Decr(key string) (int, error)

	LPush(key string, values ...string) int
	RPush(key string, values ...string) int
	LPop(key string) (string, bool)
	RPop(key string) (string, bool)
	LRange(key string, start, stop int) []string
	LLen(key string) int

	SAdd(key string, members ...string) int
	SMembers(key string) []string
	SIsMember(key string, member string) bool
	SRem(key string, members ...string) int
	SCard(key string) int

	HSet(key, field, value string) int
	HGet(key, field string) (string, bool)
	HGetAll(key string) map[string]string
	HDel(key string, fields ...string) int
	HExists(key, field string) bool
	HLen(key string) int

	Save() error
}

func newStore() DataStore {
	return &store{
		strings: make(map[string]string),
		lists:   make(map[string][]string),
		sets:    make(map[string]map[string]struct{}),
		hashes:  make(map[string]map[string]string),
	}
}

func (s *store) Set(key, value string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.strings[key] = value
}

func (s *store) Get(key string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	value, exists := s.strings[key]
	return value, exists
}

func (s *store) Exists(key string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if _, exists := s.strings[key]; exists {
		return true
	}
	if _, exists := s.lists[key]; exists {
		return true
	}
	if _, exists := s.sets[key]; exists {
		return true
	}
	if _, exists := s.hashes[key]; exists {
		return true
	}
	return false
}

func (s *store) Delete(key string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	found := false
	if _, exists := s.strings[key]; exists {
		delete(s.strings, key)
		found = true
	}
	if _, exists := s.lists[key]; exists {
		delete(s.lists, key)
		found = true
	}
	if _, exists := s.sets[key]; exists {
		delete(s.sets, key)
		found = true
	}
	if _, exists := s.hashes[key]; exists {
		delete(s.hashes, key)
		found = true
	}
	return found
}

func (s *store) Incr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.strings[key]
	if !exists {
		s.strings[key] = "1"
		return 1, nil
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	num++
	s.strings[key] = strconv.Itoa(num)
	return num, nil
}

func (s *store) Decr(key string) (int, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	value, exists := s.strings[key]
	if !exists {
		s.strings[key] = "-1"
		return -1, nil
	}

	num, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}

	num--
	s.strings[key] = strconv.Itoa(num)
	return num, nil
}

func (s *store) LPush(key string, values ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := s.lists[key]
	for i := len(values) - 1; i >= 0; i-- {
		list = append([]string{values[i]}, list...)
	}
	s.lists[key] = list
	return len(list)
}

func (s *store) RPush(key string, values ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	list := s.lists[key]
	list = append(list, values...)
	s.lists[key] = list
	return len(list)
}

func (s *store) LPop(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, exists := s.lists[key]
	if !exists || len(list) == 0 {
		return "", false
	}

	value := list[0]
	s.lists[key] = list[1:]

	if len(s.lists[key]) == 0 {
		delete(s.lists, key)
	}

	return value, true
}

func (s *store) RPop(key string) (string, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()

	list, exists := s.lists[key]
	if !exists || len(list) == 0 {
		return "", false
	}

	value := list[len(list)-1]
	s.lists[key] = list[:len(list)-1]

	if len(s.lists[key]) == 0 {
		delete(s.lists, key)
	}

	return value, true
}

func (s *store) LRange(key string, start, stop int) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list, exists := s.lists[key]
	if !exists {
		return []string{}
	}

	length := len(list)

	if start < 0 {
		start = length + start
	}
	if stop < 0 {
		stop = length + stop
	}

	if start < 0 {
		start = 0
	}
	if stop >= length {
		stop = length - 1
	}

	if start > stop || start >= length {
		return []string{}
	}

	result := make([]string, stop-start+1)
	copy(result, list[start:stop+1])
	return result
}

func (s *store) LLen(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	list, exists := s.lists[key]
	if !exists {
		return 0
	}
	return len(list)
}

func (s *store) SAdd(key string, members ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	set, exists := s.sets[key]
	if !exists {
		set = make(map[string]struct{})
		s.sets[key] = set
	}

	added := 0
	for _, member := range members {
		if _, exists := set[member]; !exists {
			set[member] = struct{}{}
			added++
		}
	}
	return added
}

func (s *store) SMembers(key string) []string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	set, exists := s.sets[key]
	if !exists {
		return []string{}
	}

	members := make([]string, 0, len(set))
	for member := range set {
		members = append(members, member)
	}
	return members
}

func (s *store) SIsMember(key string, member string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	set, exists := s.sets[key]
	if !exists {
		return false
	}

	_, isMember := set[member]
	return isMember
}

func (s *store) SRem(key string, members ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	set, exists := s.sets[key]
	if !exists {
		return 0
	}

	removed := 0
	for _, member := range members {
		if _, exists := set[member]; exists {
			delete(set, member)
			removed++
		}
	}

	if len(set) == 0 {
		delete(s.sets, key)
	}

	return removed
}

func (s *store) SCard(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	set, exists := s.sets[key]
	if !exists {
		return 0
	}
	return len(set)
}

func (s *store) HSet(key, field, value string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, exists := s.hashes[key]
	if !exists {
		hash = make(map[string]string)
		s.hashes[key] = hash
	}

	_, existed := hash[field]
	hash[field] = value

	if existed {
		return 0
	}
	return 1
}

func (s *store) HGet(key, field string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hash, exists := s.hashes[key]
	if !exists {
		return "", false
	}

	value, exists := hash[field]
	return value, exists
}

func (s *store) HGetAll(key string) map[string]string {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hash, exists := s.hashes[key]
	if !exists {
		return map[string]string{}
	}

	result := make(map[string]string, len(hash))
	for field, value := range hash {
		result[field] = value
	}
	return result
}

func (s *store) HDel(key string, fields ...string) int {
	s.mu.Lock()
	defer s.mu.Unlock()

	hash, exists := s.hashes[key]
	if !exists {
		return 0
	}

	deleted := 0
	for _, field := range fields {
		if _, exists := hash[field]; exists {
			delete(hash, field)
			deleted++
		}
	}

	if len(hash) == 0 {
		delete(s.hashes, key)
	}

	return deleted
}

func (s *store) HExists(key, field string) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hash, exists := s.hashes[key]
	if !exists {
		return false
	}

	_, exists = hash[field]
	return exists
}

func (s *store) HLen(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	hash, exists := s.hashes[key]
	if !exists {
		return 0
	}
	return len(hash)
}

func (s *store) Save() error {
	// TODO
	return nil
}
