package main

import (
	"github.com/gin-gonic/gin"

	"squalux.com/skey/lending/handler/rest"
)

func main() {
	// gin.SetMode(gin.ReleaseMode)
	r := gin.New()

	r.POST("/loan", rest.CreateLoan)

	r.GET("/outstanding", rest.GetOutstanding)

	r.POST("/payment", rest.MakePayment)

	r.GET("/delinquent", rest.IsDelinquent)

	r.Run(":8000")
}
