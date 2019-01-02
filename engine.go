package main

import (
	"context"
	"encoding/base64"
	"log"
	"time"

	"github.com/fox-one/mixin-sdk/messenger"
	"github.com/fox-one/mixin-sdk/utils"
)

type engine struct {
	*messenger.Messenger

	users map[string]*User
}

func (e engine) matchUser(ctx context.Context, userID string) (string, error) {
	return "", nil
}

func (e engine) OnMessage(ctx context.Context, msgView messenger.MessageView, userID string) error {
	log.Println("I received a msg", msgView)

	switch msgView.Category {
	case messenger.MessageCategorySystemAccountSnapshot, messenger.MessageCategorySystemConversation:
		return nil
	}

	opponentID, err := e.matchUser(ctx, msgView.UserId)
	if err != nil {
		return err
	}

	msgView.ConversationId = utils.UniqueConversationId(ClientID, opponentID)
	data, err := base64.StdEncoding.DecodeString(msgView.Data)
	if err != nil {
		return err
	}

	representativeID := utils.UUIDWithString("REPLY:" + msgView.MessageId)
	return e.SendMessage(ctx, msgView.ConversationId, opponentID, msgView.Category, string(data), representativeID)
}

func (e engine) Run(ctx context.Context) {
	for {
		if err := e.Loop(ctx, e); err != nil {
			log.Println("something is wrong", err)
			time.Sleep(1 * time.Second)
		}
	}
}
