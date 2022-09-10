package controllers

import (
	"fmt"
	"log"
	"net/http"

	"example.com/satsoverflow-backend/models"
	"github.com/gin-gonic/gin"
)

func (server *Server) GetSessionDetails(c *gin.Context) {
	session, err := server.Store.Get(c.Request, "sessionID")
	if err != nil {
		log.Fatalf("Error getting session: %v\n", err)
	}
	username, found := session.Values["username"]
	if !found {
		fmt.Printf("Not found\n")
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	var user models.User
	server.DB.Where("username = ?", username).First(&user)
	c.JSON(http.StatusOK, gin.H{
		"username":    username,
		"email":       user.Email,
		"createdDate": user.CreatedAt.String(),
	})
}
