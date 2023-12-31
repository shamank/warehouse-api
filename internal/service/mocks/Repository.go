// Code generated by mockery v2.38.0. DO NOT EDIT.

package mocks

import (
	models "github.com/shamank/warehouse-service/internal/domain/models"
	schemas "github.com/shamank/warehouse-service/internal/domain/schemas"
	mock "github.com/stretchr/testify/mock"
)

// Repository is an autogenerated mock type for the Repository type
type Repository struct {
	mock.Mock
}

// GetProductsQuantity provides a mock function with given fields: productArticle
func (_m *Repository) GetProductsQuantity(productArticle string) ([]models.WarehouseProduct, error) {
	ret := _m.Called(productArticle)

	if len(ret) == 0 {
		panic("no return value specified for GetProductsQuantity")
	}

	var r0 []models.WarehouseProduct
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]models.WarehouseProduct, error)); ok {
		return rf(productArticle)
	}
	if rf, ok := ret.Get(0).(func(string) []models.WarehouseProduct); ok {
		r0 = rf(productArticle)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.WarehouseProduct)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(productArticle)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// GetRemainingProductsByWarehouse provides a mock function with given fields: warehouseUUID
func (_m *Repository) GetRemainingProductsByWarehouse(warehouseUUID string) ([]models.Product, error) {
	ret := _m.Called(warehouseUUID)

	if len(ret) == 0 {
		panic("no return value specified for GetRemainingProductsByWarehouse")
	}

	var r0 []models.Product
	var r1 error
	if rf, ok := ret.Get(0).(func(string) ([]models.Product, error)); ok {
		return rf(warehouseUUID)
	}
	if rf, ok := ret.Get(0).(func(string) []models.Product); ok {
		r0 = rf(warehouseUUID)
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).([]models.Product)
		}
	}

	if rf, ok := ret.Get(1).(func(string) error); ok {
		r1 = rf(warehouseUUID)
	} else {
		r1 = ret.Error(1)
	}

	return r0, r1
}

// ReleaseProducts provides a mock function with given fields: productsWithSplit
func (_m *Repository) ReleaseProducts(productsWithSplit []schemas.ProductWarehouseSplitted) error {
	ret := _m.Called(productsWithSplit)

	if len(ret) == 0 {
		panic("no return value specified for ReleaseProducts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]schemas.ProductWarehouseSplitted) error); ok {
		r0 = rf(productsWithSplit)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// ReserveProducts provides a mock function with given fields: products
func (_m *Repository) ReserveProducts(products []schemas.ProductWarehouseSplitted) error {
	ret := _m.Called(products)

	if len(ret) == 0 {
		panic("no return value specified for ReserveProducts")
	}

	var r0 error
	if rf, ok := ret.Get(0).(func([]schemas.ProductWarehouseSplitted) error); ok {
		r0 = rf(products)
	} else {
		r0 = ret.Error(0)
	}

	return r0
}

// NewRepository creates a new instance of Repository. It also registers a testing interface on the mock and a cleanup function to assert the mocks expectations.
// The first argument is typically a *testing.T value.
func NewRepository(t interface {
	mock.TestingT
	Cleanup(func())
}) *Repository {
	mock := &Repository{}
	mock.Mock.Test(t)

	t.Cleanup(func() { mock.AssertExpectations(t) })

	return mock
}
