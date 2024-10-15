# dec
Implementation of floating point decimal 64

WORK IN PROGRESS

Implements floating point decimal 64 as defined by IEEE 754-2008, using binary encoding of the
coefficient. That is, the coefficient is represented as a 16 digit unsigned integer. The whole 
floating point value is encoded in 64 bits.

The implementations is mostly complying to the standard. However, all coefficient values are 
kept normalised in the range 1e15 to 1e16 - 1. Subnormal values are not represented 
and will be rounded to zero.

The implementation encodes floating point values in range 1.000... e -383 to 9.999... e 384 
with 16 digits precision in the full range.

The implementation handles +/- zero, +/- infinity and nan values.

When needed, rounding is performed with tie-to-even in all arithmetic operations 
and type conversions.


Implements 
- operations add, subtract, multiply and divide, 
- functions for rounding to integer: round, roundeven, truncate, ceiling, floor
- type conversion from integer and string
- conversion from and to float64 through conversion to string

Operations are implemented with no use of dynamic allocation or arbitrary precision math,
except for conversion to and from strings. Conversion to and from native floating point 
is inherently difficult, as it should convert to shortest decimal representation. This
type of conversion is done via conversion to string using Go standard library format and parse.

Currently no math functions are implemented.

No dependencies outside Go standard library


