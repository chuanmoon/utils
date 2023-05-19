package collection

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type Set[T comparable] map[T]struct{}

func NewSet[T comparable](v ...T) Set[T] {
	r := make(Set[T])
	r.Append(v...)
	return r
}

func NewSetWithSize[T comparable](cardinality int) Set[T] {
	return make(Set[T], cardinality)
}

func (s Set[T]) Add(v T) bool {
	prevLen := len(s)
	s[v] = struct{}{}
	return prevLen != len(s)
}

func (s Set[T]) Append(v ...T) int {
	prevLen := len(s)
	for _, val := range v {
		(s)[val] = struct{}{}
	}
	return len(s) - prevLen
}

// private version of Add which doesn't return a value
func (s Set[T]) add(v T) {
	s[v] = struct{}{}
}

func (s Set[T]) Len() int {
	return len(s)
}

func (s Set[T]) Clear() {
	for key := range s {
		delete(s, key)
	}
}

func (s Set[T]) Clone() Set[T] {
	clonedSet := NewSetWithSize[T](s.Len())
	for elem := range s {
		clonedSet.add(elem)
	}
	return clonedSet
}

func (s Set[T]) Contains(v ...T) bool {
	for _, val := range v {
		if _, ok := s[val]; !ok {
			return false
		}
	}
	return true
}

// private version of Contains for a single element v
func (s Set[T]) contains(v T) (ok bool) {
	_, ok = s[v]
	return ok
}

func (s Set[T]) Difference(other Set[T]) Set[T] {
	o := other

	diff := NewSet[T]()
	for elem := range s {
		if !o.contains(elem) {
			diff.add(elem)
		}
	}
	return diff
}

func (s Set[T]) Each(cb func(T) bool) {
	for elem := range s {
		if cb(elem) {
			break
		}
	}
}

func (s Set[T]) Equal(other Set[T]) bool {
	o := other

	if s.Len() != other.Len() {
		return false
	}
	for elem := range s {
		if !o.contains(elem) {
			return false
		}
	}
	return true
}

func (s Set[T]) Intersect(other Set[T]) Set[T] {
	o := other

	intersection := NewSet[T]()
	// loop over smaller set
	if s.Len() < other.Len() {
		for elem := range s {
			if o.contains(elem) {
				intersection.add(elem)
			}
		}
	} else {
		for elem := range o {
			if s.contains(elem) {
				intersection.add(elem)
			}
		}
	}
	return intersection
}

func (s Set[T]) IsProperSubset(other Set[T]) bool {
	return s.Len() < other.Len() && s.IsSubset(other)
}

func (s Set[T]) IsProperSuperset(other Set[T]) bool {
	return s.Len() > other.Len() && s.IsSuperset(other)
}

func (s Set[T]) IsSubset(other Set[T]) bool {
	o := other
	if s.Len() > other.Len() {
		return false
	}
	for elem := range s {
		if !o.contains(elem) {
			return false
		}
	}
	return true
}

func (s Set[T]) IsSuperset(other Set[T]) bool {
	return other.IsSubset(s)
}

func (s Set[T]) Iter() <-chan T {
	ch := make(chan T)
	go func() {
		for elem := range s {
			ch <- elem
		}
		close(ch)
	}()

	return ch
}

// Pop returns a popped item in case set is not empty, or nil-value of T
// if set is already empty
func (s Set[T]) Pop() (v T, ok bool) {
	for item := range s {
		delete(s, item)
		return item, true
	}
	return v, false
}

func (s Set[T]) Remove(v T) {
	delete(s, v)
}

func (s Set[T]) RemoveAll(i ...T) {
	for _, elem := range i {
		delete(s, elem)
	}
}

func (s Set[T]) String() string {
	items := make([]string, 0, len(s))

	for elem := range s {
		items = append(items, fmt.Sprintf("%v", elem))
	}
	return fmt.Sprintf("Set{%s}", strings.Join(items, ", "))
}

func (s Set[T]) SymmetricDifference(other Set[T]) Set[T] {
	o := other

	sd := NewSet[T]()
	for elem := range s {
		if !o.contains(elem) {
			sd.add(elem)
		}
	}
	for elem := range o {
		if !s.contains(elem) {
			sd.add(elem)
		}
	}
	return sd
}

func (s Set[T]) ToSlice() []T {
	keys := make([]T, 0, s.Len())
	for elem := range s {
		keys = append(keys, elem)
	}

	return keys
}

func (s Set[T]) Union(other Set[T]) Set[T] {
	o := other

	n := s.Len()
	if o.Len() > n {
		n = o.Len()
	}
	unionedSet := make(Set[T], n)

	for elem := range s {
		unionedSet.add(elem)
	}
	for elem := range o {
		unionedSet.add(elem)
	}
	return unionedSet
}

// MarshalJSON creates a JSON array from the set, it marshals all elements
func (s Set[T]) MarshalJSON() ([]byte, error) {
	items := make([]string, 0, s.Len())

	for elem := range s {
		b, err := json.Marshal(elem)
		if err != nil {
			return nil, err
		}

		items = append(items, string(b))
	}

	return []byte(fmt.Sprintf("[%s]", strings.Join(items, ","))), nil
}

// UnmarshalJSON recreates a set from a JSON array, it only decodes
// primitive types. Numbers are decoded as json.Number.
func (s Set[T]) UnmarshalJSON(b []byte) error {
	var i []any

	d := json.NewDecoder(bytes.NewReader(b))
	d.UseNumber()
	err := d.Decode(&i)
	if err != nil {
		return err
	}

	for _, v := range i {
		switch t := v.(type) {
		case T:
			s.add(t)
		default:
			// anything else must be skipped.
			continue
		}
	}

	return nil
}
