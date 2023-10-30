package business

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"github.com/Tango-Rocker/batch-challenge/csv"
	"html/template"
	"log"
	"log/slog"
	"os"
	"time"
)

//go:embed *.html
var mailTemplateBytes []byte

// Summary holds all the summary information.
type Summary struct {
	TotalBalance              float64
	TransactionsByMonth       map[string]int
	TotalDebitByMonth         map[string]float64
	DebitTransactionsByMonth  map[string]int
	TotalCreditByMonth        map[string]float64
	CreditTransactionsByMonth map[string]int
}

// NewSummary creates a new instance of Summary with initialized maps.
func NewSummary() *Summary {
	return &Summary{
		TransactionsByMonth:       make(map[string]int),
		TotalDebitByMonth:         make(map[string]float64),
		DebitTransactionsByMonth:  make(map[string]int),
		TotalCreditByMonth:        make(map[string]float64),
		CreditTransactionsByMonth: make(map[string]int),
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
			month := extractMonth(trx.Date)
			sum.TransactionsByMonth[month]++
			if trx.Amount < 0 {
				sum.TotalDebitByMonth[month] += trx.Amount
				sum.DebitTransactionsByMonth[month]++
			} else {
				sum.TotalCreditByMonth[month] += trx.Amount
				sum.CreditTransactionsByMonth[month]++
			} // Simplify map operations with compound assignment.
		}
	}

	// After processing all data and before sending the email, check for the success signal.
	srv.l.Info("Waiting for success signal")
	select {
	case <-srv.successSignal:
		srv.l.Info("Sending email")
		srv.mail.Send(Mail{
			To:      findMail("stub"),
			Subject: "Stori Account Balance",
			Body:    serializeSummary(sum),
		})
	case <-time.After(time.Second * 30): // Adding a timeout to avoid hanging indefinitely.
		srv.l.Error("Timeout waiting for success signal, aborting process")
	}
}

func serializeSummary(sum *Summary) string {
	t, err := template.New("email").Funcs(template.FuncMap{"printf": fmt.Sprintf}).Parse(string(mailTemplateBytes))
	if err != nil {
		log.Println("Error parsing template:", err)
		return ""
	}

	// Calculate the average debit and credit
	var totalDebitAmount, totalCreditAmount float64
	var debitCount, creditCount int
	for _, amount := range sum.TotalDebitByMonth {
		totalDebitAmount += amount
		debitCount++
	}
	for _, amount := range sum.TotalCreditByMonth {
		totalCreditAmount += amount
		creditCount++
	}

	// Avoid division by zero
	avgDebit := 0.0
	if debitCount > 0 {
		avgDebit = totalDebitAmount / float64(debitCount)
	}
	avgCredit := 0.0
	if creditCount > 0 {
		avgCredit = totalCreditAmount / float64(creditCount)
	}

	// Prepare the data for the template
	data := struct {
		*Summary
		AvgDebit  float64
		AvgCredit float64
	}{
		Summary:   sum,
		AvgDebit:  avgDebit,
		AvgCredit: avgCredit,
	}

	var ret bytes.Buffer
	if err := t.Execute(&ret, data); err != nil {
		log.Println("Error executing template:", err)
		return ""
	}

	return ret.String()
}
func extractMonth(dateStr string) string {
	// Parse the date string into a time.Time object.
	date, err := time.Parse("2006-01-02", dateStr)
	if err != nil {
		return "temporally unavailable:"
	}
	// Return the full month name.
	return date.Month().String()
}

func findMail(id string) string {
	return os.Getenv("TARGET_MAIL") //stub
}
