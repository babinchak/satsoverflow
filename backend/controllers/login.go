package controllers

import (
	"fmt"
	"log"
	"net/http"

	"example.com/satsoverflow-backend/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func (server *Server) Register(c *gin.Context) {
	type RegisterInput struct {
		Username string `json:"username" binding:"required"`
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	input := RegisterInput{}
	if err := c.BindJSON(&input); err != nil {
		log.Println("Hit here")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// email := c.Query("email")
	// password := c.Query("password")
	fmt.Printf("Email = %s\n", input.Email)
	bytes, _ := bcrypt.GenerateFromPassword([]byte(input.Password), 14)

	newUser := models.User{Username: input.Username, Email: input.Email, Password: string(bytes)}
	result := server.DB.Create(&newUser)
	log.Println("Error:", result.Error, ", rows:", result.RowsAffected)
	// session, err := store.Get(c.Request, "sessionID")
	// if err != nil {
	// 	log.Fatalf("Error getting session: %v\n", err)
	// }
	// session.Options.HttpOnly = true
	// session.Options.MaxAge = 7 * 24 * 60 * 60
	c.JSON(http.StatusOK, gin.H{
		"email": newUser.Email,
		// "otheremail": email,
		// "password":   password,
	})
}

func (server *Server) Login(c *gin.Context) {
	type LoginInput struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}
	input := LoginInput{}
	if err := c.BindJSON(&input); err != nil {
		log.Println("Hit here")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	// session, err := store.Get(c.Request, "sessionID")
	// if err != nil {
	// 	log.Fatalf("Error getting session: %v\n", err)
	// }
	// session.Options.HttpOnly = true
	// session.Options.MaxAge = 7 * 24 * 60 * 60
	// email := c.Query("email")
	// password := c.Query("password")

	var user models.User
	res := server.DB.Where("username = ?", input.Username).First(&user)
	fmt.Printf("Rows: %d\n", res.RowsAffected)
	if res.Error != nil {
		// log.Fatalf("Error querying user by email: %v\n", res.Error)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error querying user by username: %v", res.Error))
		return
		// c.AbortWithStatus(http.StatusInternalServerError)
	}
	// user := users[0]
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password))
	if err != nil {
		// log.Fatalf("Error checking password: %v\n", err)
		c.String(http.StatusInternalServerError, fmt.Sprintf("Error checking password: %v\n", err))
		return
	}
	session, err := server.Store.Get(c.Request, "sessionID")
	if err != nil {
		log.Fatalf("Error getting session: %v\n", err)
	}
	// session.Values["Foo"] = "bar"
	session.Values["username"] = input.Username
	session.Options.HttpOnly = true
	session.Options.MaxAge = 7 * 24 * 60 * 60 // expires in 7 days
	// session.Options.Secure = true
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Fatalf("Error saving session: %v\n", err)
	}
	// c.JSON(http.StatusOK, gin.H{
	// 	"message": user.Email,
	// })

}

func (server *Server) Logout(c *gin.Context) {
	session, err := server.Store.Get(c.Request, "sessionID")
	if err != nil {
		log.Fatalf("Error getting session: %v\n", err)
	}
	delete(session.Values, "username")
	// session.Options.Secure = true
	err = session.Save(c.Request, c.Writer)
	if err != nil {
		log.Fatalf("Error saving session: %v\n", err)
	}
}
