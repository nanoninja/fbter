package main

import "fmt"

type Iterator [T any] interface {
	Valid() bool
	Current() any
	Key() int
	Next()
	Rewind()
}

type Slice[T any] struct {
	data []T
	pos int
}

func (s Slice[T]) Valid() bool {
	return s.pos < len(s.data)
}

func (s Slice[T]) Current() any {
	return s.data[s.pos]
}

func (s Slice[T]) Key() int {
	return s.pos
}

func (s *Slice[T]) Next() {
	s.pos++
}

func (s *Slice[T]) Rewind() {
	s.pos = 0
}

func Looping(iter Iterator[any]) {
	for iter.Valid() {
		fmt.Println(iter.Key(), iter.Current())
		iter.Next()
	}
}

type NyNumber uint8

type NumberConstraint interface {
	~uint8 | uint16
}

func Display[T NumberConstraint](v T) {
	fmt.Println(v)
}

func main() {
	var v NyNumber = 10
	Display(v)

	var iter Iterator[any]

	iter = &Slice[string]{
		data: []string{"Hello", "World", "Gopher"},
	}

	for iter.Valid() {
		fmt.Println(iter.Key(), iter.Current())
		iter.Next()
	}

	strs := &Slice[string]{
		data: []string{"Hello", "World", "Gopher"},
	}

	nums := &Slice[int]{
		data: []int{12, 48, 90},
	}

	Looping(strs)
	Looping(nums)
}