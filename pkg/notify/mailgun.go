package notify

import (
	"context"
	"github.com/mailgun/mailgun-go/v4"
	"github.com/samvaughton/wpcommand/v2/pkg/types"
)

type Mailgun struct {
	from   string
	client *mailgun.MailgunImpl
}

func NewMailgunFromConfig(config *types.Config) *Mailgun {
	return NewMailgun(config.Mailgun.Domain, config.Mailgun.ApiKey, config.Mailgun.BaseUrl, config.Mailgun.From)
}

func NewMailgun(domain string, apiKey string, baseUrl string, from string) *Mailgun {
	mgClient := mailgun.NewMailgun(domain, apiKey)
	mgClient.SetAPIBase(baseUrl)

	return &Mailgun{client: mgClient, from: from}
}

func (mg *Mailgun) Send(msg types.NotifyMessage) error {
	return mg.SendRaw(msg.To, msg.Subject, msg.Body)
}

func (mg *Mailgun) SendRaw(to string, subject string, body string) error {
	m := mg.client.NewMessage(
		mg.from,
		subject,
		body,
		to,
	)
	_, _, err := mg.client.Send(context.Background(), m)

	return err
}
