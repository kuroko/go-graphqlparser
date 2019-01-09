package server

import (
	"bytes"
	"errors"
	"strconv"
)

// ResponseValueKinds, based on the JSON specification, tailored to suite Go.
// See: https://tools.ietf.org/html/rfc8259#section-3
const (
	ResponseValueKindInt ResponseValueKind = iota
	ResponseValueKindFloat
	ResponseValueKindString
	ResponseValueKindBoolean
	ResponseValueKindNull
	ResponseValueKindArray
	ResponseValueKindObject
)

// ResponseValueKind represents the possible kinds of ResponseValue. A ResponseValue uses it to
// identify which field it should use for it's value.
type ResponseValueKind int

// ResponseValue is a pseudo union type that represents part of a response body.
type ResponseValue struct {
	Kind         ResponseValueKind
	IntValue     int
	FloatValue   float64
	StringValue  string
	BooleanValue bool
	ArrayValue   []ResponseValue
	ObjectValue  []ResponseObjectField
}

// ResponseObjectField represents a single field that forms part of a response object value.
type ResponseObjectField struct {
	Name  string
	Value ResponseValue
}

// MarshalGraphQL is a GraphQL Marshaller that uses case statements and
// verbosity to efficiently marshal the ResponseValueKinds, coupled with
// pseudo union types this method provides significant performance gains
// over using interfaces.
func (v ResponseValue) MarshalGraphQL(buf *bytes.Buffer) error {
	switch v.Kind {
	case ResponseValueKindInt:
		buf.WriteString(strconv.Itoa(v.IntValue))
	case ResponseValueKindFloat:
		buf.WriteString(strconv.FormatFloat(v.FloatValue, 'f', 5, 64))
	case ResponseValueKindString:
		buf.WriteString(`"`)
		buf.WriteString(v.StringValue)
		buf.WriteString(`"`)
	case ResponseValueKindBoolean:
		if v.BooleanValue {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case ResponseValueKindNull:
		buf.WriteString("null")
	case ResponseValueKindArray:
		buf.WriteString("[")

		for i, av := range v.ArrayValue {
			if err := av.MarshalGraphQL(buf); err != nil {
				return err
			}

			// TODO: Optimise.
			if i < len(v.ArrayValue)-1 {
				buf.WriteString(",") // 44 = ,
			}
		}

		buf.WriteString("]")
	case ResponseValueKindObject:
		buf.WriteString("{")

		for i, ob := range v.ObjectValue {
			// Write the name of the field...
			buf.WriteString(`"`)
			buf.WriteString(ob.Name)
			buf.WriteString(`"`)
			buf.WriteString(":")

			if err := ob.Value.MarshalGraphQL(buf); err != nil {
				return err
			}

			if i < len(v.ObjectValue)-1 {
				buf.WriteString(",")
			}
		}

		buf.WriteString("}") // 125 = }
	default:
		return errors.New("TODO")
	}

	return nil
}

// Response is a server response.
type Response struct {
	Data   ResponseValue
	Errors []error // TODO: This is not the right type.
}

// MarshalGraphQL marshals the server response.
func (r Response) MarshalGraphQL() ([]byte, error) {
	buf := &bytes.Buffer{}

	// TODO: Handle errors portion of response.

	buf.WriteString("{")
	buf.WriteString(`"data":`)

	// Data must either be an object, or null.
	// TODO: Is null even valid? I guess you can't have an empty query?
	if r.Data.Kind != ResponseValueKindNull && r.Data.Kind != ResponseValueKindObject {
		return nil, errors.New("TODO")
	}

	if err := r.Data.MarshalGraphQL(buf); err != nil {
		return nil, err
	}

	buf.WriteString("}")

	return buf.Bytes(), nil
}
