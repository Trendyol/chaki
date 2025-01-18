package rand

import (
	"math/rand"
	"sync"
	"time"
)

var (
	instance *RandomGenerator
	once     sync.Once
)

func GetInstance() *RandomGenerator {
	once.Do(func() {
		instance = &RandomGenerator{
			source: rand.New(rand.NewSource(time.Now().UnixNano())),
		}
	})
	return instance
}

type RandomGenerator struct {
	source *rand.Rand
	mu     sync.Mutex
}

func (g *RandomGenerator) Int() int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.source.Int()
}

func (g *RandomGenerator) Float64() float64 {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.source.Float64()
}

func (g *RandomGenerator) Intn(n int) int {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.source.Intn(n)
}
