package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	"github.com/Sindreln/minyr/yr"
)

func main() {

	scanner := bufio.NewScanner(os.Stdin)
	fmt.Print("Skriv inn kommando: (convert/average) ")
	scanner.Scan()
	command := scanner.Text()

	switch command {
	case "convert":
		if _, err := os.Stat("kjevik-temp-fahr-20220318-20230318.csv"); err == nil {
			// Filen finnes
			fmt.Print("Filen kjevik-temp-fahr-20220318-20230318.csv finnes allerede. Vil du generere den på nytt? (j/n): ")
			scanner.Scan()
			answer := strings.ToLower(scanner.Text())
			if answer != "j" && answer != "n" {
				log.Fatal("Ugyldig svar")
			} else if answer == "n" {
				return
			}
		}
		// Kaller funksjonen convertTemperatures for å konvertere Celsius-temperaturer til Fahrenheit
		convertedTemperatures, err := yr.ConvertTemperatures()
		if err != nil {
			log.Fatal(err)
		}

		// Skriver de konverterte temperaturene til en ny fil
		if err := yr.WriteLines(convertedTemperatures, "kjevik-temp-fahr-20220318-20230318.csv"); err != nil {
			log.Fatal(err)
		}

		fmt.Println("Temperaturer konvertert!")

	case "average":
		// Spør brukeren om enheten for gjennomsnittstemperaturen
		var unit string
		for unit != "c" && unit != "f" {
			fmt.Print("Vil du ha gjennomsnittstemperaturen i Celsius eller Fahrenheit? (c/f): ")
			scanner.Scan()
			unit = strings.ToLower(scanner.Text())
		}

		// Åpner csv-filen
		file, err := yr.OpenFile("kjevik-temp-celsius-20220318-20230318.csv")
		if err != nil {
			log.Fatal(err)
		}
		defer yr.CloseFile(file)

		// Leser linjene fra csv-filen
		lines, err := yr.ReadLines(file)
		if err != nil {
			log.Fatal(err)
		}

		// Beregner gjennomsnittstemperaturen
		var sum float64
		count := 0
		for i, line := range lines {
			if i == 0 {
				continue // ignorer header-linjen
			}
			fields := strings.Split(line, ";")
			if len(fields) != 4 {
				log.Fatalf("uventet antall felt i linje %d: %d", i, len(fields))
			}
			if fields[3] == "" {
				continue // ignorer linje med tomt temperaturfelt
			}
			if temperature, err := strconv.ParseFloat(fields[3], 64); err != nil {
				log.Fatalf("kunne ikke analysere temperaturen i linje %d: %s", i, err)
			} else {
				if unit == "f" {
					temperature = temperature*1.8 + 32 // konverter til Fahrenheit
				}
				sum += temperature
				count++
			}
		}
		average := sum / float64(count)

		// Skriv ut gjennomsnittstemperaturen i valgt enhet
		if unit == "f" {
			fmt.Printf("Gjennomsnittlig temperatur: %.2f°F\n", average)
		} else {
			fmt.Printf("Gjennomsnittlig temperatur: %.2f°C\n", average)
		}

	default:
		fmt.Println("Ugyldig kommando")
	}
}
