package solver

import (
	"fmt"
)


type Index = uint32
type FF = int64

var SolverPrime FF


// Returns (gcd, x, y) such that a * x + b * y = gcd
func GCD(a, b FF) (FF, FF, FF) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, xp, yp := GCD(b, a % b)
	return gcd, yp, (xp - yp * (a / b))
}


// Returns a^(-1) modulo p
func InverseMod(a, p FF) FF {
	gcd, inv, _ := GCD(a, p)
	if gcd == 1 {
		return inv % p
	} else if gcd == -1 {
		return -inv % p
	}
	panic(fmt.Sprintf("Error taking inverse of %v modulo %v", a, p))
}
