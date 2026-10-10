package main

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
