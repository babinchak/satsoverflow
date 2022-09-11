package controllers

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
}
