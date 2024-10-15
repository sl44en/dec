package dec64

import (
	"fmt"
	"testing"
	// "math"
	// "math/bits"
	"strconv"
)

// this function test whether the same strings can be parsed by FromStrinng and standard ParseFloat
// this is not a full test, only that accepted syuntax is teh same
// note that FromString does not acept hex floating point reprsentation, which ParseFlloat does
func TestFromString1(t *testing.T) {

	for _, s := range test_ss {

		d := FromString(s)
		// if this is nan parsing failed

		succ1 := !d.IsNaN()

		fmt.Println("fromstring", s, "->", d)

		// can the same string be parsed by strconv
		_, err := strconv.ParseFloat(s, 64)
		succ2 := err == nil

		if succ1 != succ2 {
			t.Error("s1 , s2 different", succ1, succ2, s)
		}

	}

}

// In reality, we need a list of strings, held up against expected results
// this list doesnt exist at this point

func printthosefloats2() {
	for _, f := range test_floats {

		fmt.Printf("floats %v, %f, %e, %g \n", f, f, f, f)
	}
}

// what are allowwedliterals - I think we should handle teh same

var test_floats = []float64{
	.01,
	17.,
	1_7,
	// 1._9, // no.. '_' must separate successive digits

	1.0e00_0,

	1.0e000_1,
	0003.89, // this is read as a floating point
	003.89,
	03.89,
	3.89,
	-.1,
	+.1,
	0600, // this is read as an octal!!! - i dont want octal numbers in float processing
	00.1,

	// 1e,  //  exponent has no digits
	// 1e+,  // xponent has no digits
	1e+0,
	1e+0000,
	1e+00001,
}

// so, leading zeros are allowed
// and no mandatory first digit, a dot is fine

var test_ss = []string{"1.25", "-1000", "1e5", "1.0e-5", "-1.e5", "23.", "-1e04",
	"1_000_000_000",
	"-", "-0", "-1", "--6",
	"002", "00", "0", "06800", "0.0", "0e1", "01", "1", "0.01", "0.1", // frac part intpart
	".", "0.", ".0", "0",
	"_11", "1_", "0_9", "9_9", "9__9", "-1e0_4", "-1e_04",
	"5_e4", "5_.4", "5._4", "5.4_",
	"1e", "1e+", "1e+0", "1e++1", "1e-1",
	// no good....
	"e3",
	// "0x1p-2",    // this is not supported for Dec
}
