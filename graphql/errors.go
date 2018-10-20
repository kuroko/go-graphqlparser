package graphql

import (
	"bytes"
	"strconv"
)

// Error fulfils the requirements of a GraphQL error type.
type Error struct {
	Message   string
	Locations *Locations
	Path      *PathNodes
}

// NewError returns a new Error with the given message.
func NewError(message string) Error {
	return Error{
		Message: message,
	}
}

// MarshalJSON returns this Error as some JSON bytes. If the Error is invalid, an error will be
// returned.
func (e *Error) MarshalJSON() ([]byte, error) {
	buf := &bytes.Buffer{}
	buf.WriteString(`{`)
	buf.WriteString(`"message":"`)
	buf.WriteString(e.Message)
	buf.WriteString(`"`)

	if e.Locations.Len() > 0 {
		buf.WriteString(`,"locations":[`)
		e.Locations.ForEach(func(location Location, i int) {
			if i > 0 {
				buf.WriteString(",")
			}

			buf.WriteString(`{`)
			buf.WriteString(`"line":`)
			buf.WriteString(strconv.FormatUint(uint64(location.Line), 10))
			buf.WriteString(`,"column":`)
			buf.WriteString(strconv.FormatUint(uint64(location.Column), 10))
			buf.WriteString(`}`)
		})
		buf.WriteString(`]`)
	}

	if e.Path.Len() > 0 {
		buf.WriteString(`,"path":[`)
		e.Path.ForEach(func(pathNode PathNode, i int) {
			if i > 0 {
				buf.WriteString(",")
			}

			switch pathNode.Kind {
			case PathNodeKindString:
				buf.WriteString(`"`)
				buf.WriteString(pathNode.String)
				buf.WriteString(`"`)
			case PathNodeKindInt:
				buf.WriteString(strconv.Itoa(pathNode.Int))
			}
		})
		buf.WriteString(`]`)
	}

	buf.WriteString(`}`)

	return buf.Bytes(), nil
}

// Location provides information about an error as the position in a potentially multi-line string.
type Location struct {
	Line   uint
	Column uint
}

// PathNodeKind values.
const (
	PathNodeKindString PathNodeKind = iota
	PathNodeKindInt
)

// PathNodeKind an enum type that defines the type of data stored in a PathNode.
type PathNodeKind uint8

// PathNode is an individual part of the path.
type PathNode struct {
	Kind   PathNodeKind
	String string
	Int    int
}

// NewStringPathNode returns a new PathNode with the given string as it's value.
func NewStringPathNode(s string) PathNode {
	return PathNode{
		Kind:   PathNodeKindString,
		String: s,
	}
}

// NewIntPathNode returns a new Pathnode with the given int as it's value.
func NewIntPathNode(i int) PathNode {
	return PathNode{
		Kind: PathNodeKindInt,
		Int:  i,
	}
}
