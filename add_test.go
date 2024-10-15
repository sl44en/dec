package dec64

import (
	"fmt"
	"math"
	//	"math/bits"
	"testing"
)

// have to stick to 8 digit and lower for this simple div mul test to work
var add_sub_data = []Dec64{

	FromInt(0),
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

	// this ruins subtaction - too far off in scale
	// FromExpUint(25, 1221),
}

// a * b -> c,
// then c /a === b and  c / b === c
// test data with no overflows and rounding
func TestAddSub(t *testing.T) {

	// this test only works where rounding is not involved.. and maybe even issues with negative zeros

	for _, a := range add_sub_data {
		for _, b := range add_sub_data {

			helpAddSub(t, a, b)
			helpAddSub(t, b, a)
		}
	}
}

func helpAddSub(t *testing.T, a, b Dec64) {

	c := a.Add(b)
	d := c.Sub(a) // this should be b

	// we dont yet have equals !
	if d != b {
		t.Error("addsub", a, "+", b, "=", c, "; -", a, "=", d, "expected", b)
	}
}

func TestCombinations(t *testing.T) {

	zero := 0.0
	inf := 1.0 / zero
	minf := -1.0 / zero
	nan := math.Log(-1.0)
	n1 := -0.25
	n2 := 0.26

	floats := []float64{zero, -zero, n1, n2, inf, minf, nan}
	// showfloats

	prnf := fmt.Printf
	prn := fmt.Println

	prnf("%15s", "f-e")
	for _, f := range floats {
		prnf("%15f", f)
	}
	prn()

	for _, f := range floats {
		prnf("%15f", f)
		for _, e := range floats {
			prnf("%15f", f-e)
		}
		prn()
	}

	prn()

	prnf("%15s", "f+e")
	for _, f := range floats {
		prnf("%15f", f)
	}
	prn()

	for _, f := range floats {
		prnf("%15f", f)
		for _, e := range floats {
			prnf("%15f", f+e)
		}
		prn()
	}

	dvals := []Dec64{}

	for _, f := range floats {
		d := FromFloat(f)
		dvals = append(dvals, d)
	}

	prn()

	prnf("%15s", "Sub()")
	for _, a := range dvals {
		prnf("%15v", a)
	}
	prn()

	for _, a := range dvals {
		prnf("%15v", a)
		for _, b := range dvals {
			prnf("%15v", a.Sub(b))
		}
		prn()
	}

	prn()

	prnf("%15s", "Add()")
	for _, a := range dvals {
		prnf("%15v", a)
	}
	prn()

	for _, a := range dvals {
		prnf("%15v", a)
		for _, b := range dvals {
			prnf("%15v", a.Add(b))
		}
		prn()
	}

}

func TestAddSpecifics2(t *testing.T) {

	var aa, bb Dec64

	// Add diff = 1
	const nine0 = 9000_0000_0000_0000
	const nines = 9999_9999_9999_9999
	aa = FromScaledInt(16, nine0)
	bb = FromScaledInt(15, nines)
	fmt.Println("add", aa, bb, "->", aa.Add(bb))
	aa = FromScaledInt(5, 9999)
	bb = FromScaledInt(4, 9999)
	fmt.Println("add", aa, bb, "->", aa.Add(bb))
	aa = FromScaledInt(5, 90000)
	fmt.Println("add", aa, bb, "->", aa.Add(bb))

	aa = FromScaledInt(5, 99999)
	for i := range 5 + 1 {
		bb = FromScaledInt(i, 99999)
		fmt.Println("add", aa, bb, "->", aa.Add(bb))
	}

	aa = FromScaledInt(5, 34567)
	for i := range 5 + 1 {
		bb = FromScaledInt(i, 34567)
		fmt.Println("add", aa, bb, "->", aa.Add(bb))
	}

}

func TestSubSpecifics2(t *testing.T) {

	var aa, bb Dec64

	// subtract

	// with same scale diff == 0
	aa = FromScaledInt(1, 12005)
	bb = FromScaledInt(1, 12004)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	aa, bb = bb, aa
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	aa = FromScaledInt(1, 13)
	bb = FromScaledInt(1, 12005)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	aa = FromScaledInt(1, 23)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(1, 23)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))

	// for middle sdiff
	aa = FromScaledInt(4, 1012)
	bb = FromScaledInt(2, 1122)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 1244)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 1208)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 1202)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 13)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 1299)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))

	aa = FromScaledInt(16, 1000_0000_0000_0012)
	bb = FromScaledInt(2, 1122)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 1244)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 1208)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 1202)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 13)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(2, 1299)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))

	// for low values add diff = DIGITS

	aa = FromScaledInt(16, 10000000000001)
	bb = FromScaledInt(0, 5)
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

	// for low values add diff = DIGITS+1

	aa = FromScaledInt(16, 10000000000001)
	bb = FromScaledInt(-1, 999)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	aa = FromScaledInt(16, 1000)
	bb = FromScaledInt(-1, 999)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))

	aa = FromScaledInt(16, 1)
	bb = FromScaledInt(-1, 5)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(-1, 499)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))
	bb = FromScaledInt(-1, 5001)
	fmt.Println("sub", aa, bb, "->", aa.Sub(bb))

}

func TestAddSpecifics(t *testing.T) {

	k := MaxValue()
	sum := k.Add(k)
	if !sum.IsInf() {
		t.Error("Max + max should give inf", sum)
	}

	k = FromInt(5000000000000003)
	sum = k.Add(k)
	o := FromExpInt(1, 1000000000000001)
	if !o.Equal(sum) {
		t.Error("sum fail", k, sum)
	}

	k = FromInt(9999999999999998)
	m := FromInt(3)
	n := k.Add(m)
	o = FromExpInt(1, 1000000000000000)

	if !o.Equal(n) {
		t.Error("sum fail", k, m, n, o)
	}

	k = FromInt(9999999999999998)
	m = FromInt(7)
	n = k.Add(m)
	o = FromExpInt(1, 1000000000000000)

	if !o.Equal(n) {
		t.Error("sum fail", k, m, n, o)
	}
	k = FromInt(9999999999999998)
	m = FromInt(8)
	n = k.Add(m)
	o = FromExpInt(0, 10000000000000010)

	if !o.Equal(n) {
		t.Error("sum fail", k, m, n, o)
	}

	k = FromExpInt(-1, 9999999999999999)
	m = FromExpInt(-2, 9999999999999996)
	n = k.Add(m)
	o = FromExpInt(0, 1100000000000000)

	if !o.Equal(n) {
		t.Error("sum fail", k, m, n, "compare", o)
	}

	k = FromExpInt(-1, 9999999999999999)
	m = FromExpInt(-2, 9999999999999993)
	n = k.Add(m)
	o = FromExpInt(0, 1100000000000000)

	if !o.Equal(n) {
		t.Error("sum fail", k, m, n, "compare", o)
	}

	k = FromExpInt(-1, 9999999999999994)
	m = FromExpInt(-2, 9999999999999990)
	n = k.Add(m)
	o = FromExpInt(0, 1099999999999999)

	if !o.Equal(n) {
		t.Error("sum fail", k, m, n, "compare", o)
	}

	k = FromExpInt(-1, 9999999999999996)
	m = FromExpInt(-2, 9999999999999990)
	n = k.Add(m)
	o = FromExpInt(0, 1100000000000000)

	if !o.Equal(n) {
		t.Error("sum fail", k, m, n, "compare", o)
	}

	k = FromExpInt(-1, 9999999999999986)
	m = FromExpInt(-2, 9999999999999990)
	n = k.Add(m)
	o = FromExpInt(0, 1099999999999998)

	if !o.Equal(n) {
		t.Error("sum fail", k, m, n, "compare", o)
	}

	k = FromExpInt(-1, 9999999999999986)
	m = FromExpInt(-2, 9999999999999991)
	n = k.Add(m)
	o = FromExpInt(0, 1099999999999999)

	if !o.Equal(n) {
		t.Error("sum fail", k, m, n, "compare", o)
	}

}

func TestAddSample(t *testing.T) {

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
		t.Error("mul div mismatch", d1, b)
	}

	if d2 != a {
		t.Error("mul div mismatch", d1, b)
	}

}

var testint64gen = []int64{
	9_999_999,
	9_999_000,
	9_000_000,
	9_900_000,
	1_000_000,
	1_120_000,
	1_012_000,
	1_001_200,
	1_000_120,
	1_000_012,
	1_200_000,
	1_280_000,
	1_201_000,
	1234567,
	7654321,
	2, 3, 4, 5, 6, 7, 8, 9,
}

func TestAddGeneral(t *testing.T) {

	// add subtract numbers in varying scale to exercise all difefrences in scalint
	// this does not test rounding
	// no rounding will happn with these coefficients
	// rounding is expecetd to complicate testing...

	for _, d1 := range testint64gen {
		for v1 := range 8 {
			a := FromScaledInt(v1, d1)

			for _, d2 := range testint64gen {
				for v2 := range 8 {
					b := FromScaledInt(v2, d2)

					// this should give alll combination sof scale and coefficient
					// of course there will be many duplications

					c1 := a.Add(b)
					c2 := b.Add(a)
					if !c2.Equal(c1) {
						t.Error("adding order", a, "+", b, "=", c1, " != ", c2)
					}

					f1 := c1.Sub(a) // this should be b
					if !f1.Equal(b) {
						t.Error("addsub", a, "+", b, "=", c1, "; -", a, "=", f1, "expected", b)
					}

					f2 := a.Sub(b)        // b sub a will come because loop all combinations
					f2n := b.Sub(a)       // neg f2
					f2m := b.Neg().Add(a) // f2
					if !f2.Equal(f2n.Neg()) {
						t.Error("a-b != -(b-a)", a, "-", b, "=", f2, "; !=", f2n)
					}
					if !f2m.Equal(f2) {
						t.Error("a-b != -b + a)", a, "-", b, "=", f2, "; !=", f2m)
					}

					f3 := f2.Add(b) // should be a
					f4 := f2.Sub(a) // should be -b
					if !f3.Equal(a) {
						t.Error("addsub", a, "-", b, "=", f2, "; +", b, "=", f3, "expected", a)
					}
					if !f4.Equal(b.Neg()) {
						t.Error("addsub", a, "-", b, "=", f2, "; -", a, "=", f4, "expected", b)
					}

				}
			}
		}
	}
}

// try a benchmark

var sumstat bool

// it seems to make sense with this function

// adding simple values with 0 - 3 in scalin difference
func BenchmarkAddSimple(b *testing.B) {
	// remember for this stest: the tim eis for many executions
	// currently 210000 ns for 8464 additions ~25 ns per add
	const scalediff = 4
	for i := 0; i < b.N; i++ {

		for _, d1 := range testint64gen {
			for v1 := range scalediff {
				a := FromScaledInt(v1, d1)

				for _, d2 := range testint64gen {
					for v2 := range scalediff {
						b := FromScaledInt(v2, d2)

						c := a.Add(b)

						sumstat = sumstat != c.IsZero()

					}
				}
			}
		}
	}
}

// adding simple values with 0 - 3 in scalin difference
func TestAddSimple(t *testing.T) {
	const scalediff = 4
	count := 0
	for _, d1 := range testint64gen {
		for v1 := range scalediff {
			a := FromScaledInt(v1, d1)

			for _, d2 := range testint64gen {
				for v2 := range scalediff {
					b := FromScaledInt(v2, d2)

					c := a.Add(b)
					count++
					sumstat = sumstat != c.IsZero()

				}
			}
		}
	}
	fmt.Println("count from addsimple", count)
}
