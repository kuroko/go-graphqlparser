package main

import (
	"fmt"
	"unsafe"
)

func main() {
	r1 := '4'
	r2 := 'e'
	r3 := '1'
	r4 := '6'

	res := ucptor(hexRuneToInt(r1), hexRuneToInt(r2), hexRuneToInt(r3), hexRuneToInt(r4))

	fmt.Println(res)
	fmt.Println(string(res))

	fmt.Println(rtob(rune(0x4e16)))
	fmt.Println()
}

func rtob(r rune) string {
	//bs := make([]byte, utf8.RuneLen(r))
	//_ = utf8.EncodeRune(bs, r)
	//
	// return bs

	//buf := bytes.Buffer{}
	//buf.WriteRune(r)
	//
	//return buf.Bytes()

	return string(r)
}

// btos takes the given bytes, and turns them into a string.
// Q: naming btos or bbtos? :D
// TODO(seeruk): Is this actually portable then?
func btos(bs []byte) string {
	return *(*string)(unsafe.Pointer(&bs))
}

// TODO(seeruk): Here: https://github.com/graphql/graphql-js/blob/master/src/language/lexer.js#L689
func ucptor(a, b, c, d int) rune {
	return rune(a<<12 | b<<8 | c<<4 | d<<0)
}

func hexRuneToInt(r rune) int {
	isHexNum := r >= '0' && r <= '9'
	isHexLCChar := r >= 'a' && r <= 'f'
	isHexUCChar := r >= 'A' && r <= 'F'

	if !(isHexNum || isHexLCChar || isHexUCChar) {
		return -1
	}

	if isHexLCChar {
		return int(r - 87)
	} else if isHexUCChar {
		return int(r - 55)
	} else if isHexNum {
		return int(r - 48)
	}

	return -1
}
