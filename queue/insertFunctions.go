package queue

func (q *Queue) insertFifo(elem Element) {
	q.queSlice = append([]Element{elem}, q.queSlice...)
}

func (q *Queue) insertLifo(elem Element) {
	q.queSlice = append(q.queSlice, elem)
}

func (q *Queue) insertPriorityHigh(elem Element) {
	//If the queue is empty or the new element has a higher priority than the current item with the highest priority
	//it can be appended to the slice.
	if q.numElements == 0 || (q.queSlice[q.numElements-1]).Priority() < elem.Priority() {
		q.queSlice = append(q.queSlice, elem)
		return
	}

	if (q.queSlice[q.numElements-1]).Priority() == elem.Priority() {
		q.backtrackInsertionPoint(elem)
	}

	//Default case. Iterate through full queue until the first suitable spot for the new element is found.
	for i, e := range q.queSlice {
		if e.Priority() < elem.Priority() {
			continue
		}

		//e.prio >= elem.prio
		q.queSlice = append(q.queSlice[:(i-1)], append([]Element{elem}, q.queSlice[(i-1):]...)...)
		break
	}
}

func (q *Queue) insertPriorityLow(elem Element) {
	//If the queue is empty or the new element has a lower priority than the current item with the lowest priority
	//it can be appended to the slice.
	if q.numElements == 0 || (q.queSlice[q.numElements-1]).Priority() > elem.Priority() {
		q.queSlice = append(q.queSlice, elem)
		return
	}

	if (q.queSlice[q.numElements-1]).Priority() == elem.Priority() {
		q.backtrackInsertionPoint(elem)
	}

	//Default case. Iterate through full queue until the first suitable spot for the new element is found.
	for i, e := range q.queSlice {
		if e.Priority() > elem.Priority() {
			continue
		}

		//e.prio <= elem.prio
		q.queSlice = append(q.queSlice[:(i-1)], append([]Element{elem}, q.queSlice[(i-1):]...)...)
		break
	}
}

//backtrackInsertionPoint finds a suitable insertion point for an item from the back of the queue.
//queSlice[q.numElements-1].Priority() == elem.Priority() doesn't need full iteration over queue to find the spot for insertion.
//Iteration from slice end to first element which has non equal priority than the new element is enough because of the priority invariant.
//In the worst case this will iterate over whole queue, so attach new element to front of slice.
func (q *Queue) backtrackInsertionPoint(elem Element) {

	for i := q.numElements - 1; i > -1; i-- {
		if q.queSlice[i].Priority() == elem.Priority() {
			continue
		}
		q.queSlice = append(q.queSlice[i:], append([]Element{elem}, q.queSlice[:i]...)...)
		return
	}
	q.queSlice = append([]Element{elem}, q.queSlice...)

}

func (q *Queue) insertFifoLimited(elem Element) error {
	if q.numElements == q.maxnumElements && q.maxnumElements != 0 {
		_, err := q.remove()
		if err != nil {
			return err
		}
		q.numElements--
	}
	q.insertFifo(elem)
	return nil
}
