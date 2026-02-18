package mail

import (
	"context"

	"github.com/aclgo/grpc-mail/internal/models"
	"github.com/aclgo/grpc-mail/pkg/logger"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/trace"
)

type MailUseCase interface {
	Send(context.Context, *models.MailBody) error
}

type Observer struct {
	logger logger.Logger
	Tracer trace.Tracer
	Metric metric.Meter
}

func NewObserver(logger logger.Logger, tracer trace.Tracer, metric metric.Meter) *Observer {
	return &Observer{
		logger: logger,
		Tracer: tracer,
		Metric: metric,
	}
}
