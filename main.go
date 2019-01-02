package main

import (
	"context"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"log"
	"os"

	"github.com/fox-one/mixin-sdk/messenger"
	"github.com/fox-one/mixin-sdk/mixin"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

func createDb(host, username, password, database, logfile string) *gorm.DB {
	dataSourceName := fmt.Sprintf("%s:%s@%s(%s)/%s?parseTime=True&charset=utf8mb4",
		username,
		password,
		"tcp",
		host,
		database)

	if db, err := gorm.Open("mysql", dataSourceName); err != nil {
		panic(err)
	} else {
		db.DB().SetMaxIdleConns(10)
		if len(logfile) > 0 {
			db.LogMode(true)
			if f, err := os.OpenFile(logfile, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666); err == nil {
				db.SetLogger(gorm.Logger{LogWriter: log.New(f, "\r\n", 0)})
			}
		} else {
			db.LogMode(false)
		}

		return db
	}
}

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

	dbRead := createDb(Database.HostRead, Database.Username, Database.Password,
		Database.DatabaseName, Database.ReadLogFile)
	dbWrite := createDb(Database.HostWrite, Database.Username, Database.Password,
		Database.DatabaseName, Database.WriteLogFile)

	dbWrite.AutoMigrate(User{})

	e := engine{
		Messenger: m,
		dbRead:    dbRead,
		dbWrite:   dbWrite,
		users:     make(map[string]*User),
	}
	ctx := context.Background()

	e.Run(ctx)
}
