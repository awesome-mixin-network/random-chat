package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"log"

	"github.com/fox-one/mixin-sdk/messenger"
	"github.com/fox-one/mixin-sdk/mixin"
)

func main() {
	user := &mixin.User{
		UserID:    ClientID,
		SessionID: SessionID,
		PINToken:  PINToken,
	}

	block, _ := pem.Decode([]byte(SessionKey))
	privateKey, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Panicln(err)
	}
	user.SetPrivateKey(privateKey)

	m := messenger.NewMessenger(user)
	e := engine{
		Messenger: m,
		users:     make(map[string]*User),
	}
	ctx := context.Background()

	e.Run(ctx)
}
