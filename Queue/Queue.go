package Queue

import (
	"context"
	"errors"
	"fmt"
	"sort"
	"sync"
	"time"
)

type Customer struct {
	Wid       string
	Rid       string
	Ctx       context.Context
	Priority  int32
	Timestamp time.Time
	Leno      int
	Flag      int
}

func New(wid string, rid string, ctx context.Context, priority int32, timestamp time.Time, leno int, flag int) Customer {
	return Customer{
		Wid:       wid,
		Rid:       rid,
		Ctx:       ctx,
		Priority:  priority,
		Timestamp: timestamp,
		Leno:      leno,
		Flag:      flag,
	}
}

type ByPriority struct {
	Customers []Customer
	Offset    int
}

func (a ByPriority) Len() int { return len(a.Customers) - a.Offset }
func (a ByPriority) Less(i, j int) bool {
	return a.Customers[a.Offset+i].Priority < a.Customers[a.Offset+j].Priority
}
func (a ByPriority) Swap(i, j int) {
	a.Customers[a.Offset+i], a.Customers[a.Offset+j] = a.Customers[a.Offset+j], a.Customers[a.Offset+i]
}

type Queue struct {
	Size          int
	Customers     []Customer
	OverdueOffset int
	Mutex         sync.Mutex
}

func (q *Queue) Enqueue(c Customer) (bool, error) {
	q.Mutex.Lock() // Lock for concurrency safety
	defer q.Mutex.Unlock()

	if q.Size > 0 && len(q.Customers) >= q.Size {
		return true, errors.New("Queue is full") // Queue is full
	}

	// Check special cases for priority and flag
	if c.Priority == 4 || c.Priority == 5 || c.Priority == 6 {
		// Insert at the end
		q.Customers = append(q.Customers, c)
	} else if c.Priority == 1 || c.Priority == 2 || c.Priority == 3 {
		// Insert based on flag condition, searching from behind
		insertIndex := len(q.Customers)
		for insertIndex > 0 && q.Customers[insertIndex-1].Flag != 1 {
			insertIndex--
		}
		q.Customers = append(q.Customers[:insertIndex], append([]Customer{c}, q.Customers[insertIndex:]...)...)
	}

	fmt.Printf("wid:%s pid: %d pos:%d time:%s\n", c.Wid, c.Priority, c.Leno, c.Timestamp)
	fmt.Print("Queue Sorted:")
	for i := 0; i < len(q.Customers); i++ {
		fmt.Print(q.Customers[i].Priority, " ")
	}
	fmt.Println()
	return false, nil
}

func (q *Queue) Dequeue() (string, string, context.Context, int32, time.Time, error, int) {
	q.Mutex.Lock()
	defer q.Mutex.Unlock()

	if len(q.Customers) == 0 {
		return " ", " ", nil, 0, time.Time{}, errors.New("Queue is empty"), 0
	}
	c := q.Customers[0]
	q.Customers = q.Customers[1:]
	if q.OverdueOffset > 0 && (c.Priority == 4 || c.Priority == 5 || c.Priority == 6) {
		q.OverdueOffset--
	}

	return c.Wid, c.Rid, c.Ctx, c.Priority, c.Timestamp, nil, c.Flag
}

func (q *Queue) GetLength() int {
	q.Mutex.Lock() // Lock for concurrency safety
	defer q.Mutex.Unlock()

	return len(q.Customers)
}

func (q *Queue) Display() {
	q.Mutex.Lock() // Lock for concurrency safety
	defer q.Mutex.Unlock()
	if len(q.Customers) == 0 {
		fmt.Println("Queue is empty")
		return
	}
	for _, c := range q.Customers {
		fmt.Println(c.Priority)
	}
	fmt.Println()
}

func (q *Queue) IsEmpty() bool {
	q.Mutex.Lock() // Lock for concurrency safety
	defer q.Mutex.Unlock()
	return len(q.Customers) == 0
}

func (q *Queue) SortCustomers() {
	q.Mutex.Lock() // Lock for concurrency safety
	defer q.Mutex.Unlock()

	sort.Sort(ByPriority{q.Customers, q.OverdueOffset})
}

func (q *Queue) MoveToFrontIfOverdue(id int32) {
	for {
		time.Sleep(time.Millisecond * 1)

		q.Mutex.Lock()

		for i := 0; i < len(q.Customers); i++ {
			if q.Customers[i].Flag == 0 && time.Since(q.Customers[i].Timestamp) > 19*time.Second {
				// Extract the overdue customer
				q.Customers[i].Flag = 1
				overdueCustomer := q.Customers[i]
				overdueCustomer.Flag = 1

				// Find the insertion position
				insertIndex := -1
				for j := 0; j < i; j++ {
					if q.Customers[j].Priority == 1 && q.Customers[j].Timestamp.After(overdueCustomer.Timestamp) {
						insertIndex = j
						break
					}
				}

				// Only move the customer if a valid insertion index was found
				if insertIndex != -1 {
					// Remove the overdue customer from the current position
					q.Customers = append(q.Customers[:i], q.Customers[i+1:]...)

					// Insert the overdue customer at the found position
					q.Customers = append(q.Customers[:insertIndex], append([]Customer{overdueCustomer}, q.Customers[insertIndex:]...)...)
					fmt.Print("Current Queue :")
					for _, c := range q.Customers {
						fmt.Print(c.Priority, " ")
					}
					fmt.Println()

				}

			}
		}

		q.Mutex.Unlock()
	}
}
