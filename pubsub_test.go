// Copyright 2013, Chandra Sekar S.  All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the README.md file.

package pubsub

import (
	"runtime"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestSub(t *testing.T) {
	ps := New(1)
	ch1 := ps.Sub("t1")
	ch2 := ps.Sub("t1")
	ch3 := ps.Sub("t2")

	ps.Pub("hi", "t1")
	require.Equal(t, <-ch1, "hi")
	require.Equal(t, <-ch2, "hi")

	ps.Pub("hello", "t2")
	require.Equal(t, <-ch3, "hello")

	ps.Shutdown()
	_, ok := <-ch1
	require.Equal(t, ok, false)
	_, ok = <-ch2
	require.Equal(t, ok, false)
	_, ok = <-ch3
	require.Equal(t, ok, false)
}

func TestSubOnce(t *testing.T) {
	ps := New(1)
	ch := ps.SubOnce("t1")

	ps.Pub("hi", "t1")
	require.Equal(t, <-ch, "hi")

	_, ok := <-ch
	require.Equal(t, ok, false)
	ps.Shutdown()
}

func TestAddSub(t *testing.T) {
	ps := New(1)
	ch1 := ps.Sub("t1")
	ch2 := ps.Sub("t2")

	ps.Pub("hi1", "t1")
	require.Equal(t, <-ch1, "hi1")

	ps.Pub("hi2", "t2")
	require.Equal(t, <-ch2, "hi2")

	ps.AddSub(ch1, "t2", "t3")
	ps.Pub("hi3", "t2")
	require.Equal(t, <-ch1, "hi3")
	require.Equal(t, <-ch2, "hi3")

	ps.Pub("hi4", "t3")
	require.Equal(t, <-ch1, "hi4")

	ps.Shutdown()
}

func TestUnsub(t *testing.T) {
	ps := New(1)
	ch := ps.Sub("t1")

	ps.Pub("hi", "t1")
	require.Equal(t, <-ch, "hi")

	ps.Unsub(ch, "t1")
	_, ok := <-ch
	require.Equal(t, ok, false)
	ps.Shutdown()
}

func TestUnsubAll(t *testing.T) {
	ps := New(1)
	ch1 := ps.Sub("t1", "t2", "t3")
	ch2 := ps.Sub("t1", "t3")

	ps.Unsub(ch1)

	m, ok := <-ch1
	require.Equal(t, ok, false)

	ps.Pub("hi", "t1")
	m, ok = <-ch2
	require.Equal(t, m, "hi")

	ps.Shutdown()
}

func TestClose(t *testing.T) {
	ps := New(1)
	ch1 := ps.Sub("t1")
	ch2 := ps.Sub("t1")
	ch3 := ps.Sub("t2")
	ch4 := ps.Sub("t3")

	ps.Pub("hi", "t1")
	ps.Pub("hello", "t2")
	require.Equal(t, <-ch1, "hi")
	require.Equal(t, <-ch2, "hi")
	require.Equal(t, <-ch3, "hello")

	ps.Close("t1", "t2")
	_, ok := <-ch1
	require.Equal(t, ok, false)
	_, ok = <-ch2
	require.Equal(t, ok, false)
	_, ok = <-ch3
	require.Equal(t, ok, false)

	ps.Pub("welcome", "t3")
	require.Equal(t, <-ch4, "welcome")

	ps.Shutdown()
}

func TestUnsubAfterClose(t *testing.T) {
	ps := New(1)
	ch := ps.Sub("t1")
	defer func() {
		ps.Unsub(ch, "t1")
		ps.Shutdown()
	}()

	ps.Close("t1")
	_, ok := <-ch
	require.Equal(t, ok, false)
}

func TestShutdown(t *testing.T) {
	start := runtime.NumGoroutine()
	New(10).Shutdown()
	time.Sleep(1)
	require.Equal(t, runtime.NumGoroutine()-start, 1)
}

func TestMultiSub(t *testing.T) {
	ps := New(1)
	ch := ps.Sub("t1", "t2")

	ps.Pub("hi", "t1")
	require.Equal(t, <-ch, "hi")

	ps.Pub("hello", "t2")
	require.Equal(t, <-ch, "hello")

	ps.Shutdown()
	_, ok := <-ch
	require.Equal(t, ok, false)
}

func TestMultiSubOnce(t *testing.T) {
	ps := New(1)
	ch := ps.SubOnce("t1", "t2")

	ps.Pub("hi", "t1")
	require.Equal(t, <-ch, "hi")

	ps.Pub("hello", "t2")

	_, ok := <-ch
	require.Equal(t, ok, false)
	ps.Shutdown()
}

func TestMultiPub(t *testing.T) {
	ps := New(1)
	ch1 := ps.Sub("t1")
	ch2 := ps.Sub("t2")

	ps.Pub("hi", "t1", "t2")
	require.Equal(t, <-ch1, "hi")
	require.Equal(t, <-ch2, "hi")

	ps.Shutdown()
}

func TestMultiUnsub(t *testing.T) {
	ps := New(1)
	ch := ps.Sub("t1", "t2", "t3")

	ps.Unsub(ch, "t1")

	ps.Pub("hi", "t1")

	ps.Pub("hello", "t2")
	require.Equal(t, <-ch, "hello")

	ps.Unsub(ch, "t2", "t3")
	_, ok := <-ch
	require.Equal(t, ok, false)

	ps.Shutdown()
}

func TestMultiClose(t *testing.T) {
	ps := New(1)
	ch := ps.Sub("t1", "t2")

	ps.Pub("hi", "t1")
	require.Equal(t, <-ch, "hi")

	ps.Close("t1")
	ps.Pub("hello", "t2")
	require.Equal(t, <-ch, "hello")

	ps.Close("t2")
	_, ok := <-ch
	require.Equal(t, ok, false)

	ps.Shutdown()
}
