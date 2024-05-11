package matches_server

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log/slog"
	matches "matches/protos/gen/dota_tracker.matches.v1"
)

type MatchesApi struct {
	matches.UnimplementedMatchesServerServer
	matchesUser MatchesService
}

type MatchesService interface {
	MatchesUser(ctx context.Context, id int64) ([]int64, error)
	CurrentMatch(ctx context.Context, id int64) (string, error)
}

func RegisterMatchesServer(server *grpc.Server) {
	matches.RegisterMatchesServerServer(server, &MatchesApi{})
}

// MatchesCurrentUser получаем матчи по айдишнику пользователя
func (m *MatchesApi) MatchesCurrentUser(
	ctx context.Context,
	req *matches.MatchesCurrentUserRequest) (*matches.MatchesCurrentUserResponse, error) {
	var log *slog.Logger
	log.With("MatchesCurrentUser")

	log.Debug("Запрос матчей ", req.GetIdUser())
	if req.GetIdUser() == 0 {
		log.Error("Пустой ID пользователя")
		return nil, status.Error(codes.InvalidArgument, "Пустой ID пользователя")
	}

	matchesUser, err := m.matchesUser.MatchesUser(ctx, req.GetIdUser())
	if err != nil {
		log.Error("Отсутствуют данные о матчах пользователя: ", req.GetIdUser())
		return nil, status.Error(codes.Internal,
			fmt.Sprintf("Отсутствуют данные о матчах пользователя: %w", req.GetIdUser()))
	}

	return &matches.MatchesCurrentUserResponse{MatchesId: matchesUser}, nil
}

func (m *MatchesApi) MatchesCurrentMatch(
	//TODO Когда будет готова апишка получения матчей нужно продумать как вытягивать данные
	ctx context.Context,
	r *matches.MatchesCurrentMatchRequest) (*matches.MatchesCurrentMatchResponse, error) {
	if r.GetMatchId() == 0 {
		return nil, status.Error(codes.InvalidArgument, "Отсутствует id матча")
	}

	rs, err := m.matchesUser.CurrentMatch(ctx, r.GetMatchId())
	if err != nil {
		return nil, status.Error(codes.Internal, "Матч не найден")
	}

	return &matches.MatchesCurrentMatchResponse{Info: rs}, nil
}
