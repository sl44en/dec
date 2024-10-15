package dec64

import (
// "fmt"
// "math"
// "math/bits"
)

// -----------------------------------------------------------------------------
// Rounding
// -----------------------------------------------------------------------------

// Rounds x to the nearest integral value, with halfway cases
// rounding away from zero.
// This is the usual rounding method used in commercial settings
func (x Dec64) Round() Dec64 {

	t, s, e, c := decode(x)
	if t == decNormal && e < 0 {
		// normal numbers can be divided into
		// e >= 0: values are already integer, no rounding needed
		// e < 0 && e > -DIGITS: implied decimal point between digits in coefficient
		//						Determine rounding through division with 10^(-e)
		// e <= -DIGITS: value is lower than 1, rounding up to 1 possible, depending on rounding rule

		if e > -DIGITS {
			cq, cr, sf := divPow10(c, -e)
			c = sf * roundAway(cq, cr, sf)
			// c may overflow coefficient range
			if c == E_DIGITS {
				// c = E_DIGITS/10, e++
				// It cannot overflow exponent - as e < 0
				return encodeFinal(s, e+1, E_DIGITS_1)
			}
			return encodeFinal(s, e, c)
		}

		if e == -DIGITS {
			// This is special case: numbers like 0.xxxxx
			// result is 0 or 1, only
			if 0 == roundAway(0, c, E_DIGITS) {
				return Dec64{s}
			}
			return Dec64{s | ONE_PATTERN}
		}

		// when e < -DIGITS
		// value 0.0xxx and lower
		// round to zero - true for tieAway and Even, but not for Ceil
		return Dec64{s}

	}
	return x
}

// Rounds x to the nearest integral value, with halfway cases
// rounding to even.
func (x Dec64) RoundEven() Dec64 {

	t, s, e, c := decode(x)
	if t == decNormal && e < 0 {
		// normal numbers can be divided into
		// e >= 0: values are already integer, no rounding needed
		// e < 0 && e > -DIGITS: implied decimal point between digits in coefficient
		//						Determine rounding through division with 10^(-e)
		// e <= -DIGITS: value is lower than 1, rounding up to 1 possible, depending on rounding rule

		if e > -DIGITS {
			cq, cr, sf := divPow10(c, -e)
			c = sf * roundEven(cq, cr, sf)
			if c == E_DIGITS {
				// c = E_DIGITS/10, e++
				// It cannot overflow exponent - as e < 0
				return encodeFinal(s, e+1, E_DIGITS_1)
			}
			return encodeFinal(s, e, c)
		}

		if e == -DIGITS {
			// This is special case: numbers like 0.xxxxx
			// result is 0 or 1, only
			if 0 == roundEven(0, c, E_DIGITS) {
				return Dec64{s}
			}
			return Dec64{s | ONE_PATTERN}
		}

		// when e < -DIGITS
		// value 0.0xxx and lower
		// round to zero - true for tieAway and Even, but not for Ceil
		return Dec64{s}

	}
	return x
}

func (x Dec64) Trunc() Dec64 {

	t, s, e, c := decode(x)
	if t == decNormal && e < 0 {
		// normal numbers can be divided into
		// e >= 0: values are already integer, no rounding needed
		// e < 0 && e > -DIGITS: implied decimal point between digits in coefficient
		//						Determine rounding through division with 10^(-e)
		// e <= -DIGITS: value is lower than 1, rounding up to 1 possible, depending on rounding rule

		if e > -DIGITS {
			_, cr, _ := divPow10(c, -e)
			return encodeFinal(s, e, c-cr)
		}

		// when e < -DIGITS
		// value 0.xxxx and lower  truncate to zero
		return Dec64{s}
	}
	return x
}

// round away from zero
func (x Dec64) RoundUp() Dec64 {

	t, s, e, c := decode(x)
	if t == decNormal && e < 0 {
		// normal numbers can be divided into
		// e >= 0: values are already integer, no rounding needed
		// e < 0 && e > -DIGITS: implied decimal point between digits in coefficient
		//						Determine rounding through division with 10^(-e)
		// e <= -DIGITS: value is lower than 1, rounding up to 1 possible, depending on rounding rule

		if e > -DIGITS {
			_, cr, sf := divPow10(c, -e)
			if cr == 0 {
				return x // encodeFinal(s, e, c)
			}
			c += sf - cr
			// Can trigger overflow of coefficient
			if c == E_DIGITS {
				// c = E_DIGITS/10, e++
				// It cannot overflow exponent - as e < 0
				return encodeFinal(s, e+1, E_DIGITS_1)
			}
			return encodeFinal(s, e, c)
		}

		// when e < -DIGITS
		// value 0.xxxx and lower
		// because fraction is non zero, it will always round up to one
		return Dec64{s | ONE_PATTERN}
	}
	return x
}

func (x Dec64) Ceil() Dec64 {

	if x.d&SIGN_MASK == 0 {
		// positive
		return x.RoundUp()
	}
	return x.Trunc()
}

func (x Dec64) Floor() Dec64 {

	if x.d&SIGN_MASK == 0 {
		// positive
		return x.Trunc()
	}
	return x.RoundUp()
}


// ----------------------------------------------------------------------------------
// rounding primitives
// ----------------------------------------------------------------------------------

func roundUpCond(c DecBase, r, base DecBase) bool {
	// rounding: tie-even
	// unclear how to this efeficiently and properly
	r2 := r << 1
	if r2 == base {
		return (c & 1) == 1
	}
	return r2 > base
}

// return a carry digit for rounding tie - away
func roundAway(c, r, base DecBase) uint64 {
	// rounding: tie-away
	// r2 := 2 * r. If r2 >= base there should be carry
	// if r2 > base -> r2 - base is positive, 0 sign
	// if r2 == base -> r2 - base is 0, no sign bit
	// if r2 < base -> r2 - base is negative , signbit
	return c + (^((r << 1) - base) >> 63)
}

// returns rounded coefficient, tie to even
func roundEven(c DecBase, r, base DecBase) uint64 {
	// rounding: tie-even
	r2 := r + r
	if base > r2 {
		return c
	}
	if base < r2 {
		return c + 1
	}
	return c + (c & 1)
}

// returns rounded coefficient, tie to even considering a tail
// tail is more digits after the ones used to determine rounding
// if tail is non zero, there will not be an even tie
func roundEvenTail(c DecBase, r, base, tail DecBase) uint64 {

	r2 := r + r
	if base > r2 {
		// round down
		return c
	}
	if base < r2 {
		// round up
		return c + 1
	}
	// if tail != 0 there is not a tie
	if tail == 0 {
		return c + (c & 1)
	}
	return c + 1
}
