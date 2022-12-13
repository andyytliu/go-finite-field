package main

import (
	"github.com/andyytliu/go-finite-field/solver"
	"github.com/andyytliu/go-finite-field/server"
	"flag"
	"fmt"
	"log"
	"os"
	"runtime"
)

var (
	prime int
	logFile string
)

func main() {

	flag.IntVar(&prime, "p", 46337, "Prime number to use in modular calculation")
	flag.IntVar(&server.BlockSolveSize, "block", 200, "Block size of equations to solve in parallel")
	flag.StringVar(&logFile, "log", "logs.txt", "Name for the log file")
	flag.Parse()

	solver.SolverPrime = solver.FF(prime)


	os.Remove(logFile)
	file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	log.SetOutput(file)
	log.Println("*********** Start logging")


	fmt.Println("**********************************************")
	fmt.Println("****")
	fmt.Println("****    Go Finite Field    v0.7")
	fmt.Println("****        A Linear Finite-Field Solver")
	fmt.Println("****")
	fmt.Println("****      Author: Andy Yu-Ting Liu")
	fmt.Println("****")
	fmt.Println("**********************************************")
	fmt.Println("********")
	fmt.Printf("********      Log file: %v\n", logFile)
	fmt.Printf("********      NumCPU/GOMAXPROCS: %v/%v\n",
		runtime.NumCPU(), runtime.GOMAXPROCS(0))
	fmt.Printf("********      Prime: %v\n", prime)
	fmt.Printf("********      Block size: %v\n", server.BlockSolveSize)
	fmt.Println("********")
	fmt.Println("**********************************************")

	handler := server.NewHandler()
	reply := new(server.Reply)

	for {
		var com string
		var num int
		fmt.Printf(" <Please enter command>  ")
		fmt.Scan(&com)

		if com == "help" {
			PrintHelp()
		} else if com == "status" {
			handler.Status(0, reply)
			fmt.Println("**********************************************")
			fmt.Printf(" Prime: %v\n", prime)
			fmt.Printf(" Block size: %v\n", reply.BlockSize)
			fmt.Printf(" # of equations: %v\n", reply.NumEq)
			fmt.Printf(" # of solutions: %v\n", reply.NumSol)
			fmt.Println("**********************************************")
		} else if com == "block" {
			fmt.Scan(&num)
			handler.SetBlock(num, reply)
		} else if com == "read_eq" {
			fmt.Scan(&com)
			handler.ReadEquations(com, reply)
		} else if com == "read_sol" {
			fmt.Scan(&com)
			handler.ReadSolutions(com, reply)
		} else if com == "write_sol" {
			fmt.Scan(&com)
			handler.WriteSolutions(com, reply)
		} else if com == "solve" {
			fmt.Scan(&num)
			handler.SolveEquations(num, reply)
		} else {
			fmt.Println("Please enter a valid command")
		}
	}

}

func PrintHelp() {
	fmt.Println()
	fmt.Println("**********************************************")
	fmt.Println(" Available commands:")
	fmt.Println(" help:        Print this help message")
	fmt.Println(" status:      Print status of the solver\n" +
		"              including numbers of equations\n" +
		"              left and current solutions")
	fmt.Println(" block:       Change the block size")
	fmt.Println(" read_eq:     Read equations from file")
	fmt.Println(" read_sol:    Read solutions from file")
	fmt.Println(" write_sol:   Write solutions to file")
	fmt.Println(" solve:       Solve given number of equations")
	fmt.Println("**********************************************")
	fmt.Println()
}