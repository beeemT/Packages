package queue

import (
	"sync"
)

//Queuetype is the enum type for queue invariants.
//Invariants:
//  Always structure the slice in a way that the item at len(queSlice)-1 is the item for the remove operation
//  For same main ordering property of two elements the element that is older will be removed.
//	Fifo:
//		len(queSlice)-1 is the fist inserted elem
//	Lifo:
//		len(queSlice)-1 is the last inserted elem
//	PriorityHigh:
//		len(queSlice)-1 is the elem with highest priority
//	PriorityLow:
//		len(queSlice)-1 is the elem with lowest priority
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

	//FifoLimited means that the queue has a maximum capacity. Requires extra call to set capacity.
	FifoLimited

	numQueuetypes = 5
)

//Element is the interface encapsulating all element types
type Element interface {
	Priority() float64
	SetPriority(float64)

	Content() interface{}
	SetContent(interface{})
}

//Queue is a queue of type Queuetype
type Queue struct {
	order                           Queuetype
	shrinkFactor, afterShrinkFactor float64
	lock                            sync.Mutex
	queSlice                        []Element
	numElements                     int
	maxnumElements                  int
}

//NewQueue builds a new Queue with the passed Queuetype.
//Since the queue is realized through a slice, expectedLength is the initial
//cap() value of said slice.
func NewQueue(tp Queuetype) (*Queue, error) {
	if tp < 0 || tp > numQueuetypes {
		return nil, InvalidQueuetypeError{}
	}
	return &Queue{order: tp, queSlice: make([]Element, 0), shrinkFactor: 0.5, afterShrinkFactor: 0.75}, nil
}

//NewElementF64WithPriority builds a new Element with the passed content and priority.
//You cannot work with the element directly. This return value is only meant to be passed to
//queue functions.
func NewElementF64WithPriority(c interface{}, priority float64) Element {
	return &ElementF64W{priority: priority, content: &c}
}

//NewBaseElement builds a new Element with the passed content and priority = 0.
//You cannot work with the element directly. This return value is only meant to be passed to
//queue functions.
func NewBaseElement(c interface{}) Element {
	return &BaseElement{content: &c}
}

//Len returns the number of elements in the queue.
func (q *Queue) Len() int {
	return q.numElements
}

//SetLimit sets the max capacity for the queue. Returns a InvalidQueueLimitError if limit < 0.
func (q *Queue) SetLimit(limit int) error {
	if limit < 0 {
		return InvalidQueueLimitError{}
	}
	q.maxnumElements = limit
	return nil
}

//Append literally appends the element to the queue.
//Append does not uphold the invariant of the queue defined by the Queuetype and is thus unsafe.
//Use Insert for honoring the invariant.
func (q *Queue) Append(elem Element) {
	q.lock.Lock()
	defer q.lock.Unlock()

	q.queSlice = append(q.queSlice, elem)
	q.numElements++
}

//Insert inserts the passed element into the queue, according to the Queuetype of the queue.
//Insert upholds the invariant of the Queue.
//When there are multiple elements with the same priority the oldest elem will be the first that is removed.
func (q *Queue) Insert(elem Element) error {
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
	case FifoLimited:
		return q.insertFifoLimited(elem)
	default:
		return InvalidQueuetypeError{}
	}
	q.numElements++
	return nil
}

//Remove pops the element that is meant to be removed first according to the queues order.
//When there are multiple elements with the same priority the oldest elem will be the first that is removed (FIFO).
//Returns the Element split up into its pieces.
//If the list is empty, an error is returned.
func (q *Queue) Remove() (interface{}, float64, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	elemPointer, err := q.remove()
	if err != nil {
		return nil, 0, err
	}
	q.numElements--
	return elemPointer.Content(), elemPointer.Priority(), nil
}

//RemoveElement pops the element that is meant to be removed first according to the queues order.
//When there are multiple elements with the same priority the oldest elem will be the first that is removed.
//Returns the pointer to the Element itself.
func (q *Queue) RemoveElement() (Element, error) {
	q.lock.Lock()
	defer q.lock.Unlock()

	elemP, err := q.remove()
	if err != nil {
		return nil, err
	}

	q.numElements--
	return elemP, nil
}

//DeletePointer deletes all occurrences of elem out of the queue.
//This function works independent from the queueType.
//Returns number of removed elements.
func (q *Queue) DeletePointer(elem Element) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	counter := 0
	for i, e := range q.queSlice {
		if e == elem {
			q.delete(i)
			counter++
			q.numElements--
		}
	}
	return counter
}

//DeleteElem deletes all occurrences of elem out of the queue.
//This function works independent from the queueType.
//Returns number of removed elements.
func (q *Queue) DeleteElem(elem Element) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	counter := 0
	for i, e := range q.queSlice {
		if e == elem {
			q.delete(i)
			counter++
			q.numElements--
		}
	}
	return counter
}

//UpdatePriority updates the priority of all elements with priority oldPriority to the newPriority.
//Upholds the invariant of the queue.
//Returns the number of updates.
//If performanceFlag is set, elements with the same priority will be reversed in their order for ordertypes
//PriorityHigh and PriorityLow.
func (q *Queue) UpdatePriority(oldPriority, newPriority float64, performanceFlag bool) int {
	q.lock.Lock()
	defer q.lock.Unlock()

	counter := 0

	var list []Element
	if !performanceFlag {
		list = make([]Element, 0) //for buffering elements for reinsertion
	}

	switch q.order {
	case Lifo, Fifo:
		for _, e := range q.queSlice { //O(n)
			//modifing e works because queSlice is Element
			//+ Lifo and Fifo both are not sorted after priority

			if e.Priority() == oldPriority {
				e.SetPriority(newPriority)
				counter++
			}
		}

	case PriorityHigh, PriorityLow:
		//todo: use binsearch to find first elem with priority
		var modFlag bool

		for i, e := range q.queSlice {
			if e.Priority() == oldPriority {
				q.deleteWithoutMemoryManagement(i) //delete without MemoryManagement because elements get reinserted
				e.SetPriority(newPriority)
				if performanceFlag {
					q.Insert(e) //reverses the order within elements with the same priority
				} else {
					list = append(list, e)
				}
			} else if modFlag {
				break
			}
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

//GetAllElements returns a slice of references to all elements content.
func (q *Queue) GetAllElements() []interface{} {
	ret := make([]interface{}, q.numElements)
	for _, elem := range q.queSlice {
		ret = append(ret, elem.Content())
	}
	return ret
}
