// Package ast - THIS CODE IS GENERATED, DO NOT EDIT MANUALLY
package ast

// Arguments is a linked list that contains Argument values.
type Arguments struct {
	Data Argument
	Next *Arguments
}

// Add appends a Argument to this linked list and returns this new head.
func (a *Arguments) Add(data Argument) *Arguments {
	return &Arguments{
		Data: data,
		Next: a,
	}
}

// Join attaches the tail of the reciever list "a" to the head of the otherList.
func (a *Arguments) Join(otherList *Arguments) {
	current := a

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (a *Arguments) ForEach(fn func(argument Argument, i int)) {
	if a == nil {
		return
	}

	iter := 0
	current := a

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
func (a *Arguments) Len() int {
	if a == nil {
		return 0
	}

	var length int

	current := a
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
func (a *Arguments) Reverse() *Arguments {
	current := a

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

// Add appends a Definition to this linked list and returns this new head.
func (d *Definitions) Add(data Definition) *Definitions {
	return &Definitions{
		Data: data,
		Next: d,
	}
}

// Join attaches the tail of the reciever list "d" to the head of the otherList.
func (d *Definitions) Join(otherList *Definitions) {
	current := d

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
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

// Add appends a Directive to this linked list and returns this new head.
func (d *Directives) Add(data Directive) *Directives {
	return &Directives{
		Data: data,
		Next: d,
	}
}

// Join attaches the tail of the reciever list "d" to the head of the otherList.
func (d *Directives) Join(otherList *Directives) {
	current := d

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
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

// Add appends a DirectiveLocation to this linked list and returns this new head.
func (dl *DirectiveLocations) Add(data DirectiveLocation) *DirectiveLocations {
	return &DirectiveLocations{
		Data: data,
		Next: dl,
	}
}

// Join attaches the tail of the reciever list "dl" to the head of the otherList.
func (dl *DirectiveLocations) Join(otherList *DirectiveLocations) {
	current := dl

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (dl *DirectiveLocations) ForEach(fn func(directiveLocation DirectiveLocation, i int)) {
	if dl == nil {
		return
	}

	iter := 0
	current := dl

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
func (dl *DirectiveLocations) Len() int {
	if dl == nil {
		return 0
	}

	var length int

	current := dl
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
func (dl *DirectiveLocations) Reverse() *DirectiveLocations {
	current := dl

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

// Add appends a EnumValueDefinition to this linked list and returns this new head.
func (evd *EnumValueDefinitions) Add(data EnumValueDefinition) *EnumValueDefinitions {
	return &EnumValueDefinitions{
		Data: data,
		Next: evd,
	}
}

// Join attaches the tail of the reciever list "evd" to the head of the otherList.
func (evd *EnumValueDefinitions) Join(otherList *EnumValueDefinitions) {
	current := evd

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (evd *EnumValueDefinitions) ForEach(fn func(enumValueDefinition EnumValueDefinition, i int)) {
	if evd == nil {
		return
	}

	iter := 0
	current := evd

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
func (evd *EnumValueDefinitions) Len() int {
	if evd == nil {
		return 0
	}

	var length int

	current := evd
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
func (evd *EnumValueDefinitions) Reverse() *EnumValueDefinitions {
	current := evd

	var prev *EnumValueDefinitions
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}

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

// Join attaches the tail of the reciever list "e" to the head of the otherList.
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

// FieldDefinitions is a linked list that contains FieldDefinition values.
type FieldDefinitions struct {
	Data FieldDefinition
	Next *FieldDefinitions
}

// Add appends a FieldDefinition to this linked list and returns this new head.
func (fd *FieldDefinitions) Add(data FieldDefinition) *FieldDefinitions {
	return &FieldDefinitions{
		Data: data,
		Next: fd,
	}
}

// Join attaches the tail of the reciever list "fd" to the head of the otherList.
func (fd *FieldDefinitions) Join(otherList *FieldDefinitions) {
	current := fd

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (fd *FieldDefinitions) ForEach(fn func(fieldDefinition FieldDefinition, i int)) {
	if fd == nil {
		return
	}

	iter := 0
	current := fd

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
func (fd *FieldDefinitions) Len() int {
	if fd == nil {
		return 0
	}

	var length int

	current := fd
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
func (fd *FieldDefinitions) Reverse() *FieldDefinitions {
	current := fd

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

// Add appends a InputValueDefinition to this linked list and returns this new head.
func (ivd *InputValueDefinitions) Add(data InputValueDefinition) *InputValueDefinitions {
	return &InputValueDefinitions{
		Data: data,
		Next: ivd,
	}
}

// Join attaches the tail of the reciever list "ivd" to the head of the otherList.
func (ivd *InputValueDefinitions) Join(otherList *InputValueDefinitions) {
	current := ivd

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (ivd *InputValueDefinitions) ForEach(fn func(inputValueDefinition InputValueDefinition, i int)) {
	if ivd == nil {
		return
	}

	iter := 0
	current := ivd

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
func (ivd *InputValueDefinitions) Len() int {
	if ivd == nil {
		return 0
	}

	var length int

	current := ivd
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
func (ivd *InputValueDefinitions) Reverse() *InputValueDefinitions {
	current := ivd

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

// Add appends a OperationTypeDefinition to this linked list and returns this new head.
func (otd *OperationTypeDefinitions) Add(data OperationTypeDefinition) *OperationTypeDefinitions {
	return &OperationTypeDefinitions{
		Data: data,
		Next: otd,
	}
}

// Join attaches the tail of the reciever list "otd" to the head of the otherList.
func (otd *OperationTypeDefinitions) Join(otherList *OperationTypeDefinitions) {
	current := otd

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (otd *OperationTypeDefinitions) ForEach(fn func(operationTypeDefinition OperationTypeDefinition, i int)) {
	if otd == nil {
		return
	}

	iter := 0
	current := otd

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
func (otd *OperationTypeDefinitions) Len() int {
	if otd == nil {
		return 0
	}

	var length int

	current := otd
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
func (otd *OperationTypeDefinitions) Reverse() *OperationTypeDefinitions {
	current := otd

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

// Add appends a RootOperationTypeDefinition to this linked list and returns this new head.
func (rotd *RootOperationTypeDefinitions) Add(data RootOperationTypeDefinition) *RootOperationTypeDefinitions {
	return &RootOperationTypeDefinitions{
		Data: data,
		Next: rotd,
	}
}

// Join attaches the tail of the reciever list "rotd" to the head of the otherList.
func (rotd *RootOperationTypeDefinitions) Join(otherList *RootOperationTypeDefinitions) {
	current := rotd

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (rotd *RootOperationTypeDefinitions) ForEach(fn func(rootOperationTypeDefinition RootOperationTypeDefinition, i int)) {
	if rotd == nil {
		return
	}

	iter := 0
	current := rotd

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
func (rotd *RootOperationTypeDefinitions) Len() int {
	if rotd == nil {
		return 0
	}

	var length int

	current := rotd
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
func (rotd *RootOperationTypeDefinitions) Reverse() *RootOperationTypeDefinitions {
	current := rotd

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

// Add appends a Selection to this linked list and returns this new head.
func (s *Selections) Add(data Selection) *Selections {
	return &Selections{
		Data: data,
		Next: s,
	}
}

// Join attaches the tail of the reciever list "s" to the head of the otherList.
func (s *Selections) Join(otherList *Selections) {
	current := s

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (s *Selections) ForEach(fn func(selection Selection, i int)) {
	if s == nil {
		return
	}

	iter := 0
	current := s

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
func (s *Selections) Len() int {
	if s == nil {
		return 0
	}

	var length int

	current := s
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
func (s *Selections) Reverse() *Selections {
	current := s

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

// Add appends a Type to this linked list and returns this new head.
func (t *Types) Add(data Type) *Types {
	return &Types{
		Data: data,
		Next: t,
	}
}

// Join attaches the tail of the reciever list "t" to the head of the otherList.
func (t *Types) Join(otherList *Types) {
	current := t

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (t *Types) ForEach(fn func(t Type, i int)) {
	if t == nil {
		return
	}

	iter := 0
	current := t

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
func (t *Types) Len() int {
	if t == nil {
		return 0
	}

	var length int

	current := t
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
func (t *Types) Reverse() *Types {
	current := t

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

// Add appends a VariableDefinition to this linked list and returns this new head.
func (vd *VariableDefinitions) Add(data VariableDefinition) *VariableDefinitions {
	return &VariableDefinitions{
		Data: data,
		Next: vd,
	}
}

// Join attaches the tail of the reciever list "vd" to the head of the otherList.
func (vd *VariableDefinitions) Join(otherList *VariableDefinitions) {
	current := vd

	for current.Next != nil {
		current = current.Next
	}

	current.Next = otherList
}

// ForEach applies the given map function to each item in this linked list.
func (vd *VariableDefinitions) ForEach(fn func(variableDefinition VariableDefinition, i int)) {
	if vd == nil {
		return
	}

	iter := 0
	current := vd

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
func (vd *VariableDefinitions) Len() int {
	if vd == nil {
		return 0
	}

	var length int

	current := vd
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
func (vd *VariableDefinitions) Reverse() *VariableDefinitions {
	current := vd

	var prev *VariableDefinitions
	for current != nil {
		next := current.Next
		current.Next = prev
		prev = current
		current = next
	}

	return prev
}
