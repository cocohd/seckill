package datamodels

// Message rabbitmq消息体
type Message struct {
	ProductId int64
	UserId    int64
}

// NewMessage 构造消息体
func NewMessage(productId, userId int64) *Message {
	return &Message{productId, userId}
}
