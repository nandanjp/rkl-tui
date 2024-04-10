package utils

import (
	"encoding/csv"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/nandanjp/rkl/structs"
)

func parseCSV[Res structs.City](file string, parse func(fields []string) (Res, error)) ([]Res, error) {
	f, err := os.Open(file)
	defer func() {
		if err := f.Close(); err != nil {
			fmt.Printf("unexpected error encountered: %v\n", err)
		}
	}()
	if err != nil {
		return nil, err
	}
	lines, err := csv.NewReader(f).ReadAll()
	if err != nil {
		return nil, err
	}
	objects := make([]Res, len(lines))
	for i, line := range lines {
		res, err := parse(line)
		if err != nil {
			return objects, err
		}
		objects[i] = res
	}
	return objects, nil
}

func ParseCities(file string) ([]structs.City, error) {
	return parseCSV(file, func(fields []string) (structs.City, error) {
		if len(fields) < 4 {
			return structs.City{}, fmt.Errorf("failed to parse the given fields into a city struct as there were not enough fields passed in")
		}
		id := fields[0]
		parsedId, idErr := strconv.ParseUint(id, 10, 64)
		if idErr != nil {
			return structs.City{}, idErr
		}
		population := fields[3]
		population = strings.ReplaceAll(population, " ", "")
		parsedPopulation, populationErr := strconv.ParseUint(population, 10, 64)
		if populationErr != nil {
			return structs.City{}, populationErr
		}
		return structs.City{
			Id:         parsedId,
			City:       strings.ReplaceAll(fields[1], " ", ""),
			Country:    strings.ReplaceAll(fields[2], " ", ""),
			Population: parsedPopulation,
		}, nil
	})
}
