package config

import (
	"flag"
	"errors"
	"log"
)

var Token string;
var Username string;

func CheckFlags() {
	var tokenPtr = flag.String("token","","Your telegram bot token")
	var usernamePtr = flag.String("username","","Your telegram username")
	flag.Parse()

	Token = *tokenPtr
	Username = *usernamePtr
	if Token == "" {
		err := errors.New("flag token is missing")
		Fatal(err)
	}

	if Username == "" {
		err := errors.New("flag username is missing")
		Fatal(err)
	}

}

func Fatal(err error) {
	log.Fatalf( "[lncli] %v\n", err)
}
