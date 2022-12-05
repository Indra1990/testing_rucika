package controller

import (
	"fmt"
	"go-bun-chi/dto"
	"go-bun-chi/entity"
	"go-bun-chi/utils"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type OrderController struct {
	db *gorm.DB
}

func NewOrderController(db *gorm.DB) *OrderController {
	return &OrderController{db}
}

func (c *OrderController) FindMany(ctx *gin.Context) {
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

	startDate := ctx.Query("start_date")
	var startDateString string
	if startDate != "" {
		start, startErr := time.Parse("2006-01-02", startDate)
		if startErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "Invalid Request",
				"error":  startErr.Error(),
			})
			return
		}

		startDateString = start.Format("2006-01-02") + " " + "00:00:00"
	}

	endDate := ctx.Query("end_date")
	var endDateString string
	if endDate != "" {
		end, endErr := time.Parse("2006-01-02", endDate)
		if endErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "Invalid Request",
				"error":  endErr.Error(),
			})
			return
		}

		endDateString = end.Format("2006-01-02") + " " + "23:59:00"
	}

	pagination := utils.GeneratePaginationFromRequest(ctx)
	offset := (pagination.Page - 1) * pagination.Limit
	var emptyArray = make([]string, 0)
	var totalRow int64
	var orders []entity.Orders
	var ordersDTO []dto.OrderFindManyDTO

	query := c.db.Model(&orders).Select(
		"orders.id",
		"orders.title",
		"orders.order_number",
		"orders.note",
		"orders.total",
		"date_format(orders.created_at,'%Y-%m-%d %T') as created_at",
		"date_format(orders.updated_at,'%Y-%m-%d %T') as updated_at",
		"customers.name as created_by",
	).
		Joins("left join customers on orders.customer_id = customers.id")

	if startDateString != "" {
		query.Where("orders.created_at  >= ?", startDateString)
	}

	if endDateString != "" {
		query.Where("orders.created_at  <= ?", endDateString)

	}

	query.Count(&totalRow).
		Limit(pagination.Limit).
		Offset(offset).
		Scan(&ordersDTO)

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
			"error":    "Order Not Found",
		})
		return
	}

	result := map[string]interface{}{
		"data":       ordersDTO,
		"totalRow":   totalRow,
		"pagination": offset,
	}

	responseWithJson(ctx, http.StatusOK, "Success", "Find Many Order", result)

}

func (c *OrderController) Create(ctx *gin.Context) {
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

	jwtCust := ctx.MustGet("currentCustomer").(dto.CustomerDTO)

	var totalHeader float64
	var orderRequest dto.OrderCreateRequest
	var orderDetailEnt []entity.OrderDetails

	if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"result": "Invalid Request",
			"error":  err.Error(),
		})
		return
	}

	costumerId, _ := strconv.ParseUint(jwtCust.Id, 10, 64)
	currentTime := time.Now()

	headerInsert := &entity.Orders{
		CustomerId:  costumerId,
		Title:       orderRequest.Title,
		Note:        orderRequest.Note,
		OrderNumber: "PO" + currentTime.Format("20060102150405.000000"),
	}

	for _, vDet := range orderRequest.OrderDetailCreateRequest {
		qty, qtyErr := strconv.ParseInt(vDet.Qty, 10, 64)
		if qtyErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "Invalid Request",
				"error":  qtyErr.Error(),
			})
			return
		}

		price, priceErr := strconv.ParseFloat(vDet.Price, 64)
		if priceErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "Invalid Request",
				"error":  priceErr.Error(),
			})
			return
		}

		orderDetailEnt = append(orderDetailEnt, entity.OrderDetails{
			Item:   vDet.Item,
			Qty:    qty,
			Price:  price,
			Amount: float64(qty) * price,
		})

		totalHeader += float64(qty) * price
	}

	// db transaction
	transactionErr := c.db.Transaction(func(tx *gorm.DB) error {
		headerInsert.Total = totalHeader
		if err := tx.Create(headerInsert).Error; err != nil {
			return err
		}

		for i := 0; i < len(orderRequest.OrderDetailCreateRequest); i++ {
			orderDetailEnt[i].OrderId = headerInsert.ID
		}

		if err := tx.Create(&orderDetailEnt).Error; err != nil {
			return err
		}
		return nil
	})

	if transactionErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  transactionErr.Error(),
		})
		return
	}

	responseWithJson(ctx, http.StatusOK, "Success", "Successfully Create Order", nil)
}

func (c *OrderController) FindById(ctx *gin.Context) {
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
	var orderDet entity.Orders
	intParam, errIntParam := strconv.Atoi(ctx.Param("orderId"))

	if errIntParam != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  errIntParam.Error(),
		})
		return
	}
	if err := c.db.Model(&entity.Orders{}).Preload("OrderDetails").Where("id = ?", intParam).Find(&orderDet).Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  err.Error(),
		})
		return
	}

	mapDTO := c.mapOrderDetailEntityToDTO(orderDet)
	result := map[string]interface{}{
		"data": mapDTO,
	}

	responseWithJson(ctx, http.StatusOK, "Success", "Find Detail Order", result)
}

func (c *OrderController) Updater(ctx *gin.Context) {
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

	intParam, errIntParam := strconv.Atoi(ctx.Param("orderId"))

	if errIntParam != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  errIntParam.Error(),
		})
		return
	}

	var orderDetailEnt []entity.OrderDetails
	var orderDetailEntCreUpd []entity.OrderDetails

	var orderRequest dto.OrderCreateRequest
	var totalHeader float64

	if err := ctx.ShouldBindJSON(&orderRequest); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{
			"result": "Invalid Request",
			"error":  err.Error(),
		})
		return
	}

	headerUpdate := &entity.Orders{
		Title: orderRequest.Title,
		Note:  orderRequest.Note,
	}

	for _, vDet := range orderRequest.OrderDetailCreateRequest {

		qty, qtyErr := strconv.ParseInt(vDet.Qty, 10, 64)
		if qtyErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "Invalid Request",
				"error":  qtyErr.Error(),
			})
			return
		}

		price, priceErr := strconv.ParseFloat(vDet.Price, 64)
		if priceErr != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"result": "Invalid Request",
				"error":  priceErr.Error(),
			})
			return
		}

		if vDet.OrderDetailId != "" {
			orderDetailId, orderDetailIdErr := strconv.ParseUint(vDet.OrderDetailId, 10, 64)
			if orderDetailIdErr != nil {
				ctx.JSON(http.StatusBadRequest, gin.H{
					"result": "Invalid Request",
					"error":  orderDetailIdErr.Error(),
				})
				return
			}

			orderDetailEnt = append(orderDetailEnt, entity.OrderDetails{
				ID:     orderDetailId,
				Item:   vDet.Item,
				Qty:    qty,
				Price:  price,
				Amount: float64(qty) * price,
			})
		} else {
			orderDetailEntCreUpd = append(orderDetailEntCreUpd, entity.OrderDetails{
				OrderId: uint64(intParam),
				Item:    vDet.Item,
				Qty:     qty,
				Price:   price,
				Amount:  float64(qty) * price,
			})
		}

		totalHeader += float64(qty) * price
	}

	transactionErr := c.db.Transaction(func(tx *gorm.DB) error {
		headerUpdate.Total = totalHeader
		if err := tx.Model(&entity.Orders{}).Where("id = ?", intParam).Updates(headerUpdate).Error; err != nil {
			return err
		}

		if len(orderDetailEntCreUpd) > 0 {
			if err := tx.Model(&entity.OrderDetails{}).Create(&orderDetailEntCreUpd).Error; err != nil {
				return err
			}
		}

		if len(orderDetailEnt) > 0 {
			for _, ordDet := range orderDetailEnt {
				orderDetailModels := entity.OrderDetails{}
				if err := tx.Model(&orderDetailModels).Where("id = ?", ordDet.ID).Updates(entity.OrderDetails{
					Item:   ordDet.Item,
					Qty:    ordDet.Qty,
					Price:  ordDet.Price,
					Amount: ordDet.Amount,
				}).Error; err != nil {
					return err
				}
			}
		}
		return nil
	})

	if transactionErr != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  transactionErr.Error(),
		})
		return
	}

	responseWithJson(ctx, http.StatusOK, "Success", "Successfully Order Updated", nil)
}

func (c *OrderController) Deleted(ctx *gin.Context) {
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

	intParam, errIntParam := strconv.Atoi(ctx.Param("orderId"))

	if errIntParam != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  errIntParam.Error(),
		})
		return
	}

	queryHeader := c.db.Model(&entity.Orders{}).Where("id = ?", intParam).Delete(&entity.Orders{})
	if queryHeader.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  queryHeader.Error.Error(),
		})
		return
	}

	if queryHeader.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"result": "Invalid Request",
			"error":  "Order Not Found",
		})
		return
	}

	queryDetail := c.db.Model(&entity.OrderDetails{}).Where("order_id = ?", intParam).Delete(&entity.OrderDetails{})
	if queryDetail.Error != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"result": "Invalid Request",
			"error":  queryDetail.Error.Error(),
		})
		return
	}

	if queryDetail.RowsAffected == 0 {
		ctx.JSON(http.StatusOK, gin.H{
			"result": "Invalid Request",
			"error":  "Order Not Found",
		})
		return
	}

	responseWithJson(ctx, http.StatusOK, "Success", "Successfully Order Deleted", nil)

}

func (c *OrderController) mapOrderDetailEntityToDTO(ents entity.Orders) dto.OrderFindIdDTO {
	var details []dto.OrderDetailFindOrderIdDTO
	var orderUpdateAt string

	if !ents.UpdatedAt.IsZero() {
		orderUpdateAt = ents.UpdatedAt.Format("2006-01-02 15:04:05")
	}

	header := dto.OrderFindIdDTO{
		Id:          int64(ents.ID),
		CustomerId:  int64(ents.CustomerId),
		Title:       ents.Title,
		OrderNumber: ents.OrderNumber,
		Note:        ents.Note,
		Total:       fmt.Sprintf("%.f", ents.Total),
		CreatedAt:   ents.CreatedAt.Format("2006-01-02 15:04:05"),
		UpdateAt:    orderUpdateAt,
	}

	for _, val := range ents.OrderDetails {
		qty := strconv.FormatInt(val.Qty, 10)
		var orderDetailUpdatedAt string
		if !val.UpdatedAt.IsZero() {
			orderDetailUpdatedAt = val.UpdatedAt.Format("2006-01-02 15:04:05")
		}
		details = append(details, dto.OrderDetailFindOrderIdDTO{
			OrderId:   int64(val.OrderId),
			Item:      val.Item,
			Qty:       qty,
			Price:     fmt.Sprintf("%.f", val.Price),
			Amount:    fmt.Sprintf("%.f", val.Amount),
			CreatedAt: val.CreatedAt.Format("2006-01-02 15:04:05"),
			UpdateAt:  orderDetailUpdatedAt,
		})
	}

	header.OrderDetails = details
	return header
}
