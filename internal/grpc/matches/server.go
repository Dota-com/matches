package matches_server

import (
	"context"
	"google.golang.org/grpc"
	matches "matches/protos/gen/dota_tracker.matches.v1"
)

type MatchesApi struct {
	matches.UnimplementedMatchesServerServer
}

func RegisterMatchesServer(server *grpc.Server) {
	matches.RegisterMatchesServerServer(server, &MatchesApi{})
}

func (m *MatchesApi) MatchesCurrentUser(
	ctx context.Context,
	req *matches.MatchesCurrentUserRequest) *matches.MatchesCurrentUserResponse {
	panic("")
}
