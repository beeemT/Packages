package queue

//ElementF64W encapsulates all information that is needed for the storage in the queue.
type ElementF64W struct {
	priority float64
	content  *interface{}
}

func (e ElementF64W) Priority() float64 {
	return e.priority
}

func (e *ElementF64W) SetPriority(priority float64) {
	e.priority = priority
}

func (e ElementF64W) Content() interface{} {
	return e.content
}

func (e *ElementF64W) SetContent(content interface{}) {
	e.content = &content
}

//BaseElement encapsulates all information that is needed for the storage in the queue.
type BaseElement struct {
	content *interface{}
}

func (e BaseElement) Priority() float64 {
	return 0
}

func (e *BaseElement) SetPriority(priority float64) {
}

func (e BaseElement) Content() interface{} {
	return e.content
}

func (e *BaseElement) SetContent(content interface{}) {
	e.content = &content
}
