package datamodels

type Order struct {
	ID          int64 `sql:"orderID" seckill:"id"`
	UserID      int64 `sql:"userID" seckill:"UserID"`
	ProductID   int64 `sql:"productID" seckill:"ProductID"`
	OrderStatus int64 `sql:"orderStatus" seckill:"OrderStatus"`
}

const (
	OrderWait = iota
	OrderSuccess
	OrderFailed
)
