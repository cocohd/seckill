package main

import (
	"github.com/kataras/iris"
)

func main() {
	app := iris.New()
	app.Logger().SetLevel("debug")
	template := iris.HTML("./backend/web/views", ".html").Layout(
		"shared/layout.html").Reload(true)
	app.RegisterView(template)

	// 设置模板目标
	app.HandleDir("assets", iris.Dir("./backend/web/assets"))
	// 出现异常跳转指定页面
	app.OnAnyErrorCode(func(c iris.Context) {
		c.ViewData("message", c.Values().GetStringDefault("message", "访问的页面出错！"))
		c.ViewLayout("")
		c.View("shared/error.html")
	})

	// 注册控制器

	// 启动服务
	app.Run(
		iris.Addr("localhost:8080"),
		iris.WithoutServerError(iris.ErrServerClosed),
		iris.WithOptimizations,
	)
}
