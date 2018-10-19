package language

import "unsafe"

// btos takes the given bytes, and turns them into a string.
func btos(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}
