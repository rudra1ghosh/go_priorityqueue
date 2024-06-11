package Queue

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
)

type customer struct {
	wid      string
	rid      string
	ctx      context.Context
	priority int32
}

func New(wid string, rid string, ctx context.Context, priority int32) customer {
	return customer{
		wid:      wid,
		rid:      rid,
		ctx:      ctx,
		priority: priority,
	}
}

type ByPriority []customer

func (a ByPriority) Len() int           { return len(a) }
func (a ByPriority) Less(i, j int) bool { return a[i].priority < a[j].priority }
func (a ByPriority) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }

func (q *Queue) SortCustomers() {
	sort.Sort(ByPriority(q.customers))
}

type Queue struct {
	Size      int
	customers []customer
	mutex     sync.Mutex
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

func (q *Queue) Dequeue() (string, string, context.Context, int32, error) {
	q.mutex.Lock() // Lock for concurrency safety
	defer q.mutex.Unlock()

	if len(q.customers) == 0 {
		return " ", " ", nil, 0, errors.New("Queue is empty")
	}
	c := q.customers[0]
	q.customers = q.customers[1:]

	return c.wid, c.rid, c.ctx, c.priority, nil
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
		fmt.Printf("WID: %s, Priority: %d\n", c.wid, c.priority)
	}
}

func (q *Queue) IsEmpty() bool {
	return len(q.customers) == 0
}
