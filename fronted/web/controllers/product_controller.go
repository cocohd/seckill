package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"html/template"
	"os"
	"path/filepath"
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

var (
	// 生成静态html保存地址
	htmlOutPath = "./fronted/web/htmlProductShow"
	// 静态文件模板目录
	templatePath = "./fronted/web/views/template"
)

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

// GetGenerateHtml 访问静态资源
func (p *ProductController) GetGenerateHtml() {
	productIDString := p.Ctx.URLParam("productID")
	productID, err := strconv.Atoi(productIDString)
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 获取模板
	contentsTemp, err := template.ParseFiles(filepath.Join(templatePath, "product.html"))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 获取静态html生成路径
	staticTemPath := filepath.Join(htmlOutPath, "htmlProduct.html")

	// 获取数据
	product, err := p.ProductService.GetProductByID(int64(productID))
	if err != nil {
		p.Ctx.Application().Logger().Debug(err)
	}

	// 生成静态文件
	generateStaticHtml(p.Ctx, contentsTemp, staticTemPath, product)
}

// 生成静态文件
func generateStaticHtml(ctx iris.Context, template *template.Template, staticTemName string, product *datamodels.Product) {
	// 判断静态文件是否存在
	if exist(staticTemName) {
		err := os.Remove(staticTemName)
		if err != nil {
			ctx.Application().Logger().Debug(err)
		}
	}

	//	新建静态文件
	file, err := os.OpenFile(staticTemName, os.O_CREATE|os.O_WRONLY, os.ModePerm)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
	defer file.Close()

	err = template.Execute(file, &product)
	if err != nil {
		ctx.Application().Logger().Debug(err)
	}
}

// 文件是否存在
func exist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil
}
