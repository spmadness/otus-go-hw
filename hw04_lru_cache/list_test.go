package hw04lrucache

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) {
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, 3, l.Len())

		middle := l.Front().Next // 20
		l.Remove(middle)         // [10, 30]
		require.Equal(t, 2, l.Len())

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, 7, l.Len())
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Front(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{70, 80, 60, 40, 10, 30, 50}, elems)
	})

	t.Run("push front elements order is reversed in list", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushFront(20) // [20, 10]
		l.PushFront(30) // [30, 20, 10]

		require.Equal(t, 3, l.Len())
		require.Equal(t, 30, l.Front().Value)
		require.Equal(t, 10, l.Back().Value)
	})

	t.Run("move last element to front", func(t *testing.T) {
		l := NewList()

		l.PushBack(10) // [10]
		l.PushBack(20) // [10, 20]
		l.PushBack(30) // [10, 20, 30]
		require.Equal(t, 3, l.Len())
		require.Equal(t, 10, l.Front().Value)

		l.MoveToFront(l.Back()) // [30, 10, 20]
		require.Equal(t, 30, l.Front().Value)
		require.Equal(t, 20, l.Back().Value)
	})

	t.Run("remove elements", func(t *testing.T) {
		l := NewList()

		l.PushBack(10) // [10]
		l.PushBack(20) // [10, 20]
		l.PushBack(30) // [10, 20, 30]
		l.PushBack(40) // [10, 20, 30, 40]
		l.PushBack(50) // [10, 20, 30, 40, 50]
		require.Equal(t, 5, l.Len())

		l.Remove(l.Front().Next) // [10, 30, 40, 50]
		l.Remove(l.Back().Prev)  // [10, 30, 50]

		require.Equal(t, 3, l.Len())
		require.Equal(t, 10, l.Front().Value)
		require.Equal(t, 30, l.Front().Next.Value)
		require.Equal(t, 50, l.Back().Value)

		l.Remove(l.Back()) // [10, 30]
		l.Remove(l.Back()) // [10]
		l.Remove(l.Back()) // []

		require.Equal(t, 0, l.Len())
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
}
