package store

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
