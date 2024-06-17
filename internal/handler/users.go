package handler

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"gitlab.com/matchmaker/internal/entity"
)

type ResponseMessage struct {
	Message string `json:"message"`
}

// UserToPool - puts requested user to matchmaking pool.
// The user isn't able to wait for an immediate response.
func (mh *MatchmakingHandler) UserToPool(w http.ResponseWriter, r *http.Request) {
	user := entity.User{}
	buf := make([]byte, r.ContentLength)
	_, err := r.Body.Read(buf)
	if !errors.Is(err, io.EOF) && err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()

	err = json.Unmarshal(buf, &user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = mh.matchmaker.Put(r.Context(), user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	message, err := json.Marshal(ResponseMessage{"user successfully added to matchmaking pool"})
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	_, _ = w.Write(message)
}
