package controller

import (
	"errors"
	"go-bun-chi/dto"
	"go-bun-chi/entity"
	"go-bun-chi/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type CustomerController struct {
	db *gorm.DB
}

func NewCustomerController(db *gorm.DB) *CustomerController {
	return &CustomerController{db}
}

func (c *CustomerController) FindMany(ctx *gin.Context) {
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

	var customers []entity.Customer
	var findManyCustomers []dto.CustomerDTO
	var totalRow int64

	var emptyArray = make([]string, 0)
	pagination := utils.GeneratePaginationFromRequest(ctx)
	offset := (pagination.Page - 1) * pagination.Limit

	query := c.db.Model(&customers).Select(
		"customers.id",
		"customers.name",
		"customers.email",
		"date_format(customers.created_at,'%Y-%m-%d %T') as created_at",
		"date_format(customers.updated_at,'%Y-%m-%d %T') as updated_at",
	).Count(&totalRow).Limit(pagination.Limit).
		Offset(offset).
		Scan(&findManyCustomers)

	if query.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  query.Error.Error(),
		})
		return
	}

	if query.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"result":   emptyArray,
			"totalRow": 0,
			"error":    "Customer Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result":     findManyCustomers,
		"totalRow":   totalRow,
		"pagination": offset,
	})
}

func (c *CustomerController) Create(ctx *gin.Context) {
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

	var custRequest dto.CustomerCreateRequest

	if err := ctx.ShouldBindJSON(&custRequest); err != nil {
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

	password, passwordErr := bcrypt.GenerateFromPassword([]byte(custRequest.Password), bcrypt.MinCost)
	if passwordErr != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"result": "Invalid Request",
			"error":  passwordErr.Error(),
		})
	}

	cust := entity.Customer{
		Name:     custRequest.Name,
		Email:    custRequest.Email,
		Password: string(password),
	}

	query := c.db.Create(&cust)
	if query.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  query.Error.Error(),
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": "success insert customer",
	})

}

func (c *CustomerController) CustomerById(customerId int) (dto.CustomerDTO, error) {
	var customers entity.Customer
	var findManyCustomers dto.CustomerDTO

	query := c.db.Model(&customers).Select(
		"customers.id",
		"customers.name",
		"customers.email",
		"date_format(customers.created_at,'%Y-%m-%d %T') as created_at",
		"date_format(customers.updated_at,'%Y-%m-%d %T') as updated_at",
	).Where("customers.id = ?", customerId).Scan(&findManyCustomers)

	if query.Error != nil {
		return dto.CustomerDTO{}, query.Error
	}

	if query.RowsAffected == 0 {
		return dto.CustomerDTO{}, errors.New("user not found")
	}

	return findManyCustomers, nil
}

func (c *CustomerController) FindById(ctx *gin.Context) {
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

	var customers entity.Customer
	var findManyCustomers dto.CustomerDTO
	intParam, errIntParam := strconv.Atoi(ctx.Param("customerId"))

	if errIntParam != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  errIntParam.Error(),
		})
		return
	}

	query := c.db.Model(&customers).Select(
		"customers.id",
		"customers.name",
		"customers.email",
		"date_format(customers.created_at,'%Y-%m-%d %T') as created_at",
		"date_format(customers.updated_at,'%Y-%m-%d %T') as updated_at",
	).Where("customers.id = ?", intParam).Scan(&findManyCustomers)

	if query.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  query.Error.Error(),
		})
		return
	}

	if query.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"result": "Invalid Request",
			"error":  "Customer Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": findManyCustomers,
	})
}

func (c *CustomerController) Updater(ctx *gin.Context) {
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

	var custRequest dto.CustomerUpdaterRequest
	intParam, errIntParam := strconv.Atoi(ctx.Param("customerId"))

	if errIntParam != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  errIntParam.Error(),
		})
		return
	}

	if err := ctx.ShouldBindJSON(&custRequest); err != nil {
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

	updCustomer := &entity.Customer{
		Name:  custRequest.Name,
		Email: custRequest.Email,
	}

	query := c.db.Model(&entity.Customer{}).Where("id = ?", intParam).Updates(updCustomer)
	if query.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  query.Error.Error(),
		})
		return
	}

	if query.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"result": "Invalid Request",
			"error":  "Customer Not Found",
		})
		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"result": "successfully customer updated",
	})
}

func (c *CustomerController) Deleted(ctx *gin.Context) {
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

	intParam, errIntParam := strconv.Atoi(ctx.Param("customerId"))

	if errIntParam != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  errIntParam.Error(),
		})
		return
	}

	query := c.db.Model(&entity.Customer{}).Where("id = ?", intParam).Delete(&entity.Customer{})
	if query.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  query.Error.Error(),
		})
		return
	}

	if query.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"result": "Invalid Request",
			"error":  "Customer Not Found",
		})
		return
	}
	ctx.JSON(http.StatusOK, gin.H{
		"result": "successfully customer deleted",
	})
}
