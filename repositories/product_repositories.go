package repositories

import (
	"database/sql"
	"seckill/common"
	"seckill/datamodels"
	"strconv"
)

type IProduct interface {
	Conn() error
	Insert(*datamodels.Product) (int64, error)
	Delete(int64) bool
	Update(*datamodels.Product) error
	SelectByKey(int64) (*datamodels.Product, error)
	SelectAll() ([]*datamodels.Product, error)
}

type ProductManager struct {
	table     string
	mysqlConn *sql.DB
}

func NewProductManager(table string, mysqlConn *sql.DB) IProduct {
	return &ProductManager{table: table, mysqlConn: mysqlConn}
}

func (p *ProductManager) Conn() (err error) {
	if p.mysqlConn == nil {
		p.mysqlConn, err = common.NewMysqlConn()
		if err != nil {
			return
		}
	}
	if p.table == "" {
		p.table = "product"
	}
	return
}

// 插入
func (p *ProductManager) Insert(product *datamodels.Product) (productId int64, err error) {
	if err = p.Conn(); err != nil {
		return
	}

	sql := "INSERT" + p.table + "SET productName=?,productNum=?,productImage=?,productUrl=? "
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return 0, err
	}

	res, err := stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return 0, err
	}

	productId, err = res.LastInsertId()
	return productId, err
}

func (p *ProductManager) Delete(productId int64) bool {
	if err := p.Conn(); err != nil {
		return false
	}

	sql := "DELETE FROM product Where ID=? "
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return false
	}
	// 注意这里面要换成字符串
	_, err = stmt.Exec(strconv.FormatInt(productId, 10))
	if err != nil {
		return false
	}
	return true
}

func (p *ProductManager) Update(product *datamodels.Product) error {
	if err := p.Conn(); err != nil {
		return err
	}

	sql := "update product set productName=?,productNum=?,productImage=?,productUrl=? where ID=" + strconv.FormatInt(product.ID, 10)
	stmt, err := p.mysqlConn.Prepare(sql)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(product.ProductName, product.ProductNum, product.ProductImage, product.ProductUrl)
	if err != nil {
		return err
	}
	return nil
}

func (p *ProductManager) SelectByKey(productID int64) (productResult *datamodels.Product, err error) {
	if err = p.Conn(); err != nil {
		return
	}

	sql := "select * from" + p.table + "where ID=" + strconv.FormatInt(productID, 10)

	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return
	}

	res := common.GetResultRow(rows)
	if len(res) == 0 {
		return
	}

	common.DataToStructByTagSql(res, productResult)
	return
}

func (p *ProductManager) SelectAll() (productArray []*datamodels.Product, errProduct error) {
	if err := p.Conn(); err != nil {
		return nil, err
	}

	sql := "select * from" + p.table

	rows, err := p.mysqlConn.Query(sql)
	defer rows.Close()
	if err != nil {
		return nil, err
	}

	res := common.GetResultRows(rows)
	if len(res) == 0 {
		return nil, err
	}

	for _, v := range res {
		product := &datamodels.Product{}
		common.DataToStructByTagSql(v, product)
		productArray = append(productArray, product)
	}
	return
}
