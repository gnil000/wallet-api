package handlers_test

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"wallet-api/src/handlers"
	"wallet-api/src/models"
	"wallet-api/src/services"

	"github.com/goccy/go-json"
	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

func TestWalletHandler_FindById(t *testing.T) {
	mockService := new(services.WalletServiceMock)
	h := handlers.NewWalletHandler(mockService)

	validID := uuid.New().String()

	t.Run("wallet found", func(t *testing.T) {
		resp := models.GetBalanceResponse{Balance: 1000}
		mockService.On("GetWalletByID", mock.Anything, validID).Return(resp, nil)

		req := httptest.NewRequest(http.MethodGet, "/wallets/"+validID, nil)
		w := httptest.NewRecorder()

		h.FindById(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)

		var result models.GetBalanceResponse
		json.NewDecoder(w.Body).Decode(&result)
		assert.Equal(t, int64(1000), result.Balance)

		mockService.AssertExpectations(t)
	})

	t.Run("wallet not found", func(t *testing.T) {
		errResp := &models.ErrorResponse{Code: http.StatusNotFound, Message: "wallet not found"}
		mockService.On("GetWalletByID", mock.Anything, validID).Return(models.GetBalanceResponse{}, errResp)

		req := httptest.NewRequest(http.MethodGet, "/wallets/"+validID, nil)
		w := httptest.NewRecorder()

		h.FindById(w, req)

		assert.Equal(t, http.StatusNotFound, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})
}

func TestWalletHandler_ChangeBalance(t *testing.T) {
	mockService := new(services.WalletServiceMock)
	h := handlers.NewWalletHandler(mockService)

	walletID := uuid.New()
	validReq := models.ChangeBalanceRequest{
		ID:            walletID,
		Balance:       1000,
		OperationType: models.Operation_type_deposit,
	}

	t.Run("invalid json body", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/wallet", strings.NewReader("{bad json}"))
		w := httptest.NewRecorder()

		h.ChangeBalance(w, req)

		assert.Equal(t, http.StatusUnprocessableEntity, w.Result().StatusCode)
	})

	t.Run("invalid operation type", func(t *testing.T) {
		reqBody := `{"id":"` + walletID.String() + `","balance":1000,"operation_type":"INVALID"}`
		req := httptest.NewRequest(http.MethodPost, "/wallet", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		h.ChangeBalance(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("amount zero", func(t *testing.T) {
		reqBody := `{"id":"` + walletID.String() + `","balance":0,"operation_type":"DEPOSIT"}`
		req := httptest.NewRequest(http.MethodPost, "/wallet", strings.NewReader(reqBody))
		w := httptest.NewRecorder()

		h.ChangeBalance(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Result().StatusCode)
	})

	t.Run("service returns error", func(t *testing.T) {
		errResp := &models.ErrorResponse{Code: http.StatusInternalServerError, Message: "internal error"}
		mockService.On("ChangeWalletBalance", mock.Anything, validReq).Return(errResp)

		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest(http.MethodPost, "/wallet", strings.NewReader(string(body)))
		w := httptest.NewRecorder()

		h.ChangeBalance(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Result().StatusCode)
		mockService.AssertExpectations(t)
	})

	t.Run("success", func(t *testing.T) {
		mockService.On("ChangeWalletBalance", mock.Anything, validReq).Return(nil)

		body, _ := json.Marshal(validReq)
		req := httptest.NewRequest(http.MethodPost, "/wallet", strings.NewReader(string(body)))
		w := httptest.NewRecorder()

		h.ChangeBalance(w, req)

		assert.Equal(t, http.StatusOK, w.Result().StatusCode)
	})
}
