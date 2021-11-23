package controllers

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/gorilla/mux"
	"github.com/gregoflash05/gradely/models"
	"github.com/gregoflash05/gradely/utils"
	"github.com/joho/godotenv"
)

var (
	TestBuyerEmail     = "test@gmail.com"
	TestBuyerUserName  = "test@gmail.com"
	TestsellerEmail    = "test1@gmail.com"
	TestsellerUserName = "test1@gmail.com"
	TestIsVerified     = true
	TestFullName       = "testing"
	TestPhone          = "09032094355"
	TestPassword       = "testing123"
	TestDepositAmount  = 1000
	TestProductId      uint
	TestToken          string
	TestSToken         string
)

func getRouter() *mux.Router {
	router := mux.NewRouter().StrictSlash(true)
	return router
}

// Helper function to process a request and test its response
func getHTTPResponse(t *testing.T, r *mux.Router, req *http.Request) *httptest.ResponseRecorder {

	// Create a response recorder
	w := httptest.NewRecorder()

	// Create the service and process the above request.
	r.ServeHTTP(w, req)

	return w
}

func assertStatusCode(t *testing.T, got, expected int) {
	if got != expected {
		t.Errorf("got status %d expected status %d", got, expected)
	}
}

func assertResponseMessage(t *testing.T, got, expected string) {
	if got != expected {
		t.Errorf("got message: %q expected: %q", got, expected)
	}
}

func parseResponse(w *httptest.ResponseRecorder) map[string]interface{} {
	res := make(map[string]interface{})
	json.NewDecoder(w.Body).Decode(&res)
	return res
}

func TestMain(m *testing.M) {
	// load .env file if it exists
	err := godotenv.Load("../.env")
	if err != nil {
		log.Fatalf("Error loading .env file: %v", err)
	}

	fmt.Println("Environment variables successfully loaded. Starting application...")

	_, err = utils.ConnectToDB(os.Getenv("TEST_SQL_DATABASE_URL"))
	if err != nil {
		log.Fatalf("Error loading .env file: %v", errors.New("could not connect to MySql DB"))
	}
	fmt.Println("database connected")

	utils.DropTables()
	fmt.Println("tables dropped successfully")
	utils.Migrate()
	fmt.Println("migration complete")

	TestToken, TestSToken, err = setUpUserAccount()
	if err != nil {
		log.Fatal(err.Error())
	}

	TestProductId, err = setupProduct()
	if err != nil {
		log.Fatal(err.Error())
	}

	exitVal := m.Run()

	// drop tables after running all tests
	utils.DropTables()
	fmt.Println("tables dropped successfully")

	os.Exit(exitVal)
}

func setUpUserAccount() (string, string, error) {
	pass, _ := GenerateHashPassword(TestPassword)
	buyerUser := models.User{
		Email:      TestBuyerEmail,
		IsVerified: TestIsVerified,
		FullName:   TestFullName,
		UserName:   TestBuyerUserName,
		Phone:      TestPhone,
		Password:   pass,
		Role:       "buyer",
		Deposit:    TestDepositAmount,
	}

	sellerUser := models.User{
		Email:      TestsellerEmail,
		IsVerified: TestIsVerified,
		FullName:   TestFullName,
		UserName:   TestsellerUserName,
		Phone:      TestPhone,
		Password:   pass,
		Role:       "seller",
		Deposit:    TestDepositAmount,
	}

	var checkUser models.User

	result := utils.GetItemsByField(&checkUser, "email", buyerUser.Email)
	if result.RowsAffected < 1 {
		res := utils.CreateItem(&buyerUser)
		if res.RowsAffected < 1 {
			return "", "", fmt.Errorf("Buyer account not created")
		}
	}

	result = utils.GetItemsByField(&checkUser, "email", sellerUser.Email)
	if result.RowsAffected < 1 {
		res := utils.CreateItem(&sellerUser)
		if res.RowsAffected < 1 {
			return "", "", fmt.Errorf("seller account not created")
		}
	}

	btoken, err := CreateToken(strconv.FormatUint(uint64(buyerUser.ID), 10))
	if err != nil {
		return "", "", fmt.Errorf("error generating token")
	}
	stoken, err := CreateToken(strconv.FormatUint(uint64(sellerUser.ID), 10))
	if err != nil {
		return "", "", fmt.Errorf("error generating token")
	}

	return btoken, stoken, nil
}

func setupProduct() (uint, error) {

	var checkUser models.User

	result := utils.GetItemsByField(&checkUser, "email", TestsellerEmail)
	if result.RowsAffected < 1 {
		return 0, fmt.Errorf("seller does not exist")
	}

	product := models.Product{
		Cost:        50,
		ProductName: "Test Product",
		SellerId:    checkUser.ID,
	}

	res := utils.CreateItem(&product)
	if res.RowsAffected < 1 {
		return 0, fmt.Errorf("product not created")
	}

	return product.ID, nil
}
