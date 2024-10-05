package gin

import "github.com/gin-gonic/gin"

func (h *Handler) HealthHandler(ctx *gin.Context) {
	ctx.JSON(200, gin.H{"status": "ok"})
}
