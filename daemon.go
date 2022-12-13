package main

import (
	"github.com/andyytliu/go-finite-field/solver"
	"github.com/andyytliu/go-finite-field/server"
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
)

func main() {

	flag.StringVar(&port, "port", "8080", "Port number for the server to listen to")
	flag.StringVar(&logFile, "log", "logs.txt", "Name for the log file")
	flag.IntVar(&prime, "p", 46337, "Prime number to use in modular calculation")
	flag.IntVar(&server.BlockSolveSize, "block", 200, "Block size of equations to solve in parallel")
	flag.Parse()

	solver.SolverPrime = solver.FF(prime)


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
	log.Printf("*********** > Block size: %v\n", server.BlockSolveSize)
	log.Println("***********************************")
	

	handler := server.NewHandler()
	rpc.Register(handler)

	listener, err := net.Listen("tcp", ":" + port)
	if err != nil {
		log.Println(">>>>>>>>>>> errorr in server: " + err.Error())
	}
	defer listener.Close()
	
	rpc.Accept(listener)
}
