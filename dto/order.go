package dto

type OrderCreateRequest struct {
	Title                    string                      `json:"title" binding:"required"`
	Note                     string                      `json:"note"  binding:"required"`
	OrderDetailCreateRequest []*OrderDetailCreateRequest `json:"orderDetails"  binding:"required"`
}

type OrderDetailCreateRequest struct {
	Item  string `json:"item" validate:"required"`
	Qty   string `json:"qty" validate:"required"`
	Price string `json:"price" validate:"required"`
}

type OrderFindManyDTO struct {
	Id          int64  `json:"id"`
	Title       string `json:"title"`
	OrderNumber string `json:"orderNumber"`
	Note        string `json:"note"`
	Total       string `json:"total"`
	CreatedBy   string `json:"createdBy"`
	CreatedAt   string `json:"createdAt"`
	UpdateAt    string `json:"updateAt"`
}

type OrderFindIdDTO struct {
	Id           int64                       `json:"id"`
	CustomerId   int64                       `json:"cudtomerId"`
	Title        string                      `json:"title"`
	OrderNumber  string                      `json:"orderNumber"`
	Note         string                      `json:"note"`
	Total        string                      `json:"total"`
	CreatedAt    string                      `json:"createdAt"`
	UpdateAt     string                      `json:"updateAt"`
	OrderDetails []OrderDetailFindOrderIdDTO `json:"orderDetails"`
}

type OrderDetailFindOrderIdDTO struct {
	OrderId   int64  `json:"orderId"`
	Item      string `json:"item"`
	Qty       string `json:"qty"`
	Price     string `json:"price"`
	Amount    string `json:"amount"`
	CreatedAt string `json:"createdAt"`
	UpdateAt  string `json:"updateAt"`
}
