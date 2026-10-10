package store

import "strconv"

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
