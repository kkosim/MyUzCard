package main

import (
	"MyUzcardTransfer/models"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

var token string

func main() {
	rAlpha := gin.Default()

	// Обработка запроса от API Alpha
	rAlpha.POST("/login", func(c *gin.Context) {
		var credentials models.RequestLogin

		err := c.BindJSON(&credentials)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		response, err := sendLoginToOmega(credentials)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Omega"})
			return
		}
		token = response.Results.ResultLogin.Token
		c.JSON(http.StatusOK, response)

	})

	rAlpha.POST("/check", func(c *gin.Context) {
		var req models.RequestCheck

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		response, err := sendCheckToOmega(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Omega"})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	rAlpha.POST("/pay", func(c *gin.Context) {
		var req models.TransactionNumber

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		response, err := sendPayToOmega(req)
		if err != nil {
			fmt.Printf("zat", err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Omega"})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	rAlpha.POST("/get-status", func(c *gin.Context) {
		var req models.TransactionNumber

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		response, err := sendGetStatusToOmega(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Omega"})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	rAlpha.POST("/reverse", func(c *gin.Context) {
		var req models.TransactionNumber

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		response, err := sendReverseToOmega(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Omega"})
			return
		}
		c.JSON(http.StatusOK, response)
	})

	rAlpha.POST("/p2pinfowrap", func(c *gin.Context) {
		var req models.CardNumber

		err := c.BindJSON(&req)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		response, err := sendP2pinfowrapToOmega(req)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Omega"})
			return
		}
		c.JSON(http.StatusOK, response)
	})
	err := rAlpha.Run(":8080")
	if err != nil {
		log.Println("Error while running a server")
		return
	}
}

func sendLoginToOmega(alphaData models.RequestLogin) (models.Response, error) {

	omegaData := models.RequestLogin{
		Username: alphaData.Username,
		Password: alphaData.Password,
	}
	omegaJSON, err := json.Marshal(omegaData)
	if err != nil {
		fmt.Println("Error marshaling Omega data:", err)
		return models.Response{}, nil
	}
	omegaURL := "http://localhost:8081/api/Authorization/login"

	client := &http.Client{}
	req, err := http.NewRequest("POST", omegaURL, bytes.NewBuffer(omegaJSON))
	if err != nil {
		return models.Response{}, err
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		return models.Response{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error while closing a ?")
		}
	}(resp.Body)

	var response models.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return models.Response{}, err
	}

	return response, nil
}

func sendCheckToOmega(alphaData models.RequestCheck) (models.Response, error) {

	omegaData := models.RequestCheck{
		Amount:                   alphaData.Amount,
		CardNumber:               alphaData.CardNumber,
		Currency:                 alphaData.Currency,
		PartnerClientId:          alphaData.PartnerClientId,
		PartnerTransactionNumber: alphaData.PartnerTransactionNumber,
		SenderFirstName:          alphaData.SenderFirstName,
		SenderMiddleName:         alphaData.SenderMiddleName,
		SenderLastName:           alphaData.SenderLastName,
		SenderAddress:            alphaData.SenderAddress,
		Sender: &models.Sender{
			CardNumber: alphaData.Sender.CardNumber,
			CardIssuer: alphaData.Sender.CardIssuer,
			Country:    alphaData.Sender.Country,
			SenderDocument: &models.SenderDocument{
				Nationality: alphaData.Sender.SenderDocument.Nationality,
				DocType:     alphaData.Sender.SenderDocument.DocType,
				DocSerial:   alphaData.Sender.SenderDocument.DocSerial,
				DocNumber:   alphaData.Sender.SenderDocument.DocNumber,
				ValidTo:     alphaData.Sender.SenderDocument.ValidTo,
				BirthDate:   alphaData.Sender.SenderDocument.BirthDate,
			},
		},
	}
	omegaJSON, err := json.Marshal(omegaData)
	if err != nil {
		fmt.Println("Error marshaling Omega data:", err)
		return models.Response{}, nil
	}
	omegaURL := "http://localhost:8081/api/Authorization/check"

	client := &http.Client{}
	req, err := http.NewRequest("POST", omegaURL, bytes.NewBuffer(omegaJSON))
	if err != nil {
		return models.Response{}, err
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	resp, err := client.Do(req)
	if err != nil {
		return models.Response{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
			log.Println("error while closing a ?")
		}
	}(resp.Body)

	var response models.Response
	if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
		return models.Response{}, err
	}

	return response, nil
}

func sendPayToOmega(alphaData models.TransactionNumber) (models.Response, error) {

	omegaData := models.TransactionNumber{
		TransactionNumber: alphaData.TransactionNumber,
	}
	omegaJSON, err := json.Marshal(omegaData)
	if err != nil {
		fmt.Println("Error marshaling Omega data:", err)
		return models.Response{}, nil
	}
	omegaURL := "http://localhost:8081/api/Authorization/pay"

	client := &http.Client{}

	req, err := http.NewRequest("POST", omegaURL, bytes.NewBuffer(omegaJSON))
	if err != nil {
		return models.Response{}, err
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json; charset=utf-8"},
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Bearer " + token},
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.Response{}, err
	}

	defer resp.Body.Close()

	var response models.Response

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, err
	}

	if err = json.Unmarshal(bodyBytes, &response); err != nil {
		return models.Response{}, err
	}

	return response, nil
}

func sendGetStatusToOmega(alphaData models.TransactionNumber) (models.Response, error) {

	omegaData := models.TransactionNumber{
		TransactionNumber: alphaData.TransactionNumber,
	}
	omegaJSON, err := json.Marshal(omegaData)
	if err != nil {
		fmt.Println("Error marshaling Omega data:", err)
		return models.Response{}, nil
	}
	omegaURL := "http://localhost:8081/api/Authorization/get-status"

	client := &http.Client{}

	req, err := http.NewRequest("POST", omegaURL, bytes.NewBuffer(omegaJSON))
	if err != nil {
		return models.Response{}, err
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json; charset=utf-8"},
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Bearer " + token},
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.Response{}, err
	}

	defer resp.Body.Close()

	var response models.Response

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, err
	}

	if err = json.Unmarshal(bodyBytes, &response); err != nil {
		return models.Response{}, err
	}

	return response, nil
}

func sendReverseToOmega(alphaData models.TransactionNumber) (models.Response, error) {

	omegaData := models.TransactionNumber{
		TransactionNumber: alphaData.TransactionNumber,
	}
	omegaJSON, err := json.Marshal(omegaData)
	if err != nil {
		fmt.Println("Error marshaling Omega data:", err)
		return models.Response{}, nil
	}
	omegaURL := "http://localhost:8081/api/Authorization/reverse"

	client := &http.Client{}

	req, err := http.NewRequest("POST", omegaURL, bytes.NewBuffer(omegaJSON))
	if err != nil {
		return models.Response{}, err
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json; charset=utf-8"},
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Bearer " + token},
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.Response{}, err
	}

	defer resp.Body.Close()

	var response models.Response

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, err
	}

	if err = json.Unmarshal(bodyBytes, &response); err != nil {
		return models.Response{}, err
	}

	return response, nil
}

func sendP2pinfowrapToOmega(alphaData models.CardNumber) (models.Response, error) {

	omegaData := models.CardNumber{
		CardNumber: alphaData.CardNumber,
	}
	omegaJSON, err := json.Marshal(omegaData)
	if err != nil {
		fmt.Println("Error marshaling Omega data:", err)
		return models.Response{}, nil
	}
	omegaURL := "http://localhost:8081/api/Authorization/p2pinfowrap"

	client := &http.Client{}

	req, err := http.NewRequest("POST", omegaURL, bytes.NewBuffer(omegaJSON))
	if err != nil {
		return models.Response{}, err
	}

	req.Header = http.Header{
		"Content-Type":  []string{"application/json; charset=utf-8"},
		"Accept":        []string{"application/json"},
		"Authorization": []string{"Bearer " + token},
	}

	resp, err := client.Do(req)
	if err != nil {
		return models.Response{}, err
	}

	defer resp.Body.Close()

	var response models.Response

	var bodyBytes []byte
	bodyBytes, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, err
	}

	if err = json.Unmarshal(bodyBytes, &response); err != nil {
		return models.Response{}, err
	}

	return response, nil
}
