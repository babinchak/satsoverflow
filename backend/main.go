package main

import (
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
	"github.com/lightningnetwork/lnd/channeldb"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/lnwire"
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
	Paid      bool
	Hash      string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// type PendingQuestion struct {
// 	ID        uint
// 	Title     string
// 	Body      string
// 	Bounty    uint
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

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
// type MacaroonC struct {

// }

func main() {
	dsn := "host=localhost user=postgres password=dogsandcats123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Yay")
	}
	db.AutoMigrate(&User{}, &Question{}, &Answer{})
	// var opts []grpc.DialOption

	// tls_creds, err := credentials.NewClientTLSFromFile("invoicer/tls.cert", "localhost")
	// if err != nil {
	// 	log.Fatalf("Error getting tls creds: %v\n", err)
	// }
	// opts = append(opts, grpc.WithTransportCredentials(tls_creds))

	// macBytes, err := ioutil.ReadFile("invoices/admin.macaroon")
	// if err != nil {
	// 	log.Fatalf("Error loading macaroon file: %v\n", err)
	// }

	// func metadata_callback(ctx context.Context, uri ...string) (map[string]string, error) {

	// }
	// mac_creds := credentials.PerRPCCredentials{}

	// conn, err := grpc.Dial("192.168.68.54:10009", opts...)
	// if err != nil {
	// 	log.Fatalf("fail to dial: %v", err)
	// }
	// defer conn.Close()
	// client := pb.NewLightningClient()
	// lightning_client, err := lndclient.NewBasicClient("192.168.68.54:10009", "invoicer/tls.cert", "invoicer", "mainnet")
	// if err != nil {
	// 	log.Fatalf("Error setting up client: %v\n", err)
	// }
	// invoice_client, err := lndclient.NewInvoicesClient("192.168.68.54:10009", "invoicer/tls.cert", "invoicer", "mainnet")
	// if err != nil {
	// 	log.Fatalf("Error setting up client: %v\n", err)
	// }
	lndcfg := lndclient.LndServicesConfig{
		LndAddress:  "192.168.68.54:10009",
		Network:     "mainnet",
		MacaroonDir: "invoicer",
		TLSPath:     "invoicer/tls.cert",
	}

	lndservs, err := lndclient.NewLndServices(&lndcfg)
	if err != nil {
		log.Fatalf("Error getting lnd grpc services: %v\n", err)
	}

	router := gin.Default()
	router.POST("/api/question", func(c *gin.Context) {
		input := CreateQuestionInput{}
		if err := c.BindJSON(&input); err != nil {
			log.Println("Hit here")
			c.AbortWithError(http.StatusBadRequest, err)
			return
		}

		// fmt.Printf("Got message with title %s\n", v.Title)
		// conn.WriteMessage(t, msg)
		msats := lnwire.MilliSatoshi(input.Bounty * 1000)
		hash, payaddr, err := lndservs.Client.AddInvoice(c.Request.Context(), &invoicesrpc.AddInvoiceData{Memo: input.Title, Value: msats})
		if err != nil {
			log.Fatalf("Error adding invoice: %v\n", err)
		}
		fmt.Printf("Invoice added with payment_request: %s\n payment_addr: %x\n", hash.String(), payaddr)
		// conn.WriteMessage(websocket.TextMessage, []byte(payaddr))

		post := Question{Title: input.Title, Body: input.Body, Bounty: input.Bounty, Paid: false, Hash: hash.String()}
		result := db.Create(&post)
		log.Println("ID:", post.ID, ", error:", result.Error, ", rows:", result.RowsAffected)
		c.JSON(http.StatusOK, gin.H{
			"payment_request": payaddr,
			"hash":            hash.String(),
		})
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

	router.GET("/api/waitInvoicePaid", func(c *gin.Context) {
		parameters := c.Request.URL.Query()
		hash_str := parameters["hash"][0]
		hash, err := lntypes.MakeHashFromStr(hash_str)
		if err != nil {
			log.Fatalf("Error making hash from string: %v\n", err)
		}
		resp, errChan, err := lndservs.Invoices.SubscribeSingleInvoice(c.Request.Context(), hash)
		if err != nil {
			log.Fatalf("Error setting up subscribeSingleInvoice stream: %v\n", err)
		}
		for {
			fmt.Printf("Inside channel loop")
			select {
			case err := <-errChan:
				if err != nil {
					log.Fatalf("Error during subscribeSingleInvoice stream: %v\n", err)
				}
			case update := <-resp:
				state := update.State.String()
				fmt.Printf("State = %s\n", state)
				if update.State == channeldb.ContractSettled {
					// conn.WriteMessage(websocket.TextMessage, []byte("Settled"))
					c.JSON(http.StatusOK, gin.H{
						"status": "settled",
					})
					return
				}
			case <-c.Request.Context().Done():
				fmt.Printf("Closing channel while in select")
				return
			}
		}
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

	// router.GET("/api/invoice", func(c *gin.Context) {

	// })

	router.GET("/api/invoice/ws", func(c *gin.Context) {
		// wshandler(c.Writer, c.Request)
		fmt.Println("In websocket")
		// <-c.Request.Context().Done()
		// fmt.Println("After done signal")
		ctx := c.Copy()

		// return
		conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
		// done := make(chan interface{})

		// conn.SetCloseHandler()
		// <-c.Request.Context().Done()
		// fmt.Println("After done signal")
		// return
		if err != nil {
			fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
			return
		}
		defer conn.Close()
		defer fmt.Println("Closing websocket connection")

		for {
			// _, msg, err := conn.ReadMessage()

			var v struct {
				Body   string `json:"body"`
				Bounty uint64 `json:"bounty"`
				Title  string `json:"title"`
			}
			err := conn.ReadJSON(&v)
			if err != nil {
				log.Fatalf("Error getting json: %v\n", err)
				break
			}
			// fmt.Printf("Got message %v", v)
			fmt.Printf("Got message with title %s\n", v.Title)
			// conn.WriteMessage(t, msg)
			msats := lnwire.MilliSatoshi(v.Bounty * 1000)
			hash, payaddr, err := lndservs.Client.AddInvoice(c.Request.Context(), &invoicesrpc.AddInvoiceData{Memo: v.Title, Value: msats})
			if err != nil {
				log.Fatalf("Error adding invoice: %v\n", err)
			}
			fmt.Printf("Invoice added with payment_request: %s\n payment_addr: %x\n", hash.String(), payaddr)
			conn.WriteMessage(websocket.TextMessage, []byte(payaddr))

			resp, errChan, err := lndservs.Invoices.SubscribeSingleInvoice(c.Request.Context(), hash)
			if err != nil {
				log.Fatalf("Error setting up subscribeSingleInvoice stream: %v\n", err)
			}
			for {
				// err := <-errChan
				// if err != nil {
				// 	log.Fatalf("Error during subscribeSingleInvoice stream: %v\n", err)
				// }
				// update := <-resp
				fmt.Printf("Inside channel loop")
				select {
				case err := <-errChan:
					if err != nil {
						log.Fatalf("Error during subscribeSingleInvoice stream: %v\n", err)
					}
				case update := <-resp:
					state := update.State.String()
					fmt.Printf("State = %s\n", state)
					if update.State == channeldb.ContractSettled {
						conn.WriteMessage(websocket.TextMessage, []byte("Settled"))
						return
					}
				case <-ctx.Done():
					fmt.Printf("Closing channel while in select")
					return
				}

			}
		}
	})

	// Route every path not hitting api to nextjs
	router.NoRoute(ReverseProxy)
	router.Run(":8080")
}
