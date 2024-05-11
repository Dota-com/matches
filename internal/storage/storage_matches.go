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
	"time"
)

var (
	MatchesNotExists = errors.New("матчи не найдены")
)

type Storage struct {
	db *sql.DB
}

func New(storagePath config.DB) (*Storage, error) {
	connDb := fmt.Sprintf(storagePath.UserDb + ":" + storagePath.PassDb + "@" + storagePath.Host + ":" + storagePath.PortDb + "/" + storagePath.DbName)

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

func (s *Storage) AllMatches(ctx context.Context, log *slog.Logger, id int64) (domain.MatchesIds, error) {
	stmt, err := s.db.Prepare("SELECT matches_id FROM matches WHERE user_id = ?")
	if err != nil {
		log.Error("Ошибка выполнения запроса пользователя", id)
		return domain.MatchesIds{}, nil
	}

	row := stmt.QueryRowContext(ctx, id)

	var matchesId domain.MatchesIds

	err = row.Scan(&matchesId)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			log.Error("Ошибка получения полей запроса", id)
			return domain.MatchesIds{}, fmt.Errorf("query matches for user %d", id)
		}
		log.Error("Ошибка получения полей запроса", id)
		return domain.MatchesIds{}, fmt.Errorf("Ошибка получения полей запроса", id)
	}

	return domain.MatchesIds{IdsMatches: matchesId.IdsMatches}, nil
}

func (s *Storage) MatchInfo(ctx context.Context, id int64) (domain.Match, error) {
	query := `SELECT * FROM matches WHERE id=$1;`

	dbCtx, cancel := context.WithTimeout(ctx, time.Second*5)
	defer cancel()

	var match domain.Match

	err := s.db.QueryRowContext(dbCtx, query, id).Scan(&match.Id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return domain.Match{}, fmt.Errorf("match: %w", err)
		}
		return domain.Match{}, fmt.Errorf("ошибка запроса: %w", err)
	}

	return match, nil
}
