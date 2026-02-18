package ses

import (
	"context"

	"github.com/aclgo/grpc-mail/config"
	"github.com/aclgo/grpc-mail/internal/models"
)

type Ses struct {
	config *config.Config
}

func NewSes(config *config.Config) *Ses {
	return &Ses{
		config: config,
	}
}

func (s *Ses) Connect(data *models.MailBody) error {
	return nil
}

func (s *Ses) Send(ctx context.Context, data *models.MailBody) error {
	return nil
}
