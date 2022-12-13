# Go Finite Field

**Go Finite Field** is a linear solver over finite fields specifically optimized for solving sparse linear systems.

## Installation and Running the Program

Go Finite Field is built using the [Go](https://go.dev) programming language; therefore a distribution of Go is required to build the program.

To clone the source code to a local folder, run
```
> git clone https://github.com/andyytliu/go-finite-field.git
> cd go-finite-field
```
The program can be started with
```
> go run main.go
```
Alternatively, a binary file can be built first and then be used,
```
> go build main.go
> ./main
```
The program has several flags that can be used to set global parameters.
For their usage, see
```
go run main.g --help
```

Once the program is started, the `help` command,
```
<Please enter command>  help
```
will return a list of available commands and their descriptions.

For further details, please see the accompanying notes (to appear.)
