package csv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
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
	l   *slog.Logger
}

func NewCSVParser(schemaPath string, l *slog.Logger) Parser {
	f, err := os.Open(schemaPath)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	var def Schema
	if err := json.NewDecoder(f).Decode(&def); err != nil {
		panic(err)
	}

	return &csvParser{
		def: def,
		l:   l,
	}
}

//TODO: do not break the loop on error, just log it and continue

func (p *csvParser) Consume(input io.Reader, records chan<- []byte) error {
	p.l.Info("starting to consume csv file")
	reader := csv.NewReader(input)

	if p.def.SkipHeader {
		p.l.Info("skipping header")
		if _, err := reader.Read(); err != nil {
			return err
		}
	}

	line := 1
	offset := 0
	for {
		record, err := reader.Read()
		if err == io.EOF {
			p.l.Info("reached end of file")
			break
		}
		if err != nil {
			p.l.Error("error reading record", formatError(line, offset, err))
			continue
		}

		validatedRecord, err := validateAndTransformRecord(record, p.def)
		if err != nil {
			p.l.Error("error validating record", formatError(line, offset, err))
			continue
		}

		jsonRecord, err := json.Marshal(validatedRecord)
		if err != nil {
			p.l.Error("error marshalling record", formatError(line, offset, err))
			continue
		}

		offset += len(jsonRecord)
		line++
		records <- jsonRecord
	}

	p.l.Info("finished consuming csv file")
	p.l.Info("Validated %d records\n", line)
	close(records)
	return nil
}

func validateAndTransformRecord(record []string, def Schema) (map[string]string, error) {
	validatedRecord := make(map[string]string)
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
			// Assign the transformed value to the corresponding column name
			validatedRecord[colDef.Name] = transformedValue
		} else {
			return nil, fmt.Errorf("unsupported column type: %s", colDef.Type)
		}
	}

	return validatedRecord, nil
}

func formatError(line, offset int, err error) error {
	return fmt.Errorf("line %d, offset %d: %s", line, offset, err.Error())
}
