package business

import (
	"encoding/json"
	"fmt"
	"github.com/Tango-Rocker/batch-challenge/csv"
	"github.com/Tango-Rocker/batch-challenge/model"
)

type WriterConfig struct {
	BufferSize   int `env:"BUFFER_SIZE" envDefault:"1000"`
	FlushTimeout int `env:"FLUSH_TIMEOUT_MS" envDefault:"1000"`
}

type MailConfig struct {
	Host     string `env:"MAIL_SERVER_HOST"`
	Port     int    `env:"MAIL_SERVER_PORT"`
	Account  string `env:"MAIL_ACCOUNT"`
	Password string `env:"MAIL_PASSWORD"`
}

func mapRecordToTransaction(record *csv.Record) (*model.Transaction, error) {
	valuesJSON, err := json.Marshal(record.Values)
	if err != nil {
		return nil, fmt.Errorf("error marshaling values: %v", err)
	}

	var transaction model.Transaction
	err = json.Unmarshal(valuesJSON, &transaction)
	if err != nil {
		return nil, fmt.Errorf("error unmarshaling into Transaction: %v", err)
	}

	//TODO: we need to find a better way to handle this
	transaction.ExternalId = *transaction.Id
	transaction.Id = nil

	transaction.ExecutionId = record.ExecutionId
	transaction.Line = int(record.Line)     // Converting uint64 to int
	transaction.Offset = int(record.Offset) // Converting uint64 to int

	return &transaction, nil
}
