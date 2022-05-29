package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/femibiwoye/go-test/models"
	"github.com/femibiwoye/go-test/utils"
)

var (
	possibleDepositAmounts = []int{5, 10, 20, 50, 100}
)

// Deposit handles the deposit request. It checks if the user has enough money to buy the product.
func Deposit(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintID, _ := (strconv.ParseUint(userID, 10, 64))

	var user models.User

	utils.GetItemByPrimaryKey(&user, uint(uintID))

	if strings.ToLower(user.Role) != "buyer" {
		utils.GetError(fmt.Errorf("user is not a buyer"), http.StatusNotAcceptable, response)
		return
	}

	var depositRequest models.DepositRequest
	utils.ParseJSONFromRequest(request, &depositRequest)

	if !Contains(depositRequest.Amount, possibleDepositAmounts) {
		utils.GetError(errors.New("you can only deposit, 5, 10, 20, 50, 100 coins"), http.StatusBadRequest, response)
		return
	}

	updateMap := map[string]interface{}{}
	updateMap["deposit"] = depositRequest.Amount + user.Deposit

	result := utils.Db.Table("users").Where("id = ?", uint(uintID)).Updates(updateMap)

	if result.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("deposit failed"), http.StatusInternalServerError, response)
		return
	}

	utils.GetSuccess("deposit successful", nil, response)
}

// DepositReset handles the reset request.
func DepositReset(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintID, _ := (strconv.ParseUint(userID, 10, 64))

	var user models.User

	tx := utils.GetItemByPrimaryKey(&user, uint(uintID))
	if tx.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("user not found"), http.StatusUnauthorized, response)
		return
	}

	if strings.ToLower(user.Role) != "buyer" {
		utils.GetError(fmt.Errorf("user is not a buyer"), http.StatusNotAcceptable, response)
		return
	}

	updateMap := map[string]interface{}{}
	updateMap["deposit"] = 0

	result := utils.Db.Table("users").Where("id = ?", uint(uintID)).Updates(updateMap)

	if result.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("reset failed"), http.StatusInternalServerError, response)
		return
	}

	utils.GetSuccess("Reset successful", nil, response)
}

// BuyProduct handles the buy product request. It checks if the user has enough money to buy the product.
// TODO: add validation
// TODO: add transaction
func BuyProduct(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintID, _ := (strconv.ParseUint(userID, 10, 64))

	var user models.User

	utils.GetItemByPrimaryKey(&user, uint(uintID))

	if strings.ToLower(user.Role) != "buyer" {
		utils.GetError(fmt.Errorf("user is not a buyer"), http.StatusNotAcceptable, response)
		return
	}

	var buyRequest models.BuyRequest
	utils.ParseJSONFromRequest(request, &buyRequest)

	var product models.Product

	tx := utils.GetItemByPrimaryKey(&product, uint(buyRequest.ProductID))
	if tx.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("product not found"), http.StatusUnauthorized, response)
		return
	}

	totalCost := product.Cost * buyRequest.Quantity

	if totalCost > user.Deposit {
		utils.GetError(fmt.Errorf("insufficient funds"), http.StatusNotAcceptable, response)
		return
	}

	newBalance := user.Deposit - totalCost

	updateMap := map[string]interface{}{}
	updateMap["deposit"] = newBalance

	result := utils.Db.Table("users").Where("id = ?", uint(uintID)).Updates(updateMap)

	if result.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("purchase failed, try again latyer"), http.StatusInternalServerError, response)
		return
	}

	buyResponse := models.BuyResponse{
		ProductID:         buyRequest.ProductID,
		QuantityPurchased: buyRequest.Quantity,
		AmountSpent:       totalCost,
		Change:            newBalance,
	}

	utils.GetSuccess("purchase successful", buyResponse, response)

}

func Contains(v int, a []int) bool {
	for _, i := range a {
		if i == v {
			return true
		}
	}

	return false
}
