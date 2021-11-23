package controllers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/gregoflash05/gradely/models"
)

func TestDeposit(t *testing.T) {

	t.Run("test no user token", func(t *testing.T) {
		testData := models.DepositRequest{
			Amount: 100,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/deposit", Deposit).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/deposit", buf)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusUnauthorized)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "token Invalid")
	})

	t.Run("test invalid token", func(t *testing.T) {
		testData := models.DepositRequest{
			Amount: 100,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/deposit", Deposit).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/deposit", buf)
		req.Header.Add("Authorization", "Bearer ubuobda8buiwwr")

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusUnauthorized)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "token Invalid")
	})

	t.Run("test user not a buyer", func(t *testing.T) {
		testData := models.DepositRequest{
			Amount: 100,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/deposit", Deposit).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/deposit", buf)
		req.Header.Add("Authorization", "Bearer "+TestSToken)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusNotAcceptable)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "user is not a buyer")
	})

	t.Run("test amount deposited", func(t *testing.T) {
		testData := []byte(`{"amount": 60}`)
		buf := bytes.NewBuffer(testData)

		r := getRouter()
		r.HandleFunc("/v1/deposit", Deposit).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/deposit", buf)
		req.Header.Add("Authorization", "Bearer "+TestToken)

		response := getHTTPResponse(t, r, req)
		fmt.Println(response.Code)

		assertStatusCode(t, response.Code, http.StatusBadRequest)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "you can only deposit, 50, 100, 200, 500, 1000 coins")
	})

	t.Run("test deposited successfully", func(t *testing.T) {
		testData := models.DepositRequest{
			Amount: 100,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/deposit", Deposit).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/deposit", buf)
		req.Header.Add("Authorization", "Bearer "+TestToken)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, 200)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "deposit successful")
	})

}

func TestBuy(t *testing.T) {

	t.Run("test no user token", func(t *testing.T) {
		testData := models.BuyRequest{
			ProductID: 1,
			Quantity:  4,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/buy", BuyProduct).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/buy", buf)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusUnauthorized)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "token Invalid")
	})

	t.Run("test invalid token", func(t *testing.T) {
		testData := models.BuyRequest{
			ProductID: 1,
			Quantity:  4,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/buy", BuyProduct).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/buy", buf)
		req.Header.Add("Authorization", "Bearer ubuobda8buiwwr")

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusUnauthorized)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "token Invalid")
	})

	t.Run("test user not a buyer", func(t *testing.T) {
		testData := models.BuyRequest{
			ProductID: 1,
			Quantity:  4,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/buy", BuyProduct).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/buy", buf)
		req.Header.Add("Authorization", "Bearer "+TestSToken)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusNotAcceptable)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "user is not a buyer")
	})

	t.Run("test product not found", func(t *testing.T) {
		testData := models.BuyRequest{
			ProductID: 1000,
			Quantity:  4,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/buy", BuyProduct).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/buy", buf)
		req.Header.Add("Authorization", "Bearer "+TestToken)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusUnauthorized)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "product not found")
	})

	t.Run("test insuffient funds", func(t *testing.T) {
		testData := models.BuyRequest{
			ProductID: 1,
			Quantity:  47788,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/buy", BuyProduct).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/buy", buf)
		req.Header.Add("Authorization", "Bearer "+TestToken)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusNotAcceptable)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "insufficient funds")
	})

	t.Run("test purchase successful", func(t *testing.T) {
		testData := models.BuyRequest{
			ProductID: 1,
			Quantity:  1,
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/buy", BuyProduct).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/buy", buf)
		req.Header.Add("Authorization", "Bearer "+TestToken)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, 200)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "purchase successful")
	})

}
func TestProductCreate(t *testing.T) {

	t.Run("test no user token", func(t *testing.T) {
		testData := models.Product{
			Cost:        100,
			ProductName: "Test Product",
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/products", ProductCreate).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/products", buf)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusUnauthorized)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "token Invalid")
	})

	t.Run("test invalid token", func(t *testing.T) {
		testData := models.Product{
			Cost:        100,
			ProductName: "Test Product",
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/products", ProductCreate).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/products", buf)
		req.Header.Add("Authorization", "Bearer ubuobda8buiwwr")

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusUnauthorized)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "token Invalid")
	})

	t.Run("test request data", func(t *testing.T) {
		testData := []byte(`{"stem": "60"}`)
		buf := bytes.NewBuffer(testData)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/products", ProductCreate).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/products", buf)
		req.Header.Add("Authorization", "Bearer "+TestToken)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, http.StatusNotAcceptable)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "user is not a seller")
	})

	t.Run("test product added successfully", func(t *testing.T) {
		testData := models.Product{
			Cost:        100,
			ProductName: "Test Product",
		}
		buf := new(bytes.Buffer)
		json.NewEncoder(buf).Encode(testData)

		r := getRouter()
		r.HandleFunc("/v1/products", ProductCreate).Methods("POST")
		req, _ := http.NewRequest("POST", "/v1/products", buf)
		req.Header.Add("Authorization", "Bearer "+TestSToken)

		response := getHTTPResponse(t, r, req)

		assertStatusCode(t, response.Code, 200)
		assertResponseMessage(t, parseResponse(response)["message"].(string), "product added successfully")
	})

}
