package main

import (
	context "context"
	"fmt"
	"log"

	"github.com/lightninglabs/lndclient"
	"github.com/lightningnetwork/lnd/lnrpc"
)

func main() {
	fmt.Println("Hello")
	// var opts []grpc.DialOption
	// cert_creds, err := credentials.NewClientTLSFromFile("tls.cert", "localhost")
	// if err != nil {
	// 	log.Fatalf("Error getting TLS credentials: %v\n", err)
	// }
	// opts = append(opts, grpc.WithTransportCredentials(cert_creds))

	// mac_bytes, err := ioutil.ReadFile("admin.macaroon")
	// if err != nil {
	// 	log.Fatalf("Error reading macaroon file: %v\n", err)
	// }

	// auth_creds, err := grpc.Meta

	// conn, err := grpc.Dial("192.168.68.54:10009", opts...)
	// if err != nil {
	// 	log.Fatalf("Error dialing to lnd grpc server: %v\n", err)
	// }
	// defer conn.Close()

	// client := pb.NewLightningClient(conn)
	client, err := lndclient.NewBasicClient("192.168.68.54:10009", "./tls.cert", ".", "mainnet")
	if err != nil {
		log.Fatalf("Error setting up client: %v\n", err)
	}
	resp, err := client.AddInvoice(context.Background(), &lnrpc.Invoice{Memo: "testfromgo", Value: 123})
	if err != nil {
		log.Fatalf("Error adding invoice: %v\n", err)
	}
	fmt.Printf("Invoice added with payment_request: %s\n payment_addr: %x\n", resp.PaymentRequest, resp.PaymentAddr)
}
