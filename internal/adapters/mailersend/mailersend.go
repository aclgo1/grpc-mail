package mailersend

import (
	"context"
	"fmt"

	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/internal/models"
	mailer "github.com/mailersend/mailersend-go"
)

type MailerSend struct {
	ApiKey string
}

func NewMailerSend(cfg *config.Config) mail.MailUseCase {
	return &MailerSend{
		ApiKey: cfg.ApiKey,
	}
}

func (m *MailerSend) Send(ctx context.Context, data *models.MailBody) error {

	ms := mailer.NewMailersend(m.ApiKey)

	from := mailer.From{
		Name:  "",
		Email: data.From,
	}

	recipients := []mailer.Recipient{
		{
			Name:  "",
			Email: data.To,
		},
	}

	message := ms.Email.NewMessage()

	message.SetFrom(from)
	message.SetRecipients(recipients)
	message.SetSubject(data.Subject)
	message.SetHTML(data.Template)
	message.SetText(data.Body)

	_, err := ms.Email.Send(ctx, message)

	if err != nil {
		return fmt.Errorf("ms.Email.Send: %w", err)
	}

	return nil
}
