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

	rAlpha.POST("/login", func(c *gin.Context) {
		var credentials models.RequestLogin

		err := c.BindJSON(&credentials)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
			return
		}

		response, err := sendRequestToOmega("/api/Authorization/login", credentials)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Omega"})
			return
		}
		token = response.Results.ResultLogin.Token
		c.JSON(http.StatusOK, response)

	})

	routes := []string{"/check", "/pay", "/get-status", "/reverse", "/p2pinfowrap"}
	for _, route := range routes {
		rAlpha.POST(route, func(c *gin.Context) {
			var req interface{}

			err := c.BindJSON(&req)
			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
				return
			}

			response, err := sendRequestToOmega("/api/Authorization"+route, req)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send data to Omega"})
				return
			}
			fmt.Println(&response)
			c.JSON(http.StatusOK, response)
		})
	}

	err := rAlpha.Run(":8080")
	if err != nil {
		log.Println("Error while running a server")
		return
	}
}

func sendRequestToOmega(endpoint string, data interface{}) (models.Response, error) {
	omegaURL := "http://localhost:8081" + endpoint

	omegaJSON, err := json.Marshal(data)
	if err != nil {
		fmt.Println("Error marshaling Omega data:", err)
		return models.Response{}, err
	}

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
			log.Println("error while closing")
		}
	}(resp.Body)

	var response models.Response
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return models.Response{}, err
	}

	if err := json.Unmarshal(bodyBytes, &response); err != nil {
		return models.Response{}, err
	}
	fmt.Println(&response)
	return response, nil
}
