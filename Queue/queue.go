package queue

import (
	"errors"
	"sync"
)

type Queue struct {
	size int
	total int
	elements []string
	mu sync.Mutex
}

func (q *Queue) Enqueue(url string) {
	q.mu.Lock()
	defer q.mu.Unlock()

	q.elements = append(q.elements, url)
	q.size++
	q.total++
}

func (q *Queue) Deque() (string, error) {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.elements) < 1 {
		return "", errors.New("Empty queue")
	}

	url := q.elements[0]
	q.elements = q.elements[1:]
	q.size--

	return url, nil
}

func (q *Queue) Size() int {
	return q.size
}
