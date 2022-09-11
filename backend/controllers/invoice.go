package controllers

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"example.com/satsoverflow-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/channeldb"
	"github.com/lightningnetwork/lnd/lnrpc"
	"github.com/lightningnetwork/lnd/lnrpc/invoicesrpc"
	"github.com/lightningnetwork/lnd/lntypes"
	"github.com/lightningnetwork/lnd/lnwire"
)

func (server *Server) WaitInvoicePaid(c *gin.Context) {
	parameters := c.Request.URL.Query()
	hash_str := parameters["hash"][0]
	hash, err := lntypes.MakeHashFromStr(hash_str)
	if err != nil {
		log.Fatalf("Error making hash from string: %v\n", err)
	}
	resp, errChan, err := server.LndServices.Invoices.SubscribeSingleInvoice(c.Request.Context(), hash)
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
}

func (server *Server) AddFunds(c *gin.Context) {
	session, err := server.Store.Get(c.Request, "sessionID")
	if err != nil {
		log.Fatalf("Error getting session: %v\n", err)
	}

	username, found := session.Values["username"].(string)
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please log in"})
		return
	}

	type AddFundsInput struct {
		Sats uint `json:"sats" binding:"required"`
	}
	input := AddFundsInput{}
	if err := c.BindJSON(&input); err != nil {
		log.Println("Hit here")
		c.Status(http.StatusInternalServerError)
		return
	}
	msats := lnwire.MilliSatoshi(input.Sats * 1000)
	fmt.Printf("Making invoice for %d sats", msats)
	hash, payaddr, err := server.LndServices.Client.AddInvoice(c.Request.Context(), &invoicesrpc.AddInvoiceData{Memo: "Add funds to SatsMeAnything", Value: msats, Expiry: INVOICE_EXPIRY_SECS})
	if err != nil {
		log.Fatalf("Error adding invoice: %v\n", err)
	}
	invoice := models.Invoice{Username: &username, Hash: hash.String(), Paid: false}
	server.DB.Create(&invoice)
	c.JSON(http.StatusOK, gin.H{
		"payment_request": payaddr,
		"hash":            hash.String(),
	})
}

func (server *Server) WithdrawalFunds(c *gin.Context) {
	// Get PaymentRequest from body
	type WithdrawalInput struct {
		PaymentRequest string `json:"payment_request" binding:"required"`
	}
	input := WithdrawalInput{}
	if err := c.BindJSON(&input); err != nil {
		log.Println("Hit here")
		c.Status(http.StatusInternalServerError)
		return
	}

	// Decode payment request so we can get its value
	paymentRequest, err := server.LndServices.Client.DecodePaymentRequest(c.Request.Context(), input.PaymentRequest)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	// Get user balance
	session, err := server.Store.Get(c.Request, "sessionID")
	if err != nil {
		log.Fatalf("Error getting session: %v\n", err)
	}
	username, found := session.Values["username"].(string)
	if !found {
		c.JSON(http.StatusUnauthorized, gin.H{"message": "Please log in"})
		return
	}
	var user models.User
	server.DB.Where("username = ?", username).First(&user)

	// Check that user doesn't withdrawal more than their balance
	if user.Balance < uint(paymentRequest.Value.ToSatoshis()) {
		c.Status(http.StatusForbidden)
		return
		// User doesn't have enough to withdrawal
	}

	// Send payment
	sendPaymentRequest := lndclient.SendPaymentRequest{Invoice: input.PaymentRequest, Timeout: time.Minute * 10}
	statuses, errs, err := server.LndServices.Router.SendPayment(c.Request.Context(), sendPaymentRequest)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		log.Printf("Error setting up send payment: %v\n", err)
		return
	}
	for {
		select {
		case err = <-errs:
			if err != nil {
				c.Status(http.StatusInternalServerError)
				log.Printf("Error in send payment error channel: %v\n", err)
			}
		case status := <-statuses:
			if status.State == lnrpc.Payment_SUCCEEDED {
				// conn.WriteMessage(websocket.TextMessage, []byte("Settled"))
				c.JSON(http.StatusOK, gin.H{
					"status": "succeeded",
				})

				server.DB.Model(&models.User{}).Where("username = ?", user.Username).Update("balance", user.Balance-uint(paymentRequest.Value.ToSatoshis()))
				return
			} else if status.State == lnrpc.Payment_FAILED {
				c.JSON(http.StatusOK, gin.H{
					"status": "failed",
				})
			}
		case <-c.Request.Context().Done():
			fmt.Printf("Closing channel while in select")
			return
		}
	}
}
