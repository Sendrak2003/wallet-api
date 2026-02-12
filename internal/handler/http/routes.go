package http

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.Engine, h *WalletHandler) {
	api := r.Group("/api/v1")
	{
		api.POST("/wallet", h.Apply)
		api.GET("/wallets/:id", h.GetBalance)
	}
}
