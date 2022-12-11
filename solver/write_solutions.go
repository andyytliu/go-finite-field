package solver

import (
	"bufio"
	"log"
	"os"
	"strconv"
)


func WriteSolutions(file_name string,
	solutions map[Index]map[Index]FF) {

	file, err := os.Create(file_name)
	if err != nil {
		log.Println(">>>>>>>>>>> error: " + err.Error())
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	for idx1, solution := range solutions {
		_, err := writer.WriteString(strconv.FormatUint(uint64(idx1), 10) + " ")
		if err != nil {
			log.Println(">>>>>>>>>>> error: " + err.Error())
		}
		for idx2, coef := range solution {
			_, err := writer.WriteString(strconv.FormatUint(uint64(idx2), 10) + " ")
			if err != nil {
				log.Println(">>>>>>>>>>> error: " + err.Error())
			}
			_, err = writer.WriteString(strconv.FormatInt(int64(coef), 10) + " ")
			if err != nil {
				log.Println(">>>>>>>>>>> error: " + err.Error())
			}
		}
		_, err = writer.WriteString("\n")
		if err != nil {
			log.Println(">>>>>>>>>>> error: " + err.Error())
		}
		writer.Flush()
	}
}
