package main

// User user
type User struct {
	UserID     string `gorm:"TYPE:VARCHAR(36);NOT NULL;PRIMARY_KEY;" json:"user_id"`
	FullName   string `gorm:"" json:"fullname"`
	OpponentID string `gorm:"TYPE:VARCHAR(36);" json:"opponent_id,omitempty"`
	Enabled    bool   `json:"enabled"`
}
