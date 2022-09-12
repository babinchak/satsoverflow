package controllers

import (
	"os"

	"github.com/dghubble/gologin/twitter"
	"github.com/dghubble/oauth1"
	twitterOAuth1 "github.com/dghubble/oauth1/twitter"
	"github.com/gin-gonic/gin"
)

func (s *Server) initializeRoutes() {
	s.Router.POST("/api/register", s.Register)
	s.Router.POST("/api/login", s.Login)
	s.Router.POST("/api/logout", s.Logout)

	s.Router.POST("/api/answer", s.AddAnswer)
	s.Router.GET("/api/answers", s.ListAnswers)

	s.Router.POST("/api/question", s.AddQuestion)
	s.Router.GET("/api/questions", s.ListQuestions)
	s.Router.GET("/api/question", s.GetQuestion)

	s.Router.GET("/api/waitInvoicePaid", s.WaitInvoicePaid)
	s.Router.POST("/api/deposit", s.AddFunds)
	s.Router.POST("/api/withdrawal", s.WithdrawalFunds)

	s.Router.GET("/api/profile", s.GetProfile)

	config := &oauth1.Config{
		ConsumerKey:    os.Getenv("TWITTER_API_KEY"),
		ConsumerSecret: os.Getenv("TWITTER_API_SECRET"),
		CallbackURL:    "http://localhost:8080/twitter/callback",
		Endpoint:       twitterOAuth1.AuthorizeEndpoint,
	}
	s.Router.Handle("GET", "/twitter/login", gin.WrapH(twitter.LoginHandler(config, nil)))
	s.Router.Handle("GET", "/twitter/callback", gin.WrapH(twitter.CallbackHandler(config, s.issueSession(), nil)))
	// s.Router.GET("/twitter/login", s.Twitter)
	// s.Router.GET("/twitter/callback", s.TwitterCallback)
}
