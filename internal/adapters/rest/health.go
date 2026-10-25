package rest

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

type healthResponse struct {
	Status string `json:"status"`
}

func health(c *gin.Context) {
	c.JSON(http.StatusOK, healthResponse{Status: "ok"})
}
