// Package ast - THIS CODE IS GENERATED, DO NOT EDIT MANUALLY
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

// DirectiveLocations is a linked list that contains DirectiveLocation values.
type DirectiveLocations struct {
	Data DirectiveLocation
	Next *DirectiveLocations
}

// ForEach applies the given map function to each item in this linked list.
func (d *DirectiveLocations) ForEach(fn func(directiveLocation DirectiveLocation, i int)) {
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
func (d *DirectiveLocations) Len() int {
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

// Reverse reverses this linked list of DirectiveLocation. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *DirectiveLocations) Reverse() *DirectiveLocations {
	current := d

	var prev *DirectiveLocations
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// EnumValueDefinitions is a linked list that contains EnumValueDefinition values.
type EnumValueDefinitions struct {
	Data EnumValueDefinition
	Next *EnumValueDefinitions
}

// ForEach applies the given map function to each item in this linked list.
func (d *EnumValueDefinitions) ForEach(fn func(enumValueDefinition EnumValueDefinition, i int)) {
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
func (d *EnumValueDefinitions) Len() int {
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

// Reverse reverses this linked list of EnumValueDefinition. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *EnumValueDefinitions) Reverse() *EnumValueDefinitions {
	current := d

	var prev *EnumValueDefinitions
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// FieldDefinitions is a linked list that contains FieldDefinition values.
type FieldDefinitions struct {
	Data FieldDefinition
	Next *FieldDefinitions
}

// ForEach applies the given map function to each item in this linked list.
func (d *FieldDefinitions) ForEach(fn func(fieldDefinition FieldDefinition, i int)) {
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
func (d *FieldDefinitions) Len() int {
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

// Reverse reverses this linked list of FieldDefinition. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *FieldDefinitions) Reverse() *FieldDefinitions {
	current := d

	var prev *FieldDefinitions
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// InputValueDefinitions is a linked list that contains InputValueDefinition values.
type InputValueDefinitions struct {
	Data InputValueDefinition
	Next *InputValueDefinitions
}

// ForEach applies the given map function to each item in this linked list.
func (d *InputValueDefinitions) ForEach(fn func(inputValueDefinition InputValueDefinition, i int)) {
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
func (d *InputValueDefinitions) Len() int {
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

// Reverse reverses this linked list of InputValueDefinition. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *InputValueDefinitions) Reverse() *InputValueDefinitions {
	current := d

	var prev *InputValueDefinitions
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// OperationTypeDefinitions is a linked list that contains OperationTypeDefinition values.
type OperationTypeDefinitions struct {
	Data OperationTypeDefinition
	Next *OperationTypeDefinitions
}

// ForEach applies the given map function to each item in this linked list.
func (d *OperationTypeDefinitions) ForEach(fn func(operationTypeDefinition OperationTypeDefinition, i int)) {
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
func (d *OperationTypeDefinitions) Len() int {
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

// Reverse reverses this linked list of OperationTypeDefinition. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *OperationTypeDefinitions) Reverse() *OperationTypeDefinitions {
	current := d

	var prev *OperationTypeDefinitions
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

// RootOperationTypeDefinitions is a linked list that contains RootOperationTypeDefinition values.
type RootOperationTypeDefinitions struct {
	Data RootOperationTypeDefinition
	Next *RootOperationTypeDefinitions
}

// ForEach applies the given map function to each item in this linked list.
func (d *RootOperationTypeDefinitions) ForEach(fn func(rootOperationTypeDefinition RootOperationTypeDefinition, i int)) {
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
func (d *RootOperationTypeDefinitions) Len() int {
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

// Reverse reverses this linked list of RootOperationTypeDefinition. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *RootOperationTypeDefinitions) Reverse() *RootOperationTypeDefinitions {
	current := d

	var prev *RootOperationTypeDefinitions
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

// Types is a linked list that contains Type values.
type Types struct {
	Data Type
	Next *Types
}

// ForEach applies the given map function to each item in this linked list.
func (d *Types) ForEach(fn func(t Type, i int)) {
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
func (d *Types) Len() int {
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

// Reverse reverses this linked list of Type. Usually when the linked list is being 
// constructed the result will be last-to-first, so we'll want to reverse it to get it in the 
// "right" order.
func (d *Types) Reverse() *Types {
	current := d

	var prev *Types
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
