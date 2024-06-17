package handler

import "net/http"

func (mh *MatchmakingHandler) Route(mux *http.ServeMux) {
	mux.HandleFunc("POST /users/", mh.UserToPool)
}
