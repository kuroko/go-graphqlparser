package types

import (
	"bytes"
	"strconv"

	"github.com/bucketd/go-graphqlparser/ast"
)

// Error fulfils the requirements of a GraphQL error type.
type Error struct {
	Message   string
	Locations *ast.Locations
	Path      *ast.PathNodes
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
		e.Locations.ForEach(func(location ast.Location, i int) {
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
		e.Path.ForEach(func(pathNode ast.PathNode, i int) {
			if i > 0 {
				buf.WriteString(",")
			}

			switch pathNode.Kind {
			case ast.PathNodeKindString:
				buf.WriteString(`"`)
				buf.WriteString(pathNode.String)
				buf.WriteString(`"`)
			case ast.PathNodeKindInt:
				buf.WriteString(strconv.Itoa(pathNode.Int))
			}
		})
		buf.WriteString(`]`)
	}

	buf.WriteString(`}`)

	return buf.Bytes(), nil
}
