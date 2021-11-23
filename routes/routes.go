package routes

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/gregoflash05/gradely/controllers"
)

type Handler struct {
	Router *mux.Router
}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) SetupRoutes() {
	h.Router = mux.NewRouter().StrictSlash(true)

	// root
	h.Router.HandleFunc("/", VersionHandler)

	// auth
	h.Router.HandleFunc("/v1/user", controllers.UserCreate).Methods("POST")
	h.Router.HandleFunc("/v1/user", controllers.GetUser).Methods("GET")
	h.Router.HandleFunc("/v1/user", controllers.UserUpdate).Methods("PUT")
	h.Router.HandleFunc("/v1/user", controllers.UserDelete).Methods("DELETE")
	h.Router.HandleFunc("/v1/login", controllers.UserLogin).Methods("POST")
	h.Router.HandleFunc("/v1/verify-token", controllers.VerifyTokenHandler).Methods("POST")
	h.Router.HandleFunc("/v1/logout", controllers.Logout)
	h.Router.HandleFunc("/v1/logout/all", controllers.LogoutAll)

	// product
	h.Router.HandleFunc("/v1/products", controllers.ProductCreate).Methods("POST")
	h.Router.HandleFunc("/v1/products", controllers.ProductGetALL).Methods("GET")
	h.Router.HandleFunc("/v1/products/{product_id}", controllers.ProductGet).Methods("GET")
	h.Router.HandleFunc("/v1/products/{product_id}", controllers.ProductUpdate).Methods("PUT")
	h.Router.HandleFunc("/v1/products/{product_id}", controllers.ProductDelete).Methods("DELETE")

	// vending machine
	h.Router.HandleFunc("/v1/deposit", controllers.Deposit).Methods("POST")
	h.Router.HandleFunc("/v1/buy", controllers.BuyProduct).Methods("POST")
	h.Router.HandleFunc("/v1/reset", controllers.DepositReset).Methods("POST")

}

func VersionHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Pretest App - Version 0.0255\n")
}
