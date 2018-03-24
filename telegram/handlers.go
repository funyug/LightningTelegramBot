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
	"time"
)


func GetInfoHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.GetInfo(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		node_pub_key_text := "Node key: "+response.IdentityPubkey
		active_channel_text := "\nActive Channels: "+strconv.FormatInt(int64(response.NumActiveChannels),10)
		pending_channel_text := "\nPending Channels: "+strconv.FormatInt(int64(response.NumPendingChannels),10)
		total_peers_text :="\nTotal Peers: "+strconv.FormatInt(int64(response.NumPeers),10)
		block_height_text := "\nBlock Height: "+strconv.FormatInt(int64(response.BlockHeight),10)
		uri_text := "\nURIs:\n"
		for i:=1;i<=len(response.Uris);i++ {
			uri_text += strconv.FormatInt(int64(i),10)+":"
			uri_text += "\n "+response.Uris[i-1]+"\n"
		}
		b.Send(m.Sender, node_pub_key_text+active_channel_text+pending_channel_text+total_peers_text+block_height_text+uri_text)
	}
}

func GetPeersHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.ListPeers(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		peers_text := ""
		for i:=1;i<=len(response.Peers);i++ {
			peers_text += strconv.FormatInt(int64(i),10)+":"
			peers_text += "\nNode key: "+response.Peers[i-1].PubKey
			peers_text += "\nAddress: "+response.Peers[i-1].Address+"\n"
		}
		b.Send(m.Sender,peers_text)
	}
}


func NewAddressHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.GetDepositAddress(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		b.Send(m.Sender, response.Address)
	}
}

func ListChainTxnsHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.ListChainTxns(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		transactions_text := ""
		for i:=1;i<=len(response.Transactions);i++ {
			transactions_text += strconv.FormatInt(int64(i),10)+":"
			transactions_text += "\nTx hash: "+response.Transactions[i-1].TxHash
			transactions_text += "\nConfirmations: "+ strconv.FormatInt(int64(response.Transactions[i-1].NumConfirmations),10)
			transactions_text += "\nAmount: "+ strconv.FormatInt(int64(response.Transactions[i-1].Amount),10)
			transactions_text += "\nBlock Height: "+ strconv.FormatInt(int64(response.Transactions[i-1].BlockHeight),10)
			transactions_text += "\nDest addresses: "+ strings.Join(response.Transactions[i-1].DestAddresses,",")
			tm := time.Unix(response.Transactions[i-1].TimeStamp, 0)
			transactions_text += "\nTime: "+ tm.String() +"\n"
		}
		b.Send(m.Sender,transactions_text)
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


func ListInvoicesHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.ListInvoices(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		invoice_text := ""
		for i:=1;i<=len(response.Invoices);i++ {
			invoice_text += strconv.FormatInt(int64(i),10)+":"
			invoice_text += "\nMemo: "+response.Invoices[i-1].Memo
			invoice_text += "\nAmount: "+strconv.FormatInt(response.Invoices[i-1].Value,10)
			invoice_text += "\nSettled: "+strconv.FormatBool(response.Invoices[i-1].Settled)
			invoice_text += "\nRhash: "+hex.EncodeToString(response.Invoices[i-1].RHash)
			invoice_text += "\nPay_req: "+response.Invoices[i-1].PaymentRequest+"\n\n"
			tm := time.Unix(response.Invoices[i-1].CreationDate, 0)
			invoice_text += "\nCreation Date: "+tm.String()

			if response.Invoices[i-1].Settled {
				tm := time.Unix(response.Invoices[i-1].SettleDate, 0)
				invoice_text += "\nSettle Date: "+tm.String()
			}
		}
		b.Send(m.Sender,invoice_text)
	}
}

func LookupInvoice(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.LookupInvoice(c,[]byte(m.Payload))
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		invoice_text := ""
		invoice_text += "Memo: "+response.Memo
		invoice_text += "\nAmount: "+strconv.FormatInt(response.Value,10)
		invoice_text += "\nSettled: "+strconv.FormatBool(response.Settled)
		invoice_text += "\nPay_req: "+response.PaymentRequest
		tm := time.Unix(response.CreationDate, 0)
		invoice_text += "\nCreation Date: "+tm.String()
		if response.Settled {
			tm := time.Unix(response.SettleDate, 0)
			invoice_text += "\nSettle Date: "+tm.String()
		}

		b.Send(m.Sender,invoice_text)
	}
}

func ListPaymentsHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.ListPayments(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		payment_text := ""
		for i:=1;i<=len(response.Payments);i++ {
			payment_text += strconv.FormatInt(int64(i),10)+":"
			payment_text += "\nValue: "+ strconv.FormatInt(int64(response.Payments[i-1].Value),10)
			payment_text += "\nPayment Hash: "+response.Payments[i-1].PaymentHash
			payment_text += "\nFee: "+ strconv.FormatInt(int64(response.Payments[i-1].Fee),10)
			tm := time.Unix(response.Payments[i-1].CreationDate, 0)
			payment_text += "\nCreation Date: "+tm.String()+"\n\n"
		}
		b.Send(m.Sender,payment_text)
	}
}

func CloseChannelHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	data := strings.Split(m.Payload," ")
	if len(data) < 2 {
		b.Send(m.Sender,"Please provide funding_tx_id and channel index separated by a space")
		return
	}
	response,err := lnd.CloseChannel(c,data[0],data[1])
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
			case *lnrpc.CloseStatusUpdate_ClosePending:
				closingHash := update.ClosePending.Txid
				txid, err := chainhash.NewHash(closingHash)
				if err != nil {
					b.Send(m.Sender,err.Error())
					return
				}

				b.Send(m.Sender,txid.String())

			case *lnrpc.CloseStatusUpdate_ChanClose:
				return
			}
		}
	}
}

func ListChannelsHandler(b *telebot.Bot,c lnrpc.LightningClient, m *telebot.Message) {
	response,err := lnd.ListChannels(c)
	if err != nil {
		b.Send(m.Sender,err.Error())
	} else {
		channels_text := ""
		for i:=1;i<=len(response.Channels);i++ {
			channels_text += strconv.FormatInt(int64(i),10)+":"
			channels_text += "\nActive: " + strconv.FormatBool(response.Channels[i-1].Active)
			channels_text += "\nRemote Pubkey: " + response.Channels[i-1].RemotePubkey
			channels_text += "\nChannel Point: " + response.Channels[i-1].ChannelPoint
			channels_text += "\nCapacity: " + strconv.FormatInt(response.Channels[i-1].Capacity,10)
			channels_text += "\nLocal Balance: " + strconv.FormatInt(response.Channels[i-1].LocalBalance,10)
			channels_text += "\nRemote Balance: " + strconv.FormatInt(response.Channels[i-1].RemoteBalance,10)
			channels_text += "\nNumber of Updates: " + strconv.FormatUint(response.Channels[i-1].NumUpdates,10)
		}
		b.Send(m.Sender,channels_text)
	}
}

