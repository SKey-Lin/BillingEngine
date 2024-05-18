package entities

import (
	"time"

	"github.com/alpacahq/alpacadecimal"
	"gorm.io/gorm"
)

type RepaymentSchedule struct {
	gorm.Model
	LoanID        uint
	Loan          Loan                  `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Outstanding   alpacadecimal.Decimal `gorm:"type:bigint"`
	Amount        alpacadecimal.Decimal `gorm:"type:bigint"`
	Paid          alpacadecimal.Decimal `gorm:"type:bigint"`
	Installment   int32
	ScheduledDate time.Time
}
