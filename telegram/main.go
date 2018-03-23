package telegram

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/lightningnetwork/lnd/lnrpc"
	"log"
	"BitcoinTelegramBot/config"
)

func UsernameValid(m *tb.Message) bool {
	if m.Sender.Username == config.Username {
		return true
	}
	return false
}

func InitCommands(b *tb.Bot, c lnrpc.LightningClient) {

	b.Handle("/get_info", func(m *tb.Message) {
		log.Println("Received get_info command")
		if UsernameValid(m) {
			GetInfoHandler(b,c,m)
		}
	})

	b.Handle("/list_peers", func(m *tb.Message) {
		log.Println("Received get_info command")
		if UsernameValid(m) {
			GetPeersHandler(b,c,m)
		}
	})

	b.Handle("/get_address", func(m *tb.Message) {
		log.Println("Received get_address command")
		if UsernameValid(m) {
			NewAddressHandler(b,c,m)
		}
	})

	b.Handle("/wallet_balance",func(m *tb.Message) {
		log.Println("Received wallet_balance command")
		if UsernameValid(m) {
			WalletBalanceHandler(b,c,m)
		}
	})

	b.Handle("/send_coins",func(m *tb.Message) {
		log.Println("Received send_coins command")
		if UsernameValid(m) {
			SendCoinsHandler(b,c,m)
		}
	})

	b.Handle("/connect_peer",func(m *tb.Message) {
		log.Println("Received connect_peer command")
		if UsernameValid(m) {
			ConnectPeerHandler(b,c,m)
		}
	})

	b.Handle("/open_channel",func(m *tb.Message) {
		log.Println("Received open_channel command")
		if UsernameValid(m) {
			OpenChannelHandler(b,c,m)
		}
	})

	b.Handle("/channel_balance",func(m *tb.Message) {
		log.Println("Received channel_balance command")
		if UsernameValid(m) {
			ChannelBalanceHandler(b,c,m)
		}
	})

	b.Handle("/generate_invoice", func(m *tb.Message) {
		log.Println("Received generate_invoice command")
		if UsernameValid(m) {
			GenerateInvoiceHandler(b,c,m)
		}
	})

}
