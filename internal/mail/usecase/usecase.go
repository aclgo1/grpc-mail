package usecase

import (
	"context"

	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/aclgo/grpc-mail/pkg/logger"
	"github.com/pkg/errors"
)

type mailUseCase struct {
	mailUC mail.MailUseCase
	logger logger.Logger
}

func NewmailUseCase(mailUC mail.MailUseCase, logger logger.Logger) *mailUseCase {
	return &mailUseCase{
		mailUC: mailUC,
		logger: logger,
	}
}

func (u *mailUseCase) Send(ctx context.Context, data *models.MailBody) error {
	if err := u.mailUC.Send(ctx, data); err != nil {
		u.logger.Errorf("Send.mailUC.Send:%v", err)
		return errors.Wrap(err, "Send.mailUC.Send")
	}

	u.logger.Infof("%s to %s", mail.EmailSentSuccess, data.To)

	return nil
}
