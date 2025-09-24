package handlers

import (
	"net/http"
	"wallet-api/pkg/httpserver"
	"wallet-api/pkg/httpserver/utils"
	"wallet-api/src/models"
	"wallet-api/src/services"

	"github.com/goccy/go-json"
)

type WalletHandler interface {
	Register(s httpserver.Router)
	FindById(w http.ResponseWriter, r *http.Request)
	ChangeBalance(w http.ResponseWriter, r *http.Request)
}

type walletHandler struct {
	walletService services.WalletService
}

func NewWalletHandler(walletService services.WalletService) WalletHandler {
	return &walletHandler{walletService: walletService}
}

func (h *walletHandler) Register(s httpserver.Router) {
	s.POST("/wallet", h.ChangeBalance).GET("/wallets/{WALLET_UUID}", h.FindById).GET("/wallets", h.GetWallets)
}

func (h *walletHandler) FindById(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	id, err := extractIdFromPath(r)
	if err != nil {
		utils.RespondError(w, http.StatusBadRequest, "incorrect wallet id")
		return
	}
	wallet, errResp := h.walletService.GetWalletByID(r.Context(), id)
	if errResp != nil {
		utils.RespondJSON(w, errResp.Code, errResp)
		return
	}
	utils.RespondJSON(w, http.StatusOK, wallet)
}

func (h *walletHandler) ChangeBalance(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var changeBalanceReq models.ChangeBalanceRequest
	if err := json.NewDecoder(r.Body).Decode(&changeBalanceReq); err != nil {
		utils.RespondError(w, http.StatusUnprocessableEntity, "unable read request body")
		return
	}
	if !validateOperationType(changeBalanceReq.OperationType) {
		utils.RespondError(w, http.StatusBadRequest, "incorrect operation type")
		return
	}
	if !validateAmount(changeBalanceReq.Balance) {
		utils.RespondError(w, http.StatusBadRequest, "amount must be more than zero")
		return
	}
	err := h.walletService.ChangeWalletBalance(r.Context(), changeBalanceReq)
	if err != nil {
		utils.RespondJSON(w, err.Code, err)
		return
	}
	utils.RespondJSON(w, http.StatusOK, map[string]string{"message": "ok"})
}

func (h *walletHandler) GetWallets(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	wallets, errResp := h.walletService.GetWallets(r.Context())
	if errResp != nil {
		utils.RespondJSON(w, errResp.Code, errResp)
		return
	}
	utils.RespondJSON(w, http.StatusOK, wallets)
}
