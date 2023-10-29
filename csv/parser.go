package csv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"os"
)

var validatorsMap = map[string]Validator{
	"float":   FloatValidator,
	"integer": IntegerValidator,
	"date":    DateValidator,
}

type Parser interface {
	Consume(input io.Reader, records chan<- []byte) error
}

type csvParser struct {
	def Schema
}

func NewCSVParser(schemaPath string) Parser {
	f, err := os.Open(schemaPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var def Schema
	if err := json.NewDecoder(f).Decode(&def); err != nil {
		panic(err)
	}

	return &csvParser{def: def}
}

func (p *csvParser) Consume(input io.Reader, records chan<- []byte) error {
	reader := csv.NewReader(input)

	if p.def.SkipHeader {
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
		validatedRecord, err := validateAndTransformRecord(record, p.def)
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

func validateAndTransformRecord(record []string, def Schema) ([]string, error) {
	if len(record) != len(def.Columns) {
		return nil, fmt.Errorf("incorrect number of columns: got %d, want %d", len(record), len(def.Columns))
	}

	for i, value := range record {
		colDef := def.Columns[i]
		if value == "" && colDef.Required {
			return nil, fmt.Errorf("column %s is required but empty", colDef.Name)
		}

		if validator, exists := validatorsMap[colDef.Type]; exists {
			transformedValue, err := validateAndTransform(value, validator)
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
