package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Health는 단순 헬스체크용 엔드포인트입니다.
func Health(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
	})
}
