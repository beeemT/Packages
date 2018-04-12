package queue

func (q *Queue) insertFifo(elem *queueElement) {
	q.queSlice = append([]*queueElement{elem}, q.queSlice...)
}

func (q *Queue) insertLifo(elem *queueElement) {
	q.queSlice = append(q.queSlice, elem)
}

func (q *Queue) insertPriorityHigh(elem *queueElement) {
	if q.numElems == 0 || (q.queSlice[q.numElems-1]).priority < elem.priority {
		q.queSlice = append(q.queSlice, elem)
		return
	}
	for i, e := range q.queSlice {
		if e.priority < elem.priority {
			continue
		}

		//e.prio >= elem.prio
		q.queSlice = append(q.queSlice[:i], append([]*queueElement{elem}, (q.queSlice)[i:]...)...)
		break
	}
}

func (q *Queue) insertPriorityLow(elem *queueElement) {
	if q.numElems == 0 || (q.queSlice[q.numElems-1]).priority > elem.priority {
		q.queSlice = append(q.queSlice, elem)
		return
	}
	for i, e := range q.queSlice {
		if e.priority > elem.priority {
			continue
		}

		//e.prio <= elem.prio
		q.queSlice = append(q.queSlice[:i], append([]*queueElement{elem}, (q.queSlice)[i:]...)...)
		break
	}
}
