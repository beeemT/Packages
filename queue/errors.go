package queue

import (
	"fmt"
)

//EmptyListError is the error that is returned on operations that encounter an empty
//queue but need a queue with elements inside.
type EmptyListError struct {
}

//IndexOutOfBoundsError is the error that is returned on operations where an index is provided but that
//index is not within the addressable space of the queue.
type IndexOutOfBoundsError struct {
}

//InvalidQueuetypeError is returned when a nonexistent queuetype is encountered
type InvalidQueuetypeError struct {
}

func (e EmptyListError) Error() string {
	return fmt.Sprintf("%s", "Queue is empty.")
}

func (e IndexOutOfBoundsError) Error() string {
	return fmt.Sprintf("%s", "Provided Index is out of Bounds.")
}

func (e InvalidQueuetypeError) Error() string {
	return fmt.Sprintf("%s", "Provided Queuetype is not a valid Queuetype.")
}
