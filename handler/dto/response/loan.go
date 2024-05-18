package response

import (
	"github.com/alpacahq/alpacadecimal"
)

type LoanResp struct {
	Amount   alpacadecimal.Decimal `json:"amount"`
	Duration int                   `json:"duration"`
	Borrower struct {
		Name string `json:"name"`
	} `json:"borrower"`
}
