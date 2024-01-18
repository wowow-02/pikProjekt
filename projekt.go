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
	sample := "ATC"
	tolerance := 1
	time, candidatesArray := runAlgorithm(refs, sample, tolerance, AlgoritamZ)
	fmt.Println("Vrijeme za algoritam z je: ", time, " ,a kandidati su: ", candidatesArray)

}

func runAlgorithm(input []string, sample string, tolerance int, alg func([]string, string, int, *chan string)) (time.Duration, []string) {

	// Kreiraj waitgrupu za cekanje svih goroutina
	wg := new(sync.WaitGroup)

	//Kanal sa vrijednostima, tu se stavljaju kandidati
	ch := make(chan string)

	// Start timewatch
	startTime := time.Now()

	// Pozivanje 8 GO rutina
	wg.Add(1)
	go algorithmFeeder(0, 49, 749, 799, input, sample, tolerance, wg, &ch, alg)
	wg.Add(1)
	go algorithmFeeder(50, 99, 699, 749, input, sample, tolerance, wg, &ch, alg)
	wg.Add(1)
	go algorithmFeeder(100, 149, 649, 699, input, sample, tolerance, wg, &ch, alg)
	wg.Add(1)
	go algorithmFeeder(150, 199, 599, 649, input, sample, tolerance, wg, &ch, alg)
	wg.Add(1)
	go algorithmFeeder(200, 249, 549, 599, input, sample, tolerance, wg, &ch, alg)
	wg.Add(1)
	go algorithmFeeder(250, 299, 499, 549, input, sample, tolerance, wg, &ch, alg)
	wg.Add(1)
	go algorithmFeeder(300, 349, 449, 499, input, sample, tolerance, wg, &ch, alg)
	wg.Add(1)
	go algorithmFeeder(350, 399, 399, 449, input, sample, tolerance, wg, &ch, alg)

	var candidatesArray []string

	// Cekanje da se sve gorutine zavrse
	go func() {
		wg.Wait()
		close(ch)
	}()

	// Stavljanje svih kandidata u array
	for value := range ch {
		candidatesArray = append(candidatesArray, value)
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
	ch *chan string,
	alg func([]string, string, int, *chan string)) {
	defer wg.Done()
	subarray1 := input[startIndex1:endIndex1]
	copySubArray1 := make([]string, len(subarray1))
	copy(copySubArray1, subarray1)

	subarray2 := input[startIndex2:endIndex2]
	copySubArray2 := make([]string, len(subarray2))
	copy(copySubArray2, subarray2)

	mergedArray := append(copySubArray1, copySubArray2...)
	alg(mergedArray, sample, tolerance, ch)
}

func AlgoritamZ(input []string, sub string, tolerance int, ch *chan string) {
	for _, ref := range input {
		refLen, subLen := len(ref), len(sub)
		concatenated := make([]byte, refLen+subLen+2)
		copy(concatenated, sub)
		concatenated[subLen] = '$' // Specijalni separator između sub i ref
		copy(concatenated[subLen+1:], ref)
		concatenatedLen := len(concatenated)

		Z := make([]int, concatenatedLen)
		left, right := 0, 0
		candidateCount := 0

		for k := 1; k < concatenatedLen; k++ {
			// Ako je k > right, nema matcha, Z[k] se racuna:
			if k > right {
				left, right = k, k
				for right < concatenatedLen && concatenated[right-left] == concatenated[right] {
					right++
				}
				Z[k] = right - left
				right--
			} else {
				// k1 je broj koji odgovara matchu izmedu left i right
				k1 := k - left
				// Z[k1] < od preostalog intervala, onda je on Z[k]
				if Z[k1] < right-k+1 {
					Z[k] = Z[k1]
				} else {
					left = k
					for right < concatenatedLen && concatenated[right-left] == concatenated[right] {
						right++
					}
					Z[k] = right - left
					right--
				}
			}

			// Pronalazak podniza unutar ref uz uvažavanje tolerancije
			if Z[k] == subLen && k >= subLen+1 && k+subLen-1 <= refLen+subLen {
				errorCount := 0
				for i := 0; i < subLen; i++ {
					if sub[i] != ref[k+i-subLen-1] {
						errorCount++
					}
				}
				if errorCount <= tolerance {
					startIndex := k - subLen - 1
					substring := ref[startIndex : startIndex+subLen]
					*ch <- substring //TODO tu ne znam sta proslijedit, vidim da je u izvornom kodu bio substring al nema mi to smisla kao kandidat?!
					fmt.Printf("Pattern found at position: %d\n", startIndex)
					candidateCount++
				}
			}
		}
	}
}
