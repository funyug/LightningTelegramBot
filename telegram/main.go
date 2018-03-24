package telegram

import (
	tb "gopkg.in/tucnak/telebot.v2"
	"github.com/lightningnetwork/lnd/lnrpc"
	"log"
	"BitcoinTelegramBot/config"
)

func UsernameValid(message *tb.Message) bool {
	if message.Sender.Username == config.Username {
		return true
	}
	return false
}

func InitCommands(bot *tb.Bot, client lnrpc.LightningClient) {

	bot.Handle("/get_info", func(message *tb.Message) {
		log.Println("Received get_info command")
		if UsernameValid(message) {
			GetInfoHandler(bot,client,message)
		}
	})

	bot.Handle("/list_peers", func(message *tb.Message) {
		log.Println("Received get_peers command")
		if UsernameValid(message) {
			GetPeersHandler(bot,client,message)
		}
	})

	bot.Handle("/get_address", func(message *tb.Message) {
		log.Println("Received get_address command")
		if UsernameValid(message) {
			NewAddressHandler(bot,client,message)
		}
	})

	bot.Handle("/wallet_balance",func(message *tb.Message) {
		log.Println("Received wallet_balance command")
		if UsernameValid(message) {
			WalletBalanceHandler(bot,client,message)
		}
	})

	bot.Handle("/list_chain_txns",func(message *tb.Message) {
		log.Println("Received list_chain_txns command")
		if UsernameValid(message) {
			ListChainTxnsHandler(bot,client,message)
		}
	})

	bot.Handle("/send_coins",func(message *tb.Message) {
		log.Println("Received send_coins command")
		if UsernameValid(message) {
			SendCoinsHandler(bot,client,message)
		}
	})

	bot.Handle("/connect_peer",func(message *tb.Message) {
		log.Println("Received connect_peer command")
		if UsernameValid(message) {
			ConnectPeerHandler(bot,client,message)
		}
	})

	bot.Handle("/open_channel",func(message *tb.Message) {
		log.Println("Received open_channel command")
		if UsernameValid(message) {
			OpenChannelHandler(bot,client,message)
		}
	})

	bot.Handle("/channel_balance",func(message *tb.Message) {
		log.Println("Received channel_balance command")
		if UsernameValid(message) {
			ChannelBalanceHandler(bot,client,message)
		}
	})

	bot.Handle("/generate_invoice", func(message *tb.Message) {
		log.Println("Received generate_invoice command")
		if UsernameValid(message) {
			GenerateInvoiceHandler(bot,client,message)
		}
	})

	bot.Handle("/list_invoices",func(message *tb.Message) {
		log.Println("Received list_invoices command")
		if UsernameValid(message) {
			ListInvoicesHandler(bot,client,message)
		}
	})

	bot.Handle("/lookup_invoice",func(message *tb.Message) {
		log.Println("Received lookup_invoice command")
		if UsernameValid(message) {
			LookupInvoice(bot,client,message)
		}
	})

	bot.Handle("/list_payments",func(message *tb.Message) {
		log.Println("Received list_payments command")
		if UsernameValid(message) {
			ListPaymentsHandler(bot,client,message)
		}
	})

	bot.Handle("/close_channel",func(message *tb.Message) {
		log.Println("Received close_channel command")
		if UsernameValid(message) {
			CloseChannelHandler(bot,client,message)
		}
	})

	bot.Handle("/list_channels",func(message *tb.Message) {
		log.Println("Received list_channels command")
		if UsernameValid(message) {
			ListChannelsHandler(bot,client,message)
		}
	})

}
