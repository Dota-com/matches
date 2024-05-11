package service_matches

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"matches/internal/domain"
	storage_matches "matches/internal/storage"
)

type Matches struct {
	log             *slog.Logger
	matchesProvider MatchesProvider
	//matchesForKafka MatchesKafkaProvider
}

func New(
	log *slog.Logger,
	matchesUser MatchesProvider) *Matches {
	return &Matches{
		matchesProvider: matchesUser,
		log:             log,
	}
}

type MatchesProvider interface {
	AllMatches(ctx context.Context, log *slog.Logger, id int64) (domain.MatchesIds, error)
	MatchInfo(ctx context.Context, id int64) (domain.Match, error)
}

type MatchesKafkaProvider interface {
	MatchesForApi(ctx context.Context, log *slog.Logger, id int64) (bool, error)
}

// MatchesUser По айди пользователя получаем все айди.
// Проверяем на то, что айдишник передан, если он есть возвращаем списком все матчи
func (m *Matches) MatchesUser(ctx context.Context, id int64) ([]int64, error) {
	log := m.log.With("Matches User")
	log.Info("Matches User: данные по матчам пользователя: ", id)

	matchesId, err := m.matchesProvider.AllMatches(ctx, log, id)
	if err != nil {
		if errors.Is(err, storage_matches.MatchesNotExists) {
			log.Warn("Матчи пользователя ", id, " не найдены")
			return nil, fmt.Errorf("матчи пользователя %d не найдены", id)
		}
	}
	log.Info("Матчи пользователя", id, " получены из базы")

	if len(matchesId.IdsMatches) != 0 {
		return matchesId.IdsMatches, nil
	}

	//TODO Подумать насчет кафки и как реализовать
	//_, err = m.matchesForKafka.MatchesForApi(ctx, log, id)
	//if err != nil {
	//	log.Warn("Ошибка получения матчей в кафке")
	//	return nil, fmt.Errorf("ошибка получения матчей в кафке")
	//}

	return nil, err
}

func (m *Matches) CurrentMatch(ctx context.Context, id int64) (int64, error) {
	log := m.log.With(slog.String("Matches", "init"))
	log.Info("Запрос матча: ", id)

	rs, err := m.matchesProvider.MatchInfo(ctx, id)
	if err != nil {
		log.Warn("ошибка запроса матча: ", err)
		return 0, fmt.Errorf("ошибка запроса матча: %w", err)
	}

	return rs.Id, nil
}
