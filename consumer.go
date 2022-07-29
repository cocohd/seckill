package main

import (
	"fmt"
	"seckill/common"
	"seckill/rabbitmq"
	"seckill/repositories"
	"seckill/services"
)

/*rabbit 消费代码*/

func main() {
	// 初始化mysql连接
	db, err := common.NewMysqlConn()
	if err != nil {
		fmt.Println(err)
	}
	// 初始化product
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)

	// 初始化order
	order := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(order)

	// 以Simple模式消费
	rabbitmqSimple := rabbitmq.NewRabbitMQSimple("seckill")
	rabbitmqSimple.ConsumeSimple(orderService, productService)
}
