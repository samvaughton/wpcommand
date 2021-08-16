package db

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
	"github.com/uptrace/bun/driver/pgdriver"
	"time"
)

var Db *bun.DB

type QueryHook struct{}

func (h *QueryHook) BeforeQuery(ctx context.Context, event *bun.QueryEvent) context.Context {
	return ctx
}
func (h *QueryHook) AfterQuery(ctx context.Context, event *bun.QueryEvent) {
	if event.Err != nil {
		log.Info(time.Since(event.StartTime), " :: ", event.Query)
	}
}

type Model struct{}

func InitDbConnection() {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.Config.DatabaseDsn)))

	Db = bun.NewDB(sqlDb, pgdialect.New())

	Db.AddQueryHook(&QueryHook{})

	Db.RegisterModel((*types.UserAccount)(nil))
	Db.RegisterModel((*types.SiteBlueprintSet)(nil))
}

func CreateDefaultAccountAndUser() {
	accounts, err := AccountsGetAll()

	if err != nil {
		log.Error(err)

		return
	}

	if len(accounts) > 0 {
		log.Info("accounts exist")

		return
	}

	account, err := AccountCreate("Default", "default")

	if err != nil {
		log.Error(err)

		return
	}

	user, err := UserCreate("admin@admin.com", "Admin", "Admin", "password", account.Id)

	if err != nil {
		log.Error(err)

		return
	}

	log.Info(fmt.Sprintf("Default account \"%s\" has been created with user %s, password=\"%s\"", account.Key, user.Email, "password"))
}

func CreateDefaultCommands() {
	accounts, err := AccountsGetAll()

	if err != nil {
		log.Error(err)

		return
	}

	if len(accounts) > 0 {
		log.Info("accounts exist")

		return
	}

	account, err := AccountCreate("Default", "default")

	if err != nil {
		log.Error(err)

		return
	}

	user, err := UserCreate("admin@admin.com", "Admin", "Admin", "password", account.Id)

	if err != nil {
		log.Error(err)

		return
	}

	log.Info(fmt.Sprintf("Default account \"%s\" has been created with user %s, password=\"%s\"", account.Key, user.Email, "password"))
}
