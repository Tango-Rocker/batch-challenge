package main

import (
	"encoding/csv"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/Tango-Rocker/batch-challange/model"
	"github.com/Tango-Rocker/batch-challange/validation"
	"os"
)

const definitionJSON = `
	{
		"skip_header": true,
		"columns": [
			{
				"name": "id",
				"type": "integer",
				"required": true
			},
			{
				"name": "date",
				"type": "date",
				"required": true,
				"formats": ["01/2006", "January 2006"]
			},
			{
				"name": "amount",
				"type": "float",
				"required": true
			}
		]
	}`

func main() {
	path := os.Getenv("SOURCE_PATH")
	fileName := os.Getenv("FILE_NAME")

	fullPath := path + string(os.PathSeparator) + fileName

	fmt.Println("reading from source: ", fullPath)

	var def model.CSVDefinition
	err := json.Unmarshal([]byte(definitionJSON), &def)
	if err != nil {
		fmt.Println("Error:", err)
	}

	data, err := validateCSVWithDefinition(fullPath, def)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("CSV is valid!")
	}

	fmt.Println(data)
}

var ValidatorsMap = map[string]validation.Validator{
	"float":   validation.FloatValidator,
	"integer": validation.IntegerValidator,
	"date":    validation.DateValidator,
}

func validateCSVWithDefinition(csvPath string, def model.CSVDefinition) ([]map[string]string, error) {
	file, err := os.Open(csvPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	reader := csv.NewReader(file)

	if def.SkipHeader {
		_, err = reader.Read()
		if err != nil {
			return nil, err
		}
	}
	mappedRecords := make([]map[string]string, 0)
	records := make([][]string, 0)
	i := 0
	for {
		record, err := reader.Read()
		if err != nil {
			fmt.Println("Error:", err)
			break
		}

		if len(record) != len(def.Columns) {
			return nil, fmt.Errorf("row %d has an incorrect number of columns", len(records))
		}

		if len(record) != len(def.Columns) {
			return nil, fmt.Errorf("row %d has an incorrect number of columns", i+1)
		}

		for j, value := range record {
			colDef := def.Columns[j]
			if value == "" && colDef.Required {
				return nil, fmt.Errorf("row %d, Column %s is required but empty", i+1, colDef.Name)
			}

			if validator, exists := ValidatorsMap[colDef.Type]; exists {
				transformedValue, err := validation.ValidateAndTransform(value, validator)
				if err != nil {
					return nil, fmt.Errorf("row %d, Column %s: %s", i+1, colDef.Name, err.Error())
				}
				record[j] = transformedValue
			} else {
				return nil, fmt.Errorf("row %d, Column %s: Unsupported column type", i+1, colDef.Name)
			}
		}

		mappedRow, err2 := rowToMap(record, def)
		if err2 != nil {
			return nil, err2
		}
		mappedRecords = append(mappedRecords, mappedRow)
		i++
	}

	return mappedRecords, nil
}

func rowToMap(row []string, def model.CSVDefinition) (map[string]string, error) {
	mappedRow := make(map[string]string)

	if len(row) != len(def.Columns) {
		return nil, errors.New("row length doesn't match definition")
	}

	for i, value := range row {
		colDef := def.Columns[i]
		mappedRow[colDef.Name] = value
	}

	return mappedRow, nil
}
