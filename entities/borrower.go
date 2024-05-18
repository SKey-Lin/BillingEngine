package entities

import "gorm.io/gorm"

type Borrower struct {
	gorm.Model
	Name string
	Loan []Loan `gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
}
