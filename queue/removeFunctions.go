package queue

import (
	"math"
)

func (q *Queue) remove() *queueElement {
	retElem := q.queSlice[(len(q.queSlice) - 1)]
	q.queSlice[(len(q.queSlice) - 1)] = nil
	q.queSlice = q.queSlice[:(len(q.queSlice) - 1)]

	q.handleShrink()

	return retElem
}

func (q *Queue) handleShrink() {
	if float64(len(q.queSlice)) < q.shrinkFactor*float64(cap(q.queSlice)) {
		newCap := int(math.Ceil(q.afterShrinkFactor * float64(cap(q.queSlice))))
		temp := make([]*queueElement, len(q.queSlice), newCap)
		copy(temp, q.queSlice[:len(q.queSlice)])
		q.queSlice = temp
	}
}

func (q *Queue) delete(i int) {
	q.deleteWithoutMemoryManagement(i)
	q.handleShrink()
}

func (q *Queue) deleteWithoutMemoryManagement(i int) {
	copy(q.queSlice[i:], q.queSlice[i+1:])
	q.queSlice[len(q.queSlice)-1] = nil
	q.queSlice = q.queSlice[:len(q.queSlice)-1]
}
