package dec64

import (
	"fmt"
	"math"
	"math/bits"
	"testing"
)

func TestRangeSep(t *testing.T) {

	/*
	   -- test these definitions
	   const (
	   	rangeSepHi uint64 = 542101086242
	   	rangeSepLo uint64 = 13875954555633532928
	   )
	*/

	rHi, rLo := bits.Mul64(E_DIGITS, E_DIGITS_1)
	fmt.Println("calculate range separator: hi, lo", rHi, rLo)

	if rHi != rangeSepHi {
		t.Error("Range separator wrong")
	}
	if rLo != rangeSepLo {
		t.Error("Range separator wrong")
	}

}

// have to stick to 8 digit and lower for this simple div mul test to work
var muldivdata = []Dec64{

	FromInt(5000005),
	FromInt(1999998),
	FromInt(50000005),
	FromInt(19999998),

	FromInt(500005),
	FromInt(199998),
	FromInt(50005),
	FromInt(19998),
	One(1),
	FromParts(true, 0, 1),
	FromInt(9999),
	FromInt(99999999),
	FromExpUint(0, 125),
	FromExpUint(-1, 125),
	FromParts(true, -2, 125),
	FromExpUint(-3, 125),
	FromExpUint(-4, 125),

	FromExpUint(-1, 1),
	FromExpUint(-2, 1),
	FromExpUint(-3, 10002),

	FromExpUint(25, 1221),
}

// a * b -> c,
// then c /a === b and  c / b === a
// test data with no overflows and rounding
func TestMulDiv(t *testing.T) {

	for _, a := range muldivdata {
		for _, b := range muldivdata {

			c := a.Mul(b)
			c1 := b.Mul(a)
			if c != c1 {
				t.Error("mul not ombyttable", a, b, c, c1)
			}

			// prn(a, "*", b, "->", c)
			d1 := c.Div(a) // this should be b
			d2 := c.Div(b) // this should be a

			// we dont yet have equals !
			if d1 != b {
				t.Error("mul div mismatch", d1, b)
			}

			if d2 != a {
				t.Error("mul div mismatch", d2, a)
			}
		}
	}
}

// this test above, how lon

func BenchmarkMulDivSample(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for _, a := range muldivdata {
			for _, b := range muldivdata {

				c := a.Mul(b)
				c1 := b.Mul(a)
				if c != c1 {
					fmt.Println("no use")
				}

				// prn(a, "*", b, "->", c)
				d1 := c.Div(a) // this should be b
				d2 := c.Div(b) // this should be a

				// we dont yet have equals !
				if d1 != b {
					fmt.Println("no use")
				}

				if d2 != a {
					fmt.Println("no use")
				}
			}
		}
	}
}

func TestMulSpecifics(t *testing.T) {

	// mul result in low end, close to high
	var a, b, c, d Dec64

	a = FromInt(50005)
	b = FromInt(19998)
	d = FromInt(999999990)

	c = a.Mul(b)
	if !c.Equal(d) {
		t.Error("mul error a*b is ", a, b, "->", c, "expect", d)
	}

	a = FromInt(50000005)
	b = FromInt(19999998)
	d = FromInt(999999999999990)

	c = a.Mul(b)
	if !c.Equal(d) {
		t.Error("mul error a*b is ", a, b, "->", c, "expect", d)
	}

	// mul result in low end, coefficient overrun when rounding

	a = FromInt(5000_0000_0005)
	b = FromInt(1999_9999_9998)
	d = FromExpInt(8, 1000_0000_0000_0000)

	c = a.Mul(b)
	if !c.Equal(d) {
		t.Error("mul error a*b is ", a, b, "->", c, "expect", d)
	}

	// mul result in high end, still witin range, no rounding
	a = FromInt(9999_999)
	b = FromInt(8888_888)
	d = FromInt(8888_8871_1111_12)

	c = a.Mul(b)
	if !c.Equal(d) {
		t.Error("mul error a*b is ", a, b, "->", c, "expect", d)
	}

	a = FromInt(9999_9999)
	b = FromInt(8888_8888)
	d = FromInt(8888_8887_1111_1112)

	c = a.Mul(b)
	if !c.Equal(d) {
		t.Error("mul error a*b is ", a, b, "->", c, "expect", d)
	}

	a = FromInt(9999_9999_99)
	b = FromInt(8888_88)
	d = FromInt(8888_8799_9911_1112)

	c = a.Mul(b)
	if !c.Equal(d) {
		t.Error("mul error a*b is ", a, b, "->", c, "expect", d)
	}

	a = FromInt(9999_9999_9)
	b = FromInt(8888_8888_8)
	// result is mot verified!
	d = FromExpInt(2, 8888_8888_7111_1111)

	c = a.Mul(b)
	if !c.Equal(d) {
		t.Error("mul error a*b is ", a, b, "->", c, "expect", d)
	}

}

func TestMulSample(t *testing.T) {

	i := FromInt(-25)
	j := FromInt(200)
	fmt.Println("i, j", i, j)

	x := i.Mul(j)
	fmt.Println("mul i * j", x)

	y := i.Mul(Zero())
	fmt.Println("mul by zero", y)

	// array with all types
	vals := []Dec64{Zero(), FromExpInt(-2, -25), Inf(-1), NaN()}

	fmt.Println("values")

	for _, d := range vals {
		fmt.Print(d, ", ")
	}
	fmt.Println()

	fmt.Println("mul matrix")

	for _, d := range vals {
		for _, e := range vals {
			fmt.Print(d.Mul(e), ", ")

		}
		fmt.Println()
	}

	i = FromInt(2)
	j = FromInt(3)
	fmt.Println("i, j", i, j)

	fmt.Println("div i, j", i.Div(j))
	fmt.Println("div j, i", j.Div(i))

	fmt.Println("values")

	for _, d := range vals {
		fmt.Print(d, ", ")
	}
	fmt.Println()

	fmt.Println("div matrix")

	for _, d := range vals {
		for _, e := range vals {
			fmt.Print(d.Div(e), ", ")

		}
		fmt.Println()
	}

	f1 := -1.0
	f0 := 0.0
	fnan := math.Log(f1)
	finf := f1 / f0

	// array with all types
	fvals := []float64{f0, f1, finf, fnan}

	fmt.Println("what do floats do")

	for _, fd := range fvals {
		fmt.Print(fd, ", ")
	}
	fmt.Println()

	fmt.Println("div matrix")

	for _, fd := range fvals {
		for _, fe := range fvals {
			fmt.Print(fd/fe, ", ")

		}
		fmt.Println()
	}

	prn := fmt.Println

	prn("how is it rounding")

	a := FromInt(E_DIGITS_1 - 1)
	b := FromInt(2)
	prn(a, "/", b, "->", a.Div(b))

	a = FromInt(E_DIGITS_1 - 3)
	prn(a, "/", b, "->", a.Div(b))

	a = FromInt(E_DIGITS - 1)
	prn(a, "/", b, "->", a.Div(b))
	a = FromInt(E_DIGITS - 3)
	prn(a, "/", b, "->", a.Div(b))

	a = FromInt(E_DIGITS_1*5 + 25)
	b = FromInt(5)
	prn(a, "/", b, "->", a.Div(b))

	a = FromInt(E_DIGITS_1*5 + 26)
	prn(a, "/", b, "->", a.Div(b))

	a = FromInt(E_DIGITS_1*5 + 27)
	prn(a, "/", b, "->", a.Div(b))

	a = FromInt(9999999999999999)
	b = FromInt(50000005)
	prn(a, "/", b, "->", a.Div(b))
	b = FromInt(19999998)
	prn(a, "/", b, "->", a.Div(b))
	a = FromInt(999999999999999)
	b = FromInt(50000005)
	prn(a, "/", b, "->", a.Div(b))
	b = FromInt(19999998)
	prn(a, "/", b, "->", a.Div(b))

	prn("how is it rounding")
	a = FromInt(50000005)
	b = FromInt(19999998)
	prn(a, "*", b, "->", a.Mul(b))
	a = FromInt(500000005)
	b = FromInt(199999998)
	prn(a, "*", b, "->", a.Mul(b))
	a = FromInt(5000000005)
	b = FromInt(1999999998)
	prn(a, "*", b, "->", a.Mul(b))
	a = FromInt(50000000005)
	b = FromInt(19999999998)
	prn(a, "*", b, "->", a.Mul(b))
	a = FromInt(500000000005)
	b = FromInt(199999999998)
	prn(a, "*", b, "->", a.Mul(b))
	a = FromInt(5000000000005)
	b = FromInt(1999999999998)
	prn(a, "*", b, "->", a.Mul(b))
	a = FromInt(50000000000005)
	b = FromInt(19999999999998)
	prn(a, "*", b, "->", a.Mul(b))
	a = FromInt(500000000000005)
	b = FromInt(199999999999998)
	prn(a, "*", b, "->", a.Mul(b))
	a = FromInt(5000000000000005)
	b = FromInt(1999999999999998)
	prn(a, "*", b, "->", a.Mul(b))
	c := a.Mul(b)
	d1 := c
	d2 := c

	prn("Thes should be in range")

	a = FromInt(50005)
	b = FromInt(19998)
	c = a.Mul(b)
	prn(a, "*", b, "->", c)
	d1 = c.Div(a) // this should be b
	d2 = c.Div(b) // this should be a

	// we dont yet have equals !
	if d1 != b {
		t.Error("mul div mismatch a*b is ", a, b, "->", c)

		t.Error("mul div mismatch", a, b, "_>", d1, b)
	}

	if d2 != a {
		t.Error("mul div mismatch", a, b, "->", d2, a)
	}

	a = FromInt(50000005)
	b = FromInt(19999998)
	c = a.Mul(b)
	prn(a, "*", b, "->", c)
	d1 = c.Div(a) // this should be b
	d2 = c.Div(b) // this should be a

	// we dont yet have equals !
	if d1 != b {
		t.Error("mul div mismatch", d1, b)
	}

	if d2 != a {
		t.Error("mul div mismatch", d1, b)
	}

}
