package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"sync"
	"time"
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

	//Ista stvar sa ostalim algoritmima (samo copy paste)
	sample := "test"
	tolerance := 4
	time, candidatesArray := runAlgorithm(refs, sample, tolerance, testAlg)
	fmt.Println("Vrijeme za algoritam x je: ", time, " ,a kandidati su: ", candidatesArray)

}

func runAlgorithm(input []string, sample string, tolerance int, alg func([]string, string, int, *sync.WaitGroup, chan string)) (time.Duration, []string) {

	// Kreiraj waitgrupu za cekanje svih goroutina
	var wg sync.WaitGroup

	//Kanal sa vrijednostima, tu se stavljaju kandidati
	ch := make(chan string)

	// Start timewatch
	startTime := time.Now()
	/*
		// Deklaracija varijabla
		var startIndex1 int
		var endIndex1 int
		var startIndex2 int
		var endIndex2 int
		var subarray1 []string
		var subarray2 []string
		var copySubArray1 []string
		var copySubArray2 []string
		var mergedArray []string

		// Prva GO rutina
		startIndex1 = 0
		endIndex1 = 49
		subarray1 = input[startIndex1:endIndex1]
		copySubArray1 = make([]string, len(subarray1))
		copy(copySubArray1, subarray1)

		startIndex2 = 749
		endIndex2 = 799
		subarray2 = input[startIndex2:endIndex2]
		copySubArray2 = make([]string, len(subarray2))
		copy(copySubArray2, subarray2)

		mergedArray = append(copySubArray1, copySubArray2...)
		wg.Add(1)
		go alg(mergedArray, sample, tolerance, &wg, ch)
	*/
	//ako ne bude radilo na ovaj nacin mozemo probat sa ovim gore copy paste al msm da ce radit ovo
	// Pozivanje 8 GO rutina
	algorithmFeeder(0, 49, 749, 799, input, sample, tolerance, &wg, ch, alg)
	algorithmFeeder(50, 99, 699, 749, input, sample, tolerance, &wg, ch, alg)
	algorithmFeeder(100, 149, 649, 699, input, sample, tolerance, &wg, ch, alg)
	algorithmFeeder(150, 199, 599, 649, input, sample, tolerance, &wg, ch, alg)
	algorithmFeeder(200, 249, 549, 599, input, sample, tolerance, &wg, ch, alg)
	algorithmFeeder(250, 299, 499, 549, input, sample, tolerance, &wg, ch, alg)
	algorithmFeeder(300, 349, 449, 499, input, sample, tolerance, &wg, ch, alg)
	algorithmFeeder(350, 399, 399, 449, input, sample, tolerance, &wg, ch, alg)

	// Cekanje da se sve gorutine zavrse
	wg.Wait()
	close(ch)

	candidatesArray := make([]string, len(ch))

	i := 0
	// Stavljanje svih kandidata u array
	for value := range ch {
		candidatesArray[i] = value
		i = i + 1
	}

	// End timewatch
	endTime := time.Now()

	elapsedTime := endTime.Sub(startTime)
	return elapsedTime, candidatesArray
}

func algorithmFeeder(
	startIndex1 int,
	endIndex1 int,
	startIndex2 int,
	endIndex2 int,
	input []string,
	sample string,
	tolerance int,
	wg *sync.WaitGroup,
	ch chan string,
	alg func([]string, string, int, *sync.WaitGroup, chan string)) {

	subarray1 := input[startIndex1:endIndex1]
	copySubArray1 := make([]string, len(subarray1))
	copy(copySubArray1, subarray1)

	subarray2 := input[startIndex2:endIndex2]
	copySubArray2 := make([]string, len(subarray2))
	copy(copySubArray2, subarray2)

	mergedArray := append(copySubArray1, copySubArray2...)
	wg.Add(1)
	go alg(mergedArray, sample, tolerance, wg, ch)

}

func testAlg(input []string,
	sample string,
	tolerance int,
	wg *sync.WaitGroup,
	ch chan string) {
	defer wg.Done()

	//svi kandidati se stavljaju u kanal ch <- kandidat
}
