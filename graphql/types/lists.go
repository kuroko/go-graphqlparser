// Code generated by tools/listgen
// DO NOT EDIT!
package types

// Errors is a linked list that contains Error values.
type Errors struct {
	Data Error
	next *Errors
	pos  int
}

// Add appends a Error to this linked list and returns this new head.
func (es *Errors) Add(data Error) *Errors {
	var pos int

	if es != nil {
		pos = es.pos + 1
	}

	return &Errors{
		Data: data,
		next: es,
		pos:  pos,
	}
}

// ForEach applies the given map function to each item in this linked list.
func (es *Errors) ForEach(fn func(e Error, i int)) {
	if es == nil {
		return
	}

	iter := 0
	current := es

	for {
		fn(current.Data, iter)

		if current.next == nil {
			break
		}

		iter++
		current = current.next
	}
}

// ErrorsGenerator is a type used to iterate efficiently over Errors.
// @wg:ignore
type ErrorsGenerator struct {
	original *Errors
	current  *Errors
	iter     int
	length   int
}

// Next returns the current value, and it's index in the list, and sets up the next value to be
// returned.
func (g *ErrorsGenerator) Next() (Error, int) {
	if g.current == nil {
		return Error{}, -1
	}

	retv := g.current.Data
	reti := g.iter

	g.current = g.current.next
	g.iter++

	return retv, reti
}

// Reset returns this generator to it's initial state, allowing it to be used again to iterate over
// this linked list.
func (g *ErrorsGenerator) Reset() {
	g.current = g.original
	g.iter = 0
}

// Generator returns a "Generator" type for this list, allowing for much more efficient iteration
// over items within this linked list than using ForEach, though ForEach may still be more
// convenient, because ForEach is a high order function, it's slower.
func (es *Errors) Generator() ErrorsGenerator {
	return ErrorsGenerator{
		current: es,
		iter:    0,
		length:  es.Len(),
	}
}

// Insert places the Error in the position given by pos.
// The method will insert at top if pos is greater than or equal to list length.
// The method will insert at bottom if the pos is less than 0.
func (es *Errors) Insert(e Error, pos int) *Errors {
	if pos >= es.Len() || es == nil {
		return es.Add(e)
	}

	if pos < 0 {
		pos = 0
	}

	mid := es
	for mid.pos != pos {
		mid = mid.next
	}

	bot := mid.next
	mid.next = nil
	es.pos -= mid.pos

	bot = bot.Add(e)
	es.Join(bot)

	return es
}

// Join attaches the tail of the receiver list "es" to the head of the otherList.
func (es *Errors) Join(otherList *Errors) {
	if es == nil {
		return
	}

	pos := es.Len() + otherList.Len() - 1

	last := es
	for es != nil {
		es.pos = pos
		pos--
		last = es
		es = es.next
	}

	last.next = otherList
}

// Len returns the length of this linked list.
func (es *Errors) Len() int {
	if es == nil {
		return 0
	}
	return es.pos + 1
}

// Reverse reverses this linked list of Error. Usually when the linked list is being
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the
// "right" order.
func (es *Errors) Reverse() *Errors {
	current := es

	var prev *Errors
	var pos int

	for current != nil {
		current.pos = pos
		pos++

		next := current.next
		current.next = prev
		prev = current
		current = next
	}

	return prev
}

// ErrorsFromSlice returns a Errors list from a slice of Error.
func ErrorsFromSlice(sl []Error) *Errors {
	var list *Errors
	for _, v := range sl {
		list = list.Add(v)
	}
	return list.Reverse()
}
