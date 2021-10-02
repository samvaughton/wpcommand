package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"strings"
	"time"
)

func SiteExists(key string, accountId int64) bool {
	_, err := SiteGetByKey(key, accountId)

	return err != nil
}

func SiteGetAllEnabled() ([]*types.Site, error) {
	var sites = make([]*types.Site, 0)

	err := Db.NewSelect().Model(&sites).Where("enabled = ?", true).Scan(context.Background())

	if err != nil {
		return sites, err
	}

	return sites, nil
}

func SiteGetByKey(key string, accountId int64) (*types.Site, error) {
	site := new(types.Site)

	err := Db.NewSelect().Model(site).Where("key = ? and account_id = ?", key, accountId).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return site, nil
}

func SiteGetByUuid(uuid string) (*types.Site, error) {
	site := new(types.Site)

	err := Db.NewSelect().Model(site).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil, err // not found
	}

	return site, nil
}

func SiteGetByUuidSafe(uuid string, accountInt int64) (*types.Site, error) {
	site := new(types.Site)

	err := Db.NewSelect().Model(site).Where("uuid = ? and account_id = ?", uuid, accountInt).Scan(context.Background())

	if err != nil {
		return nil, err // not found
	}

	return site, nil
}

func SelectSites(siteSelector string, accountId int64) ([]*types.Site, error) {
	var err error
	var sites = make([]*types.Site, 0)

	if siteSelector == "*" || siteSelector == "all" {
		err = Db.NewSelect().Model(&sites).Where("account_id = ?", accountId).Scan(context.Background())
	} else {
		err = Db.NewSelect().Model(&sites).Where("key = ? and account_id = ?", siteSelector, accountId).Scan(context.Background())
	}

	if err != nil {
		return sites, err
	}

	return sites, nil
}

func SiteCreateFromStruct(site *types.Site, accountId int64) error {
	if site.LabelSelector == "" || site.Namespace == "" {
		return errors.New("please provide label selector & namespace when creating a site")
	}

	site.Enabled = true
	site.AccountId = accountId
	site.SiteConfig = "{}"
	site.CreatedAt = time.Now()
	site.UpdatedAt = bun.NullTime{Time: time.Now()}

	if site.Uuid == "" {
		site.Uuid = uuid.New().String()
	}

	if site.Key == "" {
		// extract key from the label selector
		parts := strings.Split(site.LabelSelector, "=")

		if len(parts) != 2 {
			return errors.New("could not generate a suitable key for the site")
		}

		site.Key = parts[1]
	}

	if site.Description == "" {
		site.Description = strings.Title(site.Key)
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

func SiteUpdate(site *types.Site) error {
	_, err := Db.NewUpdate().Model(site).WherePK().Returning("*").Exec(context.Background())

	if err != nil {
		log.Error("could not update site", err)

		return err
	}

	return nil
}
