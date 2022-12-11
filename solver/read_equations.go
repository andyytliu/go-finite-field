package solver

import (
	"bufio"
	"io"
	"log"
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
			
			coef, err := strconv.ParseInt(fields[2*i+1], 10, 32)
			if err != nil {
				log.Println(">>>>>>>>>>> error when parsing coef: " +
					fields[2*i] + ". " + err.Error())
			}
			equation[Index(idx)] = FF(coef)
		}

		if len(equation) != 0 {
			*equations = append(*equations, &equation)
		}

		if err == io.EOF {
			break
		}
	}
}
