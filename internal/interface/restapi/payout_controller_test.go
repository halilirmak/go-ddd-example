package restapi_test

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/cryptoPickle/go-ddd-example/internal/application/command"
	"github.com/cryptoPickle/go-ddd-example/internal/application/dto"
	"github.com/cryptoPickle/go-ddd-example/internal/interface/restapi"
	"github.com/cryptoPickle/go-ddd-example/internal/interface/restapi/dto/response"
	"github.com/cryptoPickle/go-ddd-example/mocks"
)

type MockPayoutService struct {
	mock.Mock
}

func (m *MockPayoutService) CreatePayouts(ctx context.Context, cmd *command.CreatePayoutCommand) (*command.CreatePayoutCommandResult, error) {
	args := m.Called(ctx, cmd)
	result, ok := args.Get(0).(*command.CreatePayoutCommandResult)
	if !ok {
		return nil, args.Error(1)
	}

	return result, args.Error(1)
}

func TestPayoutController(t *testing.T) {
	gin.SetMode(gin.TestMode)
	mockService := new(MockPayoutService)
	mockLogger := new(mocks.MockLogger)
	controller := restapi.NewPayoutController(gin.Default(), mockService, mockLogger)

	t.Run("should return error when invalid json passed", func(t *testing.T) {
		invalidJSON := `{"some: "invalidData"}`
		req, err := http.NewRequest(http.MethodPost, "/payout", bytes.NewBuffer([]byte(invalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/payout", controller.CreatePayouts)

		router.ServeHTTP(rec, req)

		assert.Equal(t, http.StatusBadRequest, rec.Code)
	})

	t.Run("should return error when there is no id", func(t *testing.T) {
		invalidJSON := `{"some": "invalidData"}`
		req, err := http.NewRequest(http.MethodPost, "/payout", bytes.NewBuffer([]byte(invalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/payout", controller.CreatePayouts)

		router.ServeHTTP(rec, req)
		body, err := io.ReadAll(rec.Body)
		assert.NoError(t, err)
		expected := `{"error":"item ids expected"}`
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, expected, string(body))
	})

	t.Run("should return error when there is no currency", func(t *testing.T) {
		invalidJSON := `{"items": ["f8a787e9-09ca-4e80-8df6-a4633e441acd"]}`
		req, err := http.NewRequest(http.MethodPost, "/payout", bytes.NewBuffer([]byte(invalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/payout", controller.CreatePayouts)

		router.ServeHTTP(rec, req)
		body, err := io.ReadAll(rec.Body)
		assert.NoError(t, err)
		expected := `{"error":"currency expected"}`
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, expected, string(body))
	})

	t.Run("should return error when there currency not supported", func(t *testing.T) {
		invalidJSON := `{"items": ["f8a787e9-09ca-4e80-8df6-a4633e441acd"], "currency": "TRY"}`
		req, err := http.NewRequest(http.MethodPost, "/payout", bytes.NewBuffer([]byte(invalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		router := gin.Default()
		router.POST("/payout", controller.CreatePayouts)

		router.ServeHTTP(rec, req)
		body, err := io.ReadAll(rec.Body)
		assert.NoError(t, err)
		expected := `{"error":"currency not available (TRY)"}`
		assert.Equal(t, http.StatusBadRequest, rec.Code)
		assert.JSONEq(t, expected, string(body))
	})

	t.Run("should return 500 when service fails", func(t *testing.T) {
		invalidJSON := `{"items": ["f8a787e9-09ca-4e80-8df6-a4633e441acd"], "currency": "GBP"}`
		req, err := http.NewRequest(http.MethodPost, "/payout", bytes.NewBuffer([]byte(invalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		mockService.On("CreatePayouts", mock.Anything, mock.Anything).Return(nil, assert.AnError).Times(1)
		router := gin.Default()
		router.POST("/payout", controller.CreatePayouts)

		router.ServeHTTP(rec, req)
		body, err := io.ReadAll(rec.Body)
		assert.NoError(t, err)
		expected := `{"error":"something unexpected happened"}`
		assert.Equal(t, http.StatusInternalServerError, rec.Code)
		assert.JSONEq(t, expected, string(body))
	})

	t.Run("should return 200 when success", func(t *testing.T) {
		invalidJSON := `{"items": ["f8a787e9-09ca-4e80-8df6-a4633e441acd"], "currency": "GBP"}`
		req, err := http.NewRequest(http.MethodPost, "/payout", bytes.NewBuffer([]byte(invalidJSON)))
		assert.NoError(t, err)
		req.Header.Set("Content-Type", "application/json")
		rec := httptest.NewRecorder()

		result := command.CreatePayoutCommandResult{
			Currency: "GBP",
			Payout: []dto.BatchPayoutDTO{
				{
					ID:              uuid.New(),
					SellerReference: "somebody",
					Payouts: []dto.PayoutDTO{
						{
							ID:       uuid.New(),
							Amount:   300.00,
							Currency: "GBP",
						},
					},
				},
			},
		}

		mockService.On("CreatePayouts", mock.Anything, mock.Anything).Return(&result, nil).Times(1)
		router := gin.Default()
		router.POST("/payout", controller.CreatePayouts)

		router.ServeHTTP(rec, req)
		body, err := io.ReadAll(rec.Body)
		assert.NoError(t, err)

		data := response.CreatePayoutResponse{}
		err = json.Unmarshal(body, &data)
		assert.NoError(t, err)

		assert.Equal(t, http.StatusCreated, rec.Code)
		assert.Equal(t, "GBP", data.Currency)
		assert.Equal(t, data.BatchPayouts[0].SellerReference, "somebody")
		assert.Equal(t, data.BatchPayouts[0].Payouts[0].Amount, float32(300.00))
	})
}
