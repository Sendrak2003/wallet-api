package http

import (
	"net/http"

	"wallet-api/internal/service"

	"github.com/gin-gonic/gin"
)

type WalletHandler struct {
	service *service.WalletService
}

type ApplyRequest struct {
	WalletID  string `json:"walletId" binding:"required" example:"11111111-1111-1111-1111-111111111111"`
	Operation string `json:"operationType" binding:"required" example:"DEPOSIT"`
	Amount    int64  `json:"amount" binding:"required" example:"1000"`
}

type BalanceResponse struct {
	Balance int64 `json:"balance" example:"1000"`
}

type ErrorResponse struct {
	Error string `json:"error" example:"invalid request"`
}

func NewWalletHandler(svc *service.WalletService) *WalletHandler {
	return &WalletHandler{service: svc}
}

// Apply godoc
// @Summary Применить операцию к кошельку
// @Description Выполняет операцию пополнения или снятия средств с кошелька
// @Tags wallet
// @Accept json
// @Produce json
// @Param request body ApplyRequest true "Данные операции"
// @Success 200 {object} BalanceResponse "Новый баланс"
// @Failure 400 {object} ErrorResponse "Ошибка валидации"
// @Failure 404 {object} ErrorResponse "Кошелек не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка"
// @Router /wallet [post]
func (h *WalletHandler) Apply(c *gin.Context) {
	var req ApplyRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request format"})
		return
	}

	balance, err := h.service.Apply(c.Request.Context(), req.WalletID, req.Operation, req.Amount)
	if err != nil {
		switch err.Error() {
		case "wallet not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "insufficient funds":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "amount must be positive":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "unknown operation type":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		case "invalid wallet id format":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}

// GetBalance godoc
// @Summary Получить баланс кошелька
// @Description Возвращает текущий баланс указанного кошелька
// @Tags wallet
// @Produce json
// @Param id path string true "ID кошелька" example(11111111-1111-1111-1111-111111111111)
// @Success 200 {object} BalanceResponse "Баланс кошелька"
// @Failure 404 {object} ErrorResponse "Кошелек не найден"
// @Failure 500 {object} ErrorResponse "Внутренняя ошибка"
// @Router /wallets/{id} [get]
func (h *WalletHandler) GetBalance(c *gin.Context) {
	walletID := c.Param("id")

	balance, err := h.service.GetBalance(c.Request.Context(), walletID)
	if err != nil {
		switch err.Error() {
		case "wallet not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "invalid wallet id format":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"balance": balance})
}
