package db

import (
	"context"
	"fmt"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"gopkg.in/guregu/null.v3"
	"time"
)

func CommandJobGetByUuid(uuid string) (*types.CommandJob, bool) {
	job := new(types.CommandJob)

	err := Db.NewSelect().Model(job).Where("uuid = ?", uuid).Scan(context.Background())

	if err != nil {
		return nil, false // not found
	}

	return job, true
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

func CreateCommandJobs(command *types.Command, sites []*types.Site, runByUserId int64) []*types.CommandJob {
	var jobs []*types.CommandJob

	for _, site := range sites {
		job := &types.CommandJob{
			Uuid:        uuid.New().String(),
			SiteId:      site.Id,
			CommandId:   command.Id,
			RunByUserId: null.IntFrom(runByUserId),
			Key:         command.Key,
			Status:      types.CommandJobStatusCreated,
			Description: fmt.Sprintf("job created via api"),
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
