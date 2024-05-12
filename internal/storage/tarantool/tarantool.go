package tarantool

import (
	"fmt"
	"github.com/tarantool/go-tarantool"
	"matches/internal/config"
)

type Tarantool struct {
	tarantool *tarantool.Connection
}

func New(cfg config.Config) (*Tarantool, error) {
	conn, err := tarantool.Connect(cfg.Tarantool.Host, tarantool.Opts{
		User: cfg.Tarantool.User,
		Pass: cfg.Tarantool.Pass,
	})

	if err != nil {
		return nil, fmt.Errorf("ошибка соединения с тарантул")
	}

	defer conn.Close()

	return &Tarantool{tarantool: conn}, nil
}

func (t *Tarantool) Select(id int64) (string, error) {
	res, err := t.tarantool.Select("", "primary", 0, 1, tarantool.IterEq, []interface{}{id})
	if err != nil {
		return "", fmt.Errorf("ошибка получения данных с тарантула: %w", err)
	}

	return res.String(), err
}
