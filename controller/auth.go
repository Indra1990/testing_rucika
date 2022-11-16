package controller

import (
	"errors"
	"go-bun-chi/database"
	"go-bun-chi/dto"
	"go-bun-chi/entity"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type AuthController struct {
	db *gorm.DB
}

func NewAuthController(db *gorm.DB) *AuthController {
	return &AuthController{db}
}
func (c *AuthController) Login(ctx *gin.Context) {
	defer func() {
		if err := recover(); err != nil {
			errRc := err.(error)
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "Invalid Request",
				"error":  errRc.Error(),
			})
			return
		}
	}()

	var authReq dto.AuthRequest
	if err := ctx.ShouldBindJSON(&authReq); err != nil {
		for _, fieldErr := range err.(validator.ValidationErrors) {
			errMsg := fieldErr.Field() + " is " + fieldErr.Tag()
			ctx.JSON(http.StatusUnprocessableEntity, gin.H{
				"result": "Invalid Request",
				"error":  errMsg,
			})
			return
		}
		return
	}
	var customer entity.Customer

	findCustomerErr := c.db.Model(&entity.Customer{}).Where("email = ?", authReq.Email).First(&customer).Error
	if findCustomerErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"result": "Invalid Request",
			"error":  "Invalid login request",
		})
		return
	}

	comparePasswordErr := bcrypt.CompareHashAndPassword([]byte(customer.Password), []byte(authReq.Password))
	if comparePasswordErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"result": "Invalid Request",
			"error":  "Invalid login request",
		})
		return
	}

	key, errKey := database.GetSecretKey()
	if errKey != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"result": "Invalid Request",
			"error":  "Invalid login request",
		})
		return
	}

	rfKey, errRf := database.GetRefreshSecretKey()
	if errRf != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"result": "Invalid Request",
			"error":  "Invalid login request",
		})
		return
	}

	atClaims := jwt.MapClaims{}
	atClaims["customer_id"] = customer.ID
	atClaims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	at := jwt.NewWithClaims(jwt.SigningMethodHS256, atClaims)
	token, err := at.SignedString([]byte(key))
	if err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"result": "Invalid Request",
			"error":  "Invalid login request",
		})
		return
	}

	rtClaims := jwt.MapClaims{}
	rtClaims["customer_id"] = customer.ID
	rtClaims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	rt := jwt.NewWithClaims(jwt.SigningMethodHS256, rtClaims)
	rfToken, err := rt.SignedString([]byte(rfKey))

	data := map[string]string{
		"acccessToken": token,
		"refreshToken": rfToken,
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": data,
	})
}

func (c *AuthController) ValidateToken(tokenEncoded string) (*jwt.Token, error) {
	key, _ := database.GetSecretKey()
	token, err := jwt.Parse(tokenEncoded, func(tokenEncoded *jwt.Token) (interface{}, error) {
		_, okey := tokenEncoded.Method.(*jwt.SigningMethodHMAC)
		if !okey {
			return nil, errors.New("invalid token")
		}
		return []byte(key), nil
	})

	if err != nil {
		return token, err
	}
	return token, nil
}
