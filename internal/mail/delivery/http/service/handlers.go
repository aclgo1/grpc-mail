package service

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/aclgo/grpc-mail/internal/mail"
	"github.com/aclgo/grpc-mail/internal/models"
	"go.opentelemetry.io/otel/metric"
)

var (
	spanServiceNameFormat     = "send-service-http-%s"
	meterServiceFormatSuccess = "send-service-http-%s-success"
	meterServiceFormatFail    = "send-service-http-%s-fail"
)

func (m *MailService) SendService(ctx context.Context) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {

		_, spanProccessing := m.observer.Tracer.Start(
			context.Background(),
			"send-service-http",
		)
		defer spanProccessing.End()

		spanProccessing.AddEvent("request-processing")

		if r.Method != http.MethodPost {

			respError := ResponseError{
				Error:      http.StatusText(http.StatusMethodNotAllowed),
				StatusCode: http.StatusMethodNotAllowed,
			}

			JSON(
				spanProccessing,
				http.StatusText(http.StatusMethodNotAllowed),
				w,
				respError,
				respError.StatusCode,
			)

			return
		}

		var data models.MailBody

		if err := json.NewDecoder(r.Body).Decode(&data); err != nil {
			m.logger.Errorf("SendService.json.NewDecoder: %v", err)

			respError := ResponseError{
				Error:      err.Error(),
				StatusCode: http.StatusBadRequest,
			}

			JSON(
				spanProccessing,
				err.Error(),
				w,
				respError,
				respError.StatusCode,
			)

			return
		}

		svc, ok := m.svcsMail[data.ServiceName]
		if !ok {
			respError := ResponseError{
				Error:      mail.ErrServiceNameNotExist.Error(),
				StatusCode: http.StatusBadRequest,
			}

			JSON(
				spanProccessing,
				mail.ErrServiceNameNotExist.Error(),
				w,
				respError,
				respError.StatusCode,
			)

			return
		}

		if err := data.Validate(); err != nil {
			m.logger.Errorf("SendService.Validate: %v", err)

			respError := ResponseError{
				Error:      err.Error(),
				StatusCode: http.StatusBadRequest,
			}

			JSON(
				spanProccessing,
				err.Error(),
				w,
				respError,
				respError.StatusCode,
			)

			return
		}

		sendSuccess, _ := m.observer.Metric.Float64Counter(
			fmt.Sprintf(meterServiceFormatSuccess, data.ServiceName),
			metric.WithUnit("0"),
		)

		sendFail, _ := m.observer.Metric.Float64Counter(
			fmt.Sprintf(meterServiceFormatFail, data.ServiceName),
			metric.WithUnit("0"),
		)

		spanProccessing.AddEvent("send-mail")

		if err := svc.Send(ctx, &data); err != nil {
			m.logger.Errorf("SendService.Send: %v", err)
			sendFail.Add(context.Background(), 1)

			respError := ResponseError{
				Error:      err.Error(),
				StatusCode: http.StatusInternalServerError,
			}

			JSON(
				spanProccessing,
				err.Error(),
				w,
				respError,
				respError.StatusCode,
			)
			return
		}

		sendSuccess.Add(context.Background(), 1)

		response := ResponsOK{
			Message: mail.EmailSentSuccess,
		}

		JSON(
			spanProccessing,
			mail.EmailSentSuccess,
			w,
			response,
			http.StatusOK,
		)
	}
}
