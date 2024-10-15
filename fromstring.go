package dec64

import (
	"fmt"
	// "math"
	// "math/bits"
	// "strconv"
)

var prn = fmt.Println

// returns 	sgn, exp, m, ok

// string has to be of form
//	   [+|-]ddd[.dddd][e|E[+|-]ddd]
// returns NaN if string is invalid
// Nan or infinity values can not be parsed

// the input should be as for decimal floating point literals, or for decimal integer literal
// no base prefixes are processed (eg 0b, 0x, oO),
// that is, eg binary, hex and octal integers not supported, neither is hex floating poit literal

// Since floating point literals accept extraneous leading zeros, it will also be allowed
// for decimal literals. However, an integer with a leading zero will still be processed as a decimal.
// This is similar behaviour as for strconv.ParseFloat()

// underscores can appear between succesive digits, such underscore do not change its value

func parseString(s string) (bool, int, uint64, bool) {

	// TODOs
	// handle overflow in mantissa - round
	// overflow in exponent - Infinity??

	var neg, negexp bool
	var intpart, fracpart bool
	var m uint64
	var ix, tenexp, exp, last int
	var d, c byte

	// for status of underscores
	const (
		nodigit = iota
		digit
		underscore
	)

	// optional sign
	if ix < len(s) {
		c = s[ix]
		if c == '-' || c == '+' {
			neg = c == '-'
			ix++
		}
	}

	// integer part
	last = nodigit
	for ; ix < len(s); ix++ {
		c = s[ix]
		d = c - '0'
		if d <= 9 {
			// digit - TODO overfloow handling
			m = m*10 + uint64(d)
			last = digit
		} else if c == '_' {
			// Underscores are allowed fillers, but only between digits
			if last != digit {
				goto FAIL
			}
			last = underscore
		} else {
			// some other character came up
			break
		}
	}
	if last == underscore {
		// underscore must go between digits
		goto FAIL
	}
	intpart = last == digit

	// we could jump to this point..
	if ix < len(s) {
		c = s[ix]
		if c == '.' {
			ix++
			goto DECIMALS
		} else if c == 'e' || c == 'E' {
			// end of coefficient
			if !intpart {
				// there were no digits
				goto FAIL
			}
			ix++
			goto EXPONENT
		} else {
			// no other acceptable
			goto FAIL
		}
	}

	if !intpart {
		// there were no digits
		goto FAIL
	}
	// there was no exponent/decimals
	return neg, 0, m, true

DECIMALS:

	// We have seen a point - collect decimals
	last = nodigit
	for ; ix < len(s); ix++ {
		c = s[ix]
		d = c - '0'
		if d <= 9 {
			// digit - TODO overfloow handling
			m = m*10 + uint64(d)
			tenexp--
			last = digit
		} else if c == '_' {
			// Underscores are allowed fillers, but only between digits
			if last != digit {
				goto FAIL
			}
			last = underscore
		} else {
			// some other character came up
			break
		}
	}
	if last == underscore {
		// underscore not between digits
		goto FAIL
	}
	fracpart = last == digit

	// at least one of int part or frac part must be present
	if !intpart && !fracpart {
		goto FAIL
	}

	if ix < len(s) {
		c = s[ix]
		if c == 'e' || c == 'E' {
			// end of coefficient
			ix++
			goto EXPONENT
		} else {
			// no other acceptable
			goto FAIL
		}
	}

	// there was no exponent
	return neg, tenexp, m, true

EXPONENT:
	// we have had the 'e'
	// exponent has an optional sign, and one or more digits.
	// it seems the exponent always has sign in outbput from FormatFloat e/g
	// [+|-] dddd

	// optional sign
	if ix < len(s) {
		c = s[ix]
		if c == '-' {
			negexp = true
			ix++
		} else if c == '+' {
			ix++
		}
	}

	// one or more exponent digits required
	last = nodigit
	for ; ix < len(s); ix++ {
		c = s[ix]
		d = c - '0'
		if d <= 9 {
			// digit - TODO overfloow handling
			exp = exp*10 + int(d)
			last = digit
		} else if c == '_' {
			// Underscores are allowed fillers, but only between digits
			if last != digit {
				goto FAIL
			}
			last = underscore
		} else {
			// some other character came up
			goto FAIL
		}
	}

	if last != digit {
		// underscore was last, or there were no exponent digits
		goto FAIL
	}

	if negexp {
		exp = -exp
	}
	return neg, tenexp + exp, m, true

FAIL:
	return false, 0, 0, false

}

// returns Dec64 value from parsing string
// string has to be of form
//	   [+|-]ddd[.dddd][e|E[+|-]ddd]
// returns NaN if string is invalid
// Nan or infinity values can not be created this way

func FromString(s string) Dec64 {
	sgn, exp, m, ok := parseString(s)
	if ok {
		return encode(boolToSign(sgn), exp, m)
	}
	return NaN()
}

// panics if string can not be parsed to Dec64
// can be used for eg initialisations
func MustParse(s string) Dec64 {
	sgn, exp, m, ok := parseString(s)
	if !ok {
		panic("String could not be parsed to dec64: " + s)
	}
	return encode(boolToSign(sgn), exp, m)
}
