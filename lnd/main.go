package lnd

import (
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
	"github.com/lightningnetwork/lnd/macaroons"
	"github.com/btcsuite/btcutil"
	"path/filepath"
	"os"
	"log"
	"google.golang.org/grpc"
	"io/ioutil"
	"github.com/lightningnetwork/lnd/lnrpc"
	"context"
	"strconv"
)

const (
	defaultTLSCertFilename  = "tls.cert"
	defaultMacaroonFilename = "admin.macaroon"
)

var (
	defaultLndDir       = btcutil.AppDataDir("lnd", false)
	defaultTLSCertPath  = filepath.Join(defaultLndDir, defaultTLSCertFilename)
	defaultMacaroonPath = filepath.Join(defaultLndDir, defaultMacaroonFilename)
)

func fatal(err error) {
	log.Print(os.Stderr, "[lncli] %v\n", err)
	os.Exit(1)
}

func Connect(conn *grpc.ClientConn) *grpc.ClientConn {
	creds, err := credentials.NewClientTLSFromFile(defaultTLSCertPath, "")
	if err != nil {
		fatal(err)
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(creds),
	}

	macBytes, err := ioutil.ReadFile(defaultMacaroonPath)
	if err != nil {
		fatal(err)
	}
	mac := &macaroon.Macaroon{}
	if err = mac.UnmarshalBinary(macBytes); err != nil {
		fatal(err)
	}

	// Now we append the macaroon credentials to the dial options.
	cred := macaroons.NewMacaroonCredential(mac)
	opts = append(opts, grpc.WithPerRPCCredentials(cred))

	conn, err = grpc.Dial("localhost:10009", opts...)
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	return conn
}

func GetInfo(c lnrpc.LightningClient) (*lnrpc.GetInfoResponse,error) {
	req := lnrpc.GetInfoRequest{}
	response, err := c.GetInfo(context.Background(),&req);
	if err != nil {
		return nil,err
		log.Println(err)
	}
	log.Println(response)
	return response, err
}

func ListPeers(c lnrpc.LightningClient) (*lnrpc.ListPeersResponse,error) {
	req := lnrpc.ListPeersRequest{}
	response, err := c.ListPeers(context.Background(),&req);
	if err != nil {
		return nil,err
		log.Println(err)
	}
	log.Println(response)
	return response, err
}



func GetDepositAddress(c lnrpc.LightningClient) (*lnrpc.NewAddressResponse,error) {
	req := lnrpc.NewWitnessAddressRequest{}
	response, err := c.NewWitnessAddress(context.Background(),&req);
	if err != nil {
		log.Println(err)
		return nil,err
	}
	log.Println(response)
	return response, err
}


func ListChainTxns(c lnrpc.LightningClient) (*lnrpc.TransactionDetails,error) {
	req := lnrpc.GetTransactionsRequest{}
	response, err := c.GetTransactions(context.Background(),&req);
	if err != nil {
		return nil,err
		log.Println(err)
	}
	log.Println(response)
	return response, err
}


func AddInvoice(c lnrpc.LightningClient, value string) (*lnrpc.AddInvoiceResponse,error) {
	amount,err := strconv.ParseInt(value,10,64)
	if err != nil {
		log.Println(err)
		return nil,err
	}
	invoice := &lnrpc.Invoice{
		Value: amount,
	}
	response, err := c.AddInvoice(context.Background(),invoice);
	if err != nil {
		log.Println(err)
		return nil,err
	}
	log.Println(response)
	return response, err
}

func WalletBalance(c lnrpc.LightningClient) (*lnrpc.WalletBalanceResponse,error) {
	req := lnrpc.WalletBalanceRequest{}
	response, err := c.WalletBalance(context.Background(),&req);
	if err != nil {
		log.Println(err)
		return nil,err
	}
	log.Println(response)
	return response, err
}

func SendCoins(c lnrpc.LightningClient, addr string, value string) (*lnrpc.SendCoinsResponse,error) {
	amount,_ := strconv.ParseInt(value,10,64)
	req := lnrpc.SendCoinsRequest{
		Addr:addr,
		Amount: amount,
		TargetConf:2,
	}
	response, err := c.SendCoins(context.Background(),&req);
	if err != nil {
		log.Println(err)
		return nil,err
	}
	log.Println(response)
	return response, err
}

func ConnectPeer(c lnrpc.LightningClient, pub_key string,host string) (*lnrpc.ConnectPeerResponse,error) {
	addr := &lnrpc.LightningAddress{
		Pubkey: pub_key,
		Host:  host,
	}
	req := &lnrpc.ConnectPeerRequest{
		Addr: addr,
		Perm: false,
	}
	response, err := c.ConnectPeer(context.Background(),req);
	if err != nil {
		log.Println(err)
		return nil,err
	}
	log.Println(response)
	return response, err
}

func OpenChannel(c lnrpc.LightningClient, pub_key_hex []byte,amount int64) (lnrpc.Lightning_OpenChannelClient,error) {
	req := &lnrpc.OpenChannelRequest{
		NodePubkey:pub_key_hex,
		LocalFundingAmount:amount,
	}

	stream, err := c.OpenChannel(context.Background(),req);
	if err != nil {
		log.Println(err)
		return nil,err
	}
	return stream, err
}

func ChannelBalance(c lnrpc.LightningClient) (*lnrpc.ChannelBalanceResponse,error) {
	req := lnrpc.ChannelBalanceRequest{}
	response, err := c.ChannelBalance(context.Background(),&req);
	if err != nil {
		log.Println(err)
		return nil,err
	}
	log.Println(response)
	return response, err
}


func ListInvoices(c lnrpc.LightningClient) (*lnrpc.ListInvoiceResponse,error) {
	req := lnrpc.ListInvoiceRequest{}
	response, err := c.ListInvoices(context.Background(),&req);
	if err != nil {
		return nil,err
		log.Println(err)
	}
	log.Println(response)
	return response, err
}

func LookupInvoice(c lnrpc.LightningClient, rHash []byte) (*lnrpc.Invoice,error) {
	req := &lnrpc.PaymentHash{
		RHash: rHash,
	}
	response, err := c.LookupInvoice(context.Background(),req);
	if err != nil {
		return nil,err
		log.Println(err)
	}
	log.Println(response)
	return response, err
}

func ListPayments(c lnrpc.LightningClient) (*lnrpc.ListPaymentsResponse,error) {
	req := lnrpc.ListPaymentsRequest{}
	response, err := c.ListPayments(context.Background(),&req);
	if err != nil {
		return nil,err
		log.Println(err)
	}
	log.Println(response)
	return response, err
}

func CloseChannel(c lnrpc.LightningClient, funding_tx_id string,index string) (lnrpc.Lightning_CloseChannelClient,error) {
	req := &lnrpc.CloseChannelRequest{}
	req.ChannelPoint.FundingTxid = &lnrpc.ChannelPoint_FundingTxidStr{
		FundingTxidStr: funding_tx_id,
	}

	channel_index,err := strconv.ParseUint(index,10,64)
	if err != nil {
		log.Println(err)
		return nil,err
	}

	req.ChannelPoint.OutputIndex = uint32(channel_index)
	stream, err := c.CloseChannel(context.Background(),req);
	if err != nil {
		log.Println(err)
		return nil,err
	}
	return stream, err
}

func ListChannels(c lnrpc.LightningClient) (*lnrpc.ListChannelsResponse,error) {
	req := lnrpc.ListChannelsRequest{}
	response, err := c.ListChannels(context.Background(),&req);
	if err != nil {
		return nil,err
		log.Println(err)
	}
	log.Println(response)
	return response, err
}