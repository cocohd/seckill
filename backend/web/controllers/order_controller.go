package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"seckill/common"
	"seckill/datamodels"
	"seckill/services"
	"strconv"
)

type OrderController struct {
	Ctx           iris.Context
	IOrderService services.IOrderService
}

func (o *OrderController) Get() mvc.View {
	orderArray, err := o.IOrderService.GetAllOrderInfo()
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	return mvc.View{
		Name: "order/view.html",
		Data: iris.Map{
			"order": orderArray,
		},
	}
}

//func (o *OrderController) GetOne() mvc.View {
//
//}

func (o *OrderController) GetDelete() {
	idString := o.Ctx.URLParam("id")
	id, err := strconv.ParseInt(idString, 10, 16)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	isDel := o.IOrderService.DeleteOrderByID(id)
	if isDel {
		o.Ctx.Application().Logger().Debug("删除商品成功，id：" + idString)
	} else {
		o.Ctx.Application().Logger().Debug("删除商品失败，id：" + idString)
	}

	o.Ctx.Redirect("/order/view.html")
}

func (o *OrderController) PostUpdate() {
	order := &datamodels.Order{}
	o.Ctx.Request().ParseForm()

	dec := common.NewDecoder(&common.DecoderOptions{TagName: "seckill"})
	if err := dec.Decode(o.Ctx.Request().Form, order); err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	err := o.IOrderService.UpdateOrder(order)
	if err != nil {
		o.Ctx.Application().Logger().Debug(err)
	}

	o.Ctx.Redirect("/order/view.html")
}
