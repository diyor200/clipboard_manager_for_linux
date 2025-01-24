package cache

import (
	"context"
	"sync"

	"golang.design/x/clipboard"
)

type Queue struct {
	Data []byte
	Kind clipboard.Format
}

var (
	q  []Queue
	mu = &sync.Mutex{}
)

func Put(ctx context.Context, data Queue) {
	mu.Lock()
	q = append(q, data)
	if len(q) > 10 {
		q = q[len(q)-10:]
	}
	mu.Unlock()
}

func ReadAll(ctx context.Context) []Queue {
	mu.Lock()
	defer mu.Unlock()
	return q
}
