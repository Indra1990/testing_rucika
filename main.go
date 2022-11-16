package main

import (
	"go-bun-chi/controller"
	"go-bun-chi/database"
	"go-bun-chi/helper"
	"log"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
)

var (
	db *gorm.DB = database.SetUpConfigMySql()
)

func main() {
	defer database.CloseDatabaseMysSqlConnection(db)

	authController := controller.NewAuthController(db)
	customerController := controller.NewCustomerController(db)
	orderController := controller.NewOrderController(db)

	r := gin.Default()
	r.POST("/login", authController.Login)
	r.Use(authMiddleware(*authController, *customerController))
	{
		// customer
		r.GET("/api/customer", customerController.FindMany)
		r.POST("/api/customer/create", customerController.Create)
		r.GET("/api/customer/:customerId", customerController.FindById)
		r.PUT("/api/customer/update/:customerId", customerController.Updater)
		r.DELETE("/api/customer/delete/:customerId", customerController.Deleted)
		// order
		r.GET("/api/order", orderController.FindMany)
		r.POST("/api/order/create", orderController.Create)
		r.GET("/api/order/:orderId", orderController.FindById)
		r.PUT("/api/order/update/:orderId", orderController.Updater)
		r.DELETE("/api/order/delete/:orderId", orderController.Deleted)

	}

	r.Run(":8081")
}

func authMiddleware(authController controller.AuthController, customerController controller.CustomerController) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if !strings.Contains(authHeader, "Bearer") {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		tokenString := ""
		arrayToken := strings.Split(authHeader, " ")
		if len(arrayToken) == 2 {
			tokenString = arrayToken[1]
		}

		token, err := authController.ValidateToken(tokenString)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok || !token.Valid {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		customerId := int(claims["customer_id"].(float64))

		log.Println("customer ID by token ", customerId)
		customerUser, err := customerController.CustomerById(customerId)
		if err != nil {
			response := helper.APIResponse("Unauthorized", http.StatusUnauthorized, "error", nil)
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, response)
			return
		}

		ctx.Set("currentCustomer", customerUser)
	}
}
