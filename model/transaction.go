package model

type Transaction struct {
	Id          *int    `json:"id" gorm:"primaryKey;autoIncrement"`
	ExternalId  int     `json:"external_id"`
	ExecutionId string  `json:"execution_id"`
	Line        int     `json:"line"`
	Offset      int     `json:"offset"`
	Amount      float64 `json:"amount"`
	Date        string  `json:"date"`
}
