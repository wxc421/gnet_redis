package queue

type Queue[T any] interface {
	Enqueue(v T)
	Dequeue() T
}
