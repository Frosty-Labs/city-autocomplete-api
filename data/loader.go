package data

import (
	"encoding/csv"
	"os"
	"strings"

	"city-autocomplete-api/models"
)

// LoadCities loads city data from CSV file
func LoadCities(filePath string) ([]models.City, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	reader.Comma = ',' // Set the delimiter to comma

	// Read all records
	records, err := reader.ReadAll()
	if err != nil {
		return nil, err
	}

	var cities []models.City

	// Skip header row if it exists
	header := true
	for _, record := range records {
		if header {
			header = false
			continue
		}

		if len(record) >= 4 {
			// Clean up the fields by trimming spaces
			name := strings.TrimSpace(record[0])
			country := strings.TrimSpace(record[1])
			subcountry := strings.TrimSpace(record[2])
			geonameid := strings.TrimSpace(record[3])

			city := models.City{
				Name:       name,
				Country:    country,
				Subcountry: subcountry,
				GeonameID:  geonameid,
			}
			cities = append(cities, city)
		}
	}

	return cities, nil
}
