package repositories

import (
	"database/sql"
	"seckill/common"
	"seckill/datamodels"
	"strconv"
)

type IOrderRepository interface {
	Conn() error
	Insert(*datamodels.Order) (int64, error)
	Delete(int64) bool
	Update(order *datamodels.Order) error
	SelectByKey(int64) (*datamodels.Order, error)
	SelectAll() ([]*datamodels.Order, error)
	SelectAllWithInfo() (map[int]map[string]string, error)
}

type OrderManagerRepository struct {
	table     string
	mysqlConn *sql.DB
}

func NewOrderManager(table string, db *sql.DB) IOrderRepository {
	return &OrderManagerRepository{table, db}
}

func (o *OrderManagerRepository) Conn() (err error) {
	if o.mysqlConn == nil {
		o.mysqlConn, err = common.NewMysqlConn()
		if err != nil {
			return err
		}
	}

	o.table = "order"
	return nil
}

func (o *OrderManagerRepository) Insert(order *datamodels.Order) (orderId int64, err error) {
	if err = o.Conn(); err != nil {
		return 0, err
	}

	sql := "insert " + o.table + " set userID=?, productID=?, orderStatus=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}
	res, err := stmt.Exec(order.UserID, order.ProductID, order.OrderStatus)
	if err != nil {
		return 0, err
	}

	return res.LastInsertId()
}

func (o *OrderManagerRepository) Delete(orderID int64) bool {
	if err := o.Conn(); err != nil {
		return false
	}

	sql := "delete from " + o.table + " where ID=?"
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	_, err = stmt.Exec(strconv.FormatInt(orderID, 10))
	if err != nil {
		return false
	}
	return true
}

func (o *OrderManagerRepository) Update(order *datamodels.Order) (err error) {
	if err = o.Conn(); err != nil {
		return err
	}

	sql := "update " + o.table + " set userID=?, orderID=?, orderStatus=? where ID=" + strconv.FormatInt(order.ID, 10)
	stmt, err := o.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}
	_, err = stmt.Exec(order.UserID, order.ProductID, order.OrderStatus)
	if err != nil {
		return err
	}
	return nil
}

func (o *OrderManagerRepository) SelectByKey(int64) (order *datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return &datamodels.Order{}, err
	}

	sql := "select * from " + o.table + " where ID=" + strconv.FormatInt(order.ID, 10)
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return &datamodels.Order{}, err
	}

	res := common.GetResultRow(rows)
	common.DataToStructByTagSql(res, &datamodels.Order{})
	return
}

func (o *OrderManagerRepository) SelectAll() (orders []*datamodels.Order, err error) {
	if err = o.Conn(); err != nil {
		return []*datamodels.Order{}, err
	}

	sql := "select * from " + o.table
	rows, err := o.mysqlConn.Query(sql)

	res := common.GetResultRows(rows)
	for _, v := range res {
		order := &datamodels.Order{}
		common.DataToStructByTagSql(v, order)
		orders = append(orders, order)
	}
	return
}

func (o *OrderManagerRepository) SelectAllWithInfo() (orderMap map[int]map[string]string, err error) {
	if err = o.Conn(); err != nil {
		return
	}

	sql := "select o.ID, p.productName, o.orderStatus from seckill.order as o left join product as p on o.productID = p.ID"
	rows, err := o.mysqlConn.Query(sql)
	if err != nil {
		return nil, err
	}

	return common.GetResultRows(rows), err
}
