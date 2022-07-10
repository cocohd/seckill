package controllers

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"net/http"
	"seckill/datamodels"
	"seckill/services"
	"strconv"
)

type UserController struct {
	Ctx         iris.Context
	UserService services.UserService
	Session     *sessions.Session
}

func (u *UserController) GetRegister() mvc.View {
	return mvc.View{
		Name: "user/register.html",
	}
}

func (u *UserController) PostRegister() {
	var (
		nickName = u.Ctx.FormValue("nickName")
		userName = u.Ctx.FormValue("userName")
		pwd      = u.Ctx.FormValue("password")
	)

	user := &datamodels.User{NickName: nickName, UserName: userName, HashedPwd: pwd}

	_, err := u.UserService.AddUser(user)
	if err != nil {
		u.Ctx.Redirect("/user/error")
	}
	u.Ctx.Redirect("/user/login")
}

func (u *UserController) GetLogin() mvc.View {
	return mvc.View{
		Name: "user/login.html",
	}
}

func (u *UserController) PostLogin() mvc.Response {
	// 获取用户提交的表单信息
	var (
		userName = u.Ctx.FormValue("userName")
		pwd      = u.Ctx.FormValue("password")
	)

	//	验证账号密码是否正确
	user, isOk := u.UserService.IsPwdSuc(userName, pwd)
	if !isOk {
		return mvc.Response{
			Path: "/user/login",
		}
	}

	// 写入用户ID到cookie中
	u.Ctx.SetCookie(&http.Cookie{Name: "uid", Value: strconv.FormatInt(user.ID, 10), Path: "/"})
	u.Session.Set("userId", strconv.FormatInt(user.ID, 10))

	return mvc.Response{
		Path: "/product/all",
	}
}
