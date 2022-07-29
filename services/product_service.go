package services

import (
	"fmt"
	"seckill/datamodels"
	"seckill/repositories"
)

type IProductService interface {
	GetProductByID(int64) (*datamodels.Product, error)
	GetAllProducts() ([]*datamodels.Product, error)
	DeleteProductByID(int64) bool
	InsertProduct(product *datamodels.Product) (int64, error)
	UpdateProduct(product *datamodels.Product) error
	// SubProductNum 秒杀成功，数量减一
	SubProductNum(int64) error
}

type ProductService struct {
	productRepository repositories.IProduct
}

// NewProductService 初始化
func NewProductService(repository repositories.IProduct) IProductService {
	return &ProductService{repository}
}

func (p *ProductService) GetProductByID(productID int64) (*datamodels.Product, error) {
	fmt.Println("********product_service**************")
	return p.productRepository.SelectByKey(productID)
}

func (p *ProductService) GetAllProducts() ([]*datamodels.Product, error) {
	return p.productRepository.SelectAll()
}

func (p *ProductService) DeleteProductByID(productID int64) bool {
	return p.productRepository.Delete(productID)
}

func (p *ProductService) InsertProduct(product *datamodels.Product) (int64, error) {
	return p.productRepository.Insert(product)
}

func (p *ProductService) UpdateProduct(product *datamodels.Product) error {
	return p.productRepository.Update(product)
}

func (p *ProductService) SubProductNum(productId int64) error {
	return p.productRepository.SubProductNum(productId)
}
