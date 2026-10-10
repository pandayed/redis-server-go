package store

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
