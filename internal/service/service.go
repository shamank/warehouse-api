package service

import (
	"errors"
	"github.com/shamank/warehouse-service/internal/domain/models"
	"github.com/shamank/warehouse-service/internal/domain/schemas"
	"github.com/shamank/warehouse-service/internal/handler"
	"log/slog"
	"sync"
)

var (
	ErrNotEnoughProducts = errors.New("not enough products in warehouses")
)

var _ handler.Service = (*Service)(nil)

//go:generate mockery --name=Repository
type Repository interface {
	GetRemainingProductsByWarehouse(warehouseUUID string) ([]models.Product, error)
	GetProductsQuantity(productArticle string) ([]models.WarehouseProduct, error)
	ReserveProducts(products []schemas.ProductWarehouseSplitted) error
	ReleaseProducts(productsWithSplit []schemas.ProductWarehouseSplitted) error
	//ReserveProduct(productArticle string, warehouseUUID string, quantity int) error
}

type Service struct {
	repo   Repository
	logger *slog.Logger
	mx     *sync.Mutex
}

func NewService(repo Repository, logger *slog.Logger) *Service {
	return &Service{
		repo:   repo,
		logger: logger,
		mx:     &sync.Mutex{},
	}
}

func (s *Service) GetRemainingProducts(warehouseUUID string) ([]schemas.Product, error) {
	products, err := s.repo.GetRemainingProductsByWarehouse(warehouseUUID)
	if err != nil {
		// TODO: обработка ошибки
		return nil, err
	}

	result := make([]schemas.Product, len(products))

	for i, product := range products {
		result[i] = schemas.Product{
			Name:     product.Name,
			Size:     product.Size,
			Code:     product.Code,
			Quantity: product.Quantity,
		}
	}

	return result, nil
}

func (s *Service) processProducts(productsToProcess []string, isRelease bool) ([]schemas.ProductWarehouseSplitted, error) {
	products := s.getProductWithCounts(productsToProcess)
	productsWithSplit := make([]schemas.ProductWarehouseSplitted, 0)

	for product, quantity := range products {
		warehouseData, err := s.processProduct(product, quantity, isRelease)
		if err != nil {
			return nil, err
		}

		productsWithSplit = append(productsWithSplit, schemas.ProductWarehouseSplitted{
			ProductArticle: product,
			WarehouseData:  warehouseData,
		})
	}

	return productsWithSplit, nil
}

func (s *Service) processProduct(product string, quantity int, isRelease bool) ([]schemas.WarehouseCounter, error) {
	warehouseData := make([]schemas.WarehouseCounter, 0)

	productInWarehouses, err := s.repo.GetProductsQuantity(product)
	if err != nil {
		return nil, err
	}

	for _, warehouseProduct := range productInWarehouses {

		warehouseQuantity := warehouseProduct.Quantity
		if isRelease {
			warehouseQuantity = warehouseProduct.ReservedQuantity
		}
		if warehouseQuantity >= quantity {
			warehouseData = append(warehouseData, schemas.WarehouseCounter{
				WarehouseUUID: warehouseProduct.WarehouseUUID,
				Count:         quantity,
			})
			quantity = 0
			break
		}

		quantity -= warehouseQuantity
		warehouseData = append(warehouseData, schemas.WarehouseCounter{
			WarehouseUUID: warehouseProduct.WarehouseUUID,
			Count:         warehouseQuantity,
		})
	}

	if quantity > 0 {
		return nil, ErrNotEnoughProducts
	}

	return warehouseData, nil
}

func (s *Service) ReserveProducts(productsToReserve []string) error {
	// т.к. во время выполнения операции может быть одновременно много запросов,
	// то для предотвращения ситуации, при которой product.quantity может уйти в минус
	// необходимо использовать мьютекс
	s.mx.Lock()
	defer s.mx.Unlock()

	productsWithSplit, err := s.processProducts(productsToReserve, false)
	if err != nil {
		s.logger.Error("error processing products", err)
		return err
	}

	err = s.repo.ReserveProducts(productsWithSplit)
	return err
}

// ReleaseProducts releases products based on the given condition
func (s *Service) ReleaseProducts(productsToRelease []string) error {
	s.mx.Lock()
	defer s.mx.Unlock()

	productsWithSplit, err := s.processProducts(productsToRelease, true)
	if err != nil {
		s.logger.Error("error processing products", err)
		return err
	}

	err = s.repo.ReleaseProducts(productsWithSplit)
	if err != nil {
		if errors.Is(err, ErrNotEnoughProducts) {
			return nil
		}
	}
	return nil
}

func (s *Service) getProductWithCounts(products []string) map[string]int {
	result := make(map[string]int)

	// т.к. товары могут повторяться, то нужно учитывать их количество
	for _, product := range products {
		result[product] += 1
	}
	return result
}
