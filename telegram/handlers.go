package telegram

import (
	"BitcoinTelegramBot/lnd"
	"gopkg.in/tucnak/telebot.v2"
	"github.com/lightningnetwork/lnd/lnrpc"
	"strconv"
	"strings"
	"encoding/hex"
	"io"
	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"fmt"
)

func NewAddressHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.GetDepositAddress(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		b.Send(m.Sender, response.Address)
	}
}

func WalletBalanceHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.WalletBalance(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		ConfirmedBalanceText := strconv.FormatInt(response.ConfirmedBalance,10)
		UnconfirmedBalanceText := strconv.FormatInt(response.UnconfirmedBalance,10)
		b.Send(m.Sender, ConfirmedBalanceText+"("+UnconfirmedBalanceText+") SAT")
	}
}

func SendCoinsHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
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
}

func ConnectPeerHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
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
}

func OpenChannelHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
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
}

func ChannelBalanceHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.ChannelBalance(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		BalanceText := strconv.FormatInt(response.Balance,10)
		b.Send(m.Sender, BalanceText+" SAT")
	}
}

func GenerateInvoiceHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.AddInvoice(c,m.Payload)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		b.Send(m.Sender, response.PaymentRequest)
	}
}


