package controllers

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lightningnetwork/lnd/channeldb"
	"github.com/lightningnetwork/lnd/lntypes"
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
