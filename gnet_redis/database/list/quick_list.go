package list

const pageSize = 1024

// QuickList is a linked list of page (which type is []T)
// QuickList has better performance than LinkedList of Add, Range and memory usage
type QuickList[T any] struct {
	linkList *LinkedList[[]T]
	size     int
}

// iterator of QuickList, move between [-1, ql.Len()]
type iterator[T any] struct {
	n      *node[[]T]
	offset int
	q      *QuickList[T]
}

// lazyInit lazily initializes a zero List value.
func (q *QuickList[T]) lazyInit() {
	if q.linkList == nil {
		q.init()
	}
}

func (q *QuickList[T]) Add(val T) {
	if q.Len() == 0 {
		q.lazyInit()
		page := q.linkList.find(0).val
		page = append(page, val)
	}
	q.linkList.ForEach(func(i int, v []T) bool {
		if len(v) < pageSize {
			v = append(v, val)
			return false
		}
		return true
	})
}

func (q *QuickList[T]) find(index int) *iterator[T] {
	if index < 0 || index >= q.size {
		return nil
	}
	var (
		pageOffset            = 0
		n          *node[[]T] = nil
		count                 = 0
	)

	if index <= q.size<<1 {
		// search from head
		n = q.linkList.first
		for {
			count += len(n.val)
			if count >= index {
				break
			}
			n = n.next
		}
	} else {
		// search from tail
		n = q.linkList.last
		for {
			count += len(n.val)
			if count >= index {
				break
			}
			n = n.prev
		}
	}

	pageOffset = index - (count - pageSize)
	return &iterator[T]{
		n:      n,
		offset: pageOffset,
		q:      q,
	}
}

func (q *QuickList[T]) Get(index int) (val T) {
	// TODO implement me
	panic("implement me")
}

func (q *QuickList[T]) Set(index int, val T) {
	// TODO implement me
	panic("implement me")
}

func (q *QuickList[T]) Insert(index int, val T) {
	// TODO implement me
	panic("implement me")
}

func (q *QuickList[T]) Remove(index int) (val T) {
	// TODO implement me
	panic("implement me")
}

func (q *QuickList[T]) RemoveLast() (val T) {
	// TODO implement me
	panic("implement me")
}

func (q *QuickList[T]) Len() int {
	return q.size
}

func (q *QuickList[T]) Range(start int, stop int) []T {
	// TODO implement me
	panic("implement me")
}

func (q *QuickList[T]) ForEach(consumer Consumer[T]) {
	// TODO implement me
	panic("implement me")
}

func (q *QuickList[T]) init() {
	page := make([]T, 0, pageSize)
	q.linkList = Make(page)
}

var _ List[any] = (*QuickList[any])(nil)
