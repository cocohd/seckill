package controllers

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/kataras/iris/v12/sessions"
	"net/http"
	"seckill/datamodels"
	"seckill/encrypt"
	"seckill/services"
	"strconv"
)

type UserController struct {
	Ctx iris.Context
	// 此处需要使用接口，不然实例化始终为nil
	UserService services.IUserService
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

	fmt.Println("*/*/*user/*/*/*", nickName, userName, pwd)

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

	// 此处未加密cookie
	u.Ctx.SetCookie(&http.Cookie{Name: "uid", Value: strconv.FormatInt(user.ID, 10), Path: "/"})
	// 加密存储cookie
	// 不能使用strconv.Itoa(int(user.ID))
	mesByte := []byte(strconv.FormatInt(user.ID, 10))
	mesEncryptedString, err := encrypt.EncodeMes(mesByte)
	if err != nil {
		u.Ctx.Application().Logger().Debug(err)
	}
	u.Ctx.SetCookie(&http.Cookie{Name: "encryptedMessage", Value: mesEncryptedString, Path: "/"})

	// 改进：不使用服务端存储session
	//u.Session.Set("userId", strconv.FormatInt(user.ID, 10))

	return mvc.Response{
		Path: "/product/all",
	}
}
