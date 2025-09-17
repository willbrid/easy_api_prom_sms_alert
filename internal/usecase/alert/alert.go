package alert

import (
	"easy-api-prom-alert-sms/config"
	"easy-api-prom-alert-sms/internal/entity"
	"easy-api-prom-alert-sms/internal/microservice"
	"easy-api-prom-alert-sms/internal/pkg/alerthelper"
	"easy-api-prom-alert-sms/internal/pkg/recipienthelper"
	"easy-api-prom-alert-sms/pkg/logger"
)

type AlertUseCase struct {
	iMsc                 microservice.IMicroservice
	recipients           []config.Recipient
	defaultRecipientName string
	simulation           bool
}

func NewAlertUseCase(iMsc microservice.IMicroservice, recipients []config.Recipient, defaultRecipientName string, simulation bool) *AlertUseCase {
	return &AlertUseCase{iMsc, recipients, defaultRecipientName, simulation}
}

func (auc *AlertUseCase) Send(alertData entity.Alert, l logger.ILogger) error {
	for _, alert := range alertData.Data.Alerts {
		alertMsg := alerthelper.GetMsgFromAlert(alert)
		recipientName := alerthelper.GetRecipientFromAlert(alert, auc.defaultRecipientName)
		members := recipienthelper.GetRecipientMembers(auc.recipients, recipientName)

		for _, member := range members {
			url, body, err := auc.iMsc.GetUrlAndBody(member, alertMsg)

			if err != nil {
				return err
			}

			if auc.simulation {
				l.Info("send request with url %s and body %s", url, body)
			} else {
				if err := auc.iMsc.Consume(url, body); err != nil {
					return err
				}
			}
		}
	}

	return nil
}
