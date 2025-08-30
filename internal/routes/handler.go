package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strings"

	"github.com/HarrisonTCodes/session-queue/internal/jwt"
	"github.com/redis/go-redis/v9"
)

type StatusResponse struct {
	Position int64 `json:"position"`
	Max      int64 `json:"max"`
}

func HandleStatus(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tknHeader := r.Header.Get("Authorization")
		if tknHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusBadRequest)
			return
		}

		tknString := strings.TrimPrefix(tknHeader, "Bearer")
		tknString = strings.TrimSpace(tknString)

		pos, err := jwt.ValidateToken(tknString)
		if err != nil {
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}

		resp := StatusResponse{
			Position: pos,
			Max:      0,
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}

type JoinResponse struct {
	Token string `json:"token"`
}

func HandleJoin(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		pos, _ := rdb.Incr(ctx, "queue:current-position").Result()

		tkn, err := jwt.CreateToken(pos)
		if err != nil {
			log.Fatal(err)
		}

		resp := JoinResponse{
			Token: tkn,
		}

		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
