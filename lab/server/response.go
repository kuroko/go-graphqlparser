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

// ResponseValue represents part of a response body.
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

func (v ResponseValue) MarshalGraphQL(buf *bytes.Buffer) error {
	switch v.Kind {
	case ResponseValueKindInt:
		buf.WriteString(strconv.Itoa(v.IntValue))
	case ResponseValueKindFloat:
		buf.WriteString(strconv.FormatFloat(v.FloatValue, 'f', 5, 64))
	case ResponseValueKindString:
		buf.Write([]byte{34})
		buf.WriteString(v.StringValue)
		buf.Write([]byte{34})
	case ResponseValueKindBoolean:
		if v.BooleanValue {
			buf.WriteString("true")
		} else {
			buf.WriteString("false")
		}
	case ResponseValueKindNull:
		buf.WriteString("null")
	case ResponseValueKindArray:
		buf.Write([]byte{91}) // 91 = [

		for i, av := range v.ArrayValue {
			if err := av.MarshalGraphQL(buf); err != nil {
				return err
			}

			// @TODO: Optimise.
			if i < len(v.ArrayValue)-1 {
				buf.Write([]byte{44}) // 44 = ,
			}
		}

		buf.Write([]byte{93}) // 93 = ]
	case ResponseValueKindObject:
		buf.Write([]byte{123}) // 123 = {

		for i, ob := range v.ObjectValue {
			// Write the name of the field...
			buf.Write([]byte{34}) // 34 = "
			buf.WriteString(ob.Name)
			buf.Write([]byte{34}) // 34 = "
			buf.Write([]byte{58}) // 58 = :

			if err := ob.Value.MarshalGraphQL(buf); err != nil {
				return err
			}

			if i < len(v.ObjectValue)-1 {
				buf.Write([]byte{44}) // 44 = ,
			}
		}

		buf.Write([]byte{125}) // 125 = }
	default:
		return errors.New("TODO")
	}

	return nil
}

type Response struct {
	Data   ResponseValue
	Errors []error // @TODO: This is not the right type.
}

func (r Response) MarshalGraphQL() ([]byte, error) {
	buf := &bytes.Buffer{}

	// @TODO: Handle errors portion of response.

	buf.Write([]byte{123}) // 123 = {
	buf.WriteString(`"data":`)

	// Data must either be an object, or null.
	// @TODO: Is null even valid? I guess you can't have an empty query?
	if r.Data.Kind != ResponseValueKindNull && r.Data.Kind != ResponseValueKindObject {
		return nil, errors.New("TODO")
	}

	if err := r.Data.MarshalGraphQL(buf); err != nil {
		return nil, err
	}

	buf.Write([]byte{125}) // 125 = }

	return buf.Bytes(), nil
}
