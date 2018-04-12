package queue

import (
	"math"
)

func (q *Queue) remove() (*queueElement, error) {
	lenQ := len(q.queSlice)
	if lenQ == 0 {
		return nil, EmptyListError{}
	}

	retElem := q.queSlice[(lenQ - 1)]
	q.queSlice[(lenQ - 1)] = nil
	q.queSlice = q.queSlice[:(lenQ - 1)]

	q.handleShrink()

	return retElem, nil
}

func (q *Queue) handleShrink() {
	lenQ := len(q.queSlice)
	if float64(lenQ) < q.shrinkFactor*float64(cap(q.queSlice)) {
		newCap := int(math.Ceil(q.afterShrinkFactor * float64(cap(q.queSlice))))
		temp := make([]*queueElement, lenQ, newCap)
		copy(temp, q.queSlice[:lenQ])
		q.queSlice = temp
	}
}

func (q *Queue) delete(i int) error {
	err := q.deleteWithoutMemoryManagement(i)
	q.handleShrink()
	return err
}

func (q *Queue) deleteWithoutMemoryManagement(i int) error {
	if q.numElems == 0 {
		return EmptyListError{}
	}
	if i < 0 || i >= q.numElems {
		return IndexOutOfBoundsError{}
	}
	lenQ := len(q.queSlice)
	copy(q.queSlice[i:], q.queSlice[i+1:])
	q.queSlice[lenQ-1] = nil
	q.queSlice = q.queSlice[:lenQ-1]
	return nil
}
