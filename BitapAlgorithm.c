#include <stdio.h> 
#include <string.h> 

#define MAX_CANDIDATES 1000 // Maksimalni broj pronađenih kandidata 

void Bitap(const char *text, const char *pattern, int tolerance, char *candidates[]) { 
    int m = strlen(pattern); 
    int n = strlen(text); 

    int R[tolerance + 1]; 
    int pattern_mask[256]; // Maskirani uzorak 

    int i, j; 

    // Inicijalizacija maskiranog uzorka 
    for (i = 0; i < 256; i++) 
        pattern_mask[i] = 0; 

    for (i = 0; i < m; i++) 
        pattern_mask[pattern[i]] |= 1 << i; 

    // Inicijalizacija polja kandidata 
    for (i = 0; i < tolerance + 1; i++) 
        R[i] = ~1; 

    int match; 
	
    for (i = 0; i < n; i++) { 
        unsigned int old_Rd1 = R[0]; 
		
        R[0] |= pattern_mask[text[i]]; 
        R[0] <<= 1; 

        for (j = 1; j <= tolerance; j++) { 
            unsigned int tmp = R[j]; 

            R[j] = (old_Rd1 & (R[j] | pattern_mask[text[i]])) << 1; 
            old_Rd1 = tmp; 
        } 

        if (R[tolerance] & (1 << m)) { 
            match = i - m + 1; 
            candidates[match] = &text[match]; 
        } 
    } 
} 

  

int main() { 
    char text[] = "ATCGATCGATCGATCG"; 
    char pattern[] = "ATCG"; 
    int tolerance = 1; 

    char *candidates[MAX_CANDIDATES]; // Polje za kandidate 
    int i; 

    // Inicijalizacija polja kandidata 
    for (i = 0; i < MAX_CANDIDATES; i++) 
        candidates[i] = NULL; 

    Bitap(text, pattern, tolerance, candidates); 

  

    // Ispis pronađenih kandidata 
    printf("Pronađeni kandidati:\n"); 
    for (i = 0; candidates[i] != NULL; i++) 
        printf("%s\n", candidates[i]); 

    return 0; 

} 