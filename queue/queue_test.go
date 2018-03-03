package queue

import (
	"testing"
)

func TestQueueLifoNoPrio(t *testing.T) {
	loops := 10000
	q := NewQueue(Lifo)
	i := 0
	for i < loops {
		elem := NewQueueElementWithoutPriority(i)
		q.Insert(elem)
		i++
	}

	l := q.Length() - 1
	for l > 0 {
		e, p := q.Remove()
		eV, ok := e.(int)
		if !ok {
			t.Log("Type not ok")
			t.Fail()
		}
		if p != 0 {
			t.Log("Priority not ok")
			t.Fail()
		}
		if eV != int(l) {
			t.Logf("Value not ok. Expected: %d but received: %d", l, eV)
			t.Fail()
		}
		l--
	}
}

func TestQueueFifoNoPrio(t *testing.T) {
	loops := 10000
	q := NewQueue(Fifo)
	i := 0
	for i < loops {
		elem := NewQueueElementWithoutPriority(i)
		q.Insert(elem)
		i++
	}

	l := q.Length()
	i = 0
	for i > l {
		e, p := q.Remove()
		eV, ok := e.(int)
		if !ok {
			t.Log("Type not ok")
			t.FailNow()
		}
		if p != 0 {
			t.Log("Priority not ok")
			t.FailNow()
		}
		if eV != i {
			t.Logf("Value not ok. Expected: %d but received: %d", i, eV)
			t.FailNow()
		}
		i++
	}

}
