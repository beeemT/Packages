package queue

import (
	"sync"
	"testing"
)

func primeQueue(numElems int, q *Queue, tp Queuetype) {
	i := 0
	for i < numElems {
		var elem *queueElement
		switch tp {
		case Fifo, Lifo:
			elem = NewQueueElementWithoutPriority(i)
		case PriorityHigh, PriorityLow:
			elem = NewQueueElementWithPriority(i, float64(i))
		default:
			panic(int(tp))
		}
		q.Insert(elem)
		i++
	}
}

func TestQueueBasicAllTypes(t *testing.T) {
	w := new(sync.WaitGroup)
	//for index := 0; index < numQueuetypes; index++ {
	go testQueue(3, w, t)
	w.Add(1)
	//}
	w.Wait()
}

func testQueue(index int, w *sync.WaitGroup, t *testing.T) {
	loops := 100000
	q, err := NewQueue(Queuetype(index))
	if err != nil {
		t.Log(err.Error())
		t.FailNow()
	}

	primeQueue(loops, q, Queuetype(index))

	if q.Length() != loops {
		t.Logf("[%d] Length not ok. Expected: %d but received %d", index, loops, q.Length())
	}

	i := 0
	expectedValue := 0
	t.Log(q.numElems)
	for i < loops {
		if index == 3 {
			t.Log(q.queSlice)
		}
		switch Queuetype(index) {
		case Fifo:
			expectedValue = i
		case Lifo:
			expectedValue = (loops - 1) - i
		case PriorityHigh:
			expectedValue = (loops - 1) - i
		case PriorityLow:
			expectedValue = i
		default:
			panic("Not all ordertypes implemented")
		}
		e, _, err := q.Remove()
		if err != nil {
			t.Logf("[%d] %s\n", index, err.Error())
			t.FailNow()
		}
		eV, ok := e.(int)
		if !ok {
			t.Logf("[%d] %s\n", index, "Type not ok")
			t.FailNow()
		}
		if eV != expectedValue {
			t.Logf("[%d] Value not ok. Expected: %d but received: %d\n", index, expectedValue, eV)
			t.FailNow()
		}
		i++
	}
	t.Logf("Queuetype %d succeded", index)
	w.Done()
}
