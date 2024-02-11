package rabbitmq

import (
	"context"

	"github.com/katana/fortuna/backend-go/internal/config/logger"
	amqp "github.com/rabbitmq/amqp091-go"
)

func (rbm *rbm_pool) SenderRb(ctx context.Context, exchange_name string, msg *Message) error {
	logger.Info("entro na FUNCAO DE PUBLICACAO")
	if rbm.channel == nil {
		logger.Info("RMB.CHANNEL E NULL")
	}
	if ctx == nil {
		logger.Info("ctx E NULL")
	}
	err := rbm.channel.PublishWithContext(ctx,
		exchange_name, // exchange amq.direct
		"",            // routing key
		false,         // mandatory
		false,         // immediate
		amqp.Publishing{
			Body:        msg.Data,
			ContentType: msg.ContentType,
		})

	if err != nil {
		logger.Error("DEU MERDA AQUI NA PUPLICACAO", err)
	}

	return nil
}
