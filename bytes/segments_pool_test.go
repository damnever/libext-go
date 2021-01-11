package bytes

import (
	"fmt"
	"math"
	"strconv"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSegmentsPool_SortedSizes(t *testing.T) {
	t.Parallel()

	p := NewSegmentsPool(SegmentsPoolSizesFrom([]int{1, 3, 5, 11, 9, 7}))
	dataset := []slicePoolTestData{
		{0, 0},
		{1, 1},
		{2, 3},
		{3, 3},
		{4, 5},
		{5, 5},
		{6, 7},
		{7, 7},
		{8, 9},
		{9, 9},
		{10, 11},
		{11, 11},
		{65535, 65535},
	}
	testSegmentsPool(t, p, dataset, "")
	testSegmentsPoolConcurrent(t, p, dataset, "")
}

func TestSegmentsPool_RangeSizes(t *testing.T) {
	t.Parallel()

	p := NewSegmentsPool(SegmentsPoolRangeSizes(2, 12, 2))
	dataset := []slicePoolTestData{
		{0, 0},
		{1, 2},
		{2, 2},
		{3, 4},
		{4, 4},
		{5, 6},
		{6, 6},
		{7, 8},
		{8, 8},
		{9, 10},
		{10, 10},
		{11, 12},
		{12, 12},
		{13, 13},
		{65536, 65536},
	}
	testSegmentsPool(t, p, dataset, "")
	testSegmentsPoolConcurrent(t, p, dataset, "")
}

func TestSegmentsPool_ExponentialSizesBase2(t *testing.T) {
	t.Parallel()

	p := NewSegmentsPool(SegmentsPoolExponentialSizes(512, 2097152, 2))
	dataset := []slicePoolTestData{
		{100, 512},
		{1023, 1024},
		{511, 512},
		{16333, 16384},
		{10000, 16384},
		{20000, 32768},
		{32666, 32768},
		{32769, 65536},
		{2000, 2048},
		{4000, 4096},
		{3000, 4096},
		{5000, 8192},
		{8189, 8192},
		{100000, 131072},
		{129999, 131072},
		{200000, 262144},
		{255555, 262144},
		{524287, 524288},
		{2097151, 2097152},
		{2097152, 2097152},
		{2097153, 2097153},
		{2100000, 2100000},
	}
	testSegmentsPool(t, p, dataset, "")
	testSegmentsPoolConcurrent(t, p, dataset, "")
}

func TestSegmentsPool_ExponentialSizesBase3(t *testing.T) {
	t.Parallel()

	p := NewSegmentsPool(SegmentsPoolExponentialSizes(200, 1594323, 3))
	dataset := []slicePoolTestData{
		{100, 243},
		{323, 729},
		{811, 2187},
		{6503, 6561},
		{10000, 19683},
		{42769, 59049},
		{111111, 177147},
		{444444, 531441},
		{1594322, 1594323},
		{1594323, 1594323},
		{2000000, 2000000},
	}
	testSegmentsPool(t, p, dataset, "")
	testSegmentsPoolConcurrent(t, p, dataset, "")
}

func TestSegmentsPool_ExponentialSizes(t *testing.T) {
	t.Parallel()

	const START = 64    // 64B
	const END = 1 << 26 // 64MiB
	for base := 2; base < 10; base++ {
		var dataset []slicePoolTestData
		min := int(math.Pow(float64(base), math.Ceil(mathLogx(float64(START), float64(base)))))
		size := min
		for ; size < END; size *= base {
			dataset = append(dataset, slicePoolTestData{size: size - 1, capacity: size},
				slicePoolTestData{size: size, capacity: size})
		}
		dataset = append(dataset, slicePoolTestData{size: size + 1, capacity: size + 1},
			slicePoolTestData{size: size + 2, capacity: size + 2})
		p := NewSegmentsPool(SegmentsPoolExponentialSizes(START, END, base))
		testSegmentsPool(t, p, dataset, "base"+strconv.Itoa(base)+": ")
	}
}

type slicePoolTestData struct {
	size     int
	capacity int
}

func testSegmentsPool(t *testing.T, p *SegmentsPool, dataset []slicePoolTestData, msgpfx string) {
	t.Helper()

	for _, v := range dataset {
		b := p.Get(v.size)
		require.Equal(t, v.capacity, cap(b), fmt.Sprintf("%s%+v", msgpfx, v))
		p.Put(b)
	}
}

func testSegmentsPoolConcurrent(t *testing.T, p *SegmentsPool, dataset []slicePoolTestData, msgpfx ...string) {
	t.Helper()

	wg := sync.WaitGroup{}
	for _, v := range dataset {
		wg.Add(1)
		go func(v slicePoolTestData) {
			defer wg.Done()

			b := p.Get(v.size)
			defer p.Put(b)
			require.Equal(t, v.capacity, cap(b), fmt.Sprintf("%s%+v", msgpfx, v))
		}(v)
	}
	wg.Wait()
}
