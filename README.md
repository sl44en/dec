# dec
Implementation of floating point decimal 64

WORK IN PROGRESS

Implements floating point decimal 64 as defined by IEEE 754-2008, using binary encoding of the
coefficient. Values are encoded in a 64 bit unsigned integer. 

The implementations is mostly complying to the standard. However, all coefficient values are 
kept normalised in range. Subnormal values are not represented and will be rounded to zero.

Encodes floating point values in range 1.000... e -383 to 9.999... e 384 with 16 digits precision 
in the full range.

The implementation handles +/- zero, +/- infinity and nan values.

Rounding is 

Implements 
- operations add, subtract, multiply and divide, 
- functions for rounding to integer: round, roundeven, trancate, ceil, floor
- type conversion from integer and string
- conversion from and to float64 throuhg conversion to string

Operations are implemented with no use of dynamic allcation or arbitrary math,
except for conversion to and from strings. Conversion to and from native floating point 
is inherently difficult, as it should convert to shortest decimal representation. This
type of conversion is done via conversion to string using Go standard library format and parse.

Currently no math functions are implemented.



