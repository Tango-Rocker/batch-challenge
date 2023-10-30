package csv

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"time"
)

var validatorsMap = map[string]Validator{
	"float":   FloatValidator,
	"integer": IntegerValidator,
	"date":    DateValidator,
}

type Record struct {
	Line        uint64
	Offset      uint64
	ExecutionId string
	Values      map[string]interface{}
}

type Parser interface {
	Consume(executionID string, input io.Reader, output chan []byte) error
}

type csvParser struct {
	def    Schema
	l      *slog.Logger
	notify chan bool
}

func NewCSVParser(schemaPath string, notify chan bool, l *slog.Logger) Parser {
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
		def:    def,
		l:      l,
		notify: notify,
	}
}
func (p *csvParser) Consume(executionId string, input io.Reader, output chan []byte) error {
	defer close(output)

	p.l.Info("starting to consume csv file")
	reader := csv.NewReader(input)

	if p.def.SkipHeader {
		p.l.Info("skipping header")
		if _, err := reader.Read(); err != nil {
			return err
		}
	}

	line := uint64(1)
	offset := uint64(0)
	for {
		record, err := reader.Read()
		if err == io.EOF {
			p.l.Info("reached end of file")
			break
		}
		if err != nil {
			p.signalError()
			p.l.Error("error reading record", formatError(line, offset, err))
			continue
		}

		validatedRecord, size, err := validateAndTransformRecord(record, p.def)
		if err != nil {
			p.signalError()
			p.l.Error("error validating record", formatError(line, offset, err))
			continue
		}

		offset += uint64(size)
		line++

		r := Record{
			Line:        line,
			Offset:      offset,
			ExecutionId: executionId,
			Values:      validatedRecord,
		}

		rbytes, err := json.Marshal(r)
		if err != nil {
			p.notify <- false
			p.l.Error("error marshalling record", formatError(line, offset, err))
			continue
		}

		select {
		case output <- rbytes:
		case <-time.After(1 * time.Second):
			err = fmt.Errorf("timeout sending record to output channel")
			p.l.Error("error sending record to output channel", formatError(line, offset, err))
		}
	}

	p.l.Info("finished consuming csv file")
	p.l.Info(fmt.Sprintf("Validated %d records\n", line))
	p.signalSuccess()
	return nil
}

func (p *csvParser) signalError() {
	if !nonBlockingSend(p.notify, false) {
		p.l.Error("error sending signal to notify channel")
	}
}

func (p *csvParser) signalSuccess() {
	if !nonBlockingSend(p.notify, true) {
		p.l.Error("error sending signal to notify channel")
	}
}

// TODO: review these 2 fucntions, they kinda suck
func validateAndTransformRecord(record []string, def Schema) (map[string]interface{}, int, error) {
	validatedRecord := make(map[string]interface{})
	if len(record) != len(def.Columns) {
		return nil, -1, fmt.Errorf("incorrect number of columns: got %d, want %d", len(record), len(def.Columns))
	}

	size := 0
	for i, value := range record {
		size += len(value)
		colDef := def.Columns[i]
		if value == "" && colDef.Required {
			return nil, -1, fmt.Errorf("column %s is required but empty", colDef.Name)
		}

		if validator, exists := validatorsMap[colDef.Type]; exists {
			transformedValue, err := validateAndTransform(value, validator)
			if err != nil {
				return nil, -1, fmt.Errorf("column %s: %s", colDef.Name, err.Error())
			}
			// Assign the transformed value to the corresponding column name
			validatedRecord[colDef.Name] = transformedValue
		} else {
			return nil, -1, fmt.Errorf("unsupported column type: %s", colDef.Type)
		}
	}

	return validatedRecord, size, nil
}

func nonBlockingSend(ch chan<- bool, data bool) bool {
	select {
	case ch <- data:
		return true
	default:
		return false
	}
}

func formatError(line, offset uint64, err error) error {
	return fmt.Errorf("line %d, offset %d: %s", line, offset, err.Error())
}
