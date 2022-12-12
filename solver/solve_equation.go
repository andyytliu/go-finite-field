package solver

import (
	"errors"
	"log"
	"sync"
)

// Plug existing solutions into equations
func UpdateEquations(
	equations []*map[Index]FF,
	solutions map[Index]map[Index]FF) {

	semLen := len(equations)  // Number of signals to get back from goroutines
	sem := make(chan struct{}, semLen)

	for id, eq := range equations {
		go func(id int, eq *map[Index]FF) {
			newEq := make(map[Index]FF)
			for idx, coef := range (*eq) {
				if coef == 0 {
					continue
				}
				if sol, ok := solutions[idx]; ok {
					for k, v := range sol {
						if v == 0 {
							continue
						}
						var newCoef FF = Mod(coef * v)
						if oldCoef, ok2 := newEq[k]; ok2 {
							newCoef = Mod(newCoef + oldCoef) 
							if newCoef == 0 {
								// Delete zero terms
								delete(newEq, k)
							} else {
								newEq[k] = newCoef
							}
						} else {
							newEq[k] = newCoef
						}
					}
				} else {
					if oldCoef, ok2 := newEq[idx]; ok2 {
						coef = Mod(coef + oldCoef) 
						if coef == 0 {
							// Delete zero terms
							delete(newEq, idx)
						} else {
							newEq[idx] = coef
						}
					} else {
						newEq[idx] = coef
					}
				}
			}
			*equations[id] = newEq
			sem <- struct{}{}
		}(id, eq)
	}
	//////
	// Receive signals from goroutines
	for i := 0; i < semLen; i++ {
		<- sem
	}
	//////
}


func SolveEquation(
	equation map[Index]FF,  // single eq to solve
	equations []*map[Index]FF,  // block equations to update
	invMap map[FF]FF,
	solutions map[Index]map[Index]FF,
	solTransIndex map[Index]map[Index]bool) error {
	

	newEq := equation

	// If the equation is zero then nothing needs to be done
	if (len(newEq) == 0) {
		return nil
	}

	// Find index to solve
	var idxToSolve Index = 0
	for k, _ := range newEq {
		if k == 0 {
			continue
		} else {
			idxToSolve = k
			break
		}
	}


	if idxToSolve == 0 {
		log.Println(">>>>>>>>>>> error: no solution for equation")
		logEquation(newEq)
		return errors.New("No solution")
	}


	// set to (-1/c) where c is the coef of idxToSolve
	var invCoef FF  = invMap[newEq[idxToSolve]]

	// Add new solution
	newSol := make(map[Index]FF)
	for k, v := range newEq {
		if k == idxToSolve {
			continue
		}

		newSol[k] = Mod(invCoef * v)

		// And modify trans indexing accordingly
		if idxs, ok := solTransIndex[k]; ok {
			idxs[idxToSolve] = true
		} else {
			solTransIndex[k] = make(map[Index]bool)
			solTransIndex[k][idxToSolve] = true
		}
	}

	

	semLen := 0  // Number of signals to get back from goroutines
	
	// # of new equations to update
	numEqToUpdate := len(equations)
	// numEqToUpdate := 100
	// if len(equations) < numEqToUpdate {
	// 	numEqToUpdate = len(equations)
	// }
	semLen += numEqToUpdate

	// # of solutions to update
	idxLen := 0
	var idxsToUpdate map[Index]bool
	if idxs, ok := solTransIndex[idxToSolve]; ok {
		idxsToUpdate = idxs
		idxLen = len(idxs)
	}
	semLen += idxLen

	
	// Make channel to receive signals from goroutines
	sem := make(chan struct{}, semLen)

	
	// Update new equations
	for _, eq := range equations[:numEqToUpdate] {

		go func(eq map[Index]FF) {
			if coef, ok := eq[idxToSolve]; ok {
				for k, v := range newSol {
					var newCoef FF = Mod(coef * v)
					if oldCoef, ok := eq[k]; ok {
						newCoef = Mod(newCoef + oldCoef)
						if newCoef == 0 {
							// Delete zero terms
							delete(eq, k)
						} else {
							eq[k] = newCoef
						}
					} else {
						eq[k] = newCoef
					}
				}
			}
			delete(eq, idxToSolve)
			sem <- struct{}{}
		}(*eq)
	}

	// Update existing solutions and modify trans indexing accordingly
	if idxsToUpdate != nil {
		//////
		mutex := sync.Mutex{}
		//////
		
		for i, _ := range idxsToUpdate {

			go func(idx Index) {
				idxCoef := solutions[idx][idxToSolve]
				for k, v := range newSol {
					var newCoef FF = Mod(idxCoef * v)

					if oldCoef, ok := solutions[idx][k]; ok {
						newCoef = Mod(newCoef + oldCoef)
						if newCoef == 0 {
							// Delete zero terms
							delete(solutions[idx], k)
							mutex.Lock()
							delete(solTransIndex[k], idx)
							mutex.Unlock()
						} else {
							solutions[idx][k] = newCoef
						}
					} else {
						solutions[idx][k] = newCoef
						// Update transpose indexing
						mutex.Lock()
						if idxs, ok := solTransIndex[k]; ok {
							idxs[idx] = true
						} else {
							solTransIndex[k] = make(map[Index]bool)
							solTransIndex[k][idx] = true
						}
						mutex.Unlock()
					}
				}
				delete(solutions[idx], idxToSolve)
				sem <- struct{}{}
			} (i)
		}
	}

	//////
	// Receive signals from goroutines
	for i := 0; i < semLen; i++ {
		<- sem
	}
	//////


	delete(solTransIndex, idxToSolve)
	solutions[idxToSolve] = newSol

	return nil
}

func logEquation(equation map[Index]FF) {
	log.Println("===> equation")
	for k,v := range equation {
		log.Printf("(%v,%v) ", k, v)
	}
	log.Println("")
}
