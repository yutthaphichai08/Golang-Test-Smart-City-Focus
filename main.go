package main

import (
	"encoding/csv"
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"
)

type AirQualityReading struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	PM25      float64   `json:"pm25"`
	CO2       float64   `json:"co2"`
}

func parseCSVReadings(filename string) ([]AirQualityReading, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, errors.New("error opening CSV file: " + err.Error())
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return nil, errors.New("error reading CSV file: " + err.Error())
	}

	var readings []AirQualityReading

	// Skip the header row and parse each record
	for i, record := range records {
		if i == 0 {
			continue // Skip the header row
		}

		// Parse timestamp
		timestamp, err := time.Parse(time.RFC3339, record[1])
		if err != nil {
			return nil, errors.New("error parsing timestamp: " + err.Error())
		}

		// Parse pm25 and co2 values
		pm25, err := strconv.ParseFloat(record[2], 64)
		if err != nil {
			return nil, errors.New("error parsing PM2.5 value: " + err.Error())
		}
		co2, err := strconv.ParseFloat(record[3], 64)
		if err != nil {
			return nil, errors.New("error parsing CO2 value: " + err.Error())
		}

		// Append the parsed record to the readings slice
		readings = append(readings, AirQualityReading{
			SensorID:  record[0],
			Timestamp: timestamp,
			PM25:      pm25,
			CO2:       co2,
		})
	}

	return readings, nil
}

func calculateAverage(readings []AirQualityReading) map[string]float64 {
	if len(readings) == 0 {
		return map[string]float64{}
	}

	totalPM25 := 0.0
	totalCO2 := 0.0

	for _, reading := range readings {
		totalPM25 += reading.PM25
		totalCO2 += reading.CO2
	}

	// Calculate averages
	averageReadings := map[string]float64{
		"pm25": totalPM25 / float64(len(readings)),
		"co2":  totalCO2 / float64(len(readings)),
	}

	return averageReadings
}

// findHighestPollutantByHour
func findHighestPollutantByHour(readings []AirQualityReading) map[int]string {
	hourlyData := make(map[int]map[string]float64)
	counts := make(map[int]map[string]int)

	// Initialize maps for each hour (0-23)
	for i := 0; i < 24; i++ {
		hourlyData[i] = map[string]float64{"pm25": 0, "co2": 0}
		counts[i] = map[string]int{"pm25": 0, "co2": 0}
	}

	for _, reading := range readings {
		hour := reading.Timestamp.Hour()

		hourlyData[hour]["pm25"] += reading.PM25
		counts[hour]["pm25"]++

		hourlyData[hour]["co2"] += reading.CO2
		counts[hour]["co2"]++
	}

	result := make(map[int]string)

	// Calculate the highest average pollutant for each hour
	for hour := 0; hour < 24; hour++ {
		maxPollutant := ""
		maxAvg := 0.0

		// Calculate average for each pollutant if there are readings
		for pollutant, sum := range hourlyData[hour] {
			if counts[hour][pollutant] > 0 {
				avg := sum / float64(counts[hour][pollutant])
				if avg > maxAvg {
					maxAvg = avg
					maxPollutant = pollutant
				}
			}
		}

		result[hour] = maxPollutant
	}

	return result
}

func main() {
	// Use the CSV file with sensor readings
	csvFilename := "readings.csv"

	// Parse the CSV data
	readings, err := parseCSVReadings(csvFilename)
	if err != nil {
		fmt.Println("Failed to parse readings:", err)
		return
	}

	// Output the parsed readings
	for _, reading := range readings {
		fmt.Printf("%+v\n", reading)
	}

	// Calculate averages
	averageValues := calculateAverage(readings)

	fmt.Println("\nAverages for all readings:")
	for pollutant, avg := range averageValues {
		fmt.Printf("%s: %.2f\n", pollutant, avg)
	}

	// Find the highest average pollutant by hour
	highestByHour := findHighestPollutantByHour(readings)
	fmt.Println("\nHighest pollutant by hour:")
	for hour, pollutant := range highestByHour {
		fmt.Printf("Hour %d: %s\n", hour, pollutant)
	}
}
