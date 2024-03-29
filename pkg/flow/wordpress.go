package flow

import (
	"fmt"
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/config"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/execution"
	"github.com/samvaughton/wpcommand/v2/pkg/pipeline"
	"github.com/samvaughton/wpcommand/v2/pkg/registry"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
	"strings"
)

func RunWpUserSync(site *types.Site, flowOpts types.FlowOptions) error {
	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "JOB_WP_USER_SYNC",
		"Detail": "",
	}).Debug("started")

	executor, err := execution.NewCommandExecutor(site, config.Config)

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "JOB_WP_USER_SYNC",
			"Detail": "INIT_EXECUTOR",
		}).Error(err)

		return err
	}
	p := pipeline.SiteCommandPipeline{
		Site:     site,
		Executor: executor,
		Config:   config.Config,
		Options:  pipeline.ExecuteOptions{},
		Commands: []pipeline.SiteCommand{
			registry.GetUserSyncCommand(site),
		},
	}

	p.Run()

	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "JOB_WP_USER_SYNC",
		"Detail": "",
	}).Debug("finished")

	return nil
}

func RunJobBasedWpUserSync(flowOpts types.FlowOptions) {
	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "WP_USER_SYNC",
		"Detail": "",
	}).Info("user sync starting")

	// locate the user sync command
	command, err := db.CommandGetByKey(registry.CmdWpSiteUserSync)

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "WP_USER_SYNC",
			"Detail": "GET_COMMAND",
		}).Error(err)

		return
	}

	// we need to create a command job for each site to sync the users

	sites, err := db.SiteGetAllEnabled()

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "WP_USER_SYNC",
			"Detail": "LIST_SITES",
		}).Error(err)

		return
	}

	db.CreateCommandJobs(command, sites, db.CreateCommandJobContext{
		Description: fmt.Sprintf("job created via %s", flowOpts.LogSource),
	}) // system
}

func RunWpUserCreate(wpUser *types.CreateWpUserPayload, site *types.Site, flowOpts types.FlowOptions) error {
	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "WP_USER_CREATE",
		"Detail": "",
	}).Info("user create starting")

	executor, err := execution.NewCommandExecutor(site, config.Config)

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "WP_USER_CREATE",
			"Detail": "INIT_EXECUTOR",
		}).Error(err)

		return err
	}

	args := fmt.Sprintf("wp user create %s %s --role=%s --user_pass=%s", wpUser.Username, wpUser.Email, wpUser.Role, wpUser.Password)

	p := pipeline.SiteCommandPipeline{
		Site:     site,
		Executor: executor,
		Config:   config.Config,
		Options:  pipeline.ExecuteOptions{},
		Commands: []pipeline.SiteCommand{
			&pipeline.SimpleCommand{Args: []string{args}},
		},
	}

	p.Run()

	return nil
}

func sanitizeUrlOutput(url string) string {
	return strings.Trim(url, "\n\r ")
}

func RunWpCreateUserLogin(site *types.Site, userId string, flowOpts types.FlowOptions) (string, error) {
	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "WP_USER_LOGIN_CREATE",
		"Detail": "",
	}).Info("user login create starting")

	executor, err := execution.NewCommandExecutor(site, config.Config)

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "WP_USER_LOGIN_CREATE",
			"Detail": "INIT_EXECUTOR",
		}).Error(err)

		return "", err
	}

	args := fmt.Sprintf("wp login create %s --url-only", userId)

	log.WithFields(log.Fields{
		"Source": flowOpts.LogSource,
		"Action": "WP_USER_LOGIN_CREATE",
		"Detail": "GENERATE_ARGS",
	}).Debug(args)

	p := pipeline.SiteCommandPipeline{
		Site:     site,
		Executor: executor,
		Config:   config.Config,
		Options:  pipeline.ExecuteOptions{},
		Commands: []pipeline.SiteCommand{
			&pipeline.SimpleCommand{Args: []string{args}},
		},
	}

	p.Run()

	if p.PreviousResult == nil {
		return "", errors.New("no result returned")
	}

	return sanitizeUrlOutput(p.PreviousResult.Output), nil
}

func RunWpUserUpdate(wpUserId int, wpUser *types.UpdateWpUserPayload, site *types.Site, flowOpts types.FlowOptions) error {
	executor, err := execution.NewCommandExecutor(site, config.Config)

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "WP_USER_CREATE",
			"Detail": "INIT_EXECUTOR",
		}).Error(err)

		return err
	}

	updatePassword := ""
	if len(wpUser.Password) > 0 {
		updatePassword = fmt.Sprintf("  --user_pass=%s", wpUser.Password)
	}

	args := fmt.Sprintf("wp user update %v --user_email=%s --role=%s%s", wpUserId, wpUser.Email, wpUser.Role, updatePassword)

	p := pipeline.SiteCommandPipeline{
		Site:     site,
		Executor: executor,
		Config:   config.Config,
		Options:  pipeline.ExecuteOptions{},
		Commands: []pipeline.SiteCommand{
			&pipeline.SimpleCommand{Args: []string{args}},
		},
	}

	p.Run()

	return nil
}

func RunWpUserDelete(wpUserId int, site *types.Site, flowOpts types.FlowOptions) error {
	executor, err := execution.NewCommandExecutor(site, config.Config)

	if err != nil {
		log.WithFields(log.Fields{
			"Source": flowOpts.LogSource,
			"Action": "WP_USER_DELETE",
			"Detail": "INIT_EXECUTOR",
		}).Error(err)

		return err
	}

	args := fmt.Sprintf("wp user delete %v --yes", wpUserId)

	p := pipeline.SiteCommandPipeline{
		Site:     site,
		Executor: executor,
		Config:   config.Config,
		Options:  pipeline.ExecuteOptions{},
		Commands: []pipeline.SiteCommand{
			&pipeline.SimpleCommand{Args: []string{args}},
		},
	}

	p.Run()

	return nil
}
