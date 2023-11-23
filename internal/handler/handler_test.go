package handler

import (
	"bytes"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/shamank/warehouse-service/internal/domain/schemas"
	"github.com/shamank/warehouse-service/internal/handler/mocks"
	"github.com/stretchr/testify/assert"
	"log/slog"
	"net/http/httptest"
	"testing"
)

func TestGetRemainingProducts(t *testing.T) {
	type Args struct {
		input  string
		output []schemas.Product
		error  error
	}

	type TestCase struct {
		url  string
		args Args

		expectedStatusCode int
		expectedResult     string
	}

	testCases := []TestCase{
		{
			url: "/getRemainingProducts?warehouse_uuid=e4aa0556-aec5-41d4-8280-885865842719",
			args: Args{
				input: "e4aa0556-aec5-41d4-8280-885865842719",
				output: []schemas.Product{
					{
						Name:     "nike",
						Size:     "XL",
						Code:     "asd-xsdad",
						Quantity: 1,
					},
				},
				error: nil,
			},
			expectedStatusCode: 200,
			expectedResult:     `[{"name":"nike","size":"XL","code":"asd-xsdad","quantity":1}]`,
		},
		{
			url:                "/getRemainingProducts",
			expectedStatusCode: 400,
			expectedResult:     `{"error":"warehouse_uuid is required"}`,
		},
	}

	for _, testCase := range testCases {
		service := mocks.NewService(t)
		service.On("GetRemainingProducts", testCase.args.input).Return(testCase.args.output, nil).Maybe()

		handler := NewHandler(service, slog.Default())

		r := gin.New()

		r.GET("/getRemainingProducts", handler.getRemainingProducts)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("GET", testCase.url, nil)

		r.ServeHTTP(w, req)

		assert.Equal(t, testCase.expectedStatusCode, w.Code)
		assert.Equal(t, testCase.expectedResult, w.Body.String())
	}
}

func TestReserveProducts(t *testing.T) {
	type Args struct {
		input []string
		error error
	}

	type TestCase struct {
		url  string
		body string
		args Args

		expectedStatusCode int
		expectedResult     string
	}

	testCases := []TestCase{
		{
			url:  "/reserveProducts",
			body: `["a1as1","xd123ed12fg"]`,
			args: Args{
				input: []string{"a1as1", "xd123ed12fg"},
				error: nil,
			},
			expectedStatusCode: 200,
			expectedResult:     `{"message":"OK"}`,
		},
		{
			url:                "/reserveProducts",
			body:               "[12312, 3231, true]",
			expectedStatusCode: 400,
			expectedResult:     `{"error":"json: cannot unmarshal number into Go value of type string"}`,
		},
		{
			url:  "/reserveProducts",
			body: `["a1as1","xd123ed12fg"]`,
			args: Args{
				input: []string{"a1as1", "xd123ed12fg"},
				error: errors.New("not enough products in warehouses"),
			},
			expectedStatusCode: 400,
			expectedResult:     `{"error":"not enough products in warehouses"}`,
		},
	}

	for _, testCase := range testCases {
		service := mocks.NewService(t)
		service.On("ReserveProducts", testCase.args.input).Return(testCase.args.error).Maybe()

		handler := NewHandler(service, slog.Default())

		r := gin.New()

		r.POST("/reserveProducts", handler.reserveProducts)

		w := httptest.NewRecorder()

		req := httptest.NewRequest("POST", testCase.url, bytes.NewBufferString(testCase.body))

		r.ServeHTTP(w, req)

		assert.Equal(t, testCase.expectedStatusCode, w.Code)
		assert.Equal(t, testCase.expectedResult, w.Body.String())

	}
}

func TestReleaseProducts(t *testing.T) {
	type Args struct {
		input []string
		error error
	}

	type TestCase struct {
		url  string
		body string
		args Args

		expectedStatusCode int
		expectedResult     string
	}

	testCases := []TestCase{
		{
			url:  "/releaseProducts",
			body: `["a1as1","xd123ed12fg"]`,
			args: Args{
				input: []string{"a1as1", "xd123ed12fg"},
				error: nil,
			},
			expectedStatusCode: 200,
			expectedResult:     `{"message":"OK"}`,
		},
		{
			url:                "/releaseProducts",
			body:               "[true, 3231, true]",
			expectedStatusCode: 400,
			expectedResult:     `{"error":"json: cannot unmarshal bool into Go value of type string"}`,
		},
		{
			url:  "/releaseProducts",
			body: `["a1as1","xd123ed12fg"]`,
			args: Args{
				input: []string{"a1as1", "xd123ed12fg"},
				error: errors.New("not enough products in warehouses"),
			},
			expectedStatusCode: 400,
			expectedResult:     `{"error":"not enough products in warehouses"}`,
		},
	}

	for _, testCase := range testCases {
		service := mocks.NewService(t)
		service.On("ReleaseProducts", testCase.args.input).Return(testCase.args.error).Maybe()

		handler := NewHandler(service, slog.Default())

		r := gin.New()

		r.POST("/releaseProducts", handler.releaseProducts)

		w := httptest.NewRecorder()
		req := httptest.NewRequest("POST", testCase.url, bytes.NewBufferString(testCase.body))

		r.ServeHTTP(w, req)

		assert.Equal(t, testCase.expectedStatusCode, w.Code)
		assert.Equal(t, testCase.expectedResult, w.Body.String())

	}

}
