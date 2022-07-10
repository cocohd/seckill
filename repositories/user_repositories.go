package repositories

import (
	"database/sql"
	"errors"
	"seckill/common"
	"seckill/datamodels"
)

type IUserRepository interface {
	Conn() error
	Insert(*datamodels.User) (int64, error)
	SelectByUserName(string) (*datamodels.User, error)
}

type UserManager struct {
	table string
	db    *sql.DB
}

func NewUserManager(table string, db *sql.DB) IUserRepository {
	return &UserManager{table: table, db: db}
}

func (u *UserManager) Conn() (err error) {
	if u.db == nil {
		u.db, err = common.NewMysqlConn()
		if err != nil {
			return err
		}
	}

	if u.table == "" {
		u.table = "user"
	}
	return nil
}

func (u *UserManager) Insert(user *datamodels.User) (userId int64, err error) {
	if err = u.Conn(); err != nil {
		return 0, err
	}

	sql := "insert into " + u.table + " set nickName=?, userName=?, hashedPwd=?"
	stmt, err := u.db.Prepare(sql)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(user.NickName, user.UserName, user.HashedPwd)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (u *UserManager) SelectByUserName(userName string) (user *datamodels.User, err error) {
	if userName == "" {
		return &datamodels.User{}, errors.New("用户名不能为空")
	}

	sql := "select * from " + u.table + " where userName=?"
	rows, err := u.db.Query(sql, userName)

	res := common.GetResultRow(rows)
	if len(res) == 0 {
		return &datamodels.User{}, err
	}

	user = &datamodels.User{}
	common.DataToStructByTagSql(res, user)
	return
}
