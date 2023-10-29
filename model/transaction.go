package model

type Transaction struct {
	Id         int     `json:"id"`
	ExternalId int     `json:"external_id"`
	BatchId    string  `json:"batch_id"`
	Line       int     `json:"line"`
	Offset     int     `json:"offset"`
	Amount     float64 `json:"amount"`
	Date       string  `json:"date"`
}
