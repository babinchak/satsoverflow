package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type User struct {
	ID           uint
	Name         string
	Email        *string
	Age          uint8
	Birthday     time.Time
	MemberNumber sql.NullString
	ActivatedAt  sql.NullTime
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

type Question struct {
	ID uint
	// Hash      string
	Title     string
	Body      string
	Bounty    uint
	CreatedAt time.Time
	UpdatedAt time.Time
}

type CreateQuestionInput struct {
	Title  string `json:"title" binding:"required"`
	Body   string `json:"body"`
	Bounty uint   `json:"bounty"`
}

type Answer struct {
	ID         uint
	Body       string
	Bounty     uint
	QuestionID uint
	CreatedAt  time.Time
	UpdatedAt  time.Time
}

type CreateAnswerInput struct {
	Body       string `json:"body" binding:"required"`
	QuestionID uint   `json:"question_id" binding:"required"`
}

func ReverseProxy(c *gin.Context) {
	remote, _ := url.Parse("http://localhost:3000")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request.Header
		req.Host = remote.Host
		req.URL = c.Request.URL
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}

	proxy.ServeHTTP(c.Writer, c.Request)
}

var wsupgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// func wshandler(w http.ResponseWriter, r *http.Request) {
// 	conn, err := wsupgrader.Upgrade(w, r, nil)
// 	if err != nil {
// 		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
// 		return
// 	}

// 	for {
// 		t, msg, err := conn.ReadMessage()
// 		if err != nil {
// 			break
// 		}
// 		fmt.Printf("Got message %s", msg)
// 		conn.WriteMessage(t, msg)
// 	}
// }

func main() {
	dsn := "host=localhost user=postgres password=dogsandcats123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Yay")
	}
	db.AutoMigrate(&User{}, &Question{}, &Answer{})
	client, err := lndclient.NewBasicClient("192.168.68.54:10009", "invoicer/tls.cert", "invoicer", "mainnet")
	if err != nil {
		log.Fatalf("Error setting up client: %v\n", err)
	}

	router := gin.Default()
	router.POST("/api/question", func(c *gin.Context) {
		input := CreateQuestionInput{}
		if err := c.BindJSON(&input); err != nil {
			log.Println("Hit here")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		post := Question{Title: input.Title, Body: input.Body, Bounty: input.Bounty}
		result := db.Create(&post)
		log.Println("ID:", post.ID, ", error:", result.Error, ", rows:", result.RowsAffected)
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.POST("/api/answer", func(c *gin.Context) {
		input := CreateAnswerInput{}
		if err := c.BindJSON(&input); err != nil {
			log.Println("Hit here")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}
		answer := Answer{Body: input.Body, QuestionID: input.QuestionID}
		result := db.Create(&answer)
		log.Println("ID:", answer.ID, ", error:", result.Error, ", rows:", result.RowsAffected)
		c.JSON(http.StatusOK, gin.H{"message": "pong"})
	})

	router.GET("/api/answers", func(c *gin.Context) {
		var Answers []Answer
		parameters := c.Request.URL.Query()
		id := parameters["question_id"][0]
		db.Where("question_id = ?", id).Order("id desc").Limit(5).Find(&Answers)
		answers := make([]map[string]interface{}, len(Answers))
		for i, post := range Answers {
			answers[i] = map[string]interface{}{"body": post.Body}
		}

		c.JSON(http.StatusOK, gin.H{
			"answers": answers,
		})
	})

	router.GET("/api/questions", func(c *gin.Context) {
		var Questions []Question
		db.Order("id desc").Limit(10).Find(&Questions)
		messages := make([]map[string]interface{}, len(Questions))
		for i, post := range Questions {
			fmt.Printf("ID %d, title = %s\n", post.ID, post.Title)
			messages[i] = map[string]interface{}{"title": post.Title, "id": post.ID, "bounty": post.Bounty}
		}

		c.JSON(http.StatusOK, gin.H{
			"messages": messages,
		})
	})

	router.GET("/api/question", func(c *gin.Context) {
		parameters := c.Request.URL.Query()
		id := parameters["id"][0]
		fmt.Printf("id: %s\n", id)
		var q Question
		db.Where("id = ?", id).First(&q)
		c.JSON(http.StatusOK, gin.H{
			"title":    q.Title,
			"body":     q.Body,
			"bounty":   q.Bounty,
			"created":  q.CreatedAt,
			"modified": q.UpdatedAt,
		})
	})

	router.GET("/api/invoice/ws", func(c *gin.Context) {
		// wshandler(c.Writer, c.Request)
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		if err != nil {
			fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
			return
		}
		defer conn.Close()

		for {
			// _, msg, err := conn.ReadMessage()

			var v struct {
				Body   string `json:"body"`
				Bounty int64  `json:"bounty"`
				Title  string `json:"title"`
			}
			err := conn.ReadJSON(&v)
			if err != nil {
				break
			}
			// fmt.Printf("Got message %v", v)
			fmt.Printf("Got message with title %s\n", v.Title)
			// conn.WriteMessage(t, msg)
			resp, err := client.AddInvoice(context.Background(), &lnrpc.Invoice{Memo: v.Title, Value: v.Bounty})
			if err != nil {
				log.Fatalf("Error adding invoice: %v\n", err)
			}
			fmt.Printf("Invoice added with payment_request: %s\n payment_addr: %x\n", resp.PaymentRequest, resp.PaymentAddr)
			conn.WriteMessage(websocket.TextMessage, []byte(resp.PaymentRequest))
		}
	})

	// Route every path not hitting api to nextjs
	router.NoRoute(ReverseProxy)
	router.Run(":8080")
}
