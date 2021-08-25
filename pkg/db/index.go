package db

import (
	"context"
	"database/sql"
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
		log.Error(time.Since(event.StartTime), " :: ", event.Query)
	}
}

func InitDbConnection() {
	sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.Config.DatabaseDsn)))

	Db = bun.NewDB(sqlDb, pgdialect.New())

	Db.AddQueryHook(&QueryHook{})

	Db.RegisterModel((*types.SiteBlueprintSet)(nil))
}
