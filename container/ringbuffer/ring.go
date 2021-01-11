package ringbuffer

const (
	// 16K if item is pointers on X64..
	highWatermark = 2048
	// 1K if item is pointers on X64..
	lowWatermark = 128
	minRingCap   = 32
)

// RingBuffer is a circular slice buffer.
// It is not goroutine safe.
type RingBuffer struct {
	items []interface{}
	head  int
	tail  int
	size  int
	cap   int
}

// New creates a new RingBuffer.
func New() *RingBuffer {
	return &RingBuffer{
		items: make([]interface{}, 2, 2),
		head:  0,
		tail:  -1,
		size:  0,
		cap:   2,
	}
}

// Append appends an item to "tail" of the ring.
func (r *RingBuffer) Append(item interface{}) {
	if r.size == r.cap {
		if bufcap := r.cap * 2; bufcap <= highWatermark {
			r.cap = bufcap
		} else {
			r.cap += highWatermark
		}
		r.resize()
	}

	r.tail = r.next(r.tail)
	r.items[r.tail] = item
	r.size++
}

// Pop pops the item from "head" of the ring.
func (r *RingBuffer) Pop() interface{} {
	item := r.Peek()
	if item != nil {
		r.items[r.head] = nil
		if r.head == r.tail {
			r.head = 0
			r.tail = -1
		} else {
			r.head = r.next(r.head)
		}

		r.size--
		if halfcap := r.cap / 2; r.size < halfcap {
			if r.cap > lowWatermark {
				r.cap = halfcap
				r.resize()
			} else if r.cap > minRingCap && r.size < minRingCap {
				r.cap = minRingCap
				r.resize()
			}
		}
	}
	return item
}

// Peek peeks the item from "head" of the ring.
func (r *RingBuffer) Peek() interface{} {
	if r.size == 0 {
		return nil
	}
	return r.items[r.head]
}

// Len returns length(not capacity) of the ring.
func (r *RingBuffer) Len() int {
	return r.size
}

func (r *RingBuffer) next(i int) int {
	return (i + 1) % r.cap
}

func (r *RingBuffer) resize() {
	items := r.items
	r.items = make([]interface{}, r.cap, r.cap)
	if r.tail < r.head {
		n := copy(r.items, items[r.head:])
		copy(r.items[n:], items[:r.tail+1])
	} else {
		copy(r.items, items[r.head:r.tail+1])
	}
	r.head = 0
	r.tail = r.size - 1
}
