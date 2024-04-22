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
	MatchesDb(ctx context.Context, log *slog.Logger, id int64) (*domain.MatchesIds, error)
}

// MatchesUser По айди пользователя получаем все айди.
// Проверяем на то, что айдишник передан, если он есть возвращаем списком все матчи
func (m *Matches) MatchesUser(ctx context.Context, id int64) ([]int64, error) {
	log := m.log.With("Matches User")
	log.Info("Matches User: данные по матчам пользователя: ", id)

	matchesId, err := m.matchesProvider.MatchesDb(ctx, log, id)
	if err != nil {
		if errors.Is(err, storage_matches.MatchesNotExists) {
			log.Warn("Матчи пользователя ", id, " не найдены")
			return nil, fmt.Errorf("матчи пользователя %d не найдены", id)
		}
	}
	log.Info("Матчи пользователя", id, " получены")
	return matchesId.IdsMatches, nil
}
