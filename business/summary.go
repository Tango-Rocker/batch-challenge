package business

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/Tango-Rocker/batch-challenge/csv"
	"log/slog"
	"time"
)

type Summary struct {
	ExecutionId   string
	AccountId     string
	TotalBalance  float64
	TotalRecords  int
	AmountByMonth map[string]float64
}

func NewSummary() *Summary {
	return &Summary{
		AmountByMonth: make(map[string]float64),
	}
}

type SummaryService struct {
	mail          *EmailService
	l             *slog.Logger
	input         chan []byte
	running       bool
	successSignal chan bool
}

func NewSummaryService(mail *EmailService, ssgnl chan bool, l *slog.Logger) *SummaryService {
	return &SummaryService{
		mail:          mail,
		l:             l.With(slog.String("service", "summary")),
		input:         make(chan []byte, 1000),
		running:       false,
		successSignal: ssgnl,
	}
}

func (srv *SummaryService) GetInputChannel() chan []byte {
	return srv.input
}

func (srv *SummaryService) Launch(ctx context.Context) {
	if !srv.running {
		go srv.run(ctx)
		srv.running = true
	}
}

func (srv *SummaryService) run(ctx context.Context) {
	srv.l.Info("Starting SummaryService")
	sum := NewSummary()
	for srv.running {
		select {
		case <-ctx.Done():
			srv.running = false
			srv.l.Info("Stopping SummaryService")

		case data, ok := <-srv.input:
			if !ok {
				// If the channel is closed, break the loop to proceed to send the email.
				srv.running = false
				srv.l.Info("Stopping SummaryService")
				break
			}

			record := &csv.Record{}
			err := json.Unmarshal(data, record)
			if err != nil {
				srv.l.Error("Error unmarshalling data, aborting summary")
				continue // Use continue to skip this iteration and try the next piece of data.
			}

			trx, err := mapRecordToTransaction(record)
			if err != nil {
				srv.l.Error("Error mapping record to transaction, aborting summary")
				continue // Same here, use continue to move on to the next record.
			}

			sum.TotalBalance += trx.Amount
			sum.TotalRecords++

			month := extractMonth(trx.Date)
			sum.AmountByMonth[month] += trx.Amount // Simplify map operations with compound assignment.
		}
	}

	// After processing all data and before sending the email, check for the success signal.
	srv.l.Info("Waiting for success signal")
	select {
	case <-srv.successSignal:
		srv.l.Info("Sending email")
		srv.mail.Send(Mail{
			To:      findMail(sum.AccountId),
			Subject: "Stori Account Balance",
			Body:    serializeSummary(sum),
		})
	case <-time.After(time.Second * 30): // Adding a timeout to avoid hanging indefinitely.
		srv.l.Error("Timeout waiting for success signal, aborting process")
	}
}

var mailTemplate = `
	Dear customer,
	Here is your card summary for theses months:

	Total Balance: %f
	Total Transactions: %d	
`

var byMothTemplate = `
	%s: %f
`

// serializeSummary converts a Summary struct to human readable mail template
func serializeSummary(sum *Summary) string {
	ret := ""

	ret += fmt.Sprintf(mailTemplate, sum.TotalBalance, sum.TotalRecords)

	for k, v := range sum.AmountByMonth {
		ret += fmt.Sprintf(byMothTemplate, k, v)
	}

	return ret
}

func extractMonth(date string) string {
	return "September" //stub
}

func findMail(id string) string {
	return "alejandroevidal1@gmail.com" //stub
}
