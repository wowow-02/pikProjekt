package main
import "fmt"

const MAX_CANDIDATES = 1000

func Bitap(text string, pattern string, tolerance int, candidates []string) {
	m := len(pattern)
	n := len(text)
	
	var R [1001]int
	var patternMask [256]int
	
	//inicijalizacija maskiranog uzorka
	for i := 0; i < 256; i++ {
		patternMask[i] = 0
	}
	
	//inicijalizacija polja kandidata
	for i := 0; i < tolerance + 1; i++ {
		R[i] = ^1
	}
	
	var match int
	
	for i := 0; i < n; i++ {
		oldRd1 := R[0]
		
		R[0] |= patternMask[text[i]]
		R[0] <<= 1
		
		for j :=1; j <= tolerance; j++ {
			tmp := R[j]
			R[j] = (oldRd1 & (R[j] | patternMask[text[i]])) << 1
			oldRd1 = tmp
		}
		
		if R[tolerance]&(1<<m) != 0 {
			match = i - m + 1
			candidates[match] = text[match : match+m]
		}
	}
}

func main() {
	text := "ATCGATCGATCGATCG"
	pattern := "ATCG"
	tolerance := 1
	
	//polja za kandidate
	var candidates [MAX_CANDIDATES]string
	
	//inicijalizacija polja za kandidate
	for i := 0; i < MAX_CANDIDATES; i++ {
		candidates[i] = ""
	}
	
	//pozivanje funkcije Bitap
	Bitap(text, pattern, tolerance, candidates[:])
	
	//ispis pronadenih kandidata
	fmt.Println("Pronađeni kandidati:")
	for _, candidate := range candidates {
		if candidate != "" {
			fmt.Println(candidate)
		}
	}
}