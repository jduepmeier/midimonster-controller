package midimonster

import "sync"

type RingBuffer struct {
	// buffer
	buf []string
	// capacitiy
	cap uint64
	// id of the oldest entry
	oldest uint64
	// id of the newest entry
	newest uint64
	// index of the current element
	current uint64
	// rw mutex
	mutex sync.RWMutex
}

func NewRingBuffer(size uint64) *RingBuffer {
	return &RingBuffer{
		buf:     make([]string, size),
		cap:     size,
		oldest:  0,
		newest:  0,
		current: 0,
	}
}

func (buffer *RingBuffer) Append(line string) {
	buffer.mutex.Lock()
	defer buffer.mutex.Unlock()
	buffer.buf[buffer.current] = line
	buffer.current = (buffer.current + 1) % buffer.cap
	buffer.newest++
	if buffer.newest >= buffer.cap {
		buffer.oldest = buffer.newest - buffer.cap
	}
}

func (buffer *RingBuffer) Size() uint64 {
	buffer.mutex.RLock()
	defer buffer.mutex.RUnlock()
	return buffer.newest - buffer.oldest
}

func (buffer *RingBuffer) GetAll() []string {
	return buffer.GetFromOldest(buffer.oldest)
}

func (buffer *RingBuffer) Newest() uint64 {
	buffer.mutex.RLock()
	defer buffer.mutex.RUnlock()
	return buffer.newest
}

func (buffer *RingBuffer) GetFromOldest(oldest uint64) []string {
	buffer.mutex.RLock()
	defer buffer.mutex.RUnlock()
	if oldest >= buffer.newest {
		return []string{}
	}
	if oldest < buffer.oldest {
		oldest = buffer.oldest
	}

	size := buffer.newest - oldest
	out := make([]string, size)

	var i uint64
	offset := buffer.Size() - size
	for i = 0; i < size; i++ {
		out[i] = buffer.buf[(buffer.oldest+i+offset)%buffer.cap]
	}
	return out
}
