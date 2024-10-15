package dec64



import (
	"fmt"
	"math"
	// "math/bits"
	// "strconv"
	"testing"
)



// --------------------------------------------------------------------
// code and example in same package, for now
// --------------------------------------------------------------------

func Example() {

	fmt.Println("scaled input")

	fmt.Println("FromScaledInt(0, 1) ->", FromScaledInt(0, 1))
	fmt.Println("FromScaledInt(0, 125) ->", FromScaledInt(0, 125))
	fmt.Println("FromScaledInt(1, 125) ->", FromScaledInt(1, 125))
	fmt.Println("FromScaledInt(4, 125) ->", FromScaledInt(4, 125))

	// subtract
	aa := FromScaledInt(16, 10000000000001)
	bb := FromScaledInt(0, 5)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	aa = FromScaledInt(16, 1)
	bb = FromScaledInt(0, 5)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(0, 51)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(0, 55)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(0, 5499)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(0, 499)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(0, 8)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(0, 1)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(-1, 5)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(-1, 499)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(-1, 5001)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))

	fmt.Println("scientific input")
	sci := FromScientific(124567, 2)
	fmt.Println("FromScientific(124567, 2)", sci)
	fmt.Println("FromScientific(124567, 0)", FromScientific(124567, 0))
	fmt.Println("FromScientific(12545, -10)", FromScientific(12545, -10))
	fmt.Println("FromScientific(1, 0)", FromScientific(1, 0))
	fmt.Println("FromScientific(-12545, 6)", FromScientific(-12545, 6))
	fmt.Println("FromScientific(0, 0)", FromScientific(0, 0))
	fmt.Println("FromScientific(-1, 0)", FromScientific(-1, 0))
	fmt.Println("FromScientific(-1, 1)", FromScientific(-1, 1))

	// loops over floats to convert to Dec
	string_convert()

	// make a dec

	fmt.Println("constants")
	fmt.Println("W", W)
	fmt.Println("K", K)
	fmt.Println("P", P)
	fmt.Println("T", T)
	fmt.Println("EXP_MAX", EXP_MAX)
	fmt.Println("EMIN", EXP_MIN)
	fmt.Println("BIAS", EXP_BIAS)
	// EXP_SHIFT_T1
	fmt.Println("EXP_SHIFT_T1", EXP_SHIFT_T1)

	fmt.Println("EXP_SHIFT_T2", EXP_SHIFT_T2)

	fmt.Printf("MIN POSITIVE %x\n", MIN_POSITIVE)
	fmt.Printf("MAX POSITIVE %x\n", MAX_POSITIVE)

	fmt.Println("MIN POSITIVE", Dec64{MIN_POSITIVE})
	fmt.Println("MAX POSITIVE", Dec64{MAX_POSITIVE})

	fmt.Println("range : ")

	vv := []Dec64{
		One(1),
		FromParts(true, 0, 0),
		FromParts(true, 0, 1),
		FromInt(9999),
		FromInt(math.MinInt64),
		Inf(-1),
		NaN(),
		FromInt(-9999999999999999),
		FromInt(-99999999999999990),
		FromInt(-99999999999999995),
		FromExpUint(0, 125),
		FromExpUint(-1, 125),
		FromParts(true, -2, 125),
		FromExpUint(-3, 125),
		FromExpUint(-4, 125),

		FromExpUint(-1, 1),
		FromExpUint(-2, 1),
		FromExpUint(-18, 12345678901234567),
		FromExpUint(-3, 10002),
		FromExpUint(-1, 1234567890123456),
		FromExpUint(-15, 1234567890123456),
		FromParts(false, -1, 9988827272800001),
		FromExpUint(-2, 9988827272800001),

		// large number
		FromExpUint(10, 123456789012),
		FromExpUint(0, 1234567890123456),
		FromParts(true, 0, 12345678901234567),
		FromExpUint(25, 12),
	}
	for _, v1 := range vv {
		fmt.Println("decoded ", v1)

	}

	// oprations

	v2 := FromInt(25)
	fmt.Println("value ", v2)
	fmt.Println("abs value ", v2.Abs())
	fmt.Println("is neg ", v2.IsNegative())
	fmt.Println("is pos  ", v2.IsPositive())

	fmt.Println("neg value: ")
	v2 = v2.Neg()
	fmt.Println("value ", v2)
	fmt.Println("is neg  ", v2.IsNegative())
	fmt.Println("is pos  ", v2.IsPositive())
	fmt.Println(" abs value ", v2.Abs())

	fmt.Println("neg value again: ")
	v2 = v2.Neg()
	fmt.Println("value ", v2)

	fmt.Println("Mul by 0.01:  MulPow10(-2)")
	v2 = v2.MulPow10(-2)
	fmt.Println("value ", v2)
	v2 = v2.MulPow10(4)
	fmt.Println("value ", v2)
	// this will overflow exponent range
	v2 = v2.MulPow10(400)
	fmt.Println("value ", v2)

	v2 = FromInt(-25)

	fmt.Println("Mul by 0.01:  MulPow10(-2)")
	v2 = v2.MulPow10(-2)
	fmt.Println("value ", v2)
	v2 = v2.MulPow10(4)
	fmt.Println("value ", v2)
	v2 = v2.MulPow10(400)
	fmt.Println("value ", v2)

	// macx value should eb a constant / parameter
	m := FromExpUint(384-15, 9_999_999_999_999_999)
	fmt.Println("value ", m)
	m = m.MulPow10(1)
	fmt.Println("value mul 10 ", m)

	fmt.Println("value max ", MaxValue())
	fmt.Println("value min positive ", MinPosValue())

	// rounding on input - 17 digits??
	// this rounds to infinity
	m = FromExpUint(384-16, 9_999_999_999_999_999_5)
	fmt.Println("value ", m)
	// this rounds to max value
	m = FromExpUint(384-16, 9_999_999_999_999_999_4)
	fmt.Println("value ", m)

	// strings
	// convert_from_string()

	// show ilog10 -needs more work

	m = FromExpUint(-2, 123)
	fmt.Println("Ilog ", m, m.ILog10())
	m = FromUint(123)
	fmt.Println("Ilog ", m, m.ILog10())
	m = FromUint(100)
	fmt.Println("Ilog ", m, m.ILog10())

}


func TestFrom(t *testing.T) {
/*
	minf := Inf(-1)
	inf := Inf(1)
	nan := NaN()
	zero := Zero()
	*/
	
	var a, b Dec64

	a = FromScientific(1, -383)
	b = MinPosValue()

	if !a.Equal(b) {
		t.Error("Creation error", a, "!=", b)
	}

	a = FromScientific(9_999_999_999_999_999, 384)
	b = MaxValue()

	if !a.Equal(b) {
		t.Error("Creation error", a, "!=", b)
	}


}

// ----------------------------------------------------------------------------
// float stuff
// ----------------------------------------------------------------------------

var td = []float64{1, 0, 1234, 1.25,
	math.Inf(0), math.Inf(-1), 1 / math.Inf(-1), math.Log(-1),
	0.1, 0.025, 1.75,
	12345678901234567, 9876543210987654,
	0.001234, 0.99999,
	234.27, 12.673,
	-1.25e-10,
	math.MaxFloat64, math.SmallestNonzeroFloat64,
}

func string_convert() {

	fmt.Println()
	fmt.Println("converting floats")

	for _, f := range td {

		// fmt.Println("float ", f)

		d := FromFloat(f)

		fmt.Println("FromFloat:", f, d)

	}

}
