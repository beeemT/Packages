package queue

import (
	"sync"
)

//Queuetype is the enum type for queue invariants.
type Queuetype int

const (
	//Fifo queue.
	Fifo Queuetype = iota

	//Lifo queue.
	Lifo

	//PriorityHigh means that on remove the elem with the highest priority value is returned.
	PriorityHigh

	//PriorityLow means that on remove the elem with the lowest priority value is returned.
	PriorityLow
)

//Queue is a queue of type Queuetype
type Queue struct {
	order                           Queuetype
	shrinkFactor, afterShrinkFactor float64
	lock                            sync.Mutex
	queSlice                        []*queueElement
}

//Invariants:
//  Always structure the slice in a way that the item at len(queSlice)-1 is the item for the remove operation
//	Fifo:
//		first elem in slice is the last inserted elem
//	Lifo:
//		first elem in slice is the first inserted elem
//	PriorityHigh:
//		first elem in slice has lowest priority, last elem highest
//	PriorityLow:
//		first elem in slice has highest priority, last elem lowest

type queueElement struct {
	priority float64
	content  interface{}
}

//NewQueue builds a new Queue with the passed Queuetype.
//Since the queue is realized through a slice, expectedLength is the initial
//cap() and len() value of said slice.
func NewQueue(tp Queuetype) *Queue {
	return &Queue{order: tp, queSlice: make([]*queueElement, 0), shrinkFactor: 1.5, afterShrinkFactor: 1.25}
}

//NewQueueElementWithPriority builds a new QueueElement with the passed content and priority.
//You cannot work with the element directly. This return value is only meant to be passed to
//queue functions.
func NewQueueElementWithPriority(c interface{}, priority float64) *queueElement {
	return &queueElement{priority: priority, content: c}
}

//NewQueueElementWithoutPriority builds a new QueueElement with the passed content and priority = 0.
//You cannot work with the element directly. This return value is only meant to be passed to
//queue functions.
func NewQueueElementWithoutPriority(c interface{}) *queueElement {
	return &queueElement{content: c}
}

//Append literally appends the element to the queue.
//Append does not uphold the invariant of the queue defined by the Queuetype.
//Use Insert for honoring the invariant.
func (q *Queue) Append(elem *queueElement) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queSlice = append(q.queSlice, elem)
}

//Insert inserts the passed element into the queue, accoring to the Queuetype of the queue.
//Insert upholds the invariant of the Queue.
//When there are multiple elems with the same priority the inserted elem will be the first that is removed.
func (q *Queue) Insert(elem *queueElement) {
	q.lock.Lock()
	defer q.lock.Unlock()

	switch q.order {
	case Fifo:
		q.insertFifo(elem)
	case Lifo:
		q.insertLifo(elem)
	case PriorityHigh:
		q.insertPriorityHigh(elem)
	case PriorityLow:
		q.insertPriorityLow(elem)
	default:
		panic("Queue has unknown order type. Can not insert elem in queue of unknown order type.")
	}
}

//Remove pops the element that is meant to be removed first according to the queues order.
//Returns the queueElement split up into its pieces.
func (q *Queue) Remove() (interface{}, float64) {
	q.lock.Lock()
	defer q.lock.Unlock()

	elemPointer := q.remove()
	return elemPointer.content, elemPointer.priority
}

//RemoveElement pops the element that is meant to be removed first according to the queues order.
//Returns the pointer to the queueElement itself.
func (q *Queue) RemoveElement() *queueElement {
	q.lock.Lock()
	defer q.lock.Unlock()

	return q.remove()
}

//DeletePointer deletes all occurences of elem out of the queue.
//This function works independent from the queueType.
//Returns true if elem was found and deleted, false else.
func (q *Queue) DeletePointer(elem *queueElement) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	removeFlag := false
	for i, e := range q.queSlice {
		if e == elem {
			q.delete(i)
			removeFlag = true
		}
	}
	return removeFlag
}

//DeleteElem deletes all occurences of elem out of the queue.
//This function works independent from the queueType.
//Returns true if elem was found and deleted, false else.
func (q *Queue) DeleteElem(elem queueElement) bool {
	q.lock.Lock()
	defer q.lock.Unlock()

	removeFlag := false
	for i, e := range q.queSlice {
		if *e == elem {
			q.delete(i)
			removeFlag = true
		}
	}
	return removeFlag
}
