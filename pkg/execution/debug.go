package execution

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"strings"
)

type DebugCommandExecutor struct {
	Site       *types.Site
	MockOutput string
	CommandLog []string
}

func (e *DebugCommandExecutor) ClearLog() {
	e.CommandLog = []string{}
}

func (e *DebugCommandExecutor) ExecuteCommand(args []string) (*types.CommandResult, error) {
	command := strings.Join(args, " ")

	fields := log.Fields{
		"event":    "execute-command",
		"executor": "debug",
		"command":  command,
		"site":     e.Site.Key,
	}

	log.WithFields(fields).Infoln("success")

	e.CommandLog = append(e.CommandLog, command)

	return &types.CommandResult{
		Command: command,
		Output:  e.MockOutput,
	}, nil
}
