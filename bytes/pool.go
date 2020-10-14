package bytes

import (
	"math"
	"sync"
)

type Pool struct {
	max2exp int
	base    int
	pools   []*sync.Pool
}

func NewPool() *Pool {
	return NewPoolWith(4, 512)
}

func NewPoolWith(max2exp, base int) *Pool {
	base = int(math.Pow(2, math.Ceil(math.Log2(float64(base))))) // Round up
	if base <= 0 {
		panic("base can not be nil")
	}

	pools := make([]*sync.Pool, max2exp+1, max2exp+1)
	for i := 0; i <= max2exp; i++ {
		sz := (1 << uint(i)) * base
		pools[i] = &sync.Pool{
			New: func() interface{} {
				return make([]byte, sz, sz)
			},
		}
	}
	return &Pool{
		max2exp: max2exp,
		base:    base,
		pools:   pools,
	}
}

func (p *Pool) Get(size int) []byte {
	if idx := p.index(size); idx >= 0 && idx <= p.max2exp {
		return p.pools[idx].Get().([]byte)
	}
	return make([]byte, size, size)
}

func (p *Pool) Put(b []byte) {
	if idx := p.index(cap(b)); idx >= 0 && idx <= p.max2exp {
		p.pools[idx].Put(b) //nolint:staticcheck
	}
}

func (p *Pool) index(n int) int {
	if n < p.base {
		return 0
	}
	return int(math.Ceil(math.Log2(math.Ceil(float64(n) / float64(p.base)))))
}
