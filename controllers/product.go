package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/gregoflash05/gradely/models"
	"github.com/gregoflash05/gradely/utils"
)

func ProductCreate(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintID, _ := (strconv.ParseUint(userID, 10, 64))

	var product models.Product
	err = utils.ParseJSONFromRequest(request, &product)

	if err != nil {
		utils.GetError(err, http.StatusBadRequest, response)
		return
	}

	var user models.User

	utils.GetItemByPrimaryKey(&user, uint(uintID))

	if strings.ToLower(user.Role) != "seller" {
		utils.GetError(fmt.Errorf("user is not a seller"), http.StatusNotAcceptable, response)
		return
	}

	product.SellerId = uint(uintID)

	res := utils.CreateItem(&product)

	if res.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("error adding product"), http.StatusInternalServerError, response)
		return
	}

	respse := map[string]interface{}{
		"product_id": product.ID,
	}

	utils.GetSuccess("product added successfully", respse, response)

}

func ProductGetALL(response http.ResponseWriter, request *http.Request) {
	_, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}
	var products []models.Product

	result := utils.Db.Find(&products)
	if result.RowsAffected < 1 {
		utils.GetError(errors.New("no products found"), http.StatusNotFound, response)
		return
	}

	utils.GetSuccess("products retreived successfully", products, response)

}

func ProductGet(response http.ResponseWriter, request *http.Request) {
	productID := mux.Vars(request)["product_id"]

	_, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintProductID, _ := (strconv.ParseUint(productID, 10, 64))
	var product models.Product

	result := utils.GetItemsByField(&product, "id", uint(uintProductID))
	if result.RowsAffected < 1 {
		utils.GetError(errors.New("product not found"), http.StatusNotFound, response)
		return
	}

	utils.GetSuccess("product retreived successfully", product, response)

}

func ProductDelete(response http.ResponseWriter, request *http.Request) {
	productID := mux.Vars(request)["product_id"]

	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintID, _ := (strconv.ParseUint(userID, 10, 64))
	uintProductID, _ := (strconv.ParseUint(productID, 10, 64))

	var user models.User

	tx := utils.GetItemByPrimaryKey(&user, uint(uintID))
	if tx.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("user not found"), http.StatusUnauthorized, response)
		return
	}

	var product models.Product

	tx = utils.GetItemByPrimaryKey(&product, uint(uintProductID))
	if tx.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("product not found"), http.StatusUnauthorized, response)
		return
	}

	if strings.ToLower(user.Role) != "seller" {
		utils.GetError(fmt.Errorf("user is not a seller"), http.StatusNotAcceptable, response)
		return
	}

	if user.ID != product.SellerId {
		utils.GetError(fmt.Errorf("user not authorized to deleted product"), http.StatusUnauthorized, response)
		return
	}

	result := utils.Db.Delete(models.Product{}, "id = ?", product.ID)

	if result.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("product delete failed"), http.StatusInternalServerError, response)
		return
	}

	utils.GetSuccess("product successfully deleted", nil, response)

}

func ProductUpdate(response http.ResponseWriter, request *http.Request) {
	productID := mux.Vars(request)["product_id"]

	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintID, _ := (strconv.ParseUint(userID, 10, 64))
	uintProductID, _ := (strconv.ParseUint(productID, 10, 64))

	var user models.User

	tx := utils.GetItemByPrimaryKey(&user, uint(uintID))
	if tx.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("user not found"), http.StatusUnauthorized, response)
		return
	}

	var product models.Product

	tx = utils.GetItemByPrimaryKey(&product, uint(uintProductID))
	if tx.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("product not found"), http.StatusUnauthorized, response)
		return
	}

	if strings.ToLower(user.Role) != "seller" {
		utils.GetError(fmt.Errorf("user is not a seller"), http.StatusNotAcceptable, response)
		return
	}

	if user.ID != product.SellerId {
		utils.GetError(fmt.Errorf("user not authorized to deleted product"), http.StatusUnauthorized, response)
		return
	}

	var updateRequest models.ProductUpdate
	if err = utils.ParseJSONFromRequest(request, &updateRequest); err != nil {
		utils.GetError(errors.New("bad update data"), http.StatusBadRequest, response)
		return
	}

	updateMap := map[string]interface{}{}

	if updateRequest.Cost != 0 {
		updateMap["cost"] = updateRequest.Cost
	} else if updateRequest.ProductName != "" {
		updateMap["product_name"] = updateRequest.ProductName
	}

	if len(updateMap) == 0 {
		utils.GetError(errors.New("empty/invalid user input data"), http.StatusBadRequest, response)
		return
	}

	result := utils.Db.Table("products").Where("id = ?", product.ID).Updates(updateMap)

	if result.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("product update failed"), http.StatusInternalServerError, response)
		return
	}

	utils.GetSuccess("product successfully updated", nil, response)

}
