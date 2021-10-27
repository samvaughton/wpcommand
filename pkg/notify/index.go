package notify

import (
	"github.com/pkg/errors"
	"github.com/samvaughton/wpcommand/v2/pkg/db"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
	log "github.com/sirupsen/logrus"
)

var Client Notifier

const ImplMailgun = "MAILGUN"

type Notifier interface {
	Send(msg types.NotifyMessage) error
	SendRaw(to string, subject string, body string) error
}

func InitNotifier(implementation string, config *types.Config) {
	switch implementation {
	case ImplMailgun:
		Client = NewMailgunFromConfig(config)
	}
}

func SendToAllUsers(subject string, body string, accountId int64) error {
	users, err := db.UsersGetByAccountIdSafe(accountId)

	if err != nil {
		return err
	}

	allErrors := ""

	for _, item := range users {
		err = Client.Send(types.NotifyMessage{
			To:      item.Email,
			Subject: subject,
			Body:    body,
		})

		if err != nil {
			log.WithFields(log.Fields{
				"To":        item.Email,
				"AccountId": accountId,
				"Action":    "NotifyAllUsers",
			}).Error(err)

			allErrors += ", " + err.Error()
		}
	}

	if len(allErrors) > 0 {
		return errors.New(allErrors)
	}

	return nil
}
