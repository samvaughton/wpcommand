package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"gopkg.in/guregu/null.v3"
	"time"
)

func CommandGetByUuid(uuid string) (*types.Command, error) {
	item := new(types.Command)

	err := Db.NewSelect().Model(item).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return item, nil
}

func CommandGetByKey(key string) (*types.Command, error) {
	item := new(types.Command)

	err := Db.NewSelect().Model(item).Where("key = ?", key).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return item, nil
}

func CommandGetByIdAccountSafe(id int64, accountId int64) (*types.Command, error) {
	item := new(types.Command)

	err := Db.NewSelect().Model(item).Where("id = ?", id).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	if item.IsDefault() == false && item.AccountId.Int64 != accountId {
		return nil, errors.New("command not found via id")
	}

	return item, nil
}

func CommandGetByUuidAccountSafe(uuid string, accountId int64) (*types.Command, error) {
	item := new(types.Command)

	err := Db.NewSelect().Model(item).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	if item.IsDefault() == false && item.AccountId.Int64 != accountId {
		return nil, errors.New("command not found via uuid")
	}

	return item, nil
}

func CommandsGetRunnableForSiteSafe(siteId int64, accountId int64) ([]*types.Command, error) {
	items := make([]*types.Command, 0)

	err := Db.
		NewSelect().
		Model(&items).
		// first clause covers site specific, second clause covers account wide, and third clause covers global
		Where("(site_id = ? AND account_id = ?) OR (site_id IS NULL AND account_id = ?) OR (site_id IS NULL AND account_id IS NULL)", siteId, accountId, accountId).
		Order("description ASC").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return items, nil
}

func CommandsGetAttachedForSiteSafe(siteId int64) ([]*types.Command, error) {
	items := make([]*types.Command, 0)

	err := Db.
		NewSelect().
		Model(&items).
		Where("site_id = ?", siteId, siteId).
		Order("description ASC").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return items, nil
}

func CommandsGetRunnableForAccountId(accountId int64) ([]*types.Command, error) {
	items := make([]*types.Command, 0)

	err := Db.
		NewSelect().
		Model(&items).
		// first clause covers site specific, second clause covers account wide, and third clause covers global
		Where("account_id IS NULL OR (account_id = ? AND site_id IS NULL)", accountId).
		Order("description ASC").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return items, nil
}

func CommandsGetAttachedForAccountId(accountId int64) ([]*types.Command, error) {
	items := make([]*types.Command, 0)

	err := Db.
		NewSelect().
		Model(&items).
		Where("account_id = ? AND site_id IS NULL", accountId).
		Order("description ASC").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return items, nil
}

func CommandsGetDefault() ([]*types.Command, error) {
	items := make([]*types.Command, 0)

	err := Db.
		NewSelect().
		Model(&items).
		Where("account_id IS NULL and site_id IS NULL").
		Order("description ASC").
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return items, nil
}

func CommandCreateAccount(accountId int64, cmd *types.Command) (*types.Command, error) {
	cmd.AccountId = null.IntFrom(accountId)
	cmd.SiteId = null.Int{}

	return CommandCreate(cmd)
}

func CommandCreateSite(siteId int64, cmd *types.Command) (*types.Command, error) {
	cmd.SiteId = null.IntFrom(siteId)
	cmd.AccountId = null.Int{}

	return CommandCreate(cmd)
}

func CommandCreateDefault(description string, key string, cmdType string) (*types.Command, error) {
	return CommandCreate(&types.Command{
		Uuid:        uuid.New().String(),
		Type:        cmdType,
		Key:         key,
		Public:      false,
		Description: description,
		CreatedAt:   time.Now(),
		Config:      "{}",
	})
}

func CommandCreate(cmd *types.Command) (*types.Command, error) {
	if cmd.HttpHeaders == "" {
		cmd.HttpHeaders = "[]"
	}

	_, err := Db.NewInsert().Model(cmd).Returning("*").Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return cmd, nil
}
