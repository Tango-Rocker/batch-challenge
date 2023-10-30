package business

import "log/slog"

type Summary struct {
	ExecutionId   string
	AccountId     string
	TotalBalance  float64
	TotalRecords  int
	AmountByMonth map[string]float64
}

type SummaryService struct {
	mail *EmailService
	l    *slog.Logger
}

func NewSummaryService(mail *EmailService, l *slog.Logger) *SummaryService {
	return &SummaryService{
		mail: mail,
		l:    l,
	}
}

func
