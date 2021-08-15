package execution

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

type CommandExecutor interface {
	ExecuteCommand(args []string) (*types.CommandResult, error)
}

func NewCommandExecutor(site *types.Site) (CommandExecutor, error) {
	dce := &DebugCommandExecutor{Site: site}

	if site.TestMode {
		return dce, nil
	}

	if site.LabelSelector != "" && site.Namespace != "" {
		ce, err := NewKubernetesCommandExecutor(site)

		if err != nil {
			log.Errorf("could not initialize kubernetes executor %s", err)
			return nil, err
		}

		return ce, nil
	}

	return dce, nil
}
