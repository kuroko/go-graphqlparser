package ast

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestArguments_Add(t *testing.T) {
	var list *Arguments

	list = list.Add(Argument{Name: "0"}).Add(Argument{Name: "1"})

	assert.Equal(t, "1", list.Data.Name)
	assert.Equal(t, 1, list.pos)

	list = list.next

	assert.Equal(t, "0", list.Data.Name)
	assert.Equal(t, 0, list.pos)
}
func TestArguments_ForEach(t *testing.T) {
	var list *Arguments

	// Does nothing if no list items.
	var doesSomething bool
	list.ForEach(func(a Argument, i int) {
		doesSomething = true
	})
	assert.Equal(t, false, doesSomething)

	// Visits each list item.
	list = list.Add(Argument{Name: "one"})
	list = list.Add(Argument{Name: "two"})

	var names []string
	var indexes []int

	list.ForEach(func(a Argument, i int) {
		names = append(names, a.Name)
		indexes = append(indexes, i)
	})

	assert.Equal(t, []string{"two", "one"}, names)
	assert.Equal(t, []int{0, 1}, indexes)
}
func TestArguments_Insert(t *testing.T) {
	var list *Arguments

	// Insert into nil list.
	list = list.Insert(Argument{Name: "bar"}, 0)
	validate(t, list, []string{"bar"})

	// Insert into bottom of list.
	list = list.Insert(Argument{Name: "!"}, 0)
	validate(t, list, []string{"bar", "!"})

	// Insert into middle of list.
	list = list.Insert(Argument{Name: "baz"}, 1)
	validate(t, list, []string{"bar", "baz", "!"})

	// Insert into top of list.
	list = list.Insert(Argument{Name: "foo"}, list.Len())
	validate(t, list, []string{"foo", "bar", "baz", "!"})

	// Inserts > len insert into top.
	list = list.Insert(Argument{Name: "what's"}, 0x7FFFFFFF)
	validate(t, list, []string{"what's", "foo", "bar", "baz", "!"})

	// Inserts < 0 insert into bottom.
	list = list.Insert(Argument{Name: "?"}, -1)
	validate(t, list, []string{"what's", "foo", "bar", "baz", "!", "?"})
}
func TestArguments_Join(t *testing.T) {
	// Join a non nil list to a list of n, one and zero elements.
	var zero *Arguments

	var one *Arguments
	one = one.Add(Argument{Name: "zero"})

	var n *Arguments
	n = n.Add(Argument{Name: "one"})
	n = n.Add(Argument{Name: "two"})

	var list *Arguments
	list = list.Add(Argument{Name: "three"})
	list = list.Add(Argument{Name: "four"})

	list.Join(n)
	list.Join(one)
	list.Join(zero)
	validate(t, list, []string{"four", "three", "two", "one", "zero"})

	// Joining a nil list to a nil list, both are still nil.
	var nl1 *Arguments
	var nl2 *Arguments

	nl1.Join(nl2)
	assert.Nil(t, nl1)
	assert.Nil(t, nl2)

	// Joining a nil list to a non nil list, nil list is not changed.
	nl2 = nl2.Add(Argument{Name: "uno"})
	nl2 = nl2.Add(Argument{Name: "dos"})
	nl1.Join(nl2)
	validate(t, nl1, nil)
}
func TestArguments_Len(t *testing.T) {
	n := (*Arguments).Add(nil, Argument{}).Add(Argument{}).Len()
	assert.Equal(t, 2, n)

	one := (*Arguments).Add(nil, Argument{}).Len()
	assert.Equal(t, 1, one)

	zero := (*Arguments).Len(nil)
	assert.Equal(t, 0, zero)
}

func TestArguments_Reverse(t *testing.T) {
	var list *Arguments

	// One element can be reversed.
	list = list.Add(Argument{Name: "first"})
	one := list.Reverse()
	assert.Equal(t, "first", list.Data.Name)
	assert.Equal(t, "first", one.Data.Name)

	// reset
	one.Reverse()

	// n elements can be reversed.
	list = list.Add(Argument{Name: "second"})
	n := list.Reverse()
	assert.Equal(t, "second", list.Data.Name)
	assert.Equal(t, "first", n.Data.Name)

	// reset
	n.Reverse()

	// Can be reversed multiple times and not mutate.
	r2 := list.Reverse().Reverse()
	assert.Equal(t, list, r2)

	// data and pos are correctly reversed.
	list = list.Add(Argument{Name: "third"})
	r3 := list.Reverse()
	validate(t, r3, []string{"first", "second", "third"})

	// data and pos are correctly re-reversed.
	r4 := r3.Reverse()
	validate(t, r4, []string{"third", "second", "first"})
}

func TestArgumentsFromSlice(t *testing.T) {
	argA := Argument{Name: "a"}
	argB := Argument{Name: "b"}
	argC := Argument{Name: "c"}

	argSlice := []Argument{argA, argB, argC}
	argList := (*Arguments).Add(nil, argA).Add(argB).Add(argC).Reverse()
	afs := ArgumentsFromSlice(argSlice)

	assert.Equal(t, argList, afs)
}

func validate(t *testing.T, list *Arguments, names []string) {
	var actualNames []string
	var actualPoss []int
	var possWanted []int

	for i := list.Len() - 1; i >= 0; i-- {
		possWanted = append(possWanted, i)

		actualNames = append(actualNames, list.Data.Name)
		actualPoss = append(actualPoss, list.pos)

		list = list.next
	}

	assert.Equal(t, names, actualNames)
	assert.Equal(t, possWanted, actualPoss)
}
