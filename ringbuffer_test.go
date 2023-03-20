package midimonster

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRingBuffer(t *testing.T) {
	buf := NewRingBuffer(1024)
	assert.Equal(t, uint64(0), buf.Size())

	testElement := "test1"
	buf.Append(testElement)
	assert.Equal(t, uint64(1), buf.Size())

	elements := buf.GetAll()
	if assert.Len(t, elements, 1) {
		assert.Equal(t, testElement, elements[0])
	}
}

func TestRingBufferMoreElementsAsBuffer(t *testing.T) {
	bufCap := uint64(2)
	buf := NewRingBuffer(bufCap)

	testStrings := []string{
		"test1",
		"test2",
		"test3",
	}

	for _, e := range testStrings {
		buf.Append(e)
	}

	assert.Equal(t, bufCap, buf.Size())

	elements := buf.GetAll()
	if assert.Len(t, elements, int(bufCap)) {
		assert.Equal(t, testStrings[1:], elements)
	}

	elements = buf.GetFromOldest(2)
	if assert.Len(t, elements, 1) {
		assert.Equal(t, testStrings[2:], elements)
	}

	elements = buf.GetFromOldest(1000)
	assert.Len(t, elements, 0)
}
