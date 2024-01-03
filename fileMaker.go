package main

import (
	"bufio"
	"fmt"
	"net/http"
	"os"
)

func makeFile() {
	// GET request na link zadan u opisu
	resp, err := http.Get("https://nanopore.s3.climb.ac.uk/MAP006-1_2D_pass.fasta")
	if err != nil {
		fmt.Println("Pogreška u dohvaćanju sekvencii s Web stranice")
	}
	defer resp.Body.Close()
	scanner := bufio.NewScanner(resp.Body)

	i := 0
	fmt.Print("Broj sekvenci koji želite spremiti u datoteku: ")
	fmt.Scan(&i)

	//analiziranje GET body-a iz requesta te spremanje zadanog broja sekvenci u .txt datoteku
	cnt := 0
	restrictedBody := ""

	for scanner.Scan() {

		line := scanner.Text()
		if line[0] == '>' {
			cnt++

			if cnt > i {
				break
			}
			continue
		}

		line += "\n"
		restrictedBody += line

	}

	filepath := "sekvence.txt"
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println("Greška prilikom otvaranja datoteke")
	}

	defer file.Close()

	_, err = file.WriteString(restrictedBody)

	if err != nil {
		fmt.Println("Greska prilikom upisivanja sekvenci u datoteku")
	}
	fmt.Println("Sekvence uspješno spremljene u datoteku Sekvence.txt")
}
