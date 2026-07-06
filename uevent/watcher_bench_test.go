package uevent_test

import (
	"context"
	"sync"
	"testing"

	"git.casinomodule.org/casino27/basic-utils/v2/uevent"
)

func BenchmarkParallelWatcherDispatch(b *testing.B) {
	input := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	watcher := uevent.NewParallelWatcher(input, func(ctx context.Context, value int) {
		wg.Done()
	})
	if !watcher.Watch(ctx) {
		b.Fatal("watcher did not start")
	}

	wg.Add(b.N)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input <- i
	}
	wg.Wait()
	b.StopTimer()
	close(input)
}

func BenchmarkBroadcastWatcherDispatchOneListener(b *testing.B) {
	benchmarkBroadcastWatcherDispatch(b, 1)
}

func BenchmarkBroadcastWatcherDispatchFourListeners(b *testing.B) {
	benchmarkBroadcastWatcherDispatch(b, 4)
}

func benchmarkBroadcastWatcherDispatch(b *testing.B, listeners int) {
	input := make(chan int)
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	watcher := uevent.NewBroadcastWatcher(input)
	for i := 0; i < listeners; i++ {
		watcher.Register(func(ctx context.Context, value int) {
			wg.Done()
		})
	}
	if !watcher.Watch(ctx) {
		b.Fatal("watcher did not start")
	}

	wg.Add(b.N * listeners)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		input <- i
	}
	wg.Wait()
	b.StopTimer()
	close(input)
}
