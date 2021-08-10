package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"strings"
	"time"
)

func AccountExists(key string) bool {
	return AccountGetByKey(key) != nil
}

func AccountGetByUuid(uuid string) *types.Account {
	account := new(types.Account)

	err := Db.NewSelect().Model(account).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil // not found
	}

	return account
}

func AccountGetByKey(key string) *types.Account {
	account := new(types.Account)

	err := Db.NewSelect().Model(account).Where("key = ?", key).Scan(context.Background())

	if err != nil {
		return nil // not found
	}

	return account
}

func AccountCreate(name string, key string) *types.Account {
	account := &types.Account{
		Uuid:      uuid.New().String(),
		Name:      strings.Trim(name, " ,-_=/\n\r\t.@"),
		Key:       strings.ToLower(strings.Trim(key, " ,-_=/\n\r\t.@")),
		Enabled:   true,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	_, err := Db.NewInsert().Model(account).Returning("*").Exec(context.Background())

	if err != nil {
		log.Error("could not create a new account", err)

		return account
	}

	return account
}
