package rest

import (
	"net/http"
	"strconv"

	dto "github.com/dranikpg/dto-mapper"
	"github.com/gin-gonic/gin"

	"squalux.com/skey/lending/handler/dto/request"
	"squalux.com/skey/lending/handler/dto/response"
	"squalux.com/skey/lending/usecases/loan"
)

func CreateLoan(c *gin.Context) {
	var body request.LoanBody
	var resp response.LoanResp

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	loan, err := loan.CreateLoan(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	dto.Map(&resp, loan)
	c.JSON(http.StatusOK, resp)
}

func GetOutstanding(c *gin.Context) {
	paramLoanID := c.Query("loan_id")

	loanID, err := strconv.ParseUint(paramLoanID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	outstanding, err := loan.GetOutstanding(uint(loanID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"outstanding": outstanding,
	})
}

func MakePayment(c *gin.Context) {
	var body request.PaymentBody

	err := c.BindJSON(&body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	err = loan.MakePayment(body)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status": true,
	})
}

func IsDelinquent(c *gin.Context) {
	paramLoanID := c.Query("loan_id")

	loanID, err := strconv.ParseUint(paramLoanID, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": err.Error(),
		})
		return
	}

	isDelinquent, err := loan.IsDelinquent(uint(loanID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"is_delinquent": isDelinquent,
	})
}
