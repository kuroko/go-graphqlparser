package unicode

import "fmt"

func main() {
	r1 := '4'
	r2 := 'e'
	r3 := '1'
	r4 := '6'

	res := ucptor(hexRuneToInt(r1), hexRuneToInt(r2), hexRuneToInt(r3), hexRuneToInt(r4))

	fmt.Println(res)
	fmt.Println(string(res))
}

// TODO(seeruk): Here: https://github.com/graphql/graphql-js/blob/master/src/language/lexer.js#L689
func ucptor(a, b, c, d int) rune {
	return rune(a<<12 | b<<8 | c<<4 | d<<0)
}

func hexRuneToInt(r rune) int {
	isHexNum := r >= '0' && r <= '9'
	isHexChar := r >= 'a' && r <= 'f'

	if !(isHexNum || isHexChar) {
		return -1
	}

	if isHexChar {
		return int(r - 87)
	} else if isHexNum {
		return int(r - 48)
	}

	return -1
}