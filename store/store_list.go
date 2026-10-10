package store

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
