package list

import (
	"strconv"
	"testing"
)

func TestGet(t *testing.T) {
	list := Make[int]()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	for i := 0; i < 10; i++ {
		v := list.Get(i)
		if i != v {
			t.Error("get test fail: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(v))
		}
	}
}

func TestRemove(t *testing.T) {
	list := Make[int]()
	for i := 0; i < 10; i++ {
		list.Add(i)
	}
	for i := 9; i >= 0; i-- {
		list.Remove(i)
		if i != list.Len() {
			t.Error("remove test fail: expected size " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(list.Len()))
		}
		list.ForEach(func(i int, v int) bool {
			if v != i {
				t.Error("remove test fail: expected " + strconv.Itoa(i) + ", actual: " + strconv.Itoa(v))
			}
			return true
		})
	}
}
