package main

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"math"
	"math/rand"
	"time"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

var (
	err     error
	newline = []byte{'\n'}
	BoolGen boolGen
)

type Map map[string]interface{}
type boolGen struct {
	src       rand.Source
	cache     int64
	remaining int
}

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func JSON(c *gin.Context, code int, success bool, message string, data interface{}) {
	exitCode := 0
	if !success {
		exitCode = 0
	}

	c.JSON(code, Message{
		ExitCode: exitCode,
		Message:  message,
		Data:     data,
	})
}

func Abort(c *gin.Context, code int, success bool, message string, data interface{}) {
	exitCode := 0
	if !success {
		exitCode = 0
	}

	c.AbortWithStatusJSON(code, Message{
		ExitCode: exitCode,
		Message:  message,
		Data:     data,
	})
}

func RandomAccountID() int {
	v := rand.Intn(9999999999999-1000000000000) + 1000000000000
	return v
}

func GenerateToken(accountID string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(accountID), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(hash), nil
}

func HashPassword(password string) ([]byte, error) {
	passwordBytes := []byte(password)

	hashedPassword, err := bcrypt.GenerateFromPassword(passwordBytes, bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	return hashedPassword, nil
}

func CheckPasswordHash(passwordInput string, hashedPassword []byte) bool {
	err = bcrypt.CompareHashAndPassword(hashedPassword, []byte(passwordInput))
	if err != nil {
		return err == nil
	}

	return true
}

func BindJSON(data []byte) (Map, error) {
	var result Map
	err := json.Unmarshal(data, &result)
	if err != nil {
		return nil, errors.New("Invalid JSON")
	}

	return result, nil
}

func RandomNumber(min, max int) int {
	return rand.Intn(max-min) + min
}

func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func NewBoolGenerator() *boolGen {
	return &boolGen{
		src: rand.NewSource(time.Now().UnixNano()),
	}
}

func (b *boolGen) Bool() bool {
	if b.remaining == 0 {
		b.cache, b.remaining = b.src.Int63(), 63
	}

	result := b.cache&0x01 == 1
	b.cache >>= 1
	b.remaining--

	return result
}
