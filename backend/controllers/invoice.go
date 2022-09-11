package controllers

import (
	"fmt"
	"log"
	"net/http"

	"example.com/satsoverflow-backend/models"
	"github.com/gin-gonic/gin"
	"github.com/lightningnetwork/lnd/channeldb"
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
