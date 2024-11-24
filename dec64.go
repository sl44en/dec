package dec64

import (
//	"fmt"
//	"math"
//
// "math/bits"
// "strconv"
)

// Decimal Floating point classification
const (
	decZero   = 1
	decNormal = 2
	decInf    = 4
	decNan    = 8
)

// ---------------------------------------------------------------------------------
// type definitions: decimal floating point 64 bit
// ---------------------------------------------------------------------------------

// how to parameterize

type DecBase = uint64

// type DecBase = uint32

type Dec64 struct {
	d DecBase
}

// constructors

// return Dec64 with value i
// if i has more than P digits, result will be rounded
func FromInt(i int64) Dec64 {
	return encodeAny(usign(i), 0, uabs(i))
}

// return Dec64 with value u
// if u has more than P digits, result will be rounded
func FromUint(u uint64) Dec64 {
	return encodeAny(0, 0, u)

}

// return Dec64 with value r = 10^exp * i
// if i has more than P digits, result will be rounded
func FromExpInt(exp int, i int64) Dec64 {
	return encodeAny(usign(i), exp, uabs(i))
}

// return Dec64 from int coefficient and a scale.
// The scale determines position of the decimal point, with a scale
// of zero indicating a decimal point in front of the first digit.
// FromScaledInt(0, 1) -> 0.1
// FromScaledInt(0, 125) -> 0.125
// FromScaledInt(1, 125) -> 1.25
// FromScaledInt(4, 125) -> 1250
// if i has more than P digits, result will be rounded
func FromScaledInt(scale int, i int64) Dec64 {
	// TODO this is almost identical to "Scientific" does not add much
	// find betetr approach
	u := uabs(i)

	return encodeAny(usign(i), scale-Len10(u), u)
}

// create Dec vale from input in scientofic format.
// that is r = d.ddd * 10^e
// inpunt io sprovided as inteher mantissa with an assumed decimal poilnt after
// frits digit.
// e.g 1.2345e2 is input as fromScientific(12345, 2) and gives value: 123.45
func FromScientific(x int64, exp int) Dec64 {

	s := uint64(x) & SIGN_MASK
	// uabs...
	y := x >> 63
	u := uint64((x ^ y) - y)

	// ok, this might be play...
	// Len10 is called both here and through encode - not ideal
	return encodeAny(s, 1+exp-Len10(u), u)
}

// return Dec64 with value r = 10^exp * u
// if u has more than P digits, result will be rounded
func FromExpUint(exp int, u uint64) Dec64 {
	// TODO get rid of this, adss nothing
	return encodeAny(0, exp, u)

}

// return encoded Dec64 with value: r = (-1)^s * 10^exp * u
// s is understood to represent 0 for false, 1 for true,
// to provide positive sign for false, negative sign for true
// if u has more than P digits, result will be rounded
// examples:
//
//	FromParts(true, -2, 125) -> -1.25
//	FromParts(false, 2, 1) -> 100
func FromParts(s bool, exp int, u uint64) Dec64 {
	return encodeAny(boolToSign(s), exp, u)
}

// encode into values of Dec64
// returns zero value, normal value or inf value depedning on input
// full range of values accepted for each type
func encodeAny(sgn uint64, e int, c DecBase) Dec64 {

	// make sure only first bit - may not be needed given the way sgn is provided
	sgn &= SIGN_MASK

	if c == 0 {
		// signed zero
		return Dec64{sgn}
	}

	// if e has extreme value, operations could cause overflow,
	// therefore e is not changed and adjustments is kept in f
	f := 0

	if c < E_DIGITS_1 {
		// scale up to target normalised coefficient range
		// f should be negative, to match use later
		f = Len10(c) - DIGITS
		c = MulPow10(c, -f)
	} else {
		if c > E_DIGITS {
			// c is larger than the normalised coefficient range.
			f = Len10(c) - DIGITS
			// divide coefficient down and round
			c = roundEven(divPow10(c, f))
		}
		// rounding can in turn overflow coefficient,
		if c == E_DIGITS {
			c = E_DIGITS_1
			f++
		}
	}

	// exp too low leads to underflow -> ZERO
	// to avoid overflow from silly input values, e is not upadted
	// with adjustemnt from normalise until after check
	if e < EXP_MIN_NORM-f {
		return Dec64{sgn}
	}

	// exp too high leads to overflow -> Infinity
	if e > EXP_MAX_NORM-f {
		return Dec64{sgn | INF_PATTERN}
	}

	return encodeFinal(sgn, e+f, c)
}

// Special values ------------------------------------------------------

// Inf returns positive infinity if sign >= 0, negative infinity if sign < 0.
func Inf(sign int) Dec64 {
	return Dec64{boolToSign(sign < 0) | INF_PATTERN}

}
func NegInf() Dec64 {
	return Dec64{SIGN_MASK | INF_PATTERN}
}

func PosInf() Dec64 {
	return Dec64{INF_PATTERN}
}

// NaN returns an IEEE 754 “not-a-number” value.
func NaN() Dec64 {
	return Dec64{NAN_PATTERN}
}

// One returns dec64 with value 1 if sign >= 0, -1 if sign < 0.
func One(sign int) Dec64 {
	return Dec64{boolToSign(sign < 0) | ONE_PATTERN}
}

// Zero returns dec64 with value 0
func Zero() Dec64 {
	return Dec64{0}
}

// Max returns max normal value
// 9.999...e384
func MaxValue() Dec64 {
	return Dec64{MAX_POSITIVE}
}

// MinPos returns minimum positive value: 1.000...e-383
// Subnormal numbers are not supported
func MinPosValue() Dec64 {
	return Dec64{MIN_POSITIVE}
}

// classifications

// true if x is NaN
func (self Dec64) IsNaN() bool {
	// special bits for nan x11111xx...
	return self.d&NAN_MASK == NAN_PATTERN
}

// True if x is positive or negative infinity
func (self Dec64) IsInf() bool {
	// special bits for inf x11110xx...
	return self.d&NAN_MASK == INF_PATTERN
}

func (self Dec64) IsZero() bool {
	// There a many encodings for zero
	// a type 1 number with coefficient 0 is a zero
	// in most cases the below would be enough
	// but -0 and multiple allowed exponent values come in the way
	//     if self.d == 0 { return true }
	v0 := self.d
	return (v0&COEFF_MASK_T1) == 0 && (v0&MSB_MASK) != MSB_MASK
}

// Signbit returns true if sign bit is set
func (self Dec64) Signbit() bool {
	return 0 != (SIGN_MASK & self.d)
}

// True if x is a finite negative value or negative infinity
// corresponds to x < 0
func (self Dec64) IsNegative() bool {

	if self.d&NAN_MASK == NAN_PATTERN {
		return false
	}
	return 0 != (SIGN_MASK&self.d) && !self.IsZero()
}

// True if x is a finite positive value or positive infinity
// corresponds to x > 0
func (self Dec64) IsPositive() bool {

	if self.d&NAN_MASK == NAN_PATTERN {
		return false
	}
	return 0 == (SIGN_MASK&self.d) && !self.IsZero()
}

// ---
// math - is this math
// ------

// Ilog10 returns the 10 exponent of x as an integer.
// for zero, Inf, negative

// Ilogb(±Inf) = MaxInt32
// Ilogb(0) = MinInt32
// Ilogb(NaN) = MaxInt32
// someho all negative are invlaid for log, but they have an exponent
func (self Dec64) ILog10() int {

	_, _, e, _ := decode(self)

	// how to handle illegal input?=

	/*
	   if t == decZero {
	   return Inf(-1)
	   }

	   if s != 0 || t == decNaN {
	   return NaN()
	   }

	   if t == decInf {
	   return Inf(1)
	   }

	*/

	return e + DIGITS_1

}

// --------------------------------------------------------------------------
// operations
// --------------------------------------------------------------------------

// return negated value
func (self Dec64) Neg() Dec64 {
	// what about negative Nan?? - its still Nan of course - negative nan is allowed
	// this allows for negative nan, and negative zero ofcourse
	// maybe this should check normal values
	return Dec64{self.d ^ SIGN_MASK}
}

// return absolute value
func (self Dec64) Abs() Dec64 {
	return Dec64{self.d &^ SIGN_MASK}
}

// multiply by a power of 10
// returns value r = d * 10^p
func (self Dec64) MulPow10(p int) Dec64 {
	// this is also called Scale()

	t, s, e, v := decode(self)

	if t == decNormal {
		// e can come out of range, while coefficient remains in range
		return encodeNormalised(s, e+p, v)
	}
	return self

}

// Datatype for decimal floating point encoded according to
// IEEE Std 754-2008

// encoding follows standard for Decimal interchange floating-point format,
// using the binary encoding for the significand

// restrictions
// 
// Coefficient is always normalised
// Subnormal cooefficients are not handled and not produced as results
//     subnormal values are rounded to zero
// the encoding has space for coefficient vales above 10^P - 1. These values are not handled
//     This implementation will not produce such over-normal nunbers as result of any operation
// providing bitpatterns that encode such coefficients will have undefined behaviour
// Coefficient values provided through the api interface, fromParts, fromExpInt, FromInt will be
// scaled and rounded to supported range, with eg a higher/lower exponent value, or 0 or Inf

/* **************************
for K                  (32, 64)
SIGN_FIELD  = 1        (1, 1)
COMB_FIELD  = W + 5    (11, 13)
TRAIL_FIELD = T        (20, 50)
*/

const K = 64 // (32, 64)
// combination field bits (5)
const W = K/16 + 4     // Exponent continuatin bits (6, 8)
const T = 15*K/16 - 10 // Coefficient continuation bits (20, 50)
const P = 9*K/32 - 2   // Decimal digits (7, 16)

const DIGITS = P // 16 for Dec64
const BITS = 64  // 64 for Dec64

const EXP_RANGE = 3 * (1 << W)   // 3 * 2 ^ W Exponent range (786)
const EXP_MAX = EXP_RANGE / 2    // Largest value is 9.99... * 10 ^ Emax
const EXP_MIN = 1 - EXP_MAX      // Smallest normalized value is 1.00... * ^ 10Emin
const EXP_TINY = 2 - P - EXP_MAX // Smallest non-zero value is 1 * 10 ^ Etiny
const EXP_BIAS = EXP_MAX + P - 2

const E_TEN = 10

// const E_DIGITS E_TEN.pow(DIGITS) // not possible
const E_DIGITS = 1_0000_0000_0000_0000
const E_DIGITS_1 = E_DIGITS / 10 // lowest value for normalised coefficient
const E_DIGITS_H = E_DIGITS / 2
const E_DIGITS_M = E_DIGITS - 1 // highest value for normalised coefficient

// mask values to extract components freom reprsentation
// Maslks are defined as u32 to match the internal type
// combination flied msb
const SIGN_SHIFT = K - 1
const SIGN_MASK = 1 << SIGN_SHIFT

// Bits 2 and 3 indicate a type 2 representation or special value
const MSB_MASK = 0b11 << (BITS - 3)

// Bits 2, 3, 4 and 5 indicate a special value, nan or inf
const SPECIAL_MASK = 0b1111 << (BITS - 5)

// nanmasks teh special msbs and the indicator for type of value
const NAN_MASK = 0b11111 << (BITS - 6)

// these are duplicates of above
const NAN_PATTERN = NAN_MASK
const INF_PATTERN = 0b11110 << (BITS - 6)

// continuation field - 6 bits for Dec
// this seems wrong!!!
const CONT_FIELD_BITS = (K/32)*2 + 4 // this is W

// what was here before??
// 10 bits for Dec64
const EXP_BITS = W + 2

// w paramter, 6 for D32, there are two more exp bits from special purpose field
// const EXP_BITS: u32 = K * 2 + 4;
const EXP_SHIFT_T1 = COEFF_TRAILING_BITS + 3
const EXP_SHIFT_T2 = COEFF_TRAILING_BITS + 1
const EXP_MASK = (1 << EXP_BITS) - 1

// const MAX_EXP_BIAS: u32 = 0b10111111

// constants related to exponents.
// they are defined as signed i32 to facilitate the conversion
// between biased and unbiased exponents.

const DIGITS_1 = DIGITS - 1

//const EXP_RANGE = 48 * (2 << (2 * K))

// const EXP_MAX = 3 * 1 << (K*2 + 3) // 96
// const EXP_MIN = 1 - EXP_MAX        // -95
// const EXP_BIAS = EXP_MAX + DIGITS - 2

const EXP_MAX_NORM = EXP_MAX - DIGITS_1
const EXP_MIN_NORM = EXP_MIN - DIGITS_1

// Coefficient
const COEFF_TRAILING_BITS = T
const COEFFICIENT_BITS_T1 = COEFF_TRAILING_BITS + 3
const COEFF_MASK_T1 = (1 << COEFFICIENT_BITS_T1) - 1

const COEFFICIENT_BITS_T2 = COEFF_TRAILING_BITS + 1
const COEFF_MASK_T2 = (1 << COEFFICIENT_BITS_T2) - 1
const COEFF_LEAD_T2 = 0b100 << COEFFICIENT_BITS_T2

// duplicate of COEFF_MASK_T1
const MAX_T1_COEFFICIENT = (1 << COEFFICIENT_BITS_T1) - 1
const TOP_T1_COEFFICIENT = 1 << COEFFICIENT_BITS_T1

// encode values

// need values for max number, minpositive, in addition to zero andsimilar

// Values
const ONE_PATTERN = ((EXP_BIAS - DIGITS_1) << EXP_SHIFT_T1) | E_DIGITS_1
const MIN_POSITIVE = E_DIGITS_1
const MAX_POSITIVE = MSB_MASK | 0x2FF<<EXP_SHIFT_T2 | E_DIGITS_M&COEFF_MASK_T2

// --------------------------------------------------------------------
// Construct final Decimal representations
// --------------------------------------------------------------------
//
// many helper functions for different conditions
//
// check coefficient range and adapt
// check exponent range and retun zero / inf if under / overflow
//
// Build final bit reprrsentation of value

// helper function for normalise and construct when results
// are from arithmetic ops and results are almost normalised
//
// construct Dec value from sign (0/1), exponent -95 <= e <= 96,
// for u32 significand value
//
// Construct will return ZERO, MAX or MIN, incase of unddeerflow / overflow,
// or a finite value. It will not return NAN.
//
// construct handles only fully normalised significand values.
// exponent can have any value, it will lead to Inf or zero values.

// general case, where all coefficient values are accepted.
// coefficient will be rounded if exceeding max coefficient
// or normalsed if sub normal
// all exponent values are accepted - will return zero or
// inf in case of under- / overflow

func encode(sgn uint64, e int, c DecBase) Dec64 {

	if c == 0 {
		// allow for signed zero?
		return Dec64{0}
	}

	// we only accept up to E_DIGITS
	// NO that doid not hold - FromFloat calling here
	// debug_assert(c <= E_DIGITS)

	// e cannot be assumed to be in valid range, because frim mul and others

	var f int

	// fmt.Println(sgn, e, c)

	if c < E_DIGITS_1 {
		// scale up to target normalised coefficient range
		f = DIGITS - Len10(c)
		c = MulPow10(c, f)
		// adjust e
		e -= f

	} else {

		if c > E_DIGITS {
			f = Len10(c) - DIGITS
			c = roundEven(divPow10(c, f))
			e += f
		}

		// rounding can in turn overflow the coefficient, or it was provided overflowed
		if c == E_DIGITS {
			c = E_DIGITS_1
			e++
		}
	}

	// e values can be out of range from start!

	if e < EXP_MIN_NORM {
		// ignoring sign
		return Dec64{sgn}
	}

	if e > EXP_MAX_NORM {
		return Dec64{sgn | INF_PATTERN}
	}
	return encodeFinal(sgn, e, c)
}

// all values are in range
func encodeFinal(sgn uint64, e int, c DecBase) Dec64 {

	// only first bit must be set in sign
	debug_assert(sgn<<1 == 0)

	debug_assert(E_DIGITS_1 <= c && c < E_DIGITS)

	// biased exponent
	eb := DecBase(e + EXP_BIAS)

	// what is actually max for this?
	debug_assert(eb < EXP_RANGE)

	if c < TOP_T1_COEFFICIENT {
		// type 1: smaller coefficient, below 2^(T+3)
		return Dec64{sgn | eb<<EXP_SHIFT_T1 | c}
	} else {
		// type 2: for large coefficient
		return Dec64{sgn | MSB_MASK | eb<<EXP_SHIFT_T2 | c&COEFF_MASK_T2}
	}
}

// all values are in range
func uencodeFinal(sgn, e, c DecBase) Dec64 {

	// only first bit must be set in sign
	debug_assert(sgn<<1 == 0)
	debug_assert(E_DIGITS_1 <= c && c < E_DIGITS)

	// what is actually max for this?
	debug_assert(e < EXP_RANGE)

	if c < TOP_T1_COEFFICIENT {
		// type 1: smaller coefficient, below 2^(T+3)
		return Dec64{sgn | e<<EXP_SHIFT_T1 | c}
	} else {
		// type 2: for large coefficient
		return Dec64{sgn | MSB_MASK | e<<EXP_SHIFT_T2 | c&COEFF_MASK_T2}
	}
}

// encode parts, w guaranteed normalised coefficient (saving time to normalise)
// but provide remainder and remainder base for rounding
// and not extreme values for exponent
// useful following multiplications and divisions

/* whre is this used

func encodeRound(sgn uint64, e int, c DecBase, r, base DecBase) Dec64 {

	debug_assert(c >= E_DIGITS_1 && c < E_DIGITS)
	// rounding: even-tie-away
	if (r + r) >= base {
		c++
		if c == E_DIGITS {
			c = E_DIGITS_1
			e++
		}
	}
	return encodeNormalised(sgn, e, c)
}
*/

// this assumes coefficient is in normalised in range, but accepts overflow from rounding
// it checks exponent values and return zero/Inf if needed
func encodeNormalised(sgn uint64, e int, c DecBase) Dec64 {

	debug_assert(c >= E_DIGITS_1 && c <= E_DIGITS)

	// adjust coefficient and exponent of overflow from rounding
	if c == E_DIGITS {
		c = E_DIGITS_1
		e++
	}

	// exp too low leads to underflow -> ZERO
	if e < EXP_MIN_NORM {
		return Dec64{sgn}
	}

	// exp too high leads to overflow -> Infinity
	if e > EXP_MAX_NORM {
		return Dec64{sgn | INF_PATTERN}
	}
	return encodeFinal(sgn, e, c)
}

func boolToSign(s bool) DecBase {
	if s {
		return SIGN_MASK
	}
	return 0
}

// SignedBAse -> DecBase

// external function
// --------------------------------------------------------------------------

func c_signed_inf(s DecBase) Dec64 {
	return Dec64{s | INF_PATTERN}
}

// helper function, unpacks the bit pattern of IEEE 754 decimal binary
// integer (BID) into its parts, and determines special values of the
// representation.
// This is a full classification, according to ieee standard,
// even though the full set of reprsentations is not implemneted.
// Deconstruct does much of he same as classify, but it alse provides the
// contents of the components of the reprrsentation sign, eponent, signicand
// (and allows for representations that are not implented).
func decode(d Dec64) (int, DecBase, int, DecBase) {
	var v0, sgn, ufe, fc DecBase
	v0 = d.d
	sgn = v0 & SIGN_MASK
	// check if MSBs of combination field not set. That should be the most
	// common, and identifies normal values with a significand lower than
	// 2^23 / 2^53
	if v0&MSB_MASK != MSB_MASK {
		// type 1 representation
		fc = v0 & COEFF_MASK_T1
		/*
			// decode is too fine - should not check for non states, like subnormal
			if fc < E_DIGITS_1 {
				// subnormal or zero
				// this implementation treat subnormal as zero
				// ! the IsZero function does not check for subnormal
				// subnormal numbers cannot appear as a result of operations in this impl
				//    return decZero, sgn, 0, 0
				// else, it would need to be handled
				if fc == 0 {
					// zero value
					return decZero, sgn, 0, 0
				} else {
					// note that other parts of this impl does not handle subnormal
					// subnormal value - it could be scaled up...
					ufe = (v0 >> EXP_SHIFT_T1) & EXP_MASK
					return decNormal, sgn, int(ufe) - EXP_BIAS, fc
				}
			}
		*/
		if fc == 0 {
			// zero value
			return decZero, sgn, 0, 0
		}

		ufe = (v0 >> EXP_SHIFT_T1) & EXP_MASK
		return decNormal, sgn, int(ufe) - EXP_BIAS, fc
	}
	if v0&SPECIAL_MASK != SPECIAL_MASK {
		// type 2 representation for finite value
		fc = COEFF_LEAD_T2 | v0&COEFF_MASK_T2
		ufe = (v0 >> EXP_SHIFT_T2) & EXP_MASK
		// should we check for over-normal
		/*
			// this implementation will not produce over-normal values
			if fc >= E_DIGITS {
				// overnormal - according to standard this must be regarded as zero
				return decZero, sgn, 0, 0
			}
		*/
		return decNormal, sgn, int(ufe) - EXP_BIAS, fc
	}
	// special values, bit patterns x11111 and x11110
	if v0&NAN_MASK == NAN_PATTERN {
		return decNan, sgn, 0, 0
	} else {
		return decInf, sgn, 0, 0
	}

}
func udecode(d Dec64) (int, DecBase, DecBase, DecBase) {
	var v0, sgn, ufe, fc DecBase
	v0 = d.d
	sgn = v0 & SIGN_MASK
	// check if MSBs of combination field not set. That should be the most
	// common, and identifies normal values with a significand lower than
	// 2^23 / 2^53
	if v0&MSB_MASK != MSB_MASK {
		// type 1 representation
		fc = v0 & COEFF_MASK_T1
		/*
			// decode is too fine - should not check for non states, like subnormal
			if fc < E_DIGITS_1 {
				// subnormal or zero
				// this implementation treat subnormal as zero
				// ! the IsZero function does not check for subnormal
				// subnormal numbers cannot appear as a result of operations in this impl
				//    return decZero, sgn, 0, 0
				// else, it would need to be handled
				if fc == 0 {
					// zero value
					return decZero, sgn, 0, 0
				} else {
					// note that other parts of this impl does not handle subnormal
					// subnormal value - it could be scaled up...
					ufe = (v0 >> EXP_SHIFT_T1) & EXP_MASK
					return decNormal, sgn, ufe, fc
				}
			}
		*/
		if fc == 0 {
			// zero value
			return decZero, sgn, 0, 0
		}

		ufe = (v0 >> EXP_SHIFT_T1) & EXP_MASK
		return decNormal, sgn, ufe, fc
	}
	if v0&SPECIAL_MASK != SPECIAL_MASK {
		// type 2 representation for finite value
		fc = COEFF_LEAD_T2 | v0&COEFF_MASK_T2
		ufe = (v0 >> EXP_SHIFT_T2) & EXP_MASK
		// should we check for over-normal
		/*
			// this implementation will not produce over-normal values
			if fc >= E_DIGITS {
				// overnormal - according to standard this must be regarded as zero
				return decZero, sgn, 0, 0
			}
		*/
		return decNormal, sgn, ufe, fc
	}
	// special values, bit patterns x11111 and x11110
	if v0&NAN_MASK == NAN_PATTERN {
		return decNan, sgn, 0, 0
	} else {
		return decInf, sgn, 0, 0
	}

}
