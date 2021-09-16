package auth

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

func CommandListFilter(userAccount *types.UserAccount, list []*types.Command) []*types.Command {
	filtered := make([]*types.Command, 0)

	for _, item := range list {
		if CommandCanRun(userAccount, item) == false {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}

func CommandCanRun(userAccount *types.UserAccount, command *types.Command) bool {

	// check if it has the command run key first
	allowed, err := Enforcer.Enforce(userAccount.GetCasbinPolicyKey(), types.AuthObjectCommand, types.AuthActionRun)

	if err != nil {
		log.Error(err)
		return false
	}

	if allowed == false {
		return false
	}

	// check if this specific command type can be run
	allowed, err = Enforcer.Enforce(userAccount.GetCasbinPolicyKey(), types.AuthObjectCommandRunType, command.Type)

	if err != nil {
		log.Error(err)
		return false
	}

	if allowed == false {
		return false
	}

	return true
}
