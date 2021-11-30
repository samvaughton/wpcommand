package db

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"github.com/uptrace/bun"
	"gopkg.in/guregu/null.v3"
	"time"
)

func CommandJobGetByUuid(uuid string) (*types.CommandJob, error) {
	job := new(types.CommandJob)

	err := Db.
		NewSelect().
		Model(job).
		Relation("Site").
		Relation("Command").
		Relation("RunByUser").
		Where("\"command_job\".uuid = ?", uuid).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return job, nil
}

func CommandJobGetByUuidSafe(uuid string, accountId int64) (*types.CommandJob, error) {
	job := new(types.CommandJob)

	err := Db.
		NewSelect().
		Model(job).
		Relation("Site").
		Relation("Command").
		Relation("RunByUser").
		Where("\"site\".account_id = ?", accountId).
		Where("\"command_job\".uuid = ?", uuid).
		Scan(context.Background())

	if err != nil {
		return nil, err
	}

	return job, nil
}

func GetCreatedJobs() []types.CommandJob {
	var err error
	var items []types.CommandJob

	err = Db.
		NewSelect().
		Model(&items).
		Relation("Site").
		Relation("Command").
		Where("status = ?", types.CommandJobStatusCreated).
		Order("created_at ASC").
		Scan(context.Background())

	if err != nil {
		log.Error(err)
	}

	return items
}

type CommandJobsFilterOptions struct {
	ExcludeKeys []string
}

func CommandJobsGetAbandoned(timeoutDuration time.Duration) ([]*types.CommandJob, error) {
	var err error
	items := make([]*types.CommandJob, 0)

	query := Db.
		NewSelect().
		Model(&items).
		Relation("Site").
		Relation("Command").
		Relation("RunByUser").
		Where("\"command_job\".status = ?", types.CommandJobStatusRunning).
		Where("? >= \"command_job\".created_at", time.Now().Add(-timeoutDuration)).
		Order("created_at DESC").
		Limit(50)

	err = query.Scan(context.Background())

	if err != nil {
		return items, err
	}

	return items, nil
}

func CommandJobsGetForAccount(accountId int64, opts CommandJobsFilterOptions) ([]*types.CommandJob, error) {
	var err error
	items := make([]*types.CommandJob, 0)

	query := Db.
		NewSelect().
		Model(&items).
		Relation("Site").
		Relation("Command").
		Relation("RunByUser").
		Where("\"site\".account_id = ?", accountId).
		Order("created_at DESC").
		Limit(50)

	if len(opts.ExcludeKeys) > 0 {
		query = query.Where("\"command\".\"key\" NOT IN (?)", bun.In(opts.ExcludeKeys))
	}

	err = query.Scan(context.Background())

	if err != nil {
		return items, err
	}

	return items, nil
}

type CreateCommandJobContext struct {
	RunByUserId int64
	Description string
	Config      map[string]interface{}
}

func CreateCommandJobs(command *types.Command, sites []*types.Site, createContext CreateCommandJobContext) []*types.CommandJob {
	var jobs []*types.CommandJob

	runBy := null.NewInt(0, false)
	if createContext.RunByUserId > 0 {
		runBy = null.IntFrom(createContext.RunByUserId)
	}

	mergedConfig, err := MergeJobConfigs(command.Config, createContext.Config)

	if err != nil {
		log.Error(err)

		return jobs
	}

	for _, site := range sites {
		job := &types.CommandJob{
			Uuid:        uuid.New().String(),
			SiteId:      site.Id,
			CommandId:   command.Id,
			RunByUserId: runBy,
			Key:         command.Key,
			Status:      types.CommandJobStatusCreated,
			Description: createContext.Description,
			CreatedAt:   time.Now(),
			Config:      mergedConfig,
			ResultData:  "{}",
		}

		_, err := Db.NewInsert().Model(job).Returning("*").Exec(context.Background())

		if err != nil {
			log.Error(err)

			continue
		}

		jobs = append(jobs, job)
	}

	return jobs
}

func MergeJobConfigs(cmdConfig string, contextConfig map[string]interface{}) (string, error) {

	cc := make(map[string]interface{})

	err := json.Unmarshal([]byte(cmdConfig), &cc)

	if err != nil {
		log.Error(err)

		return "", err
	}

	for key, item := range contextConfig {
		cc[key] = item
	}

	data, err := json.Marshal(cc)

	if err != nil {
		log.Error(err)

		return "", err
	}

	return string(data), nil
}
