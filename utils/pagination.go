package utils

import (
	"strconv"

	"github.com/gin-gonic/gin"
)

type Pagination struct {
	Limit int    `json:"limit"`
	Page  int    `json:"page"`
	Sort  string `json:"sort"`
}

type Meta struct {
	TotalRow int64 `json:"totalRow"`
}

// GeneratePaginationFromRequest
func GeneratePaginationFromRequest(c *gin.Context) Pagination {
	// Initializing default
	limit := 20
	page := 1
	query := c.Request.URL.Query()

	for key, value := range query {
		queryValue := value[len(value)-1]
		switch key {
		case "limit":
			limit, _ = strconv.Atoi(queryValue)
		case "page":
			page, _ = strconv.Atoi(queryValue)
		}
	}

	return Pagination{
		Limit: limit,
		Page:  page,
	}
}
