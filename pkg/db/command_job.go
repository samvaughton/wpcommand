package db

import (
	"context"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
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

func CommandJobsGetForAccount(accountId int64) ([]*types.CommandJob, error) {
	var err error
	items := make([]*types.CommandJob, 0)

	err = Db.
		NewSelect().
		Model(&items).
		Relation("Site").
		Relation("Command").
		Relation("RunByUser").
		Where("\"site\".account_id = ?", accountId).
		Order("created_at DESC").
		Limit(50).
		Scan(context.Background())

	if err != nil {
		return items, err
	}

	return items, nil
}

func CreateCommandJobs(command *types.Command, sites []*types.Site, runByUserId int64, description string) []*types.CommandJob {
	var jobs []*types.CommandJob

	runBy := null.NewInt(0, false)
	if runByUserId > 0 {
		runBy = null.IntFrom(runByUserId)
	}

	for _, site := range sites {
		job := &types.CommandJob{
			Uuid:        uuid.New().String(),
			SiteId:      site.Id,
			CommandId:   command.Id,
			RunByUserId: runBy,
			Key:         command.Key,
			Status:      types.CommandJobStatusCreated,
			Description: description,
			CreatedAt:   time.Now(),
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
