package handlers

import (
	"context"
	"net/http/httptest"
	"testing"

	"github.com/KAnggara75/IDXStock/internal/models"
	"github.com/gofiber/fiber/v3"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

type mockBrokerUsecase struct {
	mock.Mock
}

func (m *mockBrokerUsecase) SyncBrokerActivity(ctx context.Context, token string, params models.SyncBrokerActivityParams) ([]models.BrokerActivity, error) {
	args := m.Called(ctx, token, params)
	return args.Get(0).([]models.BrokerActivity), args.Error(1)
}

func (m *mockBrokerUsecase) ManagePartitions(ctx context.Context) (*models.PartitionManagementResponse, error) {
	args := m.Called(ctx)
	return args.Get(0).(*models.PartitionManagementResponse), args.Error(1)
}

func TestSyncBrokerActivityHandler(t *testing.T) {
	app := fiber.New()
	mockUsecase := new(mockBrokerUsecase)
	handler := NewBrokerHandler(mockUsecase)

	app.Get("/broker/sync", handler.SyncBrokerActivityHandler)

	t.Run("unauthorized_missing_token", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/broker/sync", nil)
		resp, _ := app.Test(req)
		assert.Equal(t, fiber.StatusUnauthorized, resp.StatusCode)
	})

	t.Run("success_sync", func(t *testing.T) {
		mockUsecase.On("SyncBrokerActivity", mock.Anything, "Bearer test-token", mock.Anything).
			Return([]models.BrokerActivity{}, nil)

		req := httptest.NewRequest("GET", "/broker/sync", nil)
		req.Header.Set("Authorization", "Bearer test-token")
		resp, _ := app.Test(req)

		assert.Equal(t, fiber.StatusOK, resp.StatusCode)
		mockUsecase.AssertExpectations(t)
	})
}
