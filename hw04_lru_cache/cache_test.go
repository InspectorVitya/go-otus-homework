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

	t.Run("purge logic by capacity", func(t *testing.T) {
		keys := []Key{"aaa", "bbb", "ccc", "ddd"}
		c := NewCache(len(keys) - 1)

		for _, key := range keys {
			wasInCache := c.Set(key, 0)
			require.False(t, wasInCache)
		}

		val, ok := c.Get(keys[0])
		require.False(t, ok)
		require.Equal(t, nil, val)

		for _, cache := range keys[1:] {
			val, ok = c.Get(cache)
			require.True(t, ok)
			require.NotEqual(t, nil, val)
		}
	})

	t.Run("purge logic rarely used", func(t *testing.T) {
		keys := []Key{"aaa", "bbb", "ccc"}
		c := NewCache(3)

		for _, key := range keys {
			wasInCache := c.Set(key, 0)
			require.False(t, wasInCache)
		}

		_, ok := c.Get(keys[0])
		require.True(t, ok)
		_, ok = c.Get(keys[0])
		require.True(t, ok)
		ok = c.Set(keys[1], 0)
		require.True(t, ok)

		wasInCache := c.Set("ddd", 0)
		require.False(t, wasInCache)

		val, ok := c.Get(keys[2])
		require.False(t, ok)
		require.Equal(t, nil, val)
	})
}

func TestCacheMultithreading(t *testing.T) {
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