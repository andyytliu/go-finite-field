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

// Returns a^n modulo SolverPrime
func PowerMod(a, n FF) FF {
	if n == 1 {
		return a
	}
	q, r := n/2, n%2
	tmp := PowerMod(a, q)
	tmp = Mod(tmp * tmp)
	if r != 0 {
		tmp = Mod(tmp * a)
	}
	return tmp
}

// Returns all prime factors of x
func FactorInt(x FF) (factors []FF) {
	if x % FF(2) == 0 {
		factors = append(factors, FF(2))
		for x % FF(2) == 0 {
			x = x / FF(2)
		}
	}
	for i := FF(3); i*i <= x; i = i+2 {
		if x % i == 0 {
			factors = append(factors, i)
			for x % i == 0 {
				x = x / i
			}
		}
	}
	if x > 1 {
		factors = append(factors, x)
	}
	return
}
