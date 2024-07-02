package queue

import (
	"fmt"
	"gnet_redis/internal/queue/lock_free_queue"
	"gnet_redis/internal/queue/slice_queue"
	"testing"
)

func TestQueue(t *testing.T) {
	count := 100
	queues := map[string]Queue[int]{
		"lock_free_queue": lock_free_queue.NewLKQueue[int](),
		"slice_queue":     slice_queue.NewSliceQueue[int](0),
	}
	for name, queue := range queues {
		t.Run(name, func(t *testing.T) {
			t.Log(fmt.Sprintf("%v start test...", name))
			for i := 1; i <= count; i++ {
				queue.Enqueue(i)
			}
			for i := 1; i <= count; i++ {
				value := queue.Dequeue()
				if value == 0 {
					t.Fatalf("got a nil value")
				}
				if value != i {
					t.Fatalf("expect %d but got %v", i, value)
				}
			}
		})
	}
}
