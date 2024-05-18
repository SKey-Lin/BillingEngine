package entities

import (
	"github.com/alpacahq/alpacadecimal"
	"gorm.io/gorm"
)

type Loan struct {
	gorm.Model
	Amount   alpacadecimal.Decimal `gorm:"type:bigint"`
	Duration int32
	Status   string // Assuming status will always APPROVED

	BorrowerID uint
	Borrower   Borrower `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`

	RepaymentSchedule []RepaymentSchedule `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
