package controllers

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/go-playground/validator"
	"github.com/golang-jwt/jwt"
	"github.com/gregoflash05/gradely/models"
	"github.com/gregoflash05/gradely/utils"
	"golang.org/x/crypto/bcrypt"
)

var (
	validate               = validator.New()
	errEmailNotValid       = errors.New("email address is not valid")
	errHashingFailed       = errors.New("failed to hashed password")
	ErrUserNotFound        = errors.New("user not found, confirm and try again")
	ErrInvalidCredentials  = errors.New("invalid login credentials, confirm and try again")
	ErrAccountConfirmError = errors.New("your account is not verified, kindly check your email for verification code")
	ErrAccessExpired       = errors.New("error fetching user info, access token expired, kindly login again")
	ErrGeneratingToken     = errors.New("error generating token")
	ErrConfirmPassword     = errors.New("passwords do not match")
	DefaultHashCode        = 14
)

// Method to hash password.
func GenerateHashPassword(password string) (string, error) {
	cost := 14
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), cost)

	return string(bytes), err
}

func CheckPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func ExtractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	//normally Authorization the_token_xxx
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}
	return ""
}

func FetchUserByEmail(email string) *models.User {
	u := &models.User{}
	utils.GetItemByField(&u, "email", email)

	return u
}

func UserCreate(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var user models.User
	err := utils.ParseJSONFromRequest(request, &user)

	if err != nil {
		utils.GetError(err, http.StatusBadRequest, response)
		return
	}

	userEmail := strings.ToLower(user.Email)
	if !utils.IsValidEmail(userEmail) {
		utils.GetError(errEmailNotValid, http.StatusBadRequest, response)
		return
	}

	var checkUser models.User

	result := utils.GetItemsByField(&checkUser, "email", userEmail)
	if result.RowsAffected > 0 {
		utils.GetError(
			fmt.Errorf("user with email: %s already exists", userEmail),
			http.StatusBadRequest,
			response,
		)

		return
	}

	hashPassword, err := GenerateHashPassword(user.Password)
	if err != nil {
		utils.GetError(errHashingFailed, http.StatusInternalServerError, response)
		return
	}

	user.Email = userEmail
	user.UserName = userEmail
	user.Password = hashPassword
	user.IsVerified = true
	user.Role = "buyer"

	res := utils.CreateItem(&user)

	if res.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("error Creating user"), http.StatusInternalServerError, response)
		return
	}

	respse := map[string]interface{}{
		"user_id": user.ID,
	}

	utils.GetSuccess("user created", respse, response)
}

func UserLogin(response http.ResponseWriter, request *http.Request) {
	response.Header().Add("content-type", "application/json")

	var creds models.AuthCredentials
	if err := utils.ParseJSONFromRequest(request, &creds); err != nil {
		utils.GetError(err, http.StatusBadRequest, response)
		return
	}

	if err := validate.Struct(creds); err != nil {
		utils.GetError(err, http.StatusBadRequest, response)
		return
	}

	var vser models.User

	result := utils.GetItemsByField(&vser, "email", creds.Email)
	if result.RowsAffected < 1 {
		utils.GetError(ErrUserNotFound, http.StatusBadRequest, response)
		return
	}
	// check if user is verified
	if !vser.IsVerified {
		utils.GetError(ErrAccountConfirmError, http.StatusBadRequest, response)
		return
	}

	// check password
	check := CheckPassword(creds.Password, vser.Password)
	if !check {
		utils.GetError(ErrInvalidCredentials, http.StatusBadRequest, response)
		return
	}

	token, err := CreateToken(strconv.FormatUint(uint64(vser.ID), 10))
	if err != nil {
		utils.GetError(ErrGeneratingToken, http.StatusInternalServerError, response)
		return
	}

	var sessions []models.Session

	result = utils.GetItemsByField(&sessions, "user_id", vser.ID)
	if result.RowsAffected > 1 {
		utils.GetSuccess("Login successful. There is already an active session using your account", token, response)
		return
	}
	utils.GetSuccess("Login successful", token, response)
}

func CreateToken(userid string) (string, error) {
	var err error
	//Creating Access Token
	atClaims := jwt.MapClaims{}
	atClaims["authorized"] = true
	atClaims["user_id"] = userid
	atClaims["exp"] = time.Now().Add(time.Minute * 60 * 24).Unix()
	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(os.Getenv("ACCESS_SECRET")))
	if err != nil {
		return "", err
	}
	uintID, _ := strconv.ParseUint(userid, 10, 64)
	sessionTableItem := models.Session{
		UserID: uint(uintID),
		Token:  token,
	}
	utils.CreateItem(&sessionTableItem)

	return token, nil
}

func VerifyToken(r *http.Request) (*jwt.Token, error) {
	tokenString := ExtractToken(r)
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		//Make sure that the token method conform to "SigningMethodHMAC"
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("ACCESS_SECRET")), nil
	})
	if err != nil {
		utils.Db.Delete(models.Session{}, "token = ?", tokenString)
		return nil, err
	}

	var session models.Session
	result := utils.GetItemsByField(&session, "token", tokenString)
	if result.RowsAffected < 1 {
		return nil, fmt.Errorf("not authenticated")
	}

	return token, nil
}
func TokenValid(r *http.Request) (string, error) {
	token, err := VerifyToken(r)
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok && !token.Valid {
		return "", err
	}
	return fmt.Sprintf("%v", claims["user_id"]), nil
}
func DeleteMapProps(m map[string]interface{}, s []string) {
	for _, v := range s {
		delete(m, v)
	}
}

func VerifyTokenHandler(response http.ResponseWriter, request *http.Request) {
	_, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusBadRequest, response)
		return
	}

	utils.GetSuccess("token is valid", "", response)
}

func GetUser(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	var user models.User

	uintID, _ := (strconv.ParseUint(userID, 10, 64))

	result := utils.GetItemsByField(&user, "id", uint(uintID))
	if result.RowsAffected < 1 {
		utils.GetError(errors.New("user not found"), http.StatusNotFound, response)
		return
	}

	user.Password = ""

	utils.GetSuccess("user retrieved successfully", user, response)
}

func Logout(response http.ResponseWriter, request *http.Request) {
	_, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	result := utils.Db.Delete(models.Session{}, "token = ?", ExtractToken(request))

	if result.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("logout unsuccessfull"), http.StatusInternalServerError, response)
		return
	}

	utils.GetSuccess("logout successful", "", response)

}

func LogoutAll(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintID, _ := (strconv.ParseUint(userID, 10, 64))
	result := utils.Db.Delete(models.Session{}, "user_id = ?", uint(uintID))

	if result.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("logout all sessions unsuccessfull"), http.StatusInternalServerError, response)
		return
	}

	utils.GetSuccess("logout all sessions successful", "", response)
}

func UserUpdate(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintID, _ := (strconv.ParseUint(userID, 10, 64))

	var user models.UserUpdate
	if err = utils.ParseJSONFromRequest(request, &user); err != nil {
		utils.GetError(errors.New("bad update data"), http.StatusBadRequest, response)
		return
	}
	updateMap := map[string]interface{}{}

	if user.FullName != "" {
		updateMap["full_name"] = user.FullName
	} else if user.Phone != "" {
		updateMap["phone"] = user.Phone
	} else if user.Role != "" {
		updateMap["role"] = user.Role
	}

	if len(updateMap) == 0 {
		utils.GetError(errors.New("empty/invalid user input data"), http.StatusBadRequest, response)
		return
	}

	result := utils.Db.Table("users").Where("id = ?", uint(uintID)).Updates(updateMap)

	if result.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("user update failed"), http.StatusInternalServerError, response)
		return
	}

	utils.GetSuccess("user successfully updated", nil, response)
}

func UserDelete(response http.ResponseWriter, request *http.Request) {
	userID, err := TokenValid(request)
	if err != nil {
		utils.GetError(fmt.Errorf("token Invalid"), http.StatusUnauthorized, response)
		return
	}

	uintID, _ := (strconv.ParseUint(userID, 10, 64))
	result := utils.Db.Delete(models.User{}, "id = ?", uint(uintID))

	if result.RowsAffected < 1 {
		utils.GetError(fmt.Errorf("user delete failed"), http.StatusInternalServerError, response)
		return
	}

	utils.GetSuccess("user successfully deleted", nil, response)

}
