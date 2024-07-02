package queue

import (
	"gnet_redis/internal/queue/lock_free_queue"
	"gnet_redis/internal/queue/slice_queue"
	"math/rand"
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
)

var number = 1000

func BenchmarkSliceQueue(b *testing.B) {
	lkQueue := slice_queue.NewSliceQueue[int](0)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		group := sync.WaitGroup{}
		group.Add(number * 2)
		for i := 0; i < number; i++ {
			go func() {
				defer group.Done()
				lkQueue.Enqueue(1)
			}()
		}
		for i := 0; i < number; i++ {
			go func() {
				defer group.Done()
				lkQueue.Dequeue()
			}()
		}
		group.Wait()
	}
}

func BenchmarkLKQueue(b *testing.B) {
	lkQueue := lock_free_queue.NewLKQueue[int]()

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		group := sync.WaitGroup{}
		group.Add(number * 2)
		for i := 0; i < number; i++ {
			go func() {
				defer group.Done()
				lkQueue.Enqueue(1)
			}()
		}
		for i := 0; i < number; i++ {
			go func() {
				defer group.Done()
				lkQueue.Dequeue()
			}()
		}
		group.Wait()
	}
}

//
// func BenchmarkLKQueue(b *testing.B) {
//
// 	lkQueue := lock_free_queue.NewLKQueue[int]()
//
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			lkQueue.Enqueue(1)
// 			lkQueue.Dequeue()
// 		}
// 	})
// }
//
// func BenchmarkSliceQueue(b *testing.B) {
//
// 	lkQueue := slice_queue.NewSliceQueue[int](0)
//
// 	b.ResetTimer()
// 	b.RunParallel(func(pb *testing.PB) {
// 		for pb.Next() {
// 			lkQueue.Enqueue(1)
// 			lkQueue.Dequeue()
// 		}
// 	})
// }

func BenchmarkQueue(b *testing.B) {
	queues := map[string]Queue[int]{
		"lock-free queue":   lock_free_queue.NewLKQueue[int](),
		"slice-based queue": slice_queue.NewSliceQueue[int](0),
	}

	length := 1 << 12
	inputs := make([]int, length)
	for i := 0; i < length; i++ {
		inputs = append(inputs, rand.Int())
	}

	for _, cpus := range []int{1024} {
		runtime.GOMAXPROCS(cpus)
		for name, q := range queues {
			b.Run(name+"#"+strconv.Itoa(cpus), func(b *testing.B) {
				b.ResetTimer()

				var c int64
				b.RunParallel(func(pb *testing.PB) {
					for pb.Next() {
						i := int(atomic.AddInt64(&c, 1)-1) % length
						v := inputs[i]
						if v >= 0 {
							q.Enqueue(v)
						} else {
							q.Dequeue()
						}
					}
				})
			})
		}
	}
}
