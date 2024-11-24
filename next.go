package dec64

import (
// "fmt"
//
//	"math"
//
// "math/bits"
)

// -----------------------------------------------------------------------------
// Arithmetic
// -----------------------------------------------------------------------------

// next floating point value towards positive infinity
// if self.is_nan(), this returns nan;
// if self is NEG_INFINITY, this returns MIN;
// if self is -TINY, this returns -0.0;
// if self is -0.0 or +0.0, this returns TINY;
// if self is MAX or INFINITY, this returns INFINITY;
// otherwise the unique least value greater than self is returned.
func (self Dec64) NextUp() Dec64 {

	t, s, e, c := udecode(self)

	switch t {
	case decNormal:
		if s == 0 {
			c++
			if c == E_DIGITS {
				// this will return inf if overflow normal range
				return expAdjust(0, e+1, E_DIGITS_1)
			} else {
				return uencodeFinal(0, e, c)
			}
		} else {
			if c == E_DIGITS_1 {
				// this will go to zero if underflow normal range
				return expAdjust(s, e-1, E_DIGITS-1)
			} else {
				return uencodeFinal(s, e, c-1)
			}
		} 
	case decZero:
		return Dec64{MIN_POSITIVE}
	case decInf:
		if s == 0 {  // positive
			// Inf
			return self
		} else {
			// -max
			return Dec64{SIGN_MASK | MAX_POSITIVE }
		}

	default:
		return Dec64{NAN_PATTERN}
	}
}



func (self Dec64) NextDown() Dec64 {

// maybe a bit quick... efficiency ...not
	return self.Neg().NextUp().Neg()
}


func expAdjust(sgn, e, c DecBase) Dec64 {

	debug_assert(c >= E_DIGITS_1 && c < E_DIGITS)

	// what is actually max for this?
	if e < EXP_RANGE {
		// in range
		return uencodeFinal(sgn, e, c)
	}
	if e < SIGN_MASK {
		// overflow
		return Dec64{sgn | INF_PATTERN}
	}
	// we have had underflow
	return Dec64{sgn}
}


