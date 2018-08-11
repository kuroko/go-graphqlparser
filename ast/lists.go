// THIS CODE IS GENERATED, DO NOT EDIT MANUALLY
package ast

// Arguments is a linked list that contains Argument values.
type Arguments struct {
	Data Argument
	Next *Arguments
}

// ForEach applies the given map function to each item in this linked list.
func (d *Arguments) ForEach(fn func(argument Argument, i int)) {
	if d == nil {
		return
	}

	iter := 0
	current := d

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
func (d *Arguments) Len() int {
	if d == nil {
		return 0
	}

	var length int

	current := d
	for {
		length++
		if current.Next == nil {
			break
		}

		current = current.Next
	}

	return length
}

// Reverse reverses this linked list of Argument. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *Arguments) Reverse() *Arguments {
	current := d

	var prev *Arguments
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// Definitions is a linked list that contains Definition values.
type Definitions struct {
	Data Definition
	Next *Definitions
}

// ForEach applies the given map function to each item in this linked list.
func (d *Definitions) ForEach(fn func(definition Definition, i int)) {
	if d == nil {
		return
	}

	iter := 0
	current := d

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
func (d *Definitions) Len() int {
	if d == nil {
		return 0
	}

	var length int

	current := d
	for {
		length++
		if current.Next == nil {
			break
		}

		current = current.Next
	}

	return length
}

// Reverse reverses this linked list of Definition. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *Definitions) Reverse() *Definitions {
	current := d

	var prev *Definitions
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// Directives is a linked list that contains Directive values.
type Directives struct {
	Data Directive
	Next *Directives
}

// ForEach applies the given map function to each item in this linked list.
func (d *Directives) ForEach(fn func(directive Directive, i int)) {
	if d == nil {
		return
	}

	iter := 0
	current := d

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
func (d *Directives) Len() int {
	if d == nil {
		return 0
	}

	var length int

	current := d
	for {
		length++
		if current.Next == nil {
			break
		}

		current = current.Next
	}

	return length
}

// Reverse reverses this linked list of Directive. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *Directives) Reverse() *Directives {
	current := d

	var prev *Directives
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// Selections is a linked list that contains Selection values.
type Selections struct {
	Data Selection
	Next *Selections
}

// ForEach applies the given map function to each item in this linked list.
func (d *Selections) ForEach(fn func(selection Selection, i int)) {
	if d == nil {
		return
	}

	iter := 0
	current := d

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
func (d *Selections) Len() int {
	if d == nil {
		return 0
	}

	var length int

	current := d
	for {
		length++
		if current.Next == nil {
			break
		}

		current = current.Next
	}

	return length
}

// Reverse reverses this linked list of Selection. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *Selections) Reverse() *Selections {
	current := d

	var prev *Selections
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// VariableDefinitions is a linked list that contains VariableDefinition values.
type VariableDefinitions struct {
	Data VariableDefinition
	Next *VariableDefinitions
}

// ForEach applies the given map function to each item in this linked list.
func (d *VariableDefinitions) ForEach(fn func(variableDefinition VariableDefinition, i int)) {
	if d == nil {
		return
	}

	iter := 0
	current := d

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
func (d *VariableDefinitions) Len() int {
	if d == nil {
		return 0
	}

	var length int

	current := d
	for {
		length++
		if current.Next == nil {
			break
		}

		current = current.Next
	}

	return length
}

// Reverse reverses this linked list of VariableDefinition. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *VariableDefinitions) Reverse() *VariableDefinitions {
	current := d

	var prev *VariableDefinitions
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}
