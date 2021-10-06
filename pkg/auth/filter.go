package auth

import (
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

func FilterWpUserList(userAccount *types.UserAccount, list []types.WpUser) []types.WpUser {
	filtered := make([]types.WpUser, 0)

	for _, item := range list {
		if WpUserHasReadAccess(userAccount, &item) == false {
			continue
		}

		filtered = append(filtered, item)
	}

	return filtered
}

func FindInWpUserList(userAccount *types.UserAccount, list []types.WpUser, userId int) *types.WpUser {
	for _, item := range list {
		if WpUserHasReadAccess(userAccount, &item) && userId == item.ID {
			return &item
		}
	}

	return nil
}

func WpUserHasReadAccess(userAccount *types.UserAccount, wpUser *types.WpUser) bool {
	if wpUser.Roles != "administrator" {
		allowedNormal, err := Enforcer.Enforce(userAccount.GetCasbinPolicyKey(), types.AuthObjectWordpressUser, types.AuthActionRead)

		if err != nil {
			log.Error(err)
			return false
		}

		return allowedNormal
	}

	if wpUser.Roles == "administrator" {
		allowedSpecial, err := Enforcer.Enforce(userAccount.GetCasbinPolicyKey(), types.AuthObjectWordpressUser, types.AuthActionReadSpecial)

		if err != nil {
			log.Error(err)
			return false
		}

		return allowedSpecial
	}

	return false
}

func FilterCommandList(userAccount *types.UserAccount, list []*types.Command) []*types.Command {
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
