package queue

import (
	"sync"
)

//Queuetype is the enum type for queue invariants.
//Invariants:
//  Always structure the slice in a way that the item at len(queSlice)-1 is the item for the remove operation
//  For same main ordering property of two elements the element that is older will be removed.
//	Fifo:
//		first elem in slice is the last inserted elem
//	Lifo:
//		first elem in slice is the first inserted elem
//	PriorityHigh:
//		first elem in slice has lowest priority, last elem highest
//	PriorityLow:
//		first elem in slice has highest priority, last elem lowest
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

	numQueuetypes = 4
)

//Queue is a queue of type Queuetype
type Queue struct {
	order                           Queuetype
	shrinkFactor, afterShrinkFactor float64
	lock                            sync.Mutex
	queSlice                        []*queueElement
	numElems                        int
}

type queueElement struct {
	priority float64
	content  interface{}
}

//NewQueue builds a new Queue with the passed Queuetype.
//Since the queue is realized through a slice, expectedLength is the initial
//cap() value of said slice.
func NewQueue(tp Queuetype) (*Queue, error) {
	if tp < 0 || tp > numQueuetypes {
		return nil, InvalidQueuetypeError{}
	}
	return &Queue{order: tp, queSlice: make([]*queueElement, 0), shrinkFactor: 0.5, afterShrinkFactor: 0.75}, nil
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

//Length returns the number of elements in the queue.
func (q *Queue) Length() int {
	return q.numElems
}

//Append literally appends the element to the queue.
//Append does not uphold the invariant of the queue defined by the Queuetype.
//Use Insert for honoring the invariant.
func (q *Queue) Append(elem *queueElement) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queSlice = append(q.queSlice, elem)
	q.numElems++
}

//Insert inserts the passed element into the queue, accoring to the Queuetype of the queue.
//Insert upholds the invariant of the Queue.
//When there are multiple elems with the same priority the oldest elem will be the first that is removed.
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
	q.numElems++
}

//Remove pops the element that is meant to be removed first according to the queues order.
//When there are multiple elems with the same priority the oldest elem will be the first that is removed.
//Returns the queueElement split up into its pieces.
//If the list is empty, an error is returned.
func (q *Queue) Remove() (interface{}, float64, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	elemPointer, err := q.remove()
	if err != nil {
		return nil, 0, err
	}
	q.numElems--
	return elemPointer.content, elemPointer.priority, nil
}

//RemoveElement pops the element that is meant to be removed first according to the queues order.
//When there are multiple elems with the same priority the oldest elem will be the first that is removed.
//Returns the pointer to the queueElement itself.
func (q *Queue) RemoveElement() (*queueElement, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	elemP, err := q.remove()
	if err != nil {
		return nil, err
	}

	q.numElems--
	return elemP, nil
}

//DeletePointer deletes all occurences of elem out of the queue.
//This function works independent from the queueType.
//Returns number of removed elems.
func (q *Queue) DeletePointer(elem *queueElement) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	counter := 0
	for i, e := range q.queSlice {
		if e == elem {
			q.delete(i)
			counter++
			q.numElems--
		}
	}
	return counter
}

//DeleteElem deletes all occurences of elem out of the queue.
//This function works independent from the queueType.
//Returns number of removed elems.
func (q *Queue) DeleteElem(elem queueElement) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	counter := 0
	for i, e := range q.queSlice {
		if *e == elem {
			q.delete(i)
			counter++
			q.numElems--
		}
	}
	return counter
}

//UpdatePriority updates the priority of all elements having the oldPriority with the newPriority.
//Upholds the order of the queue.
//Returns the number of updates.
//Fifo, Lifo -> no reordering of elements
//PriorityHigh, PriorityLow -> no reordering of elements unless performanceFlag is set.
//If performanceFlag is set, elements with the same priority will be reversed in their order.
func (q *Queue) UpdatePriority(oldPriority, newPriority float64, performanceFlag bool) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	counter := 0
	list := make([]*queueElement, 0) //for buffering elements for reinsertion

	for i, e := range q.queSlice {
		if e.priority == oldPriority {
			switch q.order {
			case Lifo, Fifo:
				//modifing e works because queSlice is *queueElement
				//+ Lifo and Fifo both are not sorted after priority
				e.priority = newPriority

			case PriorityHigh, PriorityLow:
				q.deleteWithoutMemoryManagement(i) //delete without MemoryManagement because elements get reinserted
				e.priority = newPriority
				if performanceFlag {
					q.Insert(e) //might reverse the order within the elems with the same priority (depends on q.order)
				} else {
					list = append(list, e)
				}
			}
			counter++
		}
	}

	if (q.order == PriorityHigh || q.order == PriorityLow) && !performanceFlag {
		l := len(list)
		for i := range list {
			q.Insert(list[l-(i+1)]) //insert oldest element first
		}
	}

	return counter
}
