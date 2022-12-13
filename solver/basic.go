package solver

import (
	"fmt"
)


type Index = uint32
type FF = int64

var SolverPrime FF


func Mod(x FF) FF {
	return x % SolverPrime
}

// Returns (gcd, x, y) such that a * x + b * y = gcd
func GCD(a, b FF) (FF, FF, FF) {
	if b == 0 {
		return a, 1, 0
	}
	gcd, xp, yp := GCD(b, a % b)
	return gcd, yp, (xp - yp * (a / b))
}


// Returns a^(-1) modulo SolverPrime
func InverseMod(x FF) FF {
	gcd, inv, _ := GCD(x, SolverPrime)
	if gcd == 1 {
		return Mod(inv)
	} else if gcd == -1 {
		return -Mod(inv)
	}
	panic(fmt.Sprintf(
		"Error taking inverse of %v modulo %v", x, SolverPrime))
}
