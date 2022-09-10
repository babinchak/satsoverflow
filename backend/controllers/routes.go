package controllers

func (s *Server) initializeRoutes() {
	s.Router.POST("/api/register", s.Register)
	s.Router.POST("/api/login", s.Login)

	s.Router.POST("/api/answer", s.AddAnswer)
	s.Router.GET("/api/answers", s.ListAnswers)

	s.Router.POST("/api/question", s.AddQuestion)
	s.Router.GET("/api/questions", s.ListQuestions)
	s.Router.GET("/api/question", s.GetQuestion)

	s.Router.GET("/api/waitInvoicePaid", s.WaitInvoicePaid)

	s.Router.GET("/api/session", s.GetSessionDetails)
}
