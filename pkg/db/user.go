package db

import (
	"context"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

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

/*
 * This 'safe' function will filter out all related accounts apart from the passed accountId
 */
func UsersGetByAccountIdSafe(accountId int64) ([]*types.User, error) {
	var err error
	var items = make([]*types.User, 0)

	err = Db.
		NewSelect().
		Model(&items).
		Relation("UserAccounts").
		Relation("UserAccounts.Account").
		Relation("UserAccounts.User").
		Join("JOIN users_accounts AS ua ON \"user\".id = ua.user_id").
		Where("ua.account_id = ?", accountId).
		Scan(context.Background())

	if err != nil {
		log.Error(err)
	}

	for _, user := range items {
		user.UserAccounts = UserFilterUserAccountsByAccountId(accountId, user.UserAccounts)
	}

	return items, nil
}

func UserGetByUuidSafe(uuid string, accountId int64) (*types.User, error) {
	user := new(types.User)

	err := Db.NewSelect().Model(user).Relation("UserAccounts").Relation("UserAccounts.Account").Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil, err // not found
	}

	filteredAccounts := UserFilterUserAccountsByAccountId(accountId, user.UserAccounts)

	if len(filteredAccounts) == 0 {
		return nil, errors.New("user not found")
	}

	user.UserAccounts = filteredAccounts

	return user, nil
}

func UserFilterUserAccountsByAccountId(accountId int64, uaList []*types.UserAccount) []*types.UserAccount {
	var filteredAccounts []*types.UserAccount

	for _, ua := range uaList {
		if ua.AccountId == accountId {
			filteredAccounts = append(filteredAccounts, ua)
		}
	}

	return filteredAccounts
}
