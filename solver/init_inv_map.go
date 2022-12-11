package solver

func InitInvMap(invMap map[FF]FF) {
	for i := -(SolverPrime - 1); i < SolverPrime; i++ {
		if i == 0 {
			continue
		}
		invMap[i] = -InverseMod(i, SolverPrime)
	}
}
