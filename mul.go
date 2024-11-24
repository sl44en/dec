package dec64

import (
	//"fmt"
	//"math"
	"math/bits"
)

// -----------------------------------------------------------------------------
// Arithmetic: Multiplication and division
// -----------------------------------------------------------------------------
/*
    * There is a complex system for the type of the result depending on
    * input types for calculations
    *
    * The table below list the result types for each combination.
    * Note that this implemenattion does not implement signed zeroes,
    * so zero results will be unsigned. Below table is based on how f64 type
    * interpretes input types and generate results.
    * -------------------------------------------------------------------------------

    X * y     y =            0          2     -56.89        inf       -inf        NaN
   x =          0            0          0         -0        NaN        NaN        NaN
   x =          2            0          4    -113.78        inf       -inf        NaN
   x =     -56.89           -0    -113.78  3236.4721       -inf        inf        NaN
   x =        inf          NaN        inf       -inf        inf       -inf        NaN
   x =       -inf          NaN       -inf        inf       -inf        inf        NaN
   x =        NaN          NaN        NaN        NaN        NaN        NaN        NaN

    x / y     y =            0          2     -56.89        inf       -inf        NaN
   x =          0          NaN          0         -0          0         -0        NaN
   x =          2          inf          1 -0.035155#          0         -0        NaN
   x =     -56.89         -inf    -28.445          1         -0          0        NaN
   x =        inf          inf        inf       -inf        NaN        NaN        NaN
   x =       -inf         -inf       -inf        inf        NaN        NaN        NaN
   x =        NaN          NaN        NaN        NaN        NaN        NaN        NaN

   ----------------------------------------------------------------------------------
*/

// value separating 31 and 32 digit wide product, calculated as:
// rangeSepHi, rangeSepLo = bits.Mul64(E_DIGITS, E_DIGITS_1)
//     hi, lo: 542101086242 13875954555633532928

const (
	rangeSepHi uint64 = 542101086242
	rangeSepLo uint64 = 13875954555633532928
)

// Following mul and divide, we need rounding, after which coefficient will generally be in range
// exponents can be far off and should be checked

// this assumes coefficient is in normalised in range
// it checks exponent values and return zero/Inf if needed
func mulEncode(sgn, e, c DecBase) Dec64 {

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

// Standard bits.Div64 modified for division by constant, utilising compiler optimisation
// for division by constant
func div64_e15(hi, lo uint64) (uint64, uint64) {

	// const parameters for divisor
	const (
		divisor = 1_000_000_000_000_000
		dbitlen = 50 // bitlength of divisor
	)
	// debug_assert(dbitlen == bits.Len64(divisor))
	// debug_assert(hi < divisor)

	// environment constants
	const (
		wbitlen = 64          // word bit length
		hbitlen = wbitlen / 2 // half word bit length
		b       = 1 << hbitlen
		mask    = b - 1
	)

	// parameterisation of divisor constants
	const (
		s   = wbitlen - dbitlen // LeadingZeros64(divisor)
		yn  = divisor << s
		yn1 = yn >> hbitlen
		yn0 = yn & mask
	)

	var q1, q0, un, un10, un1, un0, rhat, c1, c2 uint64

	un = hi<<s | lo>>(wbitlen-s)
	un10 = lo << s
	un1 = un10 >> hbitlen
	un0 = un10 & mask

	q1 = un / yn1
	rhat = un % yn1
	
	// from the blog of ridiculousfish.com (libdivide)
	// we have an alternative approach to quotient correction
	c1 = q1 * yn0
	c2 = rhat*b + un1
	if c1 > c2 {
		q1--
		if c1-c2 > yn {
			q1--
		}
	}

	un = un*b + un1 - q1*yn
	q0 = un / yn1
	rhat = un % yn1

	c1 = q0 * yn0
	c2 = rhat*b + un0
	if c1 > c2 {
		q0--
		if c1-c2 > yn {
			q0--
		}
	}

	return q1*b + q0, (un*b + un0 - q0*yn) >> s
}

func div64_e16(hi, lo uint64) (uint64, uint64) {

	// const parameters for divisor
	const (
		divisor = 10_000_000_000_000_000
		dbitlen = 54 // bitlength of divisor
	)
	// debug_assert(dbitlen == bits.Len64(divisor))
	// debug_assert(hi < divisor)

	// environment constants
	const (
		wbitlen = 64          // word bit length
		hbitlen = wbitlen / 2 // half word bit length
		b       = 1 << hbitlen
		mask    = b - 1
	)

	// parameterisation of divisor constants
	const (
		s   = wbitlen - dbitlen // LeadingZeros64(divisor)
		yn  = divisor << s
		yn1 = yn >> hbitlen
		yn0 = yn & mask
	)

	var q1, q0, un, un10, un1, un0, rhat, c1, c2 uint64

	un = hi<<s | lo>>(wbitlen-s)
	un10 = lo << s
	un1 = un10 >> hbitlen
	un0 = un10 & mask

	q1 = un / yn1
	rhat = un % yn1
	
	// from the blog of ridiculousfish.com (libdivide)
	// we have an alternative approach to quotient correction
	c1 = q1 * yn0
	c2 = rhat*b + un1
	if c1 > c2 {
		q1--
		if c1-c2 > yn {
			q1--
		}
	}

	un = un*b + un1 - q1*yn
	q0 = un / yn1
	rhat = un % yn1

	c1 = q0 * yn0
	c2 = rhat*b + un0
	if c1 > c2 {
		q0--
		if c1-c2 > yn {
			q0--
		}
	}

	return q1*b + q0, (un*b + un0 - q0*yn) >> s
}




// ----------------------------------


// const EXP_BIAS = EXP_MAX + P - 2
const EXP_ZERO_BIAS = EXP_MAX - 2

func (a Dec64) Mul(b Dec64) Dec64 {

	t1, s1, e1, c1 := udecode(a)
	t2, s2, e2, c2 := udecode(b)
	s := s1 ^ s2
	e := e1 + e2 - EXP_ZERO_BIAS
	var c, rem, mhi, mlo uint64

	switch t1 | t2 {

	case decNormal:
		// perfom multiplication; result in double word
		mhi, mlo = bits.Mul64(c1, c2)

		// result is 31 or 32 decimal digits, because coefficients are normalsed to 16 digits
		// compare with 1e31, calculated as a double word, to determine where to round
		// if c1*c2 < 1e31:
		// if (mhi, mlo) < (rangeSepHi, rangeSepLo): 
		if mhi < rangeSepHi || (mhi == rangeSepHi && mlo < rangeSepLo) {
			// low range: m is 31 digits: divide by e15
			// c, rem = bits.Div64(mhi, mlo, E_DIGITS_1)
			c, rem = div64_e15(mhi, mlo)
			debug_assert(c >= E_DIGITS_1 && c < E_DIGITS)
			// round and adjust for coefficient overflow from rounding
			// eg 5005 * 1998 -> 9999990 -> 9999.990 for 4 digit mantissa, 
			// needs rounding and will overflow
			c = roundEven(c, rem, E_DIGITS_1)
			if c == E_DIGITS {
				// c ends in hi range after rounding
				// c = E_DIGITS_1, e++
				return mulEncode(s, e, E_DIGITS_1)
			}
			return mulEncode(s, e-1, c)
		}
		// high range: m is 32 digits: divide by e16
		// c, rem = bits.Div64(mhi, mlo, E_DIGITS)
		c, rem = div64_e16(mhi, mlo)
		debug_assert(c >= E_DIGITS_1 && c < E_DIGITS)
		// round; there can be no coeeficient overflow in the high range
		c = roundEven(c, rem, E_DIGITS)
		return mulEncode(s, e, c)

	case decZero, decZero | decNormal:
		// zero * zero/normal -> zero result
		return Dec64{s}

	case decNormal | decInf, decInf:
		return Dec64{s | INF_PATTERN}

	default:
		return Dec64{NAN_PATTERN}
	}
}

func (a Dec64) Div(b Dec64) Dec64 {

	t1, s1, e1, c1 := udecode(a)
	t2, s2, e2, c2 := udecode(b)
	s := s1 ^ s2
	e := EXP_ZERO_BIAS + e1 - e2
	var mhi, mlo, c, rem uint64

	switch t1 | t2 {

	case decNormal:

		if c1 >= c2 {
			// then c1 / c2  >= 1
			// result in range 1.000 - 9.999 (for a for digit mantissa, rounded)
			// scale up with 10^(P-1) to get division result in range
			mhi, mlo = bits.Mul64(c1, E_DIGITS_1)
			c, rem = bits.Div64(mhi, mlo, c2)
			debug_assert(E_DIGITS_1 <= c && c < E_DIGITS)

			// if remainder is more than half way or equal
			// to next, add carry. distance to next is c2
			// rounding will not overflow coefficient
			c = roundEven(c, rem, c2)
			return mulEncode(s, e+1, c)

		} else {
			// c1 < c2 =>  c1 / c2 < 1
			// result in range 0.1000 - 0.9999 (for a four digit mantissa, rounded)
			// scale up with 10^P to get division result inrange
			mhi, mlo = bits.Mul64(c1, E_DIGITS)
			c, rem = bits.Div64(mhi, mlo, c2)
			debug_assert(E_DIGITS_1 <= c && c < E_DIGITS)

			// if remainder is more than half way or equal
			// to next, add carry. distance to next is c2

			//return encodeOverflow(s, e1-e2-DIGITS, roundUp(c, rem, c2))
			// I dont think there can be overflow - proof outstanding
			c = roundEven(c, rem, c2)
			return mulEncode(s, e, c)
		}

	case decZero | decNormal, decZero | decInf, decNormal | decInf:
		if t1 < t2 {
			// zero result
			return Dec64{s}
		} else {
			// inf result
			return Dec64{s | INF_PATTERN}
		}

	default:
		// 0/0, Inf/Inf and anything involving a nan -> nan
		return Dec64{NAN_PATTERN}
	}
}

