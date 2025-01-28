package selectlang

type Stack[T any] struct {
	stack []T
}

func (s *Stack[T]) Push(item T) {
	s.stack = append(s.stack, item)
}

func (s *Stack[T]) Pop() T {
	if len(s.stack) == 0 {
		var zero T

		return zero
	}

	item := s.stack[len(s.stack)-1]
	s.stack = s.stack[:len(s.stack)-1]

	return item
}

// Update will either push a new item if the stack is empty or
// update the current top item.
//
// This is a shortcut of doing: `itm := Pop()`, modify, `Push(itm)`
func (s *Stack[T]) Update(f func(item T) T) {
	if len(s.stack) == 0 {
		var zero T

		s.Push(f(zero))
	}

	s.stack[len(s.stack)-1] = f(s.stack[len(s.stack)-1])
}

// Peek will return the top item without removing it from the stack.
//
// If the stack is empty it will return an empty item.
func (s *Stack[T]) Peek() T {
	if len(s.stack) == 0 {
		var zero T

		return zero
	}

	return s.stack[len(s.stack)-1]
}
