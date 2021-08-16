package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"time"
)

func CommandGetByUuid(uuid string) (*types.Command, error) {
	job := new(types.Command)

	err := Db.NewSelect().Model(job).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return job, nil
}

func CommandCreateAccount(accountId int64, cmd *types.Command) (*types.Command, error) {
	cmd.AccountId = accountId

	return CommandCreate(cmd)
}

func CommandCreateDefault(description string, key string, cmdType string) (*types.Command, error) {
	return CommandCreate(&types.Command{
		Uuid:        uuid.New().String(),
		Type:        cmdType,
		Key:         key,
		Description: description,
		CreatedAt:   time.Now(),
	})
}

func CommandCreate(cmd *types.Command) (*types.Command, error) {
	_, err := Db.NewInsert().Model(cmd).Returning("*").Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return cmd, nil
}
