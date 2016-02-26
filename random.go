package xkcd

import (
	"math/rand"
	"sync"
)

var random *rand.Rand

type lockedRandSource struct {
	lock sync.Mutex
	src  rand.Source
}

// to satisfy rand.Source interface
func (r *lockedRandSource) Int63() int64 {
	r.lock.Lock()
	ret := r.src.Int63()
	r.lock.Unlock()
	return ret
}

// to satisfy rand.Source interface
func (r *lockedRandSource) Seed(seed int64) {
	r.lock.Lock()
	r.src.Seed(seed)
	r.lock.Unlock()
}

func randomInt(begin, end int) int {
	return random.Intn(end-begin) + begin
}
