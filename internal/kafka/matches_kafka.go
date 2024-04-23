package kafka

import (
	"context"
	"log/slog"
	"strconv"
)

type Kafka struct {
	result []string
}

func (k *Kafka) MatchesForApi(ctx context.Context, log *slog.Logger, id int64) (bool, error) {
	log.Info("Запрос матчей в кафке пользователя " + strconv.FormatInt(id, 10))

	return false, nil
}
