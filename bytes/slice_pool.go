package bytes

import (
	"math"
	"sort"
	"sync"
)

// SlicePool manages variable sized byte slices by multiple sync.Pool.
type SlicePool struct {
	sizes SlicePoolSizes
	pools []sync.Pool
}

// NewSlicePool creates a new SlicePool.
func NewSlicePool(sizes SlicePoolSizes) *SlicePool {
	p := &SlicePool{sizes: sizes}
	sizes.Iterate(func(size int) bool {
		p.pools = append(p.pools, sync.Pool{ // DO NOT copy it.
			New: func() interface{} {
				return make([]byte, size, size)
			},
		})
		return true
	})
	return p
}

// Get returns a byte slice by given size, the returned slice may bigger than size.
func (p *SlicePool) Get(size int) []byte {
	if size == 0 {
		return nil
	}

	index, ok := p.sizes.Index(size)
	if ok {
		return p.pools[index].Get().([]byte)
	}
	return make([]byte, size, size)
}

// Put puts the byte slice back into pool.
func (p *SlicePool) Put(b []byte) {
	capacity := cap(b)
	if capacity == 0 {
		return
	}

	if index, ok := p.sizes.Index(capacity); ok {
		p.pools[index].Put(b) //nolint:staticcheck
	}
}

// SlicePoolSizes represents a list of sizes.
type SlicePoolSizes interface {
	// Iterate iterates over sizes, it will stop if sizeIterator returns false.
	Iterate(sizeIterator func(size int) (stop bool))
	// Index returns the index of the given size.
	Index(size int) (index int, found bool)
}

// NOTE: the length is relatively small, maybe binary search added unnecessary overhead.
type sortedSizes []int

func (sizes sortedSizes) Iterate(sizeIterator func(int) bool) {
	for _, size := range sizes {
		if !sizeIterator(size) {
			return
		}
	}
}

func (sizes sortedSizes) Index(size int) (int, bool) {
	n := len(sizes)
	index := sort.Search(n, func(i int) bool { return sizes[i] >= size })
	return index, index != n
}

// SizesFrom creates a SlicePoolSizes from normal slice. NOTE that it will
// panic if the slice contains negative value.
func SizesFrom(sizes []int) SlicePoolSizes {
	sort.Ints(sizes)
	if len(sizes) > 0 && sizes[0] < 0 {
		panic("libext-go/bytes: negative size")
	}
	return sortedSizes(sizes)
}

// RangeSizes creates a list of sizes from [start, end], increased by step.
// NOTE that it will panic if the calculated sizes contain negative value.
func RangeSizes(start, end, step int) SlicePoolSizes {
	var sizes []int
	for size := start; size <= end; size += step {
		if size < 0 {
			panic("libext-go/bytes: negative size")
		}
		sizes = append(sizes, size)
	}
	return SizesFrom(sizes)
}

type exponentialSizes struct {
	min         int
	max         int
	maxbase2exp int
}

func (sizes exponentialSizes) Iterate(sizeIterator func(int) bool) {
	for i := 0; i <= sizes.maxbase2exp; i++ {
		size := (1 << uint(i)) * sizes.min
		if !sizeIterator(size) {
			return
		}
	}
}

func (sizes exponentialSizes) Index(size int) (int, bool) {
	if size < sizes.min {
		return 0, true
	}
	if size > sizes.max {
		return 0, false
	}
	return int(math.Ceil(math.Log2(math.Ceil(float64(size) / float64(sizes.min))))), true
}

// ExponentialSizes creates a SlicePoolSizes from [min, max], increased exponentially based on 2.
// NOTE that it will panic if the min or the max is negative, or min greater than max after round up/down.
func ExponentialSizes(min, max int) SlicePoolSizes {
	if min < 0 || max < 0 {
		panic("libext-go/bytes: min or max can not be negative")
	}
	min = int(math.Pow(2, math.Ceil(math.Log2(float64(min)))))  // Round up
	max = int(math.Pow(2, math.Floor(math.Log2(float64(max))))) // Round down
	if min > max {
		panic("libext-go/bytes: min can not greater than max")
	}
	maxexp := int(math.Log2(float64(max)))
	return exponentialSizes{min: min, max: max, maxbase2exp: maxexp}
}
