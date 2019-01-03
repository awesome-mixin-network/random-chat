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

	// HostWrite HostWrite
	HostWrite = "localhost:3306"
	// HostRead HostRead
	HostRead = "localhost:3306"
	// DatabaseName DatabaseName
	DatabaseName = "random_chat"
	// Username Username
	Username = "root"
	// Password Password
	Password = ""
	// ReadLogFile ReadLogFile
	ReadLogFile = "db_read.log"
	// WriteLogFile WriteLogFile
	WriteLogFile = "db_write.log"
)
