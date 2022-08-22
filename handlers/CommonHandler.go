package handlers

import (
	"audioPhile/claims"
	"audioPhile/database"
	"audioPhile/database/helper"
	"audioPhile/models"
	"audioPhile/utilities"
	"encoding/json"
	"fmt"
	"github.com/golang-jwt/jwt"
	"github.com/jmoiron/sqlx"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"strconv"
	"strings"
	"time"
)

var JwtKey = []byte("secureSecretText")

func SignUp(writer http.ResponseWriter, request *http.Request) {
	var userDetails models.CreateUser
	errors := json.NewDecoder(request.Body).Decode(&userDetails)
	if errors != nil {
		logrus.Error("SignUp: Error in decoding json %v", errors)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	hashPassword, hashErr := utilities.HashPassword(userDetails.Password)
	if hashErr != nil {
		logrus.Error("SignUp : Error in hashing the password")
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
	userDetails.Password = hashPassword
	//records, err := utilities.ReadData("./data.csv")
	//
	//if err != nil {
	//	logrus.Error(err)
	//}
	//var userDetails models.CreateUser
	//for _, record := range records {
	//	userDetails = models.CreateUser{
	//		FirstName:     record[0],
	//		LastName:      record[1],
	//		Email:         record[2],
	//		ContactNumber: record[3],
	//		UserName:      record[4],
	//		Password:      record[5],
	//	}
	//
	//}

	txErr := database.Tx(func(tx *sqlx.Tx) error {
		userID, err := helper.CreateUser(userDetails)
		if err != nil {
			logrus.Error("SignUp : Error in adding details to user table")
			writer.WriteHeader(http.StatusBadRequest)
			return err
		}

		err = helper.CreateRole(string(models.UserRoleUser), userID, tx)
		if err != nil {
			logrus.Error("SignUp : Error in adding details to role table")
			writer.WriteHeader(http.StatusBadRequest)
			return err
		}
		return err
	})

	if txErr != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}
}
func Login(w http.ResponseWriter, r *http.Request) {
	var credentials models.Credential
	err := json.NewDecoder(r.Body).Decode(&credentials)
	if err != nil {
		logrus.Error("Login: Error in decoding json %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if credentials.Username == "" || credentials.Password == "" {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	userDetails, err := helper.UserLogin(credentials.Username)
	if err != nil {
		logrus.Error("Login: Error in getting password %v", err)
		w.WriteHeader(http.StatusUnauthorized)
	}

	if compareErr := bcrypt.CompareHashAndPassword([]byte(userDetails.Password), []byte(credentials.Password)); compareErr != nil {
		logrus.Printf("Signin : Error in comparing the passwords.")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	//if userDetails.Password != credentials.Password {
	//	w.WriteHeader(http.StatusUnauthorized)
	//}

	userRole, err := helper.GetRole(userDetails.UserId)
	if err != nil {
		logrus.Error("Login: Error in Getting Role of User %v", err)
		w.WriteHeader(http.StatusInternalServerError)
	}

	expirationTime := time.Now().Add(time.Hour * 7)
	sessionID, err := helper.CreateSession(userDetails.UserId, expirationTime)
	if err != nil {
		logrus.Error("LogIn : Error in Creating the session %v", err)
		w.WriteHeader(http.StatusUnauthorized)
	}

	mapClaim := &claims.MapClaims{
		SessionID: sessionID,
		ID:        userDetails.UserId,
		Role:      userRole,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, mapClaim)
	signedToken, err := token.SignedString([]byte("secureSecretText"))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	tokenByte, err := json.Marshal(signedToken)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	fmt.Println(signedToken)
	_, _ = w.Write(tokenByte)
}

func GetAllProducts(writer http.ResponseWriter, request *http.Request) {
	searchProduct := request.URL.Query().Get("search")
	pageNo := request.URL.Query().Get("page")
	if pageNo == strings.TrimSpace("") {
		pageNo = "0"
	}
	Page, err := strconv.Atoi(pageNo)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	limitSize := request.URL.Query().Get("limit")
	if limitSize == strings.TrimSpace("") {
		limitSize = "5"
	}
	Limit, err := strconv.Atoi(limitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	products, err := helper.GetAllProducts(Page, Limit, searchProduct)
	if err != nil {
		logrus.Error("GetAllProducts: Products  can't be fetched %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := json.NewEncoder(writer).Encode(products)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func GetAllImageDetails(writer http.ResponseWriter, request *http.Request) {
	pageNo := request.URL.Query().Get("page")
	if pageNo == strings.TrimSpace("") {
		pageNo = "0"
	}
	Page, err := strconv.Atoi(pageNo)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}
	limitSize := request.URL.Query().Get("limit")
	if limitSize == strings.TrimSpace("") {
		limitSize = "5"
	}
	Limit, err := strconv.Atoi(limitSize)
	if err != nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	products, err := helper.GetAllImageId(Page, Limit)
	if err != nil {
		logrus.Error("GetAllImageDetails: Images  can't be fetched %v", err)
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	errs := json.NewEncoder(writer).Encode(products)
	if errs != nil {
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}

func SignOut(writer http.ResponseWriter, request *http.Request) {
	uc := helper.GetContextData(request)
	if uc == nil {
		writer.WriteHeader(http.StatusBadRequest)
		return
	}

	deleteSession := helper.DeleteSession(uc.SessionID)
	if deleteSession != nil {
		logrus.Error("SignOut: Session can't be deleted %v", deleteSession)
		writer.WriteHeader(http.StatusInternalServerError)
		return
	}

}
