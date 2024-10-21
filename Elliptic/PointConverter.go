package Elliptic

import "math/big"

/*
CoordAffine
Affine Coordinates are the standard(normal) X, Y Coordinates of a point on the Ellipse
*/
type CoordAffine struct {
    AX *big.Int
    AY *big.Int
}

/*
CoordExtended
Extended coordinates represent x y as X Y Z T satisfying the following equations:

  x=X/Z
  y=Y/Z
  x*y=T/Z
There are two variants of formulas,
Variant 1: Assumes Parameter A of the curve is -1,
Variant 2: Makes no Assumption
*/
type CoordExtended struct {
    EX *big.Int
    EY *big.Int
    EZ *big.Int
    ET *big.Int
}

/*
CoordInverted
Inverted coordinates represent x y as X Y Z satisfying the following equations:

  x=Z/X
  y=Z/Y
*/
type CoordInverted struct {
    IX *big.Int
    IY *big.Int
    IZ *big.Int
}

/*
CoordProjective
Projective coordinates [more information] represent x y as X Y Z satisfying the following equations:

  x=X/Z
  y=Y/Z
*/
type CoordProjective struct {
    PX *big.Int
    PY *big.Int
    PZ *big.Int
}

//Basic Modulus Operations

// AddModulus
// Addition Modulo prime
func AddModulus(prime, a, b *big.Int) *big.Int {
    var result = new(big.Int)
    return result.Add(a, b).Mod(result, prime)
}

// SubModulus
// Subtraction Modulo prime
func SubModulus(prime, a, b *big.Int) *big.Int {
    var result = new(big.Int)
    return result.Sub(a, b).Mod(result, prime)
}

// MulModulus
// Multiplication Modulo prime
func MulModulus(prime, a, b *big.Int) *big.Int {
    var result = new(big.Int)
    return result.Mul(a, b).Mod(result, prime)
}

// QuoModulus
// Division Modulo prime
func QuoModulus(prime, a, b *big.Int) *big.Int {
    var mmi = new(big.Int)
    mmi.ModInverse(b, prime)
    return MulModulus(prime, a, mmi)
}

//Coordinate Conversion Methods

func (e *Ellipse) Affine2Extended(InputP CoordAffine) (OutputP CoordExtended) {
    OutputP.EX = InputP.AX
    OutputP.EY = InputP.AY
    OutputP.EZ = One
    OutputP.ET = e.MulModP(InputP.AX, InputP.AY)
    return OutputP
}

func (e *Ellipse) Extended2Affine(InputP CoordExtended) (OutputP CoordAffine) {
    OutputP.AX = e.QuoModP(InputP.EX, InputP.EZ)
    OutputP.AY = e.QuoModP(InputP.EY, InputP.EZ)
    return OutputP
}
