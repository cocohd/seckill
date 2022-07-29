package main

import (
	"context"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"seckill/common"
	"seckill/fronted/web/controllers"
	"seckill/rabbitmq"
	"seckill/repositories"
	"seckill/services"
)

func main() {
	//1.创建iris 实例
	app := iris.New()
	//2.设置错误模式，在mvc模式下提示错误
	app.Logger().SetLevel("debug")
	//3.注册模板
	template := iris.HTML("./fronted/web/views", ".html").Layout("shared/layout.html").Reload(true)
	app.RegisterView(template)

	//4.设置模板
	app.HandleDir("/public", "./fronted/web/public")
	//设置静态html访问
	app.HandleDir("/html", "./fronted/web/htmlProductShow")

	//出现异常跳转到指定页面
	app.OnAnyErrorCode(func(ctx iris.Context) {
		ctx.ViewData("message", ctx.Values().GetStringDefault("message", "访问的页面出错！"))
		ctx.ViewLayout("")
		ctx.View("shared/error.html")
	})
	//连接数据库
	db, err := common.NewMysqlConn()
	if err != nil {

	}
	//sess := sessions.New(sessions.Config{
	//	Cookie:  "AdminCookie",
	//	Expires: 600 * time.Minute,
	//})

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	user := repositories.NewUserManager("user", db)
	userService := services.NewUserService(user)
	userPro := mvc.New(app.Party("/user"))
	//userPro.Register(userService, ctx, sess.Start)
	userPro.Register(userService, ctx)
	userPro.Handle(new(controllers.UserController))

	// 注册商品详情管理
	product := repositories.NewProductManager("product", db)
	productService := services.NewProductService(product)
	order := repositories.NewOrderManager("order", db)
	orderService := services.NewOrderService(order)

	// 注册rabbitMQ
	rabbitmq := rabbitmq.NewRabbitMQSimple("seckill")

	productPro := app.Party("/product")
	pro := mvc.New(productPro)
	pro.Register(productService, orderService, ctx, rabbitmq)
	pro.Handle(new(controllers.ProductController))

	app.Run(
		iris.Addr("localhost:8082"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)

}
