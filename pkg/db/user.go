package db

import (
	"context"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
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
		Relation("UserAccounts").
		Join("JOIN users_accounts ua").JoinOn("\"user\".id = ua.user_id").
		Join("JOIN accounts a").JoinOn("ua.account_id = a.id").
		Where("\"user\".email = ? and a.key = ?", email, accountKey).
		Scan(context.Background())

	if err != nil {
		log.Error(err)
		return nil // not found
	}

	// fill out the user accounts
	for _, ua := range user.UserAccounts {
		ua.User = user

		acc := new(types.Account)
		err := Db.NewSelect().Model(acc).Where("id = ?", ua.AccountId).Scan(context.Background())
		if err != nil {
			log.Errorf("error hydrating account for UserAccount: %s", err)

			continue
		}

		ua.Account = acc
	}

	return user
}
