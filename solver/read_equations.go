package solver

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"math/big"
	"os"
	"strings"
	"strconv"
)

func ReadEquations(file_name string,
	equations *[]*map[Index]FF) {

	file, err := os.Open(file_name)
	if err != nil {
		log.Println(">>>>>>>>>>> error: " + err.Error())
		return
	}
	defer file.Close()

	primeBig := new(big.Int)
	primeBig.SetString(fmt.Sprint(SolverPrime), 10)

	reader := bufio.NewReader(file)
	for {
		var (
			isPrefix = true
			err error
			line, ln []byte
		)

		for isPrefix && err == nil {
			ln, isPrefix, err = reader.ReadLine()
			line = append(line, ln...)
		}
		
		if err != nil && err != io.EOF {
			log.Println(">>>>>>>>>>> error: " + err.Error())
			break
		}
		if isPrefix {
			log.Println(">>>>>>>>>>> error: line only partially read!")
			break
		}

		fields := strings.Fields(string(line))
		if len(fields) % 2 != 0 {
			log.Println(">>>>>>>>>>> error: uneven entries in equation; correct format: 'idx_1 coef_1 idx_2 coef_2 ...'")
			break
		}

		equation := make(map[Index]FF)

		for i := 0; i < len(fields) / 2; i++ {
			idx, err := strconv.ParseUint(fields[2*i], 10, 32)
			if err != nil {
				log.Println(">>>>>>>>>>> error when parsing index: " +
					fields[2*i] + ". " + err.Error())
			}
			

			var coef FF
			split := strings.Split(fields[2*i+1], "/")
			switch len(split) {
			case 1:
				raw := new(big.Int)
				raw, ok := raw.SetString(fields[2*i+1], 10)
				if !ok {
					log.Println(">>>>>>>>>>> error when parsing coef: " +
						fields[2*i+1] + ". ")
				}
				raw.Mod(raw, primeBig)
				coef = Mod(FF(raw.Int64()))
			case 2:
				num := new(big.Int)
				den := new(big.Int)

				num, ok := num.SetString(split[0], 10)
				if !ok {
					log.Println(">>>>>>>>>>> error when parsing coef: " +
						fields[2*i+1] + ". ")
				}
				num.Mod(num, primeBig)

				den, ok = den.SetString(split[1], 10)
				if !ok {
					log.Println(">>>>>>>>>>> error when parsing coef: " +
						fields[2*i+1] + ". ")
				}
				den.Mod(den, primeBig)

				coef = Mod( FF(num.Int64()) * InverseMod(FF(den.Int64())) )
			default:
				log.Println(">>>>>>>>>>> error when parsing coef: " +
						fields[2*i+1] + ". ")
			}

			equation[Index(idx)] = coef
		}

		if len(equation) != 0 {
			*equations = append(*equations, &equation)
		}

		if err == io.EOF {
			break
		}
	}
}
