package lock_free_queue

import "testing"

func TestLockFreeQueue(t *testing.T) {
	queue := NewLockFreeQueue[string]()
	queue.Enqueue("a")
}
