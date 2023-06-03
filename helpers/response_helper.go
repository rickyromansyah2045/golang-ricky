package helpers

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func WriteJsonResponse(ctx *gin.Context, status int, payload interface{}) {
	ctx.JSON(status, payload)
}

func BadRequestResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusBadRequest, gin.H{
		"error": payload,
	})
}

func NotFoundResponse(ctx *gin.Context, payload interface{}) {
	WriteJsonResponse(ctx, http.StatusNotFound, gin.H{
		"error": payload,
	})
}

func InternalServerJsonResponse(ctx *gin.Context, err interface{}) {
	WriteJsonResponse(ctx, http.StatusInternalServerError, gin.H{
		"error": err,
	})
}

func UnauthorizeJsonResponse(ctx *gin.Context, err interface{}) {
	WriteJsonResponse(ctx, http.StatusUnauthorized, gin.H{
		"error": err,
	})
}
