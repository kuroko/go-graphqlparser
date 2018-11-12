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
func TestArguments_Insert(t *testing.T) {}
func TestArguments_Join(t *testing.T) {
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

	validate(t, list, []string{"four", "three", "two", "one", "zero"}, []int{4, 3, 2, 1, 0})
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

	// Test one element can be reversed.
	list = list.Add(Argument{Name: "first"})
	one := list.Reverse()
	assert.Equal(t, "first", list.Data.Name)
	assert.Equal(t, "first", one.Data.Name)

	// reset
	one.Reverse()

	// Test n elements can be reversed.
	list = list.Add(Argument{Name: "second"})
	n := list.Reverse()
	assert.Equal(t, "second", list.Data.Name)
	assert.Equal(t, "first", n.Data.Name)

	// reset
	n.Reverse()

	// Test can be reversed multiple times and not mutate.
	r2 := list.Reverse().Reverse()
	assert.Equal(t, list, r2)

	// Test data and pos are correctly reversed.
	list = list.Add(Argument{Name: "third"})
	r3 := list.Reverse()
	validate(t, r3, []string{"first", "second", "third"}, []int{2, 1, 0})

	// Test data and pos are correctly re-reversed.
	r4 := r3.Reverse()
	validate(t, r4, []string{"third", "second", "first"}, []int{2, 1, 0})
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

func validate(t *testing.T, list *Arguments, names []string, poss []int) {
	var actualNames []string
	var actualPoss []int

	for list != nil {
		actualNames = append(actualNames, list.Data.Name)
		actualPoss = append(actualPoss, list.pos)
		list = list.next
	}

	assert.Equal(t, names, actualNames)
	assert.Equal(t, poss, actualPoss)
}
