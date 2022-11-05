package models

type Transaction struct {
	Amount    string `json:"amount" xml:"amount" form:"amount"`
	Timestamp string `json:"timestamp" xml:"timestamp" form:"timestamp"`
}