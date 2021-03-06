package main

import "context"

// User user
type User struct {
	UserID     string `gorm:"TYPE:VARCHAR(36);NOT NULL;PRIMARY_KEY;" json:"user_id"`
	FullName   string `gorm:"TYPE:VARCHAR(100);DEFAULT:'';" json:"fullname"`
	OpponentID string `gorm:"TYPE:VARCHAR(36);DEFAULT:'';" json:"opponent_id,omitempty"`
	Enabled    bool   `gorm:"DEFAULT:1;" json:"enabled"`
}

// TableName mysql table name
func (User) TableName() string {
	return "users"
}

func (e *engine) matchUser(ctx context.Context, user *User, ignoredIDs ...string) error {
	if len(user.OpponentID) > 0 {
		return nil
	}

	ignoredIDs = append(ignoredIDs, user.UserID)
	ignored := make(map[string]bool, len(ignoredIDs))
	for _, id := range ignoredIDs {
		ignored[id] = true
	}

	var opponent *User
	for id, u := range e.users {
		if len(u.OpponentID) > 0 || !u.Enabled {
			continue
		}
		if _, found := ignored[id]; found {
			continue
		}
		opponent = u
	}

	if opponent == nil {
		var u User
		if db := e.dbRead.Where("user_id NOT IN (?)", ignoredIDs).Where("enabled = ?", true).
			Where("opponent_id = ''").First(&u); db.Error != nil {
			if db.RecordNotFound() {
				return nil
			}
			return db.Error
		}
		opponent = &u
	}

	tx := e.dbWrite.Begin()
	paras := map[string]interface{}{
		"enabled":     true,
		"opponent_id": opponent.UserID,
	}
	if err := tx.Model(user).Updates(paras).Error; err != nil {
		tx.Rollback()
		return err
	}
	paras["opponent_id"] = user.UserID
	if err := tx.Model(opponent).Updates(paras).Error; err != nil {
		tx.Rollback()
		return err
	}
	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	fullname := user.FullName
	if len(fullname) == 0 {
		fullname = "User"
	}
	if err := e.Send(ctx, opponent.UserID, fullname+" join your random conversation!"); err != nil {
		return err
	}

	fullname = opponent.FullName
	if len(fullname) == 0 {
		fullname = "User"
	}
	return e.Send(ctx, user.UserID, fullname+" join your random conversation!")
}

func (e *engine) fetchUser(ctx context.Context, userID string) (*User, error) {
	if len(userID) == 0 {
		return nil, nil
	}

	if user, found := e.users[userID]; found {
		return user, nil
	}
	var user = &User{UserID: userID}
	if err := e.dbRead.FirstOrCreate(&user).Error; err != nil {
		return nil, err
	}
	e.users[userID] = user
	return user, nil
}

func (e *engine) enableUser(ctx context.Context, user *User, enable bool) error {
	if enable {
		if err := e.dbWrite.Model(user).Update("enabled", true).Error; err != nil {
			return err
		}
		return e.matchUser(ctx, user)
	}

	opponent, err := e.fetchUser(ctx, user.OpponentID)
	if err != nil {
		return err
	}

	tx := e.dbWrite.Begin()
	paras := map[string]interface{}{
		"enabled":     false,
		"opponent_id": "",
	}
	if err := tx.Model(user).Updates(paras).Error; err != nil {
		return err
	}

	if opponent != nil {
		if err := tx.Model(opponent).Updates(paras).Error; err != nil {
			tx.Rollback()
			return err
		}
	}

	if err := tx.Commit().Error; err != nil {
		tx.Rollback()
		return err
	}

	if opponent != nil {
		fullname := user.FullName
		if len(fullname) == 0 {
			fullname = "User"
		}
		return e.Send(ctx, opponent.UserID, fullname+" quit your random conversation!")
	}

	return nil
}

func (e *engine) chageFullName(ctx context.Context, user *User, fullname string) error {
	preName := user.FullName
	if len(preName) == 0 {
		preName = "User"
	}
	if err := e.dbWrite.Model(user).Update("full_name", fullname).Error; err != nil {
		return err
	}

	if len(user.OpponentID) == 0 {
		return nil
	}

	return e.Send(ctx, user.OpponentID, preName+" updates name: "+fullname)
}

func (e *engine) changeOpponent(ctx context.Context, user *User) error {
	opponentID := user.OpponentID
	user.OpponentID = ""
	if err := e.matchUser(ctx, user, opponentID); err != nil {
		user.OpponentID = opponentID
		return err
	}
	fullname := user.FullName
	if len(fullname) == 0 {
		fullname = "User"
	}
	return e.Send(ctx, opponentID, fullname+" quit your random conversation!")
}
