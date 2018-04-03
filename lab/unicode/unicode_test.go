package unicode

import (
	"testing"
)

func BenchmarkUnquote(b *testing.B) {
	for i := 0; i < b.N; i++ {
		r1 := '4'
		r2 := 'e'
		r3 := '1'
		r4 := '6'

		res := ucptor(hexRuneToInt(r1), hexRuneToInt(r2), hexRuneToInt(r3), hexRuneToInt(r4))
		_ = res
	}
}
