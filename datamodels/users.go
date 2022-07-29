package datamodels

type User struct {
	ID        int64  `json:"id" sql:"ID" form:"ID"`
	NickName  string `json:"NickName" sql:"nickName" form:"nickName"`
	UserName  string `json:"UserName" sql:"userName" form:"userName"`
	HashedPwd string `json:"HashPwd" sql:"hashedPwd" form:"hashPwd"`
}
