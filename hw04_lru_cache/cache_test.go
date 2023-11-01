package hw04lrucache

import (
	"math/rand"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCache(t *testing.T) {
	t.Run("empty cache", func(t *testing.T) {
		c := NewCache(10)

		_, ok := c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)
	})

	t.Run("simple", func(t *testing.T) {
		c := NewCache(5)

		wasInCache := c.Set("aaa", 100)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 200)
		require.False(t, wasInCache)

		val, ok := c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 100, val)

		val, ok = c.Get("bbb")
		require.True(t, ok)
		require.Equal(t, 200, val)

		wasInCache = c.Set("aaa", 300)
		require.True(t, wasInCache)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 300, val)

		val, ok = c.Get("ccc")
		require.False(t, ok)
		require.Nil(t, val)
	})

	t.Run("purge logic: cache capacity overflow", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 10)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 20)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 30)
		require.False(t, wasInCache)

		wasInCache = c.Set("ddd", 40)
		require.False(t, wasInCache)

		_, ok := c.Get("aaa")
		require.False(t, ok)
	})

	t.Run("purge logic: long-used elements removed first", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 10)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 20)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 30)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 22)
		require.True(t, wasInCache)

		wasInCache = c.Set("aaa", 11)
		require.True(t, wasInCache)

		val, ok := c.Get("ccc")
		require.True(t, ok)
		require.Equal(t, 30, val)

		wasInCache = c.Set("ddd", 40)
		require.False(t, wasInCache)

		_, ok = c.Get("bbb")
		require.False(t, ok)

		val, ok = c.Get("aaa")
		require.True(t, ok)
		require.Equal(t, 11, val)

		wasInCache = c.Set("eee", 50)
		require.False(t, wasInCache)

		_, ok = c.Get("ccc")
		require.False(t, ok)
	})

	t.Run("cache clear", func(t *testing.T) {
		c := NewCache(3)

		wasInCache := c.Set("aaa", 10)
		require.False(t, wasInCache)

		wasInCache = c.Set("bbb", 20)
		require.False(t, wasInCache)

		wasInCache = c.Set("ccc", 30)
		require.False(t, wasInCache)

		_, ok := c.Get("aaa")
		require.True(t, ok)

		_, ok = c.Get("bbb")
		require.True(t, ok)

		_, ok = c.Get("ccc")
		require.True(t, ok)

		c.Clear()

		_, ok = c.Get("aaa")
		require.False(t, ok)

		_, ok = c.Get("bbb")
		require.False(t, ok)

		_, ok = c.Get("ccc")
		require.False(t, ok)
	})
}

func TestCacheMultithreading(_ *testing.T) {
	c := NewCache(10)
	wg := &sync.WaitGroup{}
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Set(Key(strconv.Itoa(i)), i)
		}
	}()

	go func() {
		defer wg.Done()
		for i := 0; i < 1_000_000; i++ {
			c.Get(Key(strconv.Itoa(rand.Intn(1_000_000))))
		}
	}()

	wg.Wait()
}
