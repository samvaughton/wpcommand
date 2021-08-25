package db

import (
	"context"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

func UserAccountGetByUuid(uuid string) *types.UserAccount {
	item := new(types.UserAccount)

	err := Db.
		NewSelect().
		Model(item).
		Relation("Account").
		Relation("User").
		Where("\"user_account\".\"uuid\" = ?", uuid).
		Scan(context.Background())

	if err != nil {
		return nil // not found
	}

	return item
}
