package queue

import (
	"log"
)

func (q *Queue) insertFifo(elem *QueueElement) {
	q.queSlice = append([]*QueueElement{elem}, q.queSlice...)
}

func (q *Queue) insertLifo(elem *QueueElement) {
	q.queSlice = append(q.queSlice, elem)
}

func (q *Queue) insertPriorityHigh(elem *QueueElement) {
	if q.numElems == 0 || (q.queSlice[q.numElems-1]).priority < elem.priority {
		q.queSlice = append(q.queSlice, elem)
		return
	}
	for i, e := range q.queSlice {
		if e.priority < elem.priority {
			continue
		}

		//e.prio >= elem.prio
		q.queSlice = append(q.queSlice[:i], append([]*QueueElement{elem}, (q.queSlice)[i:]...)...)
		break
	}
}

func (q *Queue) insertPriorityLow(elem *QueueElement) {
	if q.numElems == 0 || (q.queSlice[q.numElems-1]).priority > elem.priority {
		q.queSlice = append(q.queSlice, elem)
		return
	}
	for i, e := range q.queSlice {
		if e.priority > elem.priority {
			continue
		}

		//e.prio <= elem.prio
		q.queSlice = append(q.queSlice[:i], append([]*QueueElement{elem}, (q.queSlice)[i:]...)...)
		break
	}
}

func (q *Queue) insertFifoLimited(elem *QueueElement) {
	if q.numElems == q.maxNumElems && q.maxNumElems != 0 {
		_, err := q.remove()
		if err != nil {
			log.Fatalln(err.Error())
		}
		q.numElems--
	}
	q.insertFifo(elem)
}
