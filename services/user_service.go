package services

import (
	"errors"
	"golang.org/x/crypto/bcrypt"
	"seckill/datamodels"
	"seckill/repositories"
)

type IUserService interface {
	AddUser(*datamodels.User) (int64, error)
	IsPwdSuc(string, string) (*datamodels.User, bool)
}

type UserService struct {
	userRepository repositories.IUserRepository
}

func NewUserService(userRepository repositories.IUserRepository) IUserService {
	return &UserService{userRepository: userRepository}
}

func (u *UserService) AddUser(user *datamodels.User) (id int64, err error) {
	pwdDecoded, err := bcrypt.GenerateFromPassword([]byte(user.HashedPwd), bcrypt.DefaultCost)
	if err != nil {
		return 0, err
	}
	user.HashedPwd = string(pwdDecoded)
	return u.userRepository.Insert(user)
}

func (u *UserService) IsPwdSuc(userName string, pwd string) (user *datamodels.User, isOk bool) {
	user, err := u.userRepository.SelectByUserName(userName)
	if err != nil {
		return &datamodels.User{}, false
	}

	isOk, _ = ValidatePwd(user.HashedPwd, pwd)
	if !isOk {
		return &datamodels.User{}, false
	}
	return
}

func ValidatePwd(hashedPwd, pwd string) (isOk bool, err error) {
	// 将 bcrypt 散列密码与其可能的明文等效密码进行比较。 成功时返回 nil，失败时返回错误。
	if err = bcrypt.CompareHashAndPassword([]byte(hashedPwd), []byte(pwd)); err != nil {
		return false, errors.New("密码错误")
	}
	return true, nil
}
