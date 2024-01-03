package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
)

func main() {
	//ovaj dio kod odkomentirati ako se želi promjeniti broj sekvenci u datoteci
	/*i := 0
	fmt.Scan(&i)
	if i != 1 {
		makeFile()
	}*/
	//inicijalizacija potrebnih varijabli/polja i otvaranje datoteke
	var fileSekv []string
	orderingMap := make(map[int][]string)
	var refs []string

	filepath := "sekvence.txt"
	file, err := os.Open(filepath)
	if err != nil {
		fmt.Println("Došlo je do greške")
	}

	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		fileSekv = append(fileSekv, scanner.Text())

	}
	//sortiranje po duljini stringa
	sort.Slice(fileSekv, func(i, j int) bool {
		return len(fileSekv[i]) < len(fileSekv[j])
	})

	//algoritam load balancinga
	isReversed := -1
	divide := 0

	for i := 0; i < len(fileSekv); i += 8 {
		isReversed = isReversed * (-1)
		for j := i; j < i+8; j++ {

			if isReversed == 1 {
				orderingMap[divide] = append(orderingMap[divide], fileSekv[j])
			} else {
				orderingMap[divide] = append(orderingMap[divide], fileSekv[(i+8-1)-(j-i)])
			}
			divide++
			if divide == 8 {
				divide = 0
			}
		}
	}
	//spajanje u refs
	for i := 0; i < 8; i++ {

		for j := 0; j < len(orderingMap[i]); j++ {
			refs = append(refs, orderingMap[i][j])
		}
	}
	fmt.Println(len(refs))
}
