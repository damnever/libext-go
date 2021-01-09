package bytes

import (
	"bytes"
	"sync"
)

// BufferPool manages a pool of bytes.Buffer, the underlying pool is sync.Pool.
type BufferPool struct {
	pool sync.Pool
}

// NewBufferPool returns a new BufferPool.
func NewBufferPool() *BufferPool {
	return &BufferPool{
		pool: sync.Pool{
			New: func() interface{} {
				return &bytes.Buffer{}
			},
		},
	}
}

// Get returns a bytes.Buffer from pool.
func (p *BufferPool) Get() *bytes.Buffer {
	return p.pool.Get().(*bytes.Buffer)
}

// Put puts the bytes.Buffer back into pool.
func (p *BufferPool) Put(buf *bytes.Buffer) {
	buf.Reset()
	p.pool.Put(buf)
}
