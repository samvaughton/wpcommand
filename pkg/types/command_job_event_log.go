package types

import (
	"github.com/uptrace/bun"
	"time"
)

const EventLogTypeInfo = "INFO"
const EventLogTypeData = "DATA"
const EventLogTypeJobStarted = "JOB_STARTED"
const EventLogTypeJobFinished = "JOB_FINISHED"

const EventLogStatusSuccess = "SUCCESS"
const EventLogStatusFailure = "FAILURE"
const EventLogStatusSkipped = "SKIPPED"

type CommandJobEventLog struct {
	Id           int64 `bun:"id,pk"`
	CommandJobId int64
	CommandJob   *CommandJob `bun:"rel:belongs-to"`
	Uuid         string
	Status       string
	Type         string
	Command      string
	Output       string
	MetaData     string
	CreatedAt    time.Time    `bun:",nullzero,notnull,default:current_timestamp"`
	ExecutedAt   bun.NullTime `bun:",nullzero,notnull,default:current_timestamp"`
}
