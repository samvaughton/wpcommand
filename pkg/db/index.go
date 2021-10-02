package db

import (
	"context"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/stdlib"
	"github.com/oiime/logrusbun"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dialect/pgdialect"
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

	dbConfig, err := pgx.ParseConfig(config.Config.DatabaseDsn)
	if err != nil {
		panic(err)
	}

	dbConfig.PreferSimpleProtocol = true
	sqlDb := stdlib.OpenDB(*dbConfig)

	//sqlDb := sql.OpenDB(pgdriver.NewConnector(pgdriver.WithDSN(config.Config.DatabaseDsn)))

	Db = bun.NewDB(sqlDb, pgdialect.New())

	//Db.AddQueryHook(&QueryHook{})
	Db.AddQueryHook(logrusbun.NewQueryHook(logrusbun.QueryHookOptions{Logger: log.StandardLogger()}))

	Db.RegisterModel((*types.SiteBlueprintSet)(nil))
	Db.RegisterModel((*types.ObjectBlueprintStorageRelation)(nil))
}
