package main

import (
	"strconv"
	"testing"
)

//func BenchmarkUnquote(b *testing.B) {
//	for i := 0; i < b.N; i++ {
//		r1 := '4'
//		r2 := 'e'
//		r3 := '1'
//		r4 := '6'
//
//		res := ucptor(hexRuneToInt(r1), hexRuneToInt(r2), hexRuneToInt(r3), hexRuneToInt(r4))
//		_ = res
//	}
//}

func BenchmarkRuneToString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		//b.Log(btos([]byte("foo\\u4e16bar")))

		str, _ := strconv.Unquote(`"foo\u4e16bar"`)
		_ = str

		//rs := []rune{'f', 'o', 'o'}
		//rs = append(rs, rune(0x4e16))
		//rs = append(rs, 'b')
		//rs = append(rs, 'a')
		//rs = append(rs, 'r')
		//
		//s := string(rs)
		//_ = s
	}
}
