package handler

type MatchmakingHandler struct {
	matchmaker MatchmakerInterface
}

func New(matchmaker MatchmakerInterface) *MatchmakingHandler {
	return &MatchmakingHandler{
		matchmaker: matchmaker,
	}
}
