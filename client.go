package main

import (
	"flag"
	"fmt"
	"net/rpc"
)

type Reply struct {
	NumEq int
	NumSol int
}

var (
	port string
	status bool
	equationsFile string
	solutionsFile string
	writeFile string
	numEqToSolve int
	blockSize int
)

func main() {

	flag.StringVar(&port, "port", "8080", "Port number the server listens to")
	flag.BoolVar(&status, "status", false, "Look up status of the solver: numbers of equations and solutions")
	flag.StringVar(&equationsFile, "read_equations", "", "Read in equations to be solved")
	flag.StringVar(&solutionsFile, "read_solutions", "", "Read in existing solutions")
	flag.StringVar(&writeFile, "write_solutions", "", "Write solutions to file")
	flag.IntVar(&numEqToSolve, "solve_equations", 0, "Solve a given number of equations")
	flag.IntVar(&blockSize, "set_block", 0, "Set the block size in equation solving")
	flag.Parse()


	client, err := rpc.Dial("tcp", "localhost:" + port)
	if err != nil {
		fmt.Println("error in client: " + err.Error())
	}


	reply := new(Reply)

	if status {
		err = client.Call("Handler.Status", 0, &reply)
		if err != nil {
			fmt.Println("error in client: " + err.Error())
		}
		fmt.Printf("# of equations: %v\n", reply.NumEq)
		fmt.Printf("# of solutions: %v\n", reply.NumSol)
	} else

	if blockSize > 0 {
		err = client.Call("Handler.SetBlock", blockSize, &reply)
		if err != nil {
			fmt.Println("error in client: " + err.Error())
		}
	} else

	if equationsFile != "" {
		err = client.Call("Handler.ReadEquations", equationsFile, &reply)
		if err != nil {
			fmt.Println("error in client: " + err.Error())
		}
	} else

	if solutionsFile != "" {
		err = client.Call("Handler.ReadSolutions", solutionsFile, &reply)
		if err != nil {
			fmt.Println("error in client: " + err.Error())
		}
	} else

	if writeFile != "" {
		err = client.Call("Handler.WriteSolutions", writeFile, &reply)
		if err != nil {
			fmt.Println("error in client: " + err.Error())
		}
	} else

	if numEqToSolve != 0 {
		err = client.Call("Handler.SolveEquations", numEqToSolve, &reply)
		if err != nil {
			fmt.Println("error in client: " + err.Error())
		}
	} else {
		fmt.Println("Please specify operation")
	}
	
}
