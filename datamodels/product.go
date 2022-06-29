package datamodels

type Product struct {
	ID           int64  `json:"id"  sql:"ID"  seckill:"id"`
	ProductName  string `json:"ProductName" sql:"productName" seckill:"ProductName"`
	ProductNum   int64  `json:"ProductNum" sql:"productNum" seckill:"ProcuctNum"`
	ProductImage string `json:"ProductImage" sql:"productImage" seckill:"ProductImage"`
	ProductUrl   string `json:"ProductUrl" sql:"productUrl" seckill:"ProductUrl"`
}
