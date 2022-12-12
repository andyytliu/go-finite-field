package solver

import (
	"fmt"
)

var (
	testFactors []FF
)

// func InitInvMap(invMap map[FF]FF) {
// 	for i := -(SolverPrime - 1); i < SolverPrime; i++ {
// 		if i == 0 {
// 			continue
// 		}
// 		invMap[i] = -InverseMod(i, SolverPrime)
// 	}
// }

func InitInvMap(invMap map[FF]FF) {
	// Primitive root
	prim := GetPrimitive()
	inv := InverseMod(prim, SolverPrime)

	k, v := prim, -inv
	invMap[k] = v
	invMap[-k] = -v

	for k != 1 {
		k = Mod(k * prim)
		v = Mod(v * inv)
		invMap[k] = v
		invMap[-k] = -v
	}

	if FF(len(invMap)) != 2 * (SolverPrime - 1) {
		panic(fmt.Sprintf(
			"Length of InvMap seems off: %v vs %v",
			len(invMap), 2 * (SolverPrime - 1)))
	}
}

// Finds the smallest primitive root of SolverPrime
func GetPrimitive() FF {
	testFactors = FactorInt(SolverPrime - FF(1))

	for i := FF(2); i < SolverPrime; i++ {
		if isPrimitive(i) {
			return i
		}
	}
	return 0
}

// Tests if x is a primitive root of SolverPrime
// Returns true if it is a primitive root
func isPrimitive(x FF) bool {
	for _, p := range testFactors {
		if PowerMod(x, (SolverPrime - FF(1)) / p) == 1 {
			return false
		}
	}
	return true
}
