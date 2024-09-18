package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

// AirQualityReading represents a sensor's air quality data.
type AirQualityReading struct {
	SensorID  string    `json:"sensor_id"`
	Timestamp time.Time `json:"timestamp"`
	PM25      float64   `json:"pm25"`
	CO2       float64   `json:"co2"`
	NO2       float64   `json:"no2,omitempty"`
	Ozone     float64   `json:"ozone,omitempty"`
}

// parseReadings takes a byte array of JSON data and unmarshals it into a slice of AirQualityReading structs.
func parseReadings(data []byte) ([]AirQualityReading, error) {
	var readings []AirQualityReading

	// Attempt to unmarshal the JSON data
	err := json.Unmarshal(data, &readings)
	if err != nil {
		return nil, errors.New("error parsing sensor readings JSON data: " + err.Error())
	}

	return readings, nil
}

func main() {
	// Sample JSON data simulating sensor readings
	jsonData := `
	[
		{
			"sensor_id": "sensor_001",
			"timestamp": "2024-09-17T08:55:00Z",
			"pm25": 35.2,
			"co2": 412.1,
			"no2": 20.5,
			"ozone": 0.07
		},
		{
			"sensor_id": "sensor_002",
			"timestamp": "2024-09-17T08:55:00Z",
			"pm25": 42.3,
			"co2": 408.3
		}
	]`

	// Parse the JSON data
	readings, err := parseReadings([]byte(jsonData))
	if err != nil {
		fmt.Println("Failed to parse readings:", err)
		return
	}

	// Output the parsed readings
	for _, reading := range readings {
		fmt.Printf("%+v\n", reading)
	}
}
