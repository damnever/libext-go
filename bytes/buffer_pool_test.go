package bytes

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestBufferPool(t *testing.T) {
	p := NewBufferPool()
	buf := p.Get()
	p.Put(buf)
	for i := 0; i < 10; i++ {
		b := p.Get()
		require.Equal(t, buf, b)
		p.Put(b)
	}
}
