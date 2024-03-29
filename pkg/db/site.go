package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"golang.org/x/crypto/sha3"
	"strings"
	"time"
)

func SiteExists(key string, accountId int64) bool {
	_, err := SiteGetByKeySafe(key, accountId)

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

func SiteGetByKeySafe(key string, accountId int64) (*types.Site, error) {
	site := new(types.Site)

	err := Db.NewSelect().Model(site).Where("key = ? and account_id = ?", key, accountId).Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return site, nil
}

func SiteGetByAccessToken(token string) (*types.Site, error) {
	site := new(types.Site)

	err := Db.NewSelect().Model(site).Where("access_token = ?", token).Scan(context.Background())

	if err != nil {
		return nil, err // not found
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

func assignAccessToken(site *types.Site) {
	if site.AccessToken == "" {
		atUuid := uuid.New().String()
		hasher := sha3.New256()
		hasher.Write([]byte(atUuid))
		hash := hasher.Sum(nil)
		site.AccessToken = fmt.Sprintf("%x", hash)
	}
}

func SiteCreateFromStruct(site *types.Site, accountId int64) error {
	if site.LabelSelector == "" || site.Namespace == "" {
		return errors.New("please provide label selector & namespace when creating a site")
	}

	assignAccessToken(site)

	site.Enabled = true
	site.AccountId = accountId
	site.SiteConfig = "{}"
	site.WpCachedData = "{}"
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

func SiteUpdate(site *types.Site) error {
	assignAccessToken(site)

	_, err := Db.NewUpdate().Model(site).WherePK().Returning("*").Exec(context.Background())

	if err != nil {
		log.Error("could not update site", err)

		return err
	}

	return nil
}

func SiteMustGetById(siteId int64) *types.Site {
	site := new(types.Site)

	err := Db.NewSelect().Model(site).Where("id = ?", siteId).Scan(context.Background())

	if err != nil {
		panic(err)
	}

	return site
}

func SiteAddBlueprintSet(siteId int64, blueprintSetId int64) error {
	sbps := &types.SiteBlueprintSet{
		SiteId:         siteId,
		BlueprintSetId: blueprintSetId,
	}

	_, err := Db.NewInsert().Model(sbps).Returning("*").Exec(context.Background())

	return err
}
