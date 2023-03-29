package service

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// This function register the particular adapter based on providerID
func RegisterHandler() func(ctx *gin.Context) {
	return func(c *gin.Context) {
		fmt.Println("Hello from Register Handler")
	}
}
