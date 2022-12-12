package main

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

type Authentication struct{}

func (a *Authentication) CreateAccount(c *gin.Context) {
	userName, userNameParameterIncluded := c.GetQuery("username")
	emailAddress, emailAddressParameterIncluded := c.GetQuery("email_address")
	emailAddressConfirmation, emailAddressConfirmationParameterIncluded := c.GetQuery("email_address_confirmation")
	password, passwordParameterIncluded := c.GetQuery("password")
	birthdayDay, birthdayDayParameterIncluded := c.GetQuery("birthday_day")
	birthdayMonth, birthdayMonthParameterIncluded := c.GetQuery("birthday_month")
	birthdayYear, birthdayYearParameterIncluded := c.GetQuery("birthday_year")

	if !(userNameParameterIncluded || emailAddressParameterIncluded || emailAddressConfirmationParameterIncluded || passwordParameterIncluded || birthdayDayParameterIncluded || birthdayMonthParameterIncluded || birthdayYearParameterIncluded) {
		JSON(c, http.StatusBadRequest, false, "Missing required parameters", nil)
		return
	}

	if len(userName) < 2 || len(userName) > 32 {
		JSON(c, http.StatusBadRequest, false, "Username must be between 2 and 32 characters", nil)
		return
	}

	usernameRegex := regexp.MustCompile("^[a-zA-Z0-9_]+$")
	if !usernameRegex.MatchString(userName) {
		JSON(c, http.StatusBadRequest, false, "Username must only contain letters, numbers, and underscores", nil)
		return
	}

	_, accountFoundWithUsername := GetUserByUsername(userName)
	if accountFoundWithUsername {
		JSON(c, http.StatusBadRequest, false, "Username already in use", nil)
		return
	}

	if emailAddress == "" || len(emailAddress) > 254 {
		JSON(c, http.StatusBadRequest, false, "Invalid email address (length)", nil)
		return
	}

	_, accountFoundWithEmail := GetUserByEmailAddress(emailAddress)
	if accountFoundWithEmail {
		JSON(c, http.StatusBadRequest, false, "Email address already in use", nil)
		return
	}

	emailRegex := regexp.MustCompile(`^[\w-\.]+@([\w-]+\.)+[\w-]{2,4}$`)
	if !emailRegex.MatchString(emailAddress) {
		JSON(c, http.StatusBadRequest, false, "Invalid email address (not a email address)", nil)
		return
	}

	if emailAddress != emailAddressConfirmation {
		JSON(c, http.StatusBadRequest, false, "Email addresses do not match", nil)
		return
	}

	if len(password) < 4 {
		JSON(c, http.StatusBadRequest, false, "Password must be at least 4 characters", nil)
		return
	}

	birthdayDayInt, err := strconv.Atoi(birthdayDay)
	if birthdayDayInt < 1 || birthdayDayInt > 31 || err != nil {
		JSON(c, http.StatusBadRequest, false, "Invalid birthday day", nil)
		return
	}

	birthdayMonthInt, err := strconv.Atoi(birthdayMonth)
	if birthdayMonthInt < 1 || birthdayMonthInt > 12 || err != nil {
		JSON(c, http.StatusBadRequest, false, "Invalid birthday month", nil)
		return
	}

	birthdayYearInt, err := strconv.Atoi(birthdayYear)
	if birthdayYearInt < 1900 || birthdayYearInt > 2100 || err != nil {
		JSON(c, http.StatusBadRequest, false, "Invalid birthday year", nil)
		return
	}

	hashedPassword, err := HashPassword(password)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "Could not hash password", nil)
		return
	}

	birthdayTimestamp, err := time.Parse("2006-01-02", fmt.Sprintf("%s-%s-%s", birthdayYear, birthdayMonth, birthdayDay))
	if err != nil {
		fmt.Println(err)
		JSON(c, http.StatusInternalServerError, false, "Could not parse birthday. (format: 2006-01-02)", nil)
		return
	}

	accountID := RandomAccountID()
	user := User{
		ID:           fmt.Sprintf("%d", accountID),
		Username:     userName,
		Email:        emailAddress,
		CreationDate: time.Now().Unix(),
		Password:     hashedPassword,
		Birthday: Birthday{
			Timestamp: birthdayTimestamp.Unix(),
		},
		Balances: Balances{
			FIAT: FIAT{
				USD: 0.0,
			},
			Crypto: Crypto{
				BTC: 0.0,
				ETH: 0.0,
			},
		},
	}

	_, err = InsertDocument("users", user)
	if err != nil {
		JSON(c, http.StatusOK, false, "An internal error has been generated, retry later", nil)
		return
	}

	token, err := AssignToken(user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Signed up succesfully.", Map{
		"token": token,
	})
}

func (a *Authentication) NewSession(c *gin.Context) {
	emailAddress, emailAddressParameterIncluded := c.GetQuery("email_address")
	password, passwordParameterIncluded := c.GetQuery("password")

	if !(emailAddressParameterIncluded && passwordParameterIncluded) {
		JSON(c, http.StatusBadRequest, false, "Missing required parameters", nil)
		return
	}

	if emailAddress == "" || len(emailAddress) > 254 {
		JSON(c, http.StatusBadRequest, false, "Invalid email address", nil)
		return
	}

	emailRegex := "(^[a-zA-Z0-9_.+-]+@[a-zA-Z0-9-]+\\.[a-zA-Z0-9-.]+$)"
	if matched, err := regexp.MatchString(emailRegex, emailAddress); !matched || err != nil {
		JSON(c, http.StatusBadRequest, false, "Invalid email address", nil)
		return
	}

	if len(password) < 4 {
		JSON(c, http.StatusBadRequest, false, "Password must be at least 4 characters", nil)
		return
	}

	user, found := GetUserByEmailAddress(emailAddress)

	if !found || !CheckPasswordHash(password, user.Password) {
		JSON(c, http.StatusBadRequest, false, "Invalid email address or password", nil)
		return
	}

	token, err := AssignToken(user.ID)
	if err != nil {
		JSON(c, http.StatusInternalServerError, false, "An internal error has been generated, retry later", nil)
		return
	}

	JSON(c, http.StatusOK, true, "Signed in succesfully.", Map{
		"token": token,
	})
}

func (a *Authentication) RequestData(c *gin.Context) {
	token := c.GetHeader("X-Auth-Token")
	if token == "" {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	session := SessionTokens[token]
	if session.Token == "" || session.AccountID == "" {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	user, found := GetUserByID(session.AccountID)
	if !found {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	publicUser := User{
		ID:           user.ID,
		Username:     user.Username,
		Email:        user.Email,
		Birthday:     user.Birthday,
		CreationDate: user.CreationDate,
		Balances:     user.Balances,
	}

	JSON(c, http.StatusOK, true, "Requested user data succesfully.", publicUser)
}

func (a *Authentication) DestroySession(c *gin.Context) {
	token := c.GetHeader("X-Auth-Token")
	if token == "" {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	session := SessionTokens[token]
	if session.Token == "" || session.AccountID == "" {
		JSON(c, http.StatusUnauthorized, false, "Unauthorized", nil)
		return
	}

	RemoveToken(token)

	JSON(c, http.StatusOK, true, "Signed out succesfully.", nil)
}
