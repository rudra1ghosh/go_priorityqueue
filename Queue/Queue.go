package Queue

import (
	"context"
	"errors"
	"fmt"
	"log"
	"sort"
	"sync"
	"time"
)

type customer struct {
	wid       string
	rid       string
	ctx       context.Context
	priority  int32
	timestamp time.Time
}

func New(wid string, rid string, ctx context.Context, priority int32, timestamp time.Time) customer {
	return customer{
		wid:       wid,
		rid:       rid,
		ctx:       ctx,
		priority:  priority,
		timestamp: timestamp,
	}
}

type ByPriority struct {
	customers []customer
	offset    int
}

func (a ByPriority) Len() int { return len(a.customers) - a.offset }
func (a ByPriority) Less(i, j int) bool {
	return a.customers[a.offset+i].priority < a.customers[a.offset+j].priority
}
func (a ByPriority) Swap(i, j int) {
	a.customers[a.offset+i], a.customers[a.offset+j] = a.customers[a.offset+j], a.customers[a.offset+i]
}

type Queue struct {
	Size               int
	customers          []customer
	overdueOffset      int
	mutex              sync.Mutex
	lastDequeued       customer
	lastDequeuedExists bool
}

func (q *Queue) EnqueueFront(c customer) (bool, error) {
	q.mutex.Lock() // Lock for concurrency safety
	defer q.mutex.Unlock()

	if q.Size > 0 && len(q.customers) >= q.Size {
		return true, errors.New("Queue is full") // Queue is full
	}
	q.customers = append([]customer{c}, q.customers...)
	q.overdueOffset++
	return false, nil
}

func (q *Queue) Enqueue(c customer) (bool, error) {
	q.mutex.Lock() // Lock for concurrency safety
	defer q.mutex.Unlock()

	if q.Size > 0 && len(q.customers) >= q.Size {
		return true, errors.New("Queue is full") // Queue is full
	}
	q.customers = append(q.customers, c)
	return false, nil
}

func (q *Queue) Dequeue() (string, string, context.Context, int32, time.Time, error) {
	q.mutex.Lock() // Lock for concurrency safety
	defer q.mutex.Unlock()

	if len(q.customers) == 0 {
		return " ", " ", nil, 0, time.Time{}, errors.New("Queue is empty")
	}
	c := q.customers[0]
	q.customers = q.customers[1:]
	if q.overdueOffset > 0 && (c.priority == 4 || c.priority == 5 || c.priority == 6) {
		q.overdueOffset--
	}
	q.lastDequeued = c
	q.lastDequeuedExists = true
	if q.lastDequeuedExists {
		log.Printf("This is the task dequed:" + q.lastDequeued.wid)
	}

	return c.wid, c.rid, c.ctx, c.priority, c.timestamp, nil
}

func (q *Queue) GetLength() int {
	q.mutex.Lock() // Lock for concurrency safety
	defer q.mutex.Unlock()
	return len(q.customers)
}

func (q *Queue) Display() {
	q.mutex.Lock() // Lock for concurrency safety
	defer q.mutex.Unlock()
	if len(q.customers) == 0 {
		fmt.Println("Queue is empty")
		return
	}
	fmt.Println("Customers in the Queue:")
	for _, c := range q.customers {
		fmt.Printf("WID: %s, Priority: %d, Timestamp: %v\n", c.wid, c.priority, c.timestamp)
	}
}

func (q *Queue) IsEmpty() bool {
	q.mutex.Lock() // Lock for concurrency safety
	defer q.mutex.Unlock()
	return len(q.customers) == 0
}

func (q *Queue) SortCustomers() {
	q.mutex.Lock() // Lock for concurrency safety
	defer q.mutex.Unlock()

	sort.Sort(ByPriority{q.customers, q.overdueOffset})
}

func (q *Queue) MoveToFrontIfOverdue(id int32) {
	for {
		time.Sleep(time.Millisecond * 500)

		q.mutex.Lock()
		for i := q.overdueOffset; i < len(q.customers); i++ {
			if q.customers[i].priority == id && time.Since(q.customers[i].timestamp) > 19*time.Second {

				c := q.customers[i]
				// Dequeue from the current position
				q.customers = append(q.customers[:i], q.customers[i+1:]...)
				// Insert at the front of the queue
				q.customers = append([]customer{c}, q.customers...)
				// Increment the overdue offset
				q.overdueOffset++
				break // Move only one customer per iteration
			}
		}
		q.mutex.Unlock()
	}
}
