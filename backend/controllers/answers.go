package controllers

import (
	"fmt"
	"log"
	"net/http"

	"example.com/satsoverflow-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lnwire"
)

const INVOICE_EXPIRY_SECS = 600

type CreateAnswerInput struct {
	Body       string `json:"body" binding:"required"`
	QuestionID uint   `json:"question_id" binding:"required"`
}

func (server *Server) AddAnswer(c *gin.Context) {
	input := CreateAnswerInput{}
	if err := c.BindJSON(&input); err != nil {
		log.Println("Hit here")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	answer := models.Answer{Body: input.Body, QuestionID: input.QuestionID}
	result := server.DB.Create(&answer)
	log.Println("ID:", answer.ID, ", error:", result.Error, ", rows:", result.RowsAffected)
	c.JSON(http.StatusOK, gin.H{"message": "pong"})
}

func (server *Server) ListAnswers(c *gin.Context) {
	var Answers []models.Answer
	parameters := c.Request.URL.Query()
	id := parameters["question_id"][0]
	server.DB.Where("question_id = ?", id).Order("id desc").Limit(5).Find(&Answers)
	answers := make([]map[string]interface{}, len(Answers))
	for i, post := range Answers {
		answers[i] = map[string]interface{}{"body": post.Body}
	}

	c.JSON(http.StatusOK, gin.H{
		"answers": answers,
	})
}

type CreateQuestionInput struct {
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body"`
	Bounty uint   `json:"bounty"`
}

func (server *Server) AddQuestion(c *gin.Context) {
	input := CreateQuestionInput{}
	if err := c.BindJSON(&input); err != nil {
		log.Println("Hit here")
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

	// fmt.Printf("Got message with title %s\n", v.Title)
	// conn.WriteMessage(t, msg)
	msats := lnwire.MilliSatoshi(input.Bounty * 1000)
	hash, payaddr, err := server.LndServices.Client.AddInvoice(c.Request.Context(), &invoicesrpc.AddInvoiceData{Memo: input.Title, Value: msats, Expiry: INVOICE_EXPIRY_SECS})
	if err != nil {
		log.Fatalf("Error adding invoice: %v\n", err)
	}
	fmt.Printf("Invoice added with payment_request: %s\n payment_addr: %x\n", hash.String(), payaddr)
	// conn.WriteMessage(websocket.TextMessage, []byte(payaddr))

	post := models.Question{Title: input.Title, Body: input.Body, Bounty: input.Bounty, Paid: false, Hash: hash.String()}
	result := server.DB.Create(&post)
	log.Println("ID:", post.ID, ", error:", result.Error, ", rows:", result.RowsAffected)
	c.JSON(http.StatusOK, gin.H{
		"payment_request": payaddr,
		"hash":            hash.String(),
	})
}

func (server *Server) GetQuestion(c *gin.Context) {
	parameters := c.Request.URL.Query()
	id := parameters["id"][0]
	fmt.Printf("id: %s\n", id)
	var q models.Question
	res := server.DB.Where("id = ?", id).First(&q)
	// fmt.Printf("res err = %v\n", res.Error)
	// fmt.Printf("id = %d, Title = %s\n", q.ID, q.Title)
	if res.Error != nil || q.Paid == false {
		c.Status(http.StatusNotFound)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"title":    q.Title,
		"body":     q.Body,
		"bounty":   q.Bounty,
		"created":  q.CreatedAt,
		"modified": q.UpdatedAt,
	})
}

func (server *Server) ListQuestions(c *gin.Context) {
	var Questions []models.Question
	server.DB.Where("paid = ?", true).Order("id desc").Limit(10).Find(&Questions)
	messages := make([]map[string]interface{}, len(Questions))
	for i, post := range Questions {
		fmt.Printf("ID %d, title = %s\n", post.ID, post.Title)
		messages[i] = map[string]interface{}{"title": post.Title, "id": post.ID, "bounty": post.Bounty}
	}

	c.JSON(http.StatusOK, gin.H{
		"messages": messages,
	})
}
