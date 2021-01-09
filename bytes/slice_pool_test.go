package bytes

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestSlicePool_RangeSizes(t *testing.T) {
	p := NewSlicePool(RangeSizes(2, 12, 2))
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
	testSlicePool(t, p, dataset)
}

func TestSlicePool_ExponentialSizes(t *testing.T) {
	p := NewSlicePool(ExponentialSizes(512, 2097152))
	dataset := []slicePoolTestData{
		{100, 512},
		{1023, 1024},
		{511, 512},
		{10000, 16384},
		{20000, 32768},
		{32769, 65536},
		{2000, 2048},
		{3000, 4096},
		{5000, 8192},
		{100000, 131072},
		{200000, 262144},
		{524287, 524288},
		{2097151, 2097152},
		{2097152, 2097152},
		{2097153, 2097153},
		{2100000, 2100000},
	}
	testSlicePool(t, p, dataset)
}

type slicePoolTestData struct {
	size     int
	capacity int
}

func testSlicePool(t *testing.T, p *SlicePool, dataset []slicePoolTestData) {
	t.Helper()

	for _, v := range dataset {
		b := p.Get(v.size)
		require.Equal(t, v.capacity, cap(b), fmt.Sprintf("%+v", v))
		p.Put(b)
	}
}
