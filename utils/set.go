package utils

// 自定义 Set, 非并发安全
type Set[T comparable] struct {
	data map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
	return &Set[T]{
		data: make(map[T]struct{}),
	}
}

func (s *Set[T]) Add(key T) {
	s.data[key] = struct{}{}
}

func (s *Set[T]) Del(key T) {
	delete(s.data, key)
}

func (s *Set[T]) Exist(key T) bool {
	_, ok := s.data[key]
	return ok
}

func (s *Set[T]) ToList() []T {
	keys := make([]T, 0, len(s.data))
	for key := range s.data {
		keys = append(keys, key)
	}
	return keys
}
