package midimonster

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

func TestRingBuffer(t *testing.T) {
	buf := NewRingBuffer(1024)
	if buf.Size() != 0 {
		t.Errorf("buffer should be empty")
	}

	testElement := "test1"
	buf.Append(testElement)
	if buf.Size() != 1 {
		t.Errorf("buffer should now contain one element")
	}

	elements := buf.GetAll()
	if len(elements) != 1 {
		t.Errorf("GetAll should return 1 element")
	} else {
		if elements[0] != testElement {
			t.Errorf("GetAll should return the %s instead of %s", testElement, elements[0])
		}
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

	if buf.Size() != bufCap {
		t.Errorf("ring buffer size should be %d instead of %d", bufCap, buf.Size())
	}

	elements := buf.GetAll()
	if len(elements) != int(bufCap) {
		t.Errorf("ring buffer elements should be %d instead of %d", bufCap, len(elements))
	} else {
		if !cmp.Equal(elements, testStrings[1:]) {
			t.Errorf("ring buffer elements wrong:\n%s", cmp.Diff(testStrings[1:], elements))
		}
	}

	elements = buf.GetFromOldest(2)
	if len(elements) != 1 {
		t.Errorf("ring buffer elements should be %d instead of %d", 1, len(elements))
	} else {
		if !cmp.Equal(elements, testStrings[2:]) {
			t.Errorf("ring buffer elements wrong:\n%s", cmp.Diff(testStrings[2:], elements))
		}
	}

	elements = buf.GetFromOldest(1000)
	if len(elements) != 0 {
		t.Errorf("ring buffer elements should be empty")
	}
}
