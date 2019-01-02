// +build template

package main

const (
	// ClientID mixin user id
	ClientID = "xxx"
	// PIN pin
	PIN = "123456"
	// SessionID session id
	SessionID = "xxx"
	// PINToken pin token
	PINToken = "xxx"
	// SessionKey private key in pem
	SessionKey = `-----BEGIN RSA PRIVATE KEY-----
xxx
-----END RSA PRIVATE KEY-----`
)

type DatabaseConfig struct {
	HostWrite    string
	HostRead     string
	DatabaseName string
	Username     string
	Password     string

	ReadLogFile  string
	WriteLogFile string
}

var Database = &DatabaseConfig{
	HostWrite:    "localhost:3306",
	HostRead:     "localhost:3306",
	DatabaseName: "random_chat",
	Username:     "charlie",
	Password:     "",
	ReadLogFile:  "db_read.log",
	WriteLogFile: "db_write.log",
}
