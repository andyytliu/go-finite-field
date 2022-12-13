package solver

import (
	"bufio"
	"io"
	"log"
	"os"
	"strings"
	"strconv"
)

func ReadSolutions(file_name string,
	solutions map[Index]map[Index]FF,
	solTransIndex map[Index]map[Index]bool) {

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

		if len(fields) == 0 {
			if err == io.EOF {
				break
			} else {
				continue
			}
		}
		if len(fields) % 2 != 1 {
			log.Println(">>>>>>>>>>> error: non-odd entries in solution; correct format: 'idx idx_1 coef_1 idx_2 coef_2 ...'")
			break
		}

		solution := make(map[Index]FF)
		idx1, err := strconv.ParseUint(fields[0], 10, 32)
		if err != nil {
			log.Println(">>>>>>>>>>> error when parsing index 1 in solution: " +
				fields[0] + ". " + err.Error())
		}

		for i := 0; i < len(fields) / 2; i++ {
			idx2, err := strconv.ParseUint(fields[2*i + 1], 10, 32)
			if err != nil {
				log.Println(">>>>>>>>>>> error when parsing index 2 in solution: " +
					fields[2*i + 1] + ". " + err.Error())
			}

			coef, err := strconv.ParseInt(fields[2*i + 2], 10, 64)
			if err != nil {
				log.Println(">>>>>>>>>>> error when parsing coef in solution: " +
					fields[2*i + 2] + ". " + err.Error())
			}

			// Update solution
			solution[Index(idx2)] = Mod(FF(coef))
			// Update solution transpose indexing
			if idxs, ok := solTransIndex[Index(idx2)]; ok {
				idxs[Index(idx1)] = true
			} else {
				solTransIndex[Index(idx2)] = make(map[Index]bool)
				solTransIndex[Index(idx2)][Index(idx1)] = true
			}
		}

		solutions[Index(idx1)] = solution

		if err == io.EOF {
			break
		}
	}

}
