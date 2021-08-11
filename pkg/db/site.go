package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"time"
)

func SiteExists(email string) bool {
	return SiteGetByEmail(email) != nil
}

func SiteGetByEmail(email string) *types.Site {
	site := new(types.Site)

	err := Db.NewSelect().Model(site).Where("email = ?", email).Scan(context.Background())

	if err != nil {
		return nil // not found
	}

	return site
}

func SiteGetByUuid(uuid string) (*types.Site, error) {
	site := new(types.Site)

	err := Db.NewSelect().Model(site).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil, err // not found
	}

	return site, nil
}

func SelectSites(siteSelector string, accountId int64) []types.Site {
	var err error
	var sites []types.Site

	if siteSelector == "*" || siteSelector == "all" {
		err = Db.NewSelect().Model(&sites).Where("account_id = ?", accountId).Scan(context.Background())
	} else {
		err = Db.NewSelect().Model(&sites).Where("key = ? and account_id = ?", siteSelector, accountId).Scan(context.Background())
	}

	if err != nil {
		log.Error(err)
	}

	return sites
}

func SiteCreateFromStruct(site *types.Site, accountId int64) error {
	site.Enabled = true
	site.AccountId = accountId
	site.SiteConfig = "{}"
	site.CreatedAt = time.Now()
	site.UpdatedAt = bun.NullTime{Time: time.Now()}

	if site.Uuid == "" {
		site.Uuid = uuid.New().String()
	}

	_, err := Db.NewInsert().Model(site).Returning("*").Exec(context.Background())

	return err
}

func SiteCreate(key string, accountId int64, description string) *types.Site {
	site := &types.Site{
		Uuid:        uuid.New().String(),
		AccountId:   accountId,
		Key:         key,
		Description: description,
		Enabled:     true,
		CreatedAt:   time.Now(),
		SiteConfig:  "{}",
		UpdatedAt:   bun.NullTime{Time: time.Now()},
	}

	_, err := Db.NewInsert().Model(site).Returning("*").Exec(context.Background())

	if err != nil {
		log.Error("could not create a new site", err)

		return nil
	}

	return site
}
