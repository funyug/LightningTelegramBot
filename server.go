package main

import (
	"google.golang.org/grpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"BitcoinTelegramBot/lnd"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
	"log"
	"BitcoinTelegramBot/telegram"
	"BitcoinTelegramBot/config"
)


func main() {

	config.CheckFlags()

	var conn *grpc.ClientConn

	conn = lnd.Connect(conn)
	defer conn.Close()

	client := lnrpc.NewLightningClient(conn)

	bot, err := tb.NewBot(tb.Settings{
		Token:  config.Token,
		Poller: &tb.LongPoller{Timeout: 10 * time.Second},
	})

	if err != nil {
		config.Fatal(err)
	}

	telegram.InitCommands(bot,client)

	log.Println("Server started..")
	bot.Start()

}


