package execution

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

type CommandExecutor interface {
	ExecuteCommand(args []string) (*types.CommandResult, error)
}

func NewCommandExecutor(site *types.Site, config *types.Config) (CommandExecutor, error) {
	dce := &DebugCommandExecutor{Site: site}

	if site.TestMode {
		return dce, nil
	}

	if site.LabelSelector != "" && site.Namespace != "" {
		return NewKubernetesCommandExecutor(site, config), nil
	}

	return nil, errors.New(fmt.Sprintf("could not locate a suitable executor for site %s", site.Description))
}
