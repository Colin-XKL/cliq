package router

import (
    "github.com/gin-gonic/gin"

    "cliq-hub-backend/internal/http/handlers"
    "cliq-hub-backend/internal/llm"
)

func New(client llm.Client) *gin.Engine {
    r := gin.Default()
    h := handlers.NewGenerateHandler(client)

    v1 := r.Group("/v1")
    tm := v1.Group("/templates")
    tm.POST("/generate", h.Handle)
    return r
}