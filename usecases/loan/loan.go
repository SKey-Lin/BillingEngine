package loan

import (
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/alpacahq/alpacadecimal"
	"gorm.io/gorm"
	"squalux.com/skey/lending/entities"
	"squalux.com/skey/lending/handler/dto/request"
	"squalux.com/skey/lending/models"
)

const APPROVED = "APPROVED"

var INTEREST = alpacadecimal.NewFromFloat(0.1) // assumming interest always constant and flat 10%
var YEAR_IN_WEEKS = alpacadecimal.NewFromInt32(52)

func toStartOfWeek(date time.Time) time.Time {
	for date.Weekday() != time.Monday { // start of week from Monday
		date = date.AddDate(0, 0, -1)
	}

	return date
}

func createRepayment(loan entities.Loan) ([]entities.RepaymentSchedule, error) {
	duration := alpacadecimal.NewFromInt32(loan.Duration)

	year := alpacadecimal.NewFromFloat32(1)
	if duration.GreaterThan(YEAR_IN_WEEKS) {
		year = duration.Div(YEAR_IN_WEEKS).Ceil()
	}

	interestVal := loan.Amount.Mul(INTEREST).Mul(year)

	totalLoan := interestVal.Add(loan.Amount)
	repaymentPerWeek := totalLoan.Div(duration).Round(0)

	// Assumming Effective Repayment are the following week
	paymentDate := toStartOfWeek(time.Now()).Truncate(24 * time.Hour)

	repaymentChan := make(chan entities.RepaymentSchedule, loan.Duration)
	wg := new(sync.WaitGroup)
	for i := range loan.Duration {
		paymentDate = paymentDate.AddDate(0, 0, 7)
		wg.Add(1)

		go func(i int32, payDate time.Time) {
			defer wg.Done()

			var repayment entities.RepaymentSchedule
			weekCount := i

			repayment.Installment = weekCount + 1
			repayment.Paid = alpacadecimal.Zero
			repayment.LoanID = loan.ID

			repayment.ScheduledDate = payDate

			repayment.Outstanding = totalLoan.Sub(repaymentPerWeek.Mul(alpacadecimal.NewFromInt32(weekCount)))

			if i == loan.Duration-1 {
				repayment.Amount = repayment.Outstanding
			} else {
				repayment.Amount = repaymentPerWeek
			}

			repaymentChan <- repayment
		}(i, paymentDate)
	}

	wg.Wait()
	close(repaymentChan)

	repaymentList := []entities.RepaymentSchedule{}
	for repayment := range repaymentChan {
		repaymentList = append(repaymentList, repayment)
	}

	sort.Slice(repaymentList, func(i, j int) bool {
		return repaymentList[i].Installment < repaymentList[j].Installment
	})

	if err := models.DB.Create(repaymentList).Error; err != nil {
		return repaymentList, err
	}

	return repaymentList, nil

}

func CreateLoan(body request.LoanBody) (entities.Loan, error) {
	var loan entities.Loan

	err := models.DB.Transaction(func(tx *gorm.DB) error {
		var borrower entities.Borrower

		if err := models.DB.First(&borrower, "name = ?", body.BorrowerName).Error; err != nil {
			borrower.Name = body.BorrowerName

			if err := models.DB.Create(&borrower).Error; err != nil {
				return err
			}
		}

		loan.BorrowerID = borrower.ID
		loan.Amount = body.Amount
		loan.Duration = body.Duration
		loan.Status = APPROVED
		loan.Borrower = borrower

		if err := models.DB.Create(&loan).Error; err != nil {
			return err
		}

		// Can be optimize to Move into Queue Async (event driven)
		if repayments, err := createRepayment(loan); err != nil {
			return err
		} else {
			loan.RepaymentSchedule = repayments
		}

		return nil
	})
	if err != nil {
		return loan, err
	}

	return loan, nil
}

func GetOutstanding(loanID uint) (int64, error) {
	// Check if loan exists
	if err := models.DB.First(&entities.Loan{}, loanID).Error; err != nil {
		return 0, err
	}

	var repayment entities.RepaymentSchedule

	if err := models.DB.Where(map[string]interface{}{
		"loan_id": loanID,
		"paid":    0,
	}).First(&repayment).Error; err != nil {
		return 0, err
	}

	return repayment.Outstanding.BigInt().Int64(), nil
}

func IsDelinquent(loanID uint) (bool, error) {
	// Check if loan exists
	if err := models.DB.First(&entities.Loan{}, loanID).Error; err != nil {
		return false, err
	}

	pastRepayment, err := getPastRepayment(loanID)
	if err != nil {
		return false, err
	}

	return len(pastRepayment) >= 2, nil
}

func getPastRepayment(loanID uint) ([]entities.RepaymentSchedule, error) {
	repayments := []entities.RepaymentSchedule{}

	if err := models.DB.Where("scheduled_date < ?", time.Now()).Where(map[string]interface{}{
		"loan_id": loanID,
		"paid":    0,
	}).Find(&repayments).Error; err != nil {
		return repayments, err
	}
	return repayments, nil
}

func MakePayment(payment request.PaymentBody) error {

	if payment.Amount.Equal(alpacadecimal.Zero) {
		return fmt.Errorf("Payment amount required a value")
	}

	// Check if loan exists
	if err := models.DB.First(&entities.Loan{}, payment.LoanID).Error; err != nil {
		return err
	}

	pastRepayment, err := getPastRepayment(payment.LoanID)
	if err != nil {
		return err
	}

	amountToPay := alpacadecimal.Zero
	if len(pastRepayment) > 0 {
		for _, delinquent := range pastRepayment {
			amountToPay = amountToPay.Add(delinquent.Amount)
		}
	}

	if alpacadecimal.Zero.Equal(amountToPay) {
		return fmt.Errorf("There are no payment yet for now")
	}

	if !payment.Amount.Equal(amountToPay) {
		return fmt.Errorf("Amount that paid not matching amount should be paid: %d", amountToPay.BigInt().Int64())
	}

	if len(pastRepayment) > 0 {
		var ids []uint
		for _, delinquent := range pastRepayment {
			ids = append(ids, delinquent.ID)
		}

		if err := models.DB.Model(&entities.RepaymentSchedule{}).Where("id IN ?", ids).Update("paid", gorm.Expr("amount")).Error; err != nil {
			return err
		}
	}

	return nil
}
