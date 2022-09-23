package models

type Order struct {
	Model
	TransactionId   string      `json:"transaction_id" gorm:"null"`
	UserId          uint        `json:"user_id"`
	Code            string      `json:"code"`
	AmbassadorEmail string      `json:"ambassador_email"`
	Name            string      `json:"name"`
	Email           string      `json:"email"`
	Address         string      `json:"address" gorm:"null"`
	City            string      `json:"city" gorm:"null"`
	Country         string      `json:"country" gorm:"null"`
	Postcode        string      `json:"postcode" gorm:"null"`
	Complete        bool        `json:"complete" gorm:"default:false"`
	Total           float64     `json:"total" gorm:"-"`
	OrderItems      []OrderItem `json:"order_items" gorm:"foreignKey:OrderId"`
}

type OrderItem struct {
	Model
	OrderId           uint    `json:"order_id"`
	ProductTitle      string  `json:"product_title"`
	Quantity          uint    `json:"quantity"`
	Price             float64 `json:"price"`
	AdminRevenue      float64 `json:"admin_revenue"`
	AmbassadorRevenue float64 `json:"ambassador_revenue"`
}

func (order *Order) GetTotal() float64 {
	var total float64 = 0

	for _, item := range order.OrderItems {
		total += item.Price * float64(item.Quantity)
	}

	return total
}
