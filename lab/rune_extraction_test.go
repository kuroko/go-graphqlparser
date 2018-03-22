package lab

import (
	"testing"
	"unicode/utf8"
)

const input = "世界 This is a string. Hello, 世界"

var inputLen = len(input)

func BenchmarkForRange(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pos := 0

		// Read until the end of the string, one rune at a time.
		for {
			if pos+1 > inputLen {
				break
			}

			var r rune
			for _, r = range input[pos:] {
				break
			}

			// Move along the length of this rune, so we start reading the
			// next whole rune, not the next byte.
			pos += utf8.RuneLen(r)

			_ = r
		}
	}
}

func BenchmarkUtf8DecodeRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		pos := 0

		inputBytes := []byte(input)

		// Read until the end of the string, one rune at a time.
		for {
			if pos+1 > inputLen {
				break
			}

			r, w := utf8.DecodeRune(inputBytes[pos:])

			// Move along the length of this rune, so we start reading the
			// next whole rune, not the next byte.
			pos += w

			_ = r
		}
	}
}
