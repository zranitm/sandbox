package model

type TransactionResponse struct {
	Period           string        `json:"period"`
	TotalIncome      int           `json:"total_income"`
	TotalExpenditure int           `json:"total_expenditure"`
	Transactions     []Transaction `json:"transactions"`
}

type Transaction struct {
	Date    string `json:"date"`
	Amount  int    `json:"amount"`
	Content string `json:"content"`
}
