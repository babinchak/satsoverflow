package controllers

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
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/channeldb"
	"gopkg.in/boj/redistore.v1"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Server struct {
	DB          *gorm.DB
	Store       *redistore.RediStore
	LndServices *lndclient.GrpcLndServices
	Router      *gin.Engine
}

func (server *Server) Initialize() {
	// Initialize postgres
	dsn := "host=localhost user=postgres password=dogsandcats123 dbname=postgres port=5432 sslmode=disable TimeZone=Asia/Shanghai"
	var err error
	server.DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Yay")
	}
	server.DB.AutoMigrate(&models.User{}, &models.Question{}, &models.Answer{})

	// Initialize redis session store
	server.Store, err = redistore.NewRediStore(10, "tcp", "localhost:6379", "", []byte("secret-key"))
	if err != nil {
		log.Fatalf("Error setting up redis store: %v\n", err)
	}
	// defer store.Close()

	// Initialize Lnd GRPC services
	lndcfg := lndclient.LndServicesConfig{
		LndAddress:  "192.168.68.54:10009",
		Network:     "mainnet",
		MacaroonDir: "invoicer",
		TLSPath:     "invoicer/tls.cert",
	}

	server.LndServices, err = lndclient.NewLndServices(&lndcfg)
	if err != nil {
		log.Fatalf("Error getting lnd grpc services: %v\n", err)
	}

	// Initialize router
	server.Router = gin.Default()
	server.initializeRoutes()

	// Initialize daemons
	go server.subscribeInvoicesDaemon()
	// go server.deleteExpiredInvoicesDaemon()
}

func reverseProxy(c *gin.Context) {
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

func (server *Server) Run() {
	server.Router.NoRoute(reverseProxy)
	server.Router.Run(":8080")
}

func (server *Server) Close() {
	server.Store.Close()
}

func (server *Server) subscribeInvoicesDaemon() {
	invoices, errs, err := server.LndServices.Client.SubscribeInvoices(context.Background(), lndclient.InvoiceSubscriptionRequest{})
	if err != nil {
		log.Fatalf("Error setting up subscribe invoices stream: %v\n", err)
	}
	for {
		select {
		case invoice := <-invoices:
			fmt.Printf("Invoice state: %s\nInvoice memo: %s\n", invoice.State.String(), invoice.Memo)
			if invoice.State == channeldb.ContractSettled {
				hash := invoice.Hash.String()
				server.DB.Model(&models.Question{}).Where("hash = ?", hash).Update("paid", true)
			}
		case err := <-errs:
			fmt.Printf("Error in subscribe invoices stream: %v\n", err)
		}
	}
}

func (server *Server) deleteExpiredInvoicesDaemon() {
	var Questions []models.Question
	for {
		now := time.Now()
		then := now.Add(-time.Minute * 10)
		fmt.Printf("then is %s\n", time.Time.String(then))
		fmt.Printf("now is %s\n", time.Time.String(now))
		server.DB.Where("paid = ? AND created_at < ?", false, then).Find(&Questions)
		for _, q := range Questions {
			fmt.Printf("ID = %d, Title = %s, Created = %s\n", q.ID, q.Title, time.Time.String(q.CreatedAt))
		}
		time.Sleep(15 * time.Minute)
	}
}
