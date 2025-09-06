package main

import "fmt"

// Queue is a simple FIFO queue
type Queue[T any] struct {
	items []T
}

// Enqueue adds an item at the end
func (q *Queue[T]) Enqueue(item T) {
	q.items = append(q.items, item)
}

// Dequeue removes and returns the first item
func (q *Queue[T]) Dequeue() (T, bool) {
	var zero T
	if len(q.items) == 0 {
		return zero, false
	}
	item := q.items[0]
	q.items = q.items[1:]
	return item, true
}

// IsEmpty checks if the queue is empty
func (q *Queue[T]) IsEmpty() bool {
	return len(q.items) == 0
}

// Size returns the number of items
func (q *Queue[T]) Size() int {
	return len(q.items)
}

// PrintAll prints all elements in the queue
func (q *Queue[T]) PrintAll() {
	if len(q.items) == 0 {
		fmt.Println("Queue is empty")
		return
	}
	fmt.Println("Queue elements:")
	for _, item := range q.items {
		fmt.Println(item)
	}
}
