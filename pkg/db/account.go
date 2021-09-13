package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

func AccountExists(key string) bool {
	return AccountGetByKey(key) != nil
}

func AccountsGetAll() ([]types.Account, error) {
	var err error
	var items []types.Account

	err = Db.NewSelect().Model(&items).Scan(context.Background())

	if err != nil {
		return []types.Account{}, err
	}

	return items, nil
}

func AccountGetByUuid(uuid string) (*types.Account, error) {
	account := new(types.Account)

	err := Db.NewSelect().Model(account).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil, err // not found
	}

	return account, nil
}

func AccountGetByKey(key string) *types.Account {
	account := new(types.Account)

	err := Db.NewSelect().Model(account).Where("key = ?", key).Scan(context.Background())

	if err != nil {
		return nil // not found
	}

	return account
}

func AccountCreate(bdb bun.IDB, name string, key string) (*types.Account, error) {
	account := &types.Account{
		Uuid:      uuid.New().String(),
		Name:      strings.Trim(name, " ,-_=/\n\r\t.@"),
		Key:       strings.ToLower(strings.Trim(key, " ,-_=/\n\r\t.@")),
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := bdb.NewInsert().Model(account).Returning("*").Exec(context.Background())

	if err != nil {
		log.Error("could not create a new account", err)

		return nil, err
	}

	return account, nil
}

func AccountUpdate(bdb bun.IDB, account *types.Account) error {
	_, err := bdb.NewUpdate().Model(account).Where("id = ?", account.Id).Returning("*").Exec(context.Background())

	if err != nil {
		log.Error("could not update account", err)

		return err
	}

	return nil
}
