package bytes

import (
	"math"
	"sort"
	"sync"
)

// SegmentsPool manages variable sized byte slices by multiple sync.Pool.
type SegmentsPool struct {
	sizes SegmentsPoolSizes
	pools []sync.Pool
}

// NewSegmentsPool creates a new SegmentsPool from given sizes.
func NewSegmentsPool(sizes SegmentsPoolSizes) *SegmentsPool {
	p := &SegmentsPool{sizes: sizes}
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
func (p *SegmentsPool) Get(size int) []byte {
	if size <= 0 {
		return nil
	}

	index, ok := p.sizes.Index(size)
	if ok {
		return p.pools[index].Get().([]byte)
	}
	return make([]byte, size, size)
}

// Put puts the byte slice back into pool.
func (p *SegmentsPool) Put(b []byte) {
	capacity := cap(b)
	if capacity == 0 {
		return
	}

	if index, ok := p.sizes.Index(capacity); ok {
		p.pools[index].Put(b) //nolint:staticcheck
	}
}

// SegmentsPoolSizes represents a list of sizes.
type SegmentsPoolSizes interface {
	// Iterate iterates over sizes, it will stop if sizeIterator returns false.
	Iterate(sizeIterator func(size int) (stop bool))
	// Index returns the index of the given size.
	Index(size int) (index int, found bool)
}

// NOTE: the length is relatively small, maybe binary search added unnecessary overhead, such as cache miss.
type sortedSegmentsPoolSizes []int

func (sizes sortedSegmentsPoolSizes) Iterate(sizeIterator func(int) bool) {
	for _, size := range sizes {
		if !sizeIterator(size) {
			return
		}
	}
}

func (sizes sortedSegmentsPoolSizes) Index(size int) (int, bool) {
	n := len(sizes)
	index := sort.Search(n, func(i int) bool { return sizes[i] >= size })
	return index, index != n
}

// SegmentsPoolSizesFrom creates a SegmentsPoolSizes from normal slice. NOTE that it will
// panic if the slice contains negative value.
func SegmentsPoolSizesFrom(sizes []int) SegmentsPoolSizes {
	sort.Ints(sizes)
	if len(sizes) > 0 && sizes[0] < 0 {
		panic("libext-go/bytes: negative size")
	}
	return sortedSegmentsPoolSizes(sizes)
}

type rangeSegmentsPoolSizes struct {
	start int
	end   int
	step  int
}

func (sizes rangeSegmentsPoolSizes) Iterate(sizeIterator func(int) bool) {
	for size := sizes.start; size <= sizes.end; size += sizes.step {
		if !sizeIterator(size) {
			return
		}
	}
}

func (sizes rangeSegmentsPoolSizes) Index(size int) (int, bool) {
	if size <= sizes.start {
		return 0, true
	}
	if size > sizes.end {
		return 0, false
	}

	size -= sizes.start
	index := size / sizes.step
	if size%sizes.step > 0 {
		index++
	}
	return index, true
}

// SegmentsPoolRangeSizes creates a list of sizes from [start, end], increased by step.
// NOTE that it will panic if the start less than the end, or either of them less than 1.
func SegmentsPoolRangeSizes(start, end, step int) SegmentsPoolSizes {
	if start > end {
		panic("libext-go/bytes: start greater than end")
	}
	if end < 1 {
		panic("libext-go/bytes: end less than 1")
	}
	if step < 1 {
		panic("libext-go/bytes: step less than 1")
	}

	end = (end-start)/step*step + start
	return rangeSegmentsPoolSizes{start: start, end: end, step: step}
}

type exponentialSegmentsPoolSizes struct {
	min    int
	minf   float64
	max    int
	maxf   float64
	base   int
	basef  float64
	maxexp int
}

func (sizes exponentialSegmentsPoolSizes) Iterate(sizeIterator func(int) bool) {
	for size := sizes.min; size <= sizes.max; size *= sizes.base {
		if !sizeIterator(size) {
			return
		}
	}
}

func (sizes exponentialSegmentsPoolSizes) Index(size int) (int, bool) {
	if size <= sizes.min {
		return 0, true
	}
	if size > sizes.max {
		return 0, false
	}
	sizef := float64(size)

	var index int
	if sizes.base == 2 {
		index = int(math.Ceil(math.Log2(sizef / sizes.minf)))
	} else {
		indexf := mathLogx(sizef/sizes.minf, sizes.basef)
		index = int(indexf)
		if indexf != float64(int(indexf)) {
			capability := int(sizes.minf * math.Pow(sizes.basef, float64(index)))
			if capability < size {
				index = int(math.Ceil(indexf))
			}
		}
	}
	return index, true
}

// SegmentsPoolExponentialSizes creates a SegmentsPoolSizes from [min, max], increased exponentially
// like that: `min * math.Pow(base, 0...)`. The min may round up and the max may round down.
//
// NOTE that it will panic if the min or the max is negative, or base less than 2,
// or the min greater than the max after round up/down.
func SegmentsPoolExponentialSizes(min, max, base int) SegmentsPoolSizes {
	if base < 2 {
		panic("libext-go/bytes: base less than 2")
	}
	if min < 0 || max < 0 || base < 2 {
		panic("libext-go/bytes: min or max is negative")
	}

	min = int(math.Pow(float64(base), math.Ceil(mathLogx(float64(min), float64(base)))))  // Round up
	max = int(math.Pow(float64(base), math.Floor(mathLogx(float64(max), float64(base))))) // Round down
	if min > max {
		panic("libext-go/bytes: min can not greater than max")
	}
	maxexp := int(mathLogx(float64(max), float64(base)))
	return exponentialSegmentsPoolSizes{
		min:    min,
		minf:   float64(min),
		max:    max,
		maxf:   float64(max),
		base:   base,
		basef:  float64(base),
		maxexp: maxexp,
	}
}

func mathLogx(v, base float64) float64 {
	return math.Log(v) / math.Log(base)
}
