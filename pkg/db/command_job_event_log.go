package db

import (
	"context"
	"encoding/json"
	"github.com/google/uuid"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	"github.com/uptrace/bun"
	"time"
)

func CreateCommandJobEvent(jobId int64, rType string, status string, command string, output string, metaData map[string]interface{}) (*types.CommandJobEventLog, error) {
	mdJson, err := json.Marshal(metaData)

	if err != nil {
		return nil, err
	}

	job := &types.CommandJobEventLog{
		Uuid:         uuid.New().String(),
		CommandJobId: jobId,
		Type:         rType,
		Status:       status,
		Command:      command,
		Output:       output,
		MetaData:     string(mdJson),
		CreatedAt:    time.Now(),
		ExecutedAt:   bun.NullTime{Time: time.Now()},
	}

	_, err = Db.NewInsert().Model(job).Returning("*").Exec(context.Background())

	if err != nil {
		return nil, err
	}

	return job, nil
}

func CommandJobEventLogGetByJob(jobId int64) ([]*types.CommandJobEventLog, error) {
	var err error
	items := make([]*types.CommandJobEventLog, 0)

	err = Db.
		NewSelect().
		Model(&items).
		Relation("CommandJob").
		Where("command_job_id = ?", jobId).
		Order("created_at ASC").
		Scan(context.Background())

	if err != nil {
		return items, err
	}

	return items, nil
}
