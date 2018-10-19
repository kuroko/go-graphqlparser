package graphql

// Error ...
type Error error

// Errors is a linked list that contains Error values.
type Errors struct {
	Data Error
	Next *Errors
}

// Add appends a Error to this linked list and returns this new head.
func (e *Errors) Add(data Error) *Errors {
	return &Errors{
		Data: data,
		Next: e,
	}
}

// Join attaches the tail of the receiver list "e" to the head of the otherList.
func (e *Errors) Join(otherList *Errors) {
	current := e

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (e *Errors) ForEach(fn func(err Error, i int)) {
	if e == nil {
		return
	}

	iter := 0
	current := e

	for {
		fn(current.Data, iter)

		if current.Next == nil {
			break
		}

		iter++
		current = current.Next
	}
}

// Len returns the length of this linked list.
func (e *Errors) Len() int {
	if e == nil {
		return 0
	}

	var length int

	current := e
	for {
		length++
		if current.Next == nil {
			break
		}

		current = current.Next
	}

	return length
}

// Reverse reverses this linked list of Error. Usually when the linked list is being
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the
// "right" order.
func (e *Errors) Reverse() *Errors {
	current := e

	var prev *Errors
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}
