package store

import "sync"

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
}

// New creates an independent in-memory store.
func New() DataStore {
	return &store{
		strings: make(map[string]string),
		lists:   make(map[string][]string),
		sets:    make(map[string]map[string]struct{}),
		hashes:  make(map[string]map[string]string),
	}
}
