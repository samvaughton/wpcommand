package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"strings"
	"time"
)

func UserExists(email string) bool {
	return UserGetByEmail(email) != nil
}

func UserGetByUuid(uuid string) *types.User {
	user := new(types.User)

	err := Db.NewSelect().Model(user).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil // not found
	}

	return user
}

func UserGetByEmail(email string) *types.User {
	user := new(types.User)

	err := Db.NewSelect().Model(user).Where("email = ?", email).Scan(context.Background())

	if err != nil {
		return nil // not found
	}

	return user
}

func UserGetByEmailAndAccountKey(email string, accountKey string) *types.User {
	user := new(types.User)

	err := Db.
		NewSelect().
		Model(user).
		Relation("Accounts").
		Join("JOIN users_accounts ua").JoinOn("\"user\".id = ua.user_id").
		Join("JOIN accounts a").JoinOn("ua.account_id = a.id").
		Where("\"user\".email = ? and a.key = ?", email, accountKey).
		Scan(context.Background())

	if err != nil {
		log.Error(err)
		return nil // not found
	}

	return user
}

func UserCreate(email string, firstName string, lastName string, passwordPlain string, accountId int64) (*types.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(passwordPlain), 10)

	if err != nil {
		log.Error("could not create user, an error occurred with password hashing", err)
	}

	user := &types.User{
		Uuid:      uuid.New().String(),
		Email:     strings.ToLower(strings.Trim(email, " ,-_=/\n\r\t.@")),
		Password:  string(hashedPassword),
		FirstName: strings.Title(strings.ToLower(firstName)),
		LastName:  strings.Title(strings.ToLower(lastName)),
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err = Db.NewInsert().Model(user).Returning("*").Exec(context.Background())

	if err != nil {
		log.Error("could not create a new user", err)

		return nil, err
	}

	// insert relation
	ua := &types.UserAccount{UserId: user.Id, AccountId: accountId}

	_, err = Db.NewInsert().Model(ua).Exec(context.Background())
	if err != nil {
		log.Error("could not create user <-> account relation", err)

		return user, err
	}

	return user, nil
}
