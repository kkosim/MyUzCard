package main

import (
	"MyUzcardTransfer/models"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strings"
	"time"
)

var token string

func main() {
	rOmega := gin.Default()

	authMiddleware := func(c *gin.Context) {
		header := c.GetHeader("Authorization")
		headerParts := strings.Split(header, " ")
		if len(headerParts) != 2 || headerParts[1] != token {
			c.JSON(http.StatusBadRequest, models.Response{
				Results: nil,
				Error: &models.Error{
					Code:    -100,
					Message: "internal server error",
				},
				Success: false,
			})
			c.Abort()
			return
		}
	}

	rOmega.POST("/api/Authorization/login", func(c *gin.Context) {
		var loginData models.RequestLogin
		if err := c.ShouldBindJSON(&loginData); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if loginData.Username == "test" && loginData.Password == "test@123" {
			token = "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
			c.JSON(http.StatusOK, models.Response{
				Results: &models.Results{
					ResultLogin: &models.ResultLogin{
						Token:      token,
						UserID:     26,
						ExpireDate: "2022-06-06T10:20:59.1529539",
						RoleID:     2,
					},
				},
				Error:   nil,
				Success: true,
			})
		} else {
			c.JSON(http.StatusBadRequest, models.Response{
				Results: nil,
				Error: &models.Error{
					Code:    -2,
					Message: "User not found",
				},
				Success: false,
			})
		}
	})

	secured := rOmega.Group("/api/Authorization", authMiddleware)
	{
		secured.POST("/check", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Results: &models.Results{
					ResultCheck: &models.ResultCheck{
						TransactionNumber: "11112222070580de84dbf6-8183-412c-9ba8-06abb4f96c6c",
						State:             0,
						CreatedDate:       "2022-11-11T22:07:06.1931708",
						Amount:            1000,
						CommissionAmount:  2,
						TotalAmount:       1002,
						CardHolder:        "MIRKURBONOV MIRGANI GULO",
						BankName:          "IpakYuliBank",
						CurrencyInfo: &models.CurrencyInfo{
							From: "USD",
							To:   "UZS",
							Rate: 11306.79,
						},
					},
				},
				Error:   nil,
				Success: true,
			})
		})

		secured.POST("/pay", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Results: &models.Results{
					ResultPay: &models.ResultPay{
						TransactionNumber: "cc6895851c10b9386c2418008b02af39",
						State:             1,
						CreatedDate:       "28.10.2020 10:00:00",
						PaymentDate:       "28.10.2020 10:01:58",
						Amount:            12.3,
						TotalAmount:       12,
						CommissionAmount:  1,
					},
				},
				Error:   nil,
				Success: true,
			})
		})

		secured.POST("/get-status", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Results: &models.Results{
					ResultGetStatus: &models.ResultGetStatus{
						State:            1,
						CreatedDate:      "28.10.2020 10:00:00",
						PaymentDate:      "28.10.2020 10:01:58",
						Amount:           12.3,
						TotalAmount:      12,
						CommissionAmount: 1,
					},
				},
				Error:   nil,
				Success: true,
			})
		})

		secured.POST("/reverse", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Results: &models.Results{
					ResultReverse: &models.ResultReverse{
						State:            1,
						CreatedDate:      "28.10.2020 10:00:00",
						PaymentDate:      "28.10.2020 10:01:58",
						Amount:           12.3,
						TotalAmount:      12,
						CommissionAmount: 1,
						ReverseDate:      time.Now(),
					},
				},
				Error:   nil,
				Success: true,
			})
		})

		secured.POST("/p2pinfowrap", func(c *gin.Context) {
			c.JSON(http.StatusOK, models.Response{
				Results: &models.Results{
					ResultP2pinfowrap: &models.ResultP2pinfowrap{
						CardHolder: "Mirgani Mirkurbonov",
						BankName:   "Ipak Yuli Bank",
					},
				},
				Error:   nil,
				Success: true,
			})
		})
	}

	log.Println("Starting server on port 8081...")
	err := rOmega.Run(":8081")
	if err != nil {
		log.Println("Error while running a server")
		return
	}
}
