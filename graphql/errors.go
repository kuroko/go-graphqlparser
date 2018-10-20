package graphql

import (
	"errors"
	"strconv"
	"unsafe"
)

// Error ...
type Error struct {
	Message   string     `json:"message"`
	Locations []Location `json:"locations"`
	Path      []PathInfo `json:"path"`
}

type Location struct {
	Line   uint `json:"line"`
	Column uint `json:"column"`
}

const (
	PathInfoKindString PathInfoKind = iota
	PathInfoKindInt
)

type PathInfoKind uint8

type PathInfo struct {
	Kind   PathInfoKind
	String string
	Int    int
}

func (p PathInfo) MarshalJSON() ([]byte, error) {
	switch p.Kind {
	case PathInfoKindString:
		return stob(`"` + p.String + `"`), nil
	case PathInfoKindInt:
		return stob(strconv.Itoa(p.Int)), nil
	}

	return []byte{}, errors.New("graphql: Invalid PathInfo")
}

// stob takes the given string, and turns them into a byte slice.
func stob(s string) []byte {
	return *(*[]byte)(unsafe.Pointer(&s))
}
