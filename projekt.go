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
	fmt.Println("Broj očitanja: ", len(refs))

	//Ista stvar sa ostalim algoritmima (samo copy paste)
	sample := "AACTGCGCA"
	tolerance := 7

	//time, candidatesArray := runAlgorithm(refs, sample, tolerance, Lavenshtein)
	//time, candidatesArray := runAlgorithm(refs, sample, tolerance, RabinKarp)
	//time, candidatesArray := runAlgorithm(refs, sample, tolerance, Naive_algorithm)
	time, candidatesArray := runAlgorithm(refs, sample, tolerance, Bitap_algorithm)

	//fmt.Println("Vrijeme za algoritam z je: ", time, " ,a kandidati su: ", candidatesArray)
	fmt.Println("Postavljena vrijednost tolerancije: ", tolerance)
	fmt.Println("Broj kandidata: ", len(candidatesArray))
	fmt.Println("Vrijeme za algoritam je: ", time)

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

// ---------------------------------------------------------------------------
func bitapSearch(ref, pattern string, tolerance int, kandidati *chan string) {
	m := len(pattern)
	mask := make([]int, 128)

	for _, char := range pattern {
		mask[char] |= 1 << (m - 1)
	}

	// Inicijalizacija tablice s nulama
	R := 0
	shiftedMask := mask

	for i := 0; i < m; i++ {
		shiftedMask[pattern[i]] |= 1 << i
	}

	for _, char := range ref {
		R = (R << 1) | 1
		R &= shiftedMask[char]

		// Ako je najviši bit nula, podudaranje pronađeno
		if R&(1<<(m-1)) == 0 {
			*kandidati <- ref
			break
		}
	}

	for i := 1; i <= tolerance; i++ {
		R = R << 1

		for j := range mask {
			R |= mask[j]

			// Ako je najviši bit nula, podudaranje pronađeno
			if R&(1<<(m-1)) == 0 {
				*kandidati <- ref
				break
			}
		}
	}
}

func Bitap_algorithm(ref []string, sub string, tolerance int, kandidati *chan string) {
	for _, refString := range ref {
		bitapSearch(refString, sub, tolerance, kandidati)
	}
}

// ---------------------------------------------------------------------------
func levenshteinDistance(a, b string) int {
	m, n := len(a), len(b)
	dp := make([][]int, m+1) //Inicijalizacija i izrada matrice Levensteinove udaljenosti

	for i := range dp {
		dp[i] = make([]int, n+1)
	}

	for i := 0; i <= m; i++ {
		for j := 0; j <= n; j++ {
			if i == 0 {
				dp[i][j] = j
			} else if j == 0 {
				dp[i][j] = i
			} else if a[i-1] == b[j-1] {
				dp[i][j] = dp[i-1][j-1]
			} else {
				dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
			}
		}
	}

	return dp[m][n]
}

func min(a, b, c int) int {
	if a < b {
		if a < c {
			return a
		}
		return c
	}

	if b < c {
		return b
	}
	return c
}

func Lavenshtein(ref []string, sub string, tolerance int, kandidati *chan string) {
	for _, refString := range ref {
		for i := 0; i <= len(refString)-len(sub); i++ {
			substring := refString[i : i+len(sub)]

			// Računaj Levenshtein udaljenost između podniza i podniza iz referentnog niza
			distance := levenshteinDistance(sub, substring)

			// Ako je udaljenost manja ili jednaka toleranciji, dodaj podniz u kanal
			if distance <= tolerance {
				*kandidati <- substring
			}
		}
	}
}

//---------------------------------------------------------------------------

const (
	ALPHABET_SIZE  = 4
	PRIME          = 101
	MAX_CANDIDATES = 1000000
)

func hash(str string, length int) int {
	hashValue := 0
	for i := 0; i < length; i++ {
		hashValue = (hashValue*ALPHABET_SIZE + int(str[i])) % PRIME
	}
	return hashValue
}

func power(x, y int) int {
	result := 1
	for i := 0; i < y; i++ {
		result = (result * x) % PRIME
	}
	return result
}

func rehash(oldHash int, oldChar, newChar byte, length int) int {
	newHash := (oldHash - int(oldChar)*power(ALPHABET_SIZE, length-1)) % PRIME
	newHash = (newHash*ALPHABET_SIZE + int(newChar)) % PRIME
	if newHash < 0 { //Ne želimo negativne vrijednosti za našu fju
		newHash += PRIME
	}
	return newHash
}

func RabinKarp(ref []string, sub string, tolerance int, kandidati *chan string) {
	subLen := len(sub)
	for _, refString := range ref {
		refLen := len(refString) //Početne hash vrijednosti koje se onda preračunavaju s rolling hashom
		refHash := hash(refString, subLen)
		subHash := hash(sub, subLen)

		for i := 0; i <= refLen-subLen; i++ {
			if refHash == subHash {
				errorCount := 0
				for j := 0; j < subLen; j++ {
					if refString[i+j] != sub[j] {
						errorCount++
					}
				}
				if errorCount <= tolerance {
					*kandidati <- refString[i : i+subLen]
					//*kandidati = append(*kandidati, refString[i:i+subLen])
					//fmt.Printf("Pattern found at position: %d in ref: %s\n", i, refString)
				}
			}
			if i < refLen-subLen { //rolling hash
				refHash = rehash(refHash, refString[i], refString[i+subLen], subLen)
			}
		}
	}
}

// ---------------------------------------------------------------------------
func Naive_algorithm(ref []string, sub string, tolerance int, kandidati *chan string) {
	refLen, subLen := len(ref), len(sub)
	for _, refString := range ref {
		for i := 0; i <= refLen-subLen; i++ {
			diffCount := 0

			for j := 0; j < subLen; j++ {
				if sub[j] != refString[i+j] {
					diffCount++
				}

				if diffCount > tolerance {
					break
				}
			}

			if diffCount <= tolerance {
				*kandidati <- refString[i : i+subLen]
				//fmt.Printf("Pattern found at position: %d\n", i)
			}
		}
	}
}

//---------------------------------------------------------------------------
