package rest

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func NewRouter() http.Handler {
	router := gin.New()
	router.HandleMethodNotAllowed = true //Return 405 instead of 404 when the path exists but the method does not.

	router.Use(gin.Recovery()) // Turn panics inside handlers into 500 responses instead of crashing the server.

	router.GET("/health", health)

	return router
}
