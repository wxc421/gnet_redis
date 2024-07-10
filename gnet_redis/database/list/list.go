package list

// Consumer traverses list.
// It receives index and value as params, returns true to continue traversal, while returns false to break
type Consumer[T any] func(i int, v T) bool

type List[T any] interface {
	Add(val T)
	Get(index int) (val T)
	Set(index int, val T)
	Insert(index int, val T)
	Remove(index int) (val T)
	RemoveLast() (val T)
	Len() int
	Range(start int, stop int) []T
	ForEach(consumer Consumer[T])
}
