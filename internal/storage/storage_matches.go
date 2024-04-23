package storage_matches

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	_ "github.com/lib/pq"
	"log"
	"log/slog"
	"matches/internal/config"
	"matches/internal/domain"
)

var (
	MatchesNotExists = errors.New("матчи не найдены")
)

type Storage struct {
	db *sql.DB
}

func New(storagePath config.DB) (*Storage, error) {
	connDb := storagePath.UserDb + ":" + storagePath.PassDb + "@" + storagePath.Host + ":" + storagePath.PortDb + "/" + storagePath.DbName

	db, err := sql.Open("postgres", connDb)
	if err != nil {
		panic("Ошибка соединения с базой " + connDb)
	}

	defer func() {
		if err := db.Close(); err != nil {
			log.Fatal("Ошибка закрытия базы ", err)
		}
	}()

	if err = db.Ping(); err != nil {
		log.Fatal("Ошибка базы ", err)
		return nil, fmt.Errorf("%s: %s", "postgre", err)
	}

	return &Storage{db: db}, nil
}

func (s *Storage) MatchesDb(ctx context.Context, log *slog.Logger, id int64) (*domain.MatchesIds, error) {
	stmt, err := s.db.Prepare("SELECT matches_id FROM matches WHERE user_id = ?")
	if err != nil {
		log.Error("Ошибка выполнения запроса пользователя", id)
		return &domain.MatchesIds{}, nil
	}

	row := stmt.QueryRowContext(ctx, id)

	var matchesId domain.MatchesIds

	err = row.Scan(&matchesId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Ошибка получения полей запроса", id)
			return &domain.MatchesIds{}, fmt.Errorf("query matches for user %d", id)
		}
		log.Error("Ошибка получения полей запроса", id)
		return &domain.MatchesIds{}, fmt.Errorf("Ошибка получения полей запроса", id)
	}

	return &domain.MatchesIds{IdsMatches: matchesId.IdsMatches}, nil
}
