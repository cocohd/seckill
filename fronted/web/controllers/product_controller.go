package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"seckill/datamodels"
	"seckill/services"
	"strconv"
)

type ProductController struct {
	Ctx            iris.Context
	ProductService services.IProductService
	OrderService   services.IOrderService
	Session        *sessions.Session
}

func (p *ProductController) GetDetail() mvc.View {
	//productIdString := p.Ctx.URLParam("productId")
	//productId, err := strconv.Atoi(productIdString)
	//if err != nil {
	//	p.Ctx.Application().Logger().Debug(err)
	//}
	//
	//product, err := p.ProductService.GetProductByID(int64(productId))
	fmt.Println("**********11111*************")
	product, err := p.ProductService.GetProductByID(2)
	fmt.Println("**********2222*************")

	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/view.html",
		Data: iris.Map{
			"product": product,
		},
	}
}

func (p *ProductController) GetOrder() mvc.View {
	userIdString := p.Ctx.GetCookie("uid")
	productIdString := p.Ctx.URLParam("productID")
	productId, err := strconv.Atoi(productIdString)

	fmt.Println("-----productId------", productId)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	product, err := p.ProductService.GetProductByID(int64(productId))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 判断商品数量是否够
	var orderID int64
	showMessage := "抢购失败！"
	if product.ProductNum > 0 {
		product.ProductNum--
		err := p.ProductService.UpdateProduct(product)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}

		//	创建订单
		userId, err := strconv.Atoi(userIdString)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		}

		orderID, err = p.OrderService.InsertOrder(&datamodels.Order{
			UserID:      int64(userId),
			ProductID:   int64(productId),
			OrderStatus: datamodels.OrderSuccess,
		})
		fmt.Println("插入订单、？？？", err)
		if err != nil {
			p.Ctx.Application().Logger().Debug(err)
		} else {
			showMessage = "抢购成功！"
		}
		fmt.Println("showMessage:", showMessage)

	}

	return mvc.View{
		Layout: "shared/productLayout.html",
		Name:   "product/result.html",
		Data: iris.Map{
			"orderID":     orderID,
			"showMessage": showMessage,
		},
	}
}
