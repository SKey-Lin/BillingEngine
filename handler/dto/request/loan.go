package request

import "github.com/alpacahq/alpacadecimal"

type LoanBody struct {
	Amount       alpacadecimal.Decimal `json:"amount"`
	Duration     int32                 `json:"duration"`
	BorrowerName string                `json:"borrower_name"`
}

type PaymentBody struct {
	Amount alpacadecimal.Decimal `json:"amount"`
	LoanID uint                  `json:"loan_id"`
}
