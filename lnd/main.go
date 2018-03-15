package lnd

import (
	"google.golang.org/grpc/credentials"
	"gopkg.in/macaroon.v2"
	"github.com/lightningnetwork/lnd/macaroons"
	"github.com/btcsuite/btcutil"
	"path/filepath"
	"fmt"
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
	fmt.Fprintf(os.Stderr, "[lncli] %v\n", err)
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

	macConstraints := []macaroons.Constraint{
		macaroons.TimeoutConstraint(60),
	}

	// Apply constraints to the macaroon.
	constrainedMac, err := macaroons.AddConstraints(mac, macConstraints...)
	if err != nil {
		fatal(err)
	}

	// Now we append the macaroon credentials to the dial options.
	cred := macaroons.NewMacaroonCredential(constrainedMac)
	opts = append(opts, grpc.WithPerRPCCredentials(cred))

	conn, err = grpc.Dial("localhost:10009", opts...)
	if err != nil {
		log.Fatalf("did not connect: %s", err)
	}
	return conn
}


func GetDepositAddress(c lnrpc.LightningClient) (*lnrpc.NewAddressResponse,error) {
	req := lnrpc.NewWitnessAddressRequest{}
	response, err := c.NewWitnessAddress(context.Background(),&req);
	if err != nil {
		return nil,err
	}
	return response, err
}

func AddInvoice(c lnrpc.LightningClient, value string) (*lnrpc.AddInvoiceResponse,error) {
	amount,err := strconv.ParseInt(value,0,64)
	if err != nil {
		return nil,err
	}
	invoice := &lnrpc.Invoice{
		Value: amount,
	}
	response, err := c.AddInvoice(context.Background(),invoice);
	if err != nil {
		return nil,err
	}
	return response, err
}
