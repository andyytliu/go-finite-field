package server

import (
	. "github.com/andyytliu/go-finite-field/solver"
	"log"
)

var (
	BlockSolveSize int = 200

 	equations []*map[Index]FF
	solutions = make(map[Index]map[Index]FF)
	solTransIndex = make(map[Index]map[Index]bool)
	InvMap = make(map[FF]FF)
)

type Handler struct {
	Channel chan struct{}
}

type Reply struct {
	NumEq int
	NumSol int
}

func (handler *Handler) Status(_ int, reply *Reply) error {
	reply.NumEq = len(equations)
	reply.NumSol = len(solutions)
	return nil
}

func (handler *Handler) SetBlock(size int, reply *Reply) error {
	go func() {
		handler.Channel <- struct{}{}
		BlockSolveSize = size
		log.Printf("*********** Set block size to: %v\n", size)
		<- handler.Channel
	}()
	return nil
}

func (handler *Handler) ReadEquations(fileName string, reply *Reply) error {
	go func() {
		handler.Channel <- struct{}{}
		ReadEquations(fileName, &equations)
		log.Printf("*********** Read %v equations\n", len(equations))
		<- handler.Channel
	}()
	return nil
}

func (handler *Handler) ReadSolutions(fileName string, reply *Reply) error {
	go func() {
		handler.Channel <- struct{}{}
		ReadSolutions(fileName, solutions, solTransIndex)
		log.Printf("*********** Read solutions from file: %v\n", fileName)
		<- handler.Channel
	}()
	return nil
}

func (handler *Handler) WriteSolutions(fileName string, reply *Reply) error {
	go func() {
		handler.Channel <- struct{}{}
		WriteSolutions(fileName, solutions)
		log.Printf("*********** Written solutions to file: %v\n", fileName)
		<- handler.Channel
	}()
	return nil
}

func (handler *Handler) SolveEquations(totalNumToSolve int, reply *Reply) error {
	if len(equations) < totalNumToSolve {
		totalNumToSolve  = len(equations)
	}
	
	numToSolve := BlockSolveSize
	log.Printf("*********** Solving %v equations in blocks of %v\n",
		totalNumToSolve, BlockSolveSize)

	var err error

	go func() {
		handler.Channel <- struct{}{}

		for totalNumToSolve > 0 {
			if totalNumToSolve < numToSolve {
				numToSolve = totalNumToSolve
			} else {
				numToSolve = BlockSolveSize			
			}
			totalNumToSolve -= numToSolve

			log.Println("Start back-substitution")
			UpdateEquations(equations[:numToSolve], solutions)
			
			log.Printf("Start solving batch of %v equations\n", numToSolve)
			for i := 0; i < numToSolve; i++ {
				err = SolveEquation(*equations[i], equations[i+1:numToSolve],
					InvMap, solutions, solTransIndex)
				if err != nil {
					break
				}
			}
			if err != nil {
				break
			}
			equations = equations[numToSolve:]
			log.Printf("Done  solving batch of %v equations\n", numToSolve)
		}

		if err != nil {
			log.Println(">>>>>>>>>>> Error in solving equations")
		} else {
			log.Println("*********** Done solving equations")
		}
		<- handler.Channel
	}()

	return nil
}