package counter

import "sync"

type Counter struct {
	m sync.Mutex
	c int
}

func New() *Counter {
	return &Counter{}
}

func (c *Counter) Count() int {
	c.m.Lock()
	defer c.m.Unlock()
	c.c += 1
	return c.c
}

func (c *Counter) Value() int {
	c.m.Lock()
	defer c.m.Unlock()
	return c.c
}
