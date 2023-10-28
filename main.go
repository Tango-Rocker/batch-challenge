package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"github.com/Tango-Rocker/batch-challange/model"
	"github.com/Tango-Rocker/batch-challange/validation"
	"io"
	"os"
	"sync"
)

var ValidatorsMap = map[string]validation.Validator{
	"float":   validation.FloatValidator,
	"integer": validation.IntegerValidator,
	"date":    validation.DateValidator,
}

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
		return
	}

	file, err := os.Open(fullPath)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	defer file.Close()

	records := make(chan []string)
	var wg sync.WaitGroup

	// Start CSV reading and validation in a goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := validateCSVWithDefinition(file, records, def); err != nil {
			fmt.Println("Validation error:", err)
			close(records)
		}
	}()

	// Start processing the records in a separate goroutine
	wg.Add(1)
	go func() {
		defer wg.Done()
		processRecords(records, os.Stdout)
	}()

	wg.Wait()
}

func validateCSVWithDefinition(input io.Reader, records chan<- []string, def model.CSVDefinition) error {
	reader := csv.NewReader(input)

	if def.SkipHeader {
		if _, err := reader.Read(); err != nil {
			return err
		}
	}

	i := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Validate and transform the record
		validatedRecord, err := validateAndTransformRecord(record, def)
		if err != nil {
			return err
		}

		// Send the valid record to the channel
		records <- validatedRecord
		i++
	}

	fmt.Printf("Validated %d records\n", i)
	close(records)
	return nil
}

func processRecords(records <-chan []string, output io.Writer) {
	i := 0
	for record := range records {
		// Process the record to construct the desired payload
		// For the sake of the example, we just print it out.
		_, err := fmt.Fprintln(output, record)
		if err != nil {
			fmt.Printf("Error writing record %d: %s\n", i, err.Error())
		}
		i++
	}
	fmt.Printf("Processed %d records\n", i)
}

func validateAndTransformRecord(record []string, def model.CSVDefinition) ([]string, error) {
	if len(record) != len(def.Columns) {
		return nil, fmt.Errorf("incorrect number of columns: got %d, want %d", len(record), len(def.Columns))
	}

	for i, value := range record {
		colDef := def.Columns[i]
		if value == "" && colDef.Required {
			return nil, fmt.Errorf("column %s is required but empty", colDef.Name)
		}

		if validator, exists := ValidatorsMap[colDef.Type]; exists {
			transformedValue, err := validation.ValidateAndTransform(value, validator)
			if err != nil {
				return nil, fmt.Errorf("column %s: %s", colDef.Name, err.Error())
			}
			record[i] = transformedValue
		} else {
			return nil, fmt.Errorf("unsupported column type: %s", colDef.Type)
		}
	}

	return record, nil
}
