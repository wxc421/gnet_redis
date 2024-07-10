package list

// LinkedList is doubly linked list
type LinkedList[T any] struct {
	first *node[T]
	last  *node[T]
	size  int
}

type node[T any] struct {
	val  T
	prev *node[T]
	next *node[T]
}

// Add adds value to the tail
func (list *LinkedList[T]) Add(val T) {
	if list == nil {
		panic("list is nil")
	}
	n := &node[T]{
		val: val,
	}
	if list.last == nil {
		// empty list
		list.first = n
		list.last = n
	} else {
		n.prev = list.last
		list.last.next = n
		list.last = n
	}
	list.size++
}

func (list *LinkedList[T]) find(index int) (n *node[T]) {
	if index < list.size/2 {
		n = list.first
		for i := 0; i < index; i++ {
			n = n.next
		}
	} else {
		n = list.last
		for i := list.size - 1; i > index; i-- {
			n = n.prev
		}
	}
	return n
}

// Get returns value at the given index
func (list *LinkedList[T]) Get(index int) (val T) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index >= list.size {
		panic("index out of bound")
	}
	return list.find(index).val
}

// Set updates value at the given index, the index should between [0, list.size]
func (list *LinkedList[T]) Set(index int, val T) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index > list.size {
		panic("index out of bound")
	}
	n := list.find(index)
	n.val = val
}

// Insert inserts value at the given index, the original element at the given index will move backward
func (list *LinkedList[T]) Insert(index int, val T) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index > list.size {
		panic("index out of bound")
	}

	if index == list.size {
		list.Add(val)
		return
	}
	// list is not empty
	pivot := list.find(index)
	n := &node[T]{
		val:  val,
		prev: pivot.prev,
		next: pivot,
	}
	if pivot.prev == nil {
		list.first = n
	} else {
		pivot.prev.next = n
	}
	pivot.prev = n
	list.size++
}

func (list *LinkedList[T]) removeNode(n *node[T]) {
	if n.prev == nil {
		list.first = n.next
	} else {
		n.prev.next = n.next
	}
	if n.next == nil {
		list.last = n.prev
	} else {
		n.next.prev = n.prev
	}

	// for gc
	n.prev = nil
	n.next = nil

	list.size--
}

// Remove removes value at the given index
func (list *LinkedList[T]) Remove(index int) (val interface{}) {
	if list == nil {
		panic("list is nil")
	}
	if index < 0 || index >= list.size {
		panic("index out of bound")
	}

	n := list.find(index)
	list.removeNode(n)
	return n.val
}

// RemoveLast removes the last element and returns its value
func (list *LinkedList[T]) RemoveLast() (val interface{}) {
	if list == nil {
		panic("list is nil")
	}
	if list.last == nil {
		// empty list
		return nil
	}
	n := list.last
	list.removeNode(n)
	return n.val
}

// Len returns the number of elements in list
func (list *LinkedList[T]) Len() int {
	if list == nil {
		panic("list is nil")
	}
	return list.size
}

// Range returns elements which index within [start, stop)
func (list *LinkedList[T]) Range(start int, stop int) []T {
	if list == nil {
		panic("list is nil")
	}
	if start < 0 || start >= list.size {
		panic("`start` out of range")
	}
	if stop < start || stop > list.size {
		panic("`stop` out of range")
	}

	sliceSize := stop - start
	slice := make([]T, sliceSize)
	n := list.first
	i := 0
	sliceIndex := 0
	for n != nil {
		if i >= start && i < stop {
			slice[sliceIndex] = n.val
			sliceIndex++
		} else if i >= stop {
			break
		}
		i++
		n = n.next
	}
	return slice
}

// ForEach visits each element in the list
// if the consumer returns false, the loop will be break
func (list *LinkedList[T]) ForEach(consumer Consumer[T]) {
	if list == nil {
		panic("list is nil")
	}
	n := list.first
	i := 0
	for n != nil {
		goNext := consumer(i, n.val)
		if !goNext {
			break
		}
		i++
		n = n.next
	}
}

// Make creates a new linked list
func Make[T any](vals ...T) *LinkedList[T] {
	list := LinkedList[T]{}
	for _, v := range vals {
		list.Add(v)
	}
	return &list
}
