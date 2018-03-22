package main

import (
	"google.golang.org/grpc"
	"github.com/lightningnetwork/lnd/lnrpc"
	"fmt"
	"BitcoinTelegramBot/lnd"
	tb "gopkg.in/tucnak/telebot.v2"
	"time"
	"strconv"
	"strings"
	"encoding/hex"
	"io"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"log"
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

	b.Handle("/get_deposit_address", func(m *tb.Message) {
		log.Println("Received get_deposit_address command")
		response,err := lnd.GetDepositAddress(c)
		if err != nil {
			b.Send(m.Sender,err.Error())
		} else {
			b.Send(m.Sender, response.Address)
		}
	})

	b.Handle("/wallet_balance",func(m *tb.Message) {
		log.Println("Received wallet_balance command")
		response,err := lnd.WalletBalance(c)
		if err != nil {
			b.Send(m.Sender,err.Error())
		} else {
			ConfirmedBalanceText := strconv.FormatInt(response.ConfirmedBalance,10)
			UnconfirmedBalanceText := strconv.FormatInt(response.UnconfirmedBalance,10)
			b.Send(m.Sender, ConfirmedBalanceText+"("+UnconfirmedBalanceText+") SAT")
		}
	})

	b.Handle("/send_coins",func(m *tb.Message) {
		log.Println("Received send_coins command")
		data := strings.Split(m.Payload," ")
		if len(data) < 2 {
			b.Send(m.Sender,"Please provide address and amount separated by a space")
			return
		}
		response,err := lnd.SendCoins(c,data[0],data[1])
		if err != nil {
			b.Send(m.Sender,err.Error())
		} else {
			b.Send(m.Sender, "Transaction_id:"+response.Txid)
		}
	})

	b.Handle("/connect_peer",func(m *tb.Message) {
		log.Println("Received connect_peer command")
		targetAddress := m.Payload
		splitAddr := strings.Split(targetAddress, "@")
		if len(splitAddr) != 2 {
			b.Send(m.Sender,"target address expected in format: pubkey@host:port")
			return
		}
		_,err := lnd.ConnectPeer(c,splitAddr[0],splitAddr[1])
		if err != nil {
			b.Send(m.Sender,err.Error())
		} else {
			b.Send(m.Sender, "Connected")
		}
	})

	b.Handle("/open_channel",func(m *tb.Message) {
		log.Println("Received open_channel command")
		data := strings.Split(m.Payload," ")
		if len(data) < 2 {
			b.Send(m.Sender,"Please provide node_pub_key and amount separated by a space")
			return
		}
		nodePubHex, err := hex.DecodeString(data[0])
		if err != nil {
			b.Send(m.Sender,"unable to decode node public key")
			return
		}
		amount,err := strconv.ParseInt(data[1],10,64)
		response,err := lnd.OpenChannel(c,nodePubHex,amount)
		if err != nil {
			b.Send(m.Sender,err.Error())
			return
		} else {
			for {
				resp, err := response.Recv()
				if err == io.EOF {
					return
				} else if err != nil {
					b.Send(m.Sender,err.Error())
					return
				}

				switch update := resp.Update.(type) {
				case *lnrpc.OpenStatusUpdate_ChanPending:
					txid, err := chainhash.NewHash(update.ChanPending.Txid)
					if err != nil {
						b.Send(m.Sender,err.Error())
						return
					}

					b.Send(m.Sender,"Channel opening initiated. Funding txid: "+txid.String())

				case *lnrpc.OpenStatusUpdate_ChanOpen:
					channelPoint := update.ChanOpen.ChannelPoint

					var txidHash []byte
					switch channelPoint.GetFundingTxid().(type) {
					case *lnrpc.ChannelPoint_FundingTxidBytes:
						txidHash = channelPoint.GetFundingTxidBytes()
					case *lnrpc.ChannelPoint_FundingTxidStr:
						s := channelPoint.GetFundingTxidStr()
						h, err := chainhash.NewHashFromStr(s)
						if err != nil {
							b.Send(m.Sender,err.Error())
							return
						}

						txidHash = h[:]
					}

					txid, err := chainhash.NewHash(txidHash)
					if err != nil {
						b.Send(m.Sender,err.Error())
						return
					}

					index := channelPoint.OutputIndex
					b.Send(m.Sender,"Channel ready. Txid:"+fmt.Sprintf("%v", txid)+" Channel Index:"+fmt.Sprintf("%v", index))
				}
			}
		}
	})

	b.Handle("/channel_balance",func(m *tb.Message) {
		log.Println("Received channel_balance command")
		response,err := lnd.ChannelBalance(c)
		if err != nil {
			b.Send(m.Sender,err.Error())
		} else {
			BalanceText := strconv.FormatInt(response.Balance,10)
			b.Send(m.Sender, BalanceText+" SAT")
		}
	})

	b.Handle("/generate_invoice", func(m *tb.Message) {
		log.Println("Received generate_invoice command")
		response,err := lnd.AddInvoice(c,m.Payload)
		if err != nil {
			b.Send(m.Sender,err.Error())
		} else {
			b.Send(m.Sender, response.PaymentRequest)
		}
	})

	log.Println("Server started..")
	response,err := lnd.GetInfo(c)
	if err != nil {
		fatal(err)
	} else {
		log.Println(response)
	}
	b.Start()

}

func fatal(err error) {
	log.Fatalf( "[lncli] %v\n", err)
}
