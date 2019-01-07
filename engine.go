package main

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"strings"
	"time"

	"github.com/fox-one/mixin-sdk/messenger"
	"github.com/fox-one/mixin-sdk/utils"
	"github.com/jinzhu/gorm"
)

type engine struct {
	*messenger.Messenger

	dbRead  *gorm.DB
	dbWrite *gorm.DB

	users map[string]*User
}

func (e *engine) OnMessage(ctx context.Context, msgView messenger.MessageView, userID string) error {
	d, err := json.Marshal(msgView)
	log.Println(string(d), err)

	switch msgView.Category {
	case messenger.MessageCategorySystemAccountSnapshot, messenger.MessageCategorySystemConversation:
		return nil
	}

	user, err := e.fetchUser(ctx, msgView.UserId)
	if err != nil {
		return err
	}

	data, err := base64.StdEncoding.DecodeString(msgView.Data)
	if err != nil {
		return err
	}

	if msgView.Category == messenger.MessageCategoryPlainText {
		switch string(data) {
		case "/start":
			return e.enableUser(ctx, user, true)

		case "/stop":
			return e.enableUser(ctx, user, false)

		case "/change":
			return e.changeOpponent(ctx, user)
		}

		if strings.HasPrefix(string(data), "/name ") {
			return e.chageFullName(ctx, user, string(data)[6:])
		}
	}

	if len(user.OpponentID) == 0 {
		msg := `/start to start;
/stop to stop;
/change to seek a new mate;
/name XXX to name your self;`
		return e.Send(ctx, user.UserID, msg)
	}

	msgView.ConversationId = utils.UniqueConversationID(ClientID, user.OpponentID)
	return e.SendMessage(ctx, msgView.ConversationId, user.OpponentID, msgView.Category, string(data), "")
}

func (e *engine) Run(ctx context.Context) {
	for {
		if err := e.Loop(ctx, e); err != nil {
			log.Println("something is wrong", err)
			time.Sleep(1 * time.Second)
		}
	}
}

func (e *engine) Send(ctx context.Context, userID, content string) error {
	msgView := messenger.MessageView{
		ConversationId: utils.UniqueConversationID(ClientID, userID),
		UserId:         userID,
	}
	return e.SendPlainText(ctx, msgView, content)
}
