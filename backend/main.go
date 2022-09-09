package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"example.com/satsoverflow-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/channeldb"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/lnwire"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/boj/redistore.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

const INVOICE_EXPIRY_SECS = 600

// type User struct {
// 	ID           uint
// 	Name         string
// 	Email        *string
// 	Age          uint8
// 	Birthday     time.Time
// 	MemberNumber sql.NullString
// 	ActivatedAt  sql.NullTime
// 	CreatedAt    time.Time
// 	UpdatedAt    time.Time
// }

// type Question struct {
// 	ID uint
// 	// Hash      string
// 	Title     string
// 	Body      string
// 	Bounty    uint
// 	Paid      bool
// 	Hash      string
// 	CreatedAt time.Time
// 	UpdatedAt time.Time
// }

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

// type Answer struct {
// 	ID         uint
// 	Body       string
// 	Bounty     uint
// 	QuestionID uint
// 	CreatedAt  time.Time
// 	UpdatedAt  time.Time
// }

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

// subscribeInvoicesDaemon gets notified of invoice state updates by the lnd
// client, and sets the questions associated with newly settled invoices to be
// paid.  Only questions that are paid for show up on the website.
func subscribeInvoicesDaemon(client lndclient.LightningClient, db *gorm.DB) {
	invoices, errs, err := client.SubscribeInvoices(context.Background(), lndclient.InvoiceSubscriptionRequest{})
	if err != nil {
		log.Fatalf("Error setting up subscribe invoices stream: %v\n", err)
	}
	for {
		select {
		case invoice := <-invoices:
			fmt.Printf("Invoice state: %s\nInvoice memo: %s\n", invoice.State.String(), invoice.Memo)
			if invoice.State == channeldb.ContractSettled {
				hash := invoice.Hash.String()
				db.Model(&models.Question{}).Where("hash = ?", hash).Update("paid", true)
			}
		case err := <-errs:
			fmt.Printf("Error in subscribe invoices stream: %v\n", err)
		}
	}
}

func deleteExpiredInvoicesDaemon(db *gorm.DB) {
	var Questions []models.Question
	for {
		now := time.Now()
		then := now.Add(-time.Minute * 10)
		fmt.Printf("then is %s\n", time.Time.String(then))
		fmt.Printf("now is %s\n", time.Time.String(now))
		db.Where("paid = ? AND created_at < ?", false, then).Find(&Questions)
		for _, q := range Questions {
			fmt.Printf("ID = %d, Title = %s, Created = %s\n", q.ID, q.Title, time.Time.String(q.CreatedAt))
		}
		time.Sleep(15 * time.Minute)
	}
}

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
	db.AutoMigrate(&models.User{}, &models.Question{}, &models.Answer{})

	store, err := redistore.NewRediStore(10, "tcp", "localhost:6379", "", []byte("secret-key"))
	if err != nil {
		log.Fatalf("Error setting up redis store: %v\n", err)
	}
	defer store.Close()
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
	go subscribeInvoicesDaemon(lndservs.Client, db)
	// go deleteExpiredInvoicesDaemon(db)

	router := gin.Default()

	router.POST("/api/login", func(c *gin.Context) {
		type LoginInput struct {
			Email    string `json:"email" binding:"required"`
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
		res := db.Where("email = ?", input.Email).First(&user)
		fmt.Printf("Rows: %d\n", res.RowsAffected)
		if res.Error != nil {
			// log.Fatalf("Error querying user by email: %v\n", res.Error)
			c.String(http.StatusInternalServerError, fmt.Sprintf("Error querying user by email: %v", res.Error))
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
		session, err := store.Get(c.Request, "sessionID")
		if err != nil {
			log.Fatalf("Error getting session: %v\n", err)
		}
		// session.Values["Foo"] = "bar"
		session.Values["email"] = input.Email
		session.Options.HttpOnly = true
		session.Options.MaxAge = 7 * 24 * 60 * 60 // expires in 7 days
		// session.Options.Secure = true
		err = session.Save(c.Request, c.Writer)
		if err != nil {
			log.Fatalf("Error saving session: %v\n", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"message": user.Email,
		})

	})

	router.POST("/api/register", func(c *gin.Context) {
		type RegisterInput struct {
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

		newUser := models.User{Email: input.Email, Password: string(bytes)}
		result := db.Create(&newUser)
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
	})

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
		hash, payaddr, err := lndservs.Client.AddInvoice(c.Request.Context(), &invoicesrpc.AddInvoiceData{Memo: input.Title, Value: msats, Expiry: INVOICE_EXPIRY_SECS})
		if err != nil {
			log.Fatalf("Error adding invoice: %v\n", err)
		}
		fmt.Printf("Invoice added with payment_request: %s\n payment_addr: %x\n", hash.String(), payaddr)
		// conn.WriteMessage(websocket.TextMessage, []byte(payaddr))

		post := models.Question{Title: input.Title, Body: input.Body, Bounty: input.Bounty, Paid: false, Hash: hash.String()}
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
		answer := models.Answer{Body: input.Body, QuestionID: input.QuestionID}
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
		var Answers []models.Answer
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
		var Questions []models.Question
		db.Where("paid = ?", true).Order("id desc").Limit(10).Find(&Questions)
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
		var q models.Question
		res := db.Where("id = ?", id).First(&q)
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
	})

	// router.GET("/api/invoice", func(c *gin.Context) {

	// })

	// router.GET("/api/invoice/ws", func(c *gin.Context) {
	// 	// wshandler(c.Writer, c.Request)
	// 	fmt.Println("In websocket")
	// 	// <-c.Request.Context().Done()
	// 	// fmt.Println("After done signal")
	// 	ctx := c.Copy()

	// 	// return
	// 	conn, err := wsupgrader.Upgrade(c.Writer, c.Request, nil)
	// 	// done := make(chan interface{})

	// 	// conn.SetCloseHandler()
	// 	// <-c.Request.Context().Done()
	// 	// fmt.Println("After done signal")
	// 	// return
	// 	if err != nil {
	// 		fmt.Printf("Failed to set websocket upgrade: %+v\n", err)
	// 		return
	// 	}
	// 	defer conn.Close()
	// 	defer fmt.Println("Closing websocket connection")

	// 	for {
	// 		// _, msg, err := conn.ReadMessage()

	// 		var v struct {
	// 			Body   string `json:"body"`
	// 			Bounty uint64 `json:"bounty"`
	// 			Title  string `json:"title"`
	// 		}
	// 		err := conn.ReadJSON(&v)
	// 		if err != nil {
	// 			log.Fatalf("Error getting json: %v\n", err)
	// 			break
	// 		}
	// 		// fmt.Printf("Got message %v", v)
	// 		fmt.Printf("Got message with title %s\n", v.Title)
	// 		// conn.WriteMessage(t, msg)
	// 		msats := lnwire.MilliSatoshi(v.Bounty * 1000)
	// 		hash, payaddr, err := lndservs.Client.AddInvoice(c.Request.Context(), &invoicesrpc.AddInvoiceData{Memo: v.Title, Value: msats})
	// 		if err != nil {
	// 			log.Fatalf("Error adding invoice: %v\n", err)
	// 		}
	// 		fmt.Printf("Invoice added with payment_request: %s\n payment_addr: %x\n", hash.String(), payaddr)
	// 		conn.WriteMessage(websocket.TextMessage, []byte(payaddr))

	// 		resp, errChan, err := lndservs.Invoices.SubscribeSingleInvoice(c.Request.Context(), hash)
	// 		if err != nil {
	// 			log.Fatalf("Error setting up subscribeSingleInvoice stream: %v\n", err)
	// 		}
	// 		for {
	// 			// err := <-errChan
	// 			// if err != nil {
	// 			// 	log.Fatalf("Error during subscribeSingleInvoice stream: %v\n", err)
	// 			// }
	// 			// update := <-resp
	// 			fmt.Printf("Inside channel loop")
	// 			select {
	// 			case err := <-errChan:
	// 				if err != nil {
	// 					log.Fatalf("Error during subscribeSingleInvoice stream: %v\n", err)
	// 				}
	// 			case update := <-resp:
	// 				state := update.State.String()
	// 				fmt.Printf("State = %s\n", state)
	// 				if update.State == channeldb.ContractSettled {
	// 					conn.WriteMessage(websocket.TextMessage, []byte("Settled"))
	// 					return
	// 				}
	// 			case <-ctx.Done():
	// 				fmt.Printf("Closing channel while in select")
	// 				return
	// 			}

	// 		}
	// 	}
	// })

	// Route every path not hitting api to nextjs
	router.NoRoute(ReverseProxy)
	router.Run(":8080")
}
