package main

import (
	. "github.com/andyytliu/go-finite-field/solver"
	"flag"
	"log"
	"net"
	"net/rpc"
	"os"
	"runtime"
)

var (
	port string
	logFile string
	prime int
	blockSolveSize int

 	equations []*map[Index]FF
	solutions = make(map[Index]map[Index]FF)
	solTransIndex = make(map[Index]map[Index]bool)
	invMap = make(map[FF]FF)
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
		blockSolveSize = size
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
	
	numToSolve := blockSolveSize
	log.Printf("*********** Solving %v equations in blocks of %v\n",
		totalNumToSolve, blockSolveSize)

	var err error

	go func() {
		handler.Channel <- struct{}{}

		for totalNumToSolve > 0 {
			if totalNumToSolve < numToSolve {
				numToSolve = totalNumToSolve
			} else {
				numToSolve = blockSolveSize			
			}
			totalNumToSolve -= numToSolve

			log.Println("Start back-substitution")
			UpdateEquations(equations[:numToSolve], solutions)
			
			log.Printf("Start solving batch of %v equations\n", numToSolve)
			for i := 0; i < numToSolve; i++ {
				err = SolveEquation(*equations[i], equations[i+1:numToSolve],
					invMap, solutions, solTransIndex)
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


func main() {

	flag.StringVar(&port, "port", "8080", "Port number for the server to listen to")
	flag.StringVar(&logFile, "log", "logs.txt", "Name for the log file")
	flag.IntVar(&prime, "p", 46337, "Prime number to use in modular calculation")
	flag.IntVar(&blockSolveSize, "block", 200, "Block size of equations to solve in parallel")
	flag.Parse()

	SolverPrime = FF(prime)
	InitInvMap(invMap)


	os.Remove(logFile)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY, 0666)
    if err != nil {
        log.Fatal(err)
    }
    log.SetOutput(file)

	log.Println("***********************************")
    log.Println("*********** Start logging")
	log.Printf("*********** > NumCPU: %v\n", runtime.NumCPU())
	log.Printf("*********** > GOMAXPROCS: %v\n", runtime.GOMAXPROCS(0))
	log.Printf("*********** > Port: %v\n", port)
	log.Printf("*********** > Log file: %v\n", logFile)
	log.Printf("*********** > Prime: %v\n", prime)
	log.Printf("*********** > Block size: %v\n", blockSolveSize)
	log.Println("***********************************")
	

	handler := new(Handler)
	handler.Channel = make(chan struct{}, 1) // chan of size 1, works as a mutex
	rpc.Register(handler)

	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Println(">>>>>>>>>>> errorr in server: " + err.Error())
	}
	defer listener.Close()
	
	rpc.Accept(listener)
}
