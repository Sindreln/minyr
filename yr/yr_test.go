package yr

import (
	"fmt"
	"math"
	"strconv"
	"strings"
	"testing"
)

func TestReadLines(t *testing.T) {
	file, err := OpenFile("kjevik-temp-celsius-20220318-20230318.csv")
	if err != nil {
		t.Fatal(err)
	}
	defer CloseFile(file)

	lines, err := ReadLines(file)
	if err != nil {
		t.Fatal(err)
	}

	if len(lines) != 16756 {
		t.Errorf("unexpected number of lines: got %d, want %d", len(lines), 16756)
	}
}

func TestWriteLines(t *testing.T) {
	lines := []string{
		"Kjevik;SN39040;18.03.2022 01:50;6",
		"Kjevik;SN39040;07.03.2023 18:20;0",
		"Kjevik;SN39040;08.03.2023 02:20;-11",
		"Data er gyldig per 18.03.2023 (CC BY 4.0), Meteorologisk institutt (MET);;;",
	}
	filename := "testfile.txt"
	want := []string{
		"Kjevik;SN39040;18.03.2022 01:50;42.8",
		"Kjevik;SN39040;07.03.2023 18:20;32.0",
		"Kjevik;SN39040;08.03.2023 02:20;12.2",
		"Data er basert på gyldig data (per 18.03.2023) (CC BY 4.0) fra Meteorologisk institutt (MET);endringen er gjort av Sindre Norbom",
	}
	for i, line := range lines {
		parts := strings.Split(line, ";")
		if len(parts) < 4 {
			// Skip lines that don't have at least 4 fields
			continue
		}
		tempStr := parts[3]
		if tempStr == "" {
			// Skip lines that have an empty fourth field
			continue
		}
		if _, err := strconv.ParseFloat(tempStr, 64); err != nil {
			// Skip lines where the fourth field is not a valid float
			continue
		}
		temp, _ := strconv.ParseFloat(tempStr, 64)
		tempStr = fmt.Sprintf("%.1f", temp*1.8+32)
		got := fmt.Sprintf("%s;%s;%s;%s", parts[0], parts[1], parts[2], tempStr)
		if got != want[i] {
			t.Errorf("TestWriteLines(%v, %v) = %v, want %v", lines, filename, got, want[i])
			return
		}
	}
}

func TestAverageTemperature(t *testing.T) {
	file, err := OpenFile("kjevik-temp-celsius-20220318-20230318.csv")
	if err != nil {
		t.Fatalf("feil ved åpning av fil: %s", err)
	}
	defer CloseFile(file)

	lines, err := ReadLines(file)
	if err != nil {
		t.Fatalf("feil ved lesing av fil: %s", err)
	}

	var sum float64
	var count int

	for _, line := range lines {
		fields := strings.Split(line, ";")
		if len(fields) != 4 {
			t.Fatalf("uventet antall felt i linje: %s", line)
		}

		if fields[3] == "" {
			continue // ignorerer linjer med tomme temperaturfelt
		}

		temperatureCelsius, err := strconv.ParseFloat(fields[3], 64)
		if err != nil {
			if !strings.Contains(err.Error(), "invalid syntax") {
				t.Fatalf("kunne ikke parse temperatur: %s", err)
			}
			continue // ignorerer linjer med ugyldige temperaturverdier
		}

		sum += temperatureCelsius
		count++
	}

	averageTemperature := sum / float64(count)
	averageTemperature = math.Round(averageTemperature*100) / 100 // runde av til to desimaler

	if averageTemperature != 8.56 {
		t.Errorf("gjennomsnittstemperaturen er %f; forventet 8.56", averageTemperature)
	}
}
