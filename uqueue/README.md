# Queues Library

This library provides generic queue data structures for Go. It includes FIFO, lock-free concurrent FIFO, and priority
queue implementations.

## Features

- Generic Implementation: FIFO and priority queues are implemented using Go generics, allowing
  you to queue items of any type.
- Thread-Safety: Queue operations are safe across multiple goroutines. `ConcurrentFIFOQueueImpl` provides a lock-free
  FIFO implementation for high-concurrency FIFO workloads.
- Timeouts for Polling: The `Poll` function allows consumers to wait for items with a timeout.

## Structures

### FIFOQueueImpl

A standard FIFO (First-In-First-Out) queue implementation.

#### Methods

- `NewFIFOQueue`: Initializes a new FIFO queue with optional initial elements.
- `Queue`: Enqueues an item to the back of the queue.
- `Fetch`: Dequeues an item from the front of the queue without waiting.
- `Poll`: Dequeues an item from the front of the queue or waits for a specified timeout.

### ConcurrentFIFOQueueImpl

A lock-free FIFO queue implementation based on the Michael-Scott queue algorithm.

#### Methods

- `NewConcurrentFIFOQueueImpl`: Initializes a new concurrent FIFO queue.
- `Queue`: Enqueues an item to the back of the queue.
- `Fetch`: Dequeues an item from the front of the queue without waiting.
- `Poll`: Dequeues an item from the front of the queue or waits for a specified timeout.

### PriorityQueueImpl

A priority queue where items can be enqueued with a priority level, and dequeuing will retrieve the highest priority
item.

#### Methods

- `NewPrioritizedPriorityQueue`: Initializes a new empty priority queue.
- `Queue`: Enqueues an item with a specified priority.
- `Fetch`: Dequeues the highest-priority item from the queue without waiting.
- `Poll`: Dequeues the highest-priority item from the queue or waits for a specified timeout.

## Dependencies

The library uses the `container/heap` package for the priority queue implementation and an `opt` package for optional
value handling.

## Usage

#### FIFO queue:

```go
package myprogram

import (
	"time"

	"git.casinomodule.org/casino27/basic-utils/v2/uqueue"
)

q := uqueue.NewFIFOQueue[int](1, 2, 3)
q.Queue(4)
item := q.Poll(5 * time.Second)

```

#### Priority queue:

```go
package myprogram

import (
	"time"

	"git.casinomodule.org/casino27/basic-utils/v2/uqueue"
)

pq := uqueue.NewPriorityQueue[int]()
pq.Queue(1, 3) // The number 3 here is the priority.
pq.Queue(2, 1)
item := pq.Poll(5 * time.Second)
```

## Author

Developed by [@kordax](mailto:dmorozov@valoru-software.com) (Dmitry Morozov)

[Valoru Software](https://valoru-software.com)
