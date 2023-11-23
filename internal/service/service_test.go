package service

import (
	"errors"
	"github.com/shamank/warehouse-service/internal/domain/models"
	"github.com/shamank/warehouse-service/internal/domain/schemas"
	"github.com/shamank/warehouse-service/internal/service/mocks"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"testing"
)

func TestService_GetRemainingProducts(t *testing.T) {
	type Args struct {
		input  string
		output []models.Product
		error  error
	}

	type TestCase struct {
		warehouseUUID   string
		args            Args
		expectedError   error
		excepctedResult []schemas.Product
	}

	testCases := []TestCase{
		{
			warehouseUUID: "uuid",
			expectedError: nil,
			args: Args{
				input: "uuid",
				output: []models.Product{
					{

						UUID:     "94df01f6-4ead-482a-8bd5-84ee7e1aef57",
						Name:     "nike",
						Size:     "XL",
						Code:     "asd-xsdad",
						Quantity: 1,
					},
				},
				error: nil,
			},
			excepctedResult: []schemas.Product{
				{
					Name:     "nike",
					Size:     "XL",
					Code:     "asd-xsdad",
					Quantity: 1,
				},
			},
		},
		{
			warehouseUUID: "uuid",
			expectedError: errors.New("some error"),
			args: Args{
				input:  "uuid",
				output: nil,
				error:  errors.New("some error"),
			},
			excepctedResult: nil,
		},
	}

	for _, testCase := range testCases {
		repo := mocks.NewRepository(t)

		repo.On("GetRemainingProductsByWarehouse", testCase.args.input).Return(testCase.args.output, testCase.args.error)

		svc := NewService(repo, slog.Default())

		result, err := svc.GetRemainingProducts(testCase.warehouseUUID)

		assert.Equal(t, err, testCase.expectedError)
		assert.Equal(t, result, testCase.excepctedResult)
	}
}

func TestService_ReserveProducts(t *testing.T) {
	type productQuantityArgs struct {
		input               string
		productInWarehouses []models.WarehouseProduct
		error               error
	}

	type reserveProductArgs struct {
		input  []schemas.ProductWarehouseSplitted
		output error
	}

	type TestCase struct {
		productsToReserve   []string
		productQuantityArgs []productQuantityArgs
		reserveProductArgs  reserveProductArgs
		expectedError       error
	}

	testCase1 := TestCase{
		productsToReserve: []string{"product1", "product2", "product1"},
		productQuantityArgs: []productQuantityArgs{
			{
				input: "product1",
				productInWarehouses: []models.WarehouseProduct{
					{
						WarehouseUUID: "d10c8d17-6d15-445e-b643-6affa59aa26c",
						Quantity:      1,
					},
					{
						WarehouseUUID: "a00518e4-be6e-4eb7-9f95-bb52cc8b8548",
						Quantity:      1,
					},
				},
				error: nil,
			},
			{
				input: "product2",
				productInWarehouses: []models.WarehouseProduct{
					{
						WarehouseUUID: "233ef39e-bdea-41dc-a5a2-31c8f0e29d6e",
						Quantity:      1,
					},
				},
				error: nil,
			},
		},
		reserveProductArgs: reserveProductArgs{
			input: []schemas.ProductWarehouseSplitted{
				{
					ProductArticle: "product1",
					WarehouseData: []schemas.WarehouseCounter{
						{
							WarehouseUUID: "d10c8d17-6d15-445e-b643-6affa59aa26c",
							Count:         1,
						},
						{
							WarehouseUUID: "a00518e4-be6e-4eb7-9f95-bb52cc8b8548",
							Count:         1,
						},
					},
				},
				{
					ProductArticle: "product2",
					WarehouseData: []schemas.WarehouseCounter{
						{
							WarehouseUUID: "233ef39e-bdea-41dc-a5a2-31c8f0e29d6e",
							Count:         1,
						},
					},
				},
			},
			output: nil,
		},
		expectedError: nil,
	}
	testCase2 := TestCase{
		productsToReserve: []string{"product1", "product2", "product1"},
		productQuantityArgs: []productQuantityArgs{
			{
				input: "product1",
				productInWarehouses: []models.WarehouseProduct{
					{
						WarehouseUUID: "d10c8d17-6d15-445e-b643-6affa59aa26c",
						Quantity:      1,
					},
				},
				error: nil,
			},
			{
				input: "product2",
				productInWarehouses: []models.WarehouseProduct{
					{
						WarehouseUUID: "233ef39e-bdea-41dc-a5a2-31c8f0e29d6e",
						Quantity:      1,
					},
				},
				error: nil,
			},
		},
		expectedError: ErrNotEnoughProducts,
	}

	repo1 := mocks.NewRepository(t)
	repo2 := mocks.NewRepository(t)

	for _, productQuantityArgs := range testCase1.productQuantityArgs {
		repo1.On("GetProductsQuantity", productQuantityArgs.input).Maybe().Return(productQuantityArgs.productInWarehouses, productQuantityArgs.error)
	}

	for _, productQuantityArgs := range testCase2.productQuantityArgs {
		repo2.On("GetProductsQuantity", productQuantityArgs.input).Maybe().Return(productQuantityArgs.productInWarehouses, productQuantityArgs.error)
	}

	repo1.On("ReserveProducts", testCase1.reserveProductArgs.input).Once().Return(testCase1.reserveProductArgs.output)

	svc1 := NewService(repo1, slog.Default())
	svc2 := NewService(repo2, slog.Default())

	err := svc1.ReserveProducts(testCase1.productsToReserve)
	assert.Equal(t, err, testCase1.expectedError)

	err = svc2.ReserveProducts(testCase2.productsToReserve)
	assert.Equal(t, err, testCase2.expectedError)
}

func TestService_ReleaseProducts(t *testing.T) {
	type productQuantityArgs struct {
		input               string
		productInWarehouses []models.WarehouseProduct
		error               error
	}

	type releaseProductArgs struct {
		input  []schemas.ProductWarehouseSplitted
		output error
	}

	type TestCase struct {
		productsToRelease   []string
		productQuantityArgs []productQuantityArgs
		releaseProductArgs  releaseProductArgs
		expectedError       error
	}

	testCase1 := TestCase{
		productsToRelease: []string{"product1", "product2", "product1"},
		productQuantityArgs: []productQuantityArgs{
			{
				input: "product1",
				productInWarehouses: []models.WarehouseProduct{
					{
						WarehouseUUID:    "d10c8d17-6d15-445e-b643-6affa59aa26c",
						ReservedQuantity: 1,
					},
					{
						WarehouseUUID:    "a00518e4-be6e-4eb7-9f95-bb52cc8b8548",
						ReservedQuantity: 1,
					},
				},
				error: nil,
			},
			{
				input: "product2",
				productInWarehouses: []models.WarehouseProduct{
					{
						WarehouseUUID:    "233ef39e-bdea-41dc-a5a2-31c8f0e29d6e",
						ReservedQuantity: 1,
					},
				},
				error: nil,
			},
		},
		releaseProductArgs: releaseProductArgs{
			input: []schemas.ProductWarehouseSplitted{
				{
					ProductArticle: "product1",
					WarehouseData: []schemas.WarehouseCounter{
						{
							WarehouseUUID: "d10c8d17-6d15-445e-b643-6affa59aa26c",
							Count:         1,
						},
						{
							WarehouseUUID: "a00518e4-be6e-4eb7-9f95-bb52cc8b8548",
							Count:         1,
						},
					},
				},
				{
					ProductArticle: "product2",
					WarehouseData: []schemas.WarehouseCounter{
						{
							WarehouseUUID: "233ef39e-bdea-41dc-a5a2-31c8f0e29d6e",
							Count:         1,
						},
					},
				},
			},
			output: nil,
		},
		expectedError: nil,
	}

	testCase2 := TestCase{
		productsToRelease: []string{"product1", "product2", "product1"},
		productQuantityArgs: []productQuantityArgs{
			{
				input: "product1",
				productInWarehouses: []models.WarehouseProduct{
					{
						WarehouseUUID:    "d10c8d17-6d15-445e-b643-6affa59aa26c",
						ReservedQuantity: 1,
					},
					{
						WarehouseUUID:    "a00518e4-be6e-4eb7-9f95-bb52cc8b8548",
						ReservedQuantity: 1,
					},
				},
				error: nil,
			},
			{
				input: "product2",
				productInWarehouses: []models.WarehouseProduct{
					{
						WarehouseUUID:    "233ef39e-bdea-41dc-a5a2-31c8f0e29d6e",
						ReservedQuantity: 0,
					},
				},
				error: nil,
			},
		},
		expectedError: ErrNotEnoughProducts,
	}

	repo1 := mocks.NewRepository(t)
	repo2 := mocks.NewRepository(t)

	for _, productQuantityArgs := range testCase1.productQuantityArgs {
		repo1.On("GetProductsQuantity", productQuantityArgs.input).Maybe().Return(productQuantityArgs.productInWarehouses, productQuantityArgs.error)
	}

	for _, productQuantityArgs := range testCase2.productQuantityArgs {
		repo2.On("GetProductsQuantity", productQuantityArgs.input).Maybe().Return(productQuantityArgs.productInWarehouses, productQuantityArgs.error)
	}

	repo1.On("ReleaseProducts", testCase1.releaseProductArgs.input).Once().Return(testCase1.releaseProductArgs.output)

	svc1 := NewService(repo1, slog.Default())
	svc2 := NewService(repo2, slog.Default())

	err := svc1.ReleaseProducts(testCase1.productsToRelease)
	assert.Equal(t, err, testCase1.expectedError)

	err = svc2.ReleaseProducts(testCase2.productsToRelease)
	assert.Equal(t, err, testCase2.expectedError)
}
