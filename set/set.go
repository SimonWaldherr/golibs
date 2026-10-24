// Package set provides small generic sets using Go equality for membership.
package set

// Set stores unique comparable values. Its zero value is ready to use.
// Sets are not safe for concurrent mutation. Copying a Set shares any initialized
// storage; use Clone for an independent set. Keys themselves are copied shallowly.
type Set[T comparable] struct {
	values map[T]struct{}
}

// New returns a set containing values, ignoring repeated equal values.
func New[T comparable](values ...T) Set[T] {
	var result Set[T]
	result.Add(values...)
	return result
}

// Add inserts values. An empty call does not allocate storage.
func (s *Set[T]) Add(values ...T) {
	if len(values) == 0 {
		return
	}
	if s.values == nil {
		s.values = make(map[T]struct{}, len(values))
	}
	for _, value := range values {
		s.values[value] = struct{}{}
	}
}

// Remove deletes values. Missing values are ignored.
func (s *Set[T]) Remove(values ...T) {
	for _, value := range values {
		delete(s.values, value)
	}
}

// Has reports whether value belongs to the set.
func (s Set[T]) Has(value T) bool {
	_, exists := s.values[value]
	return exists
}

// Len returns the number of stored values.
func (s Set[T]) Len() int { return len(s.values) }

// Values returns an independent slice in unspecified order, or nil for an empty set.
func (s Set[T]) Values() []T {
	var result []T
	if len(s.values) > 0 {
		result = make([]T, 0, len(s.values))
		for value := range s.values {
			result = append(result, value)
		}
	}
	return result
}

// Clone returns an independent set with the same keys.
func (s Set[T]) Clone() Set[T] {
	var result Set[T]
	if len(s.values) > 0 {
		result.values = make(map[T]struct{}, len(s.values))
		for value := range s.values {
			result.values[value] = struct{}{}
		}
	}
	return result
}

// Union returns an independent set containing values in either set.
func (s Set[T]) Union(other Set[T]) Set[T] {
	if len(s.values) == 0 {
		return other.Clone()
	}
	result := s.Clone()
	for value := range other.values {
		result.Add(value)
	}
	return result
}

// Intersection returns an independent set containing values present in both sets.
func (s Set[T]) Intersection(other Set[T]) Set[T] {
	if s.Len() > other.Len() {
		s, other = other, s
	}
	var result Set[T]
	for value := range s.values {
		if other.Has(value) {
			result.Add(value)
		}
	}
	return result
}

// Difference returns an independent set containing values in s but not other.
func (s Set[T]) Difference(other Set[T]) Set[T] {
	if len(other.values) == 0 {
		return s.Clone()
	}
	var result Set[T]
	for value := range s.values {
		if !other.Has(value) {
			result.Add(value)
		}
	}
	return result
}
