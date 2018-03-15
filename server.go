package main

import (
	"google.golang.org/grpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"fmt"
	"os"
	"BitcoinTelegramBot/lnd"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
)



func main() {
	var conn *grpc.ClientConn

	conn = lnd.Connect(conn)
	defer conn.Close()

	c := lnrpc.NewLightningClient(conn)

	b, err := tb.NewBot(tb.Settings{
		Token:  "TOKEN_HERE",
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		fatal(err)
	}

	b.Handle("/deposit", func(m *tb.Message) {
		response,err := lnd.GetDepositAddress(c)
		if err != nil {
			b.Send(m.Sender,err)
		} else {
			b.Send(m.Sender, response.Address)
		}
	})

	b.Handle("/generate_invoice", func(m *tb.Message) {
		response,err := lnd.AddInvoice(c,m.Payload)
		if err != nil {
			b.Send(m.Sender,err)
		} else {
			b.Send(m.Sender, response.PaymentRequest)
		}
	})

	fmt.Printf("Server started..")
	b.Start()

}

func fatal(err error) {
	fmt.Fprintf(os.Stderr, "[lncli] %v\n", err)
	os.Exit(1)
}
