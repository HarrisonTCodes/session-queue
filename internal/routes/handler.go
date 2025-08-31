package routes

import (
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/HarrisonTCodes/session-queue/internal/jwt"
	"github.com/redis/go-redis/v9"
)

type StatusResponse struct {
	Position int64  `json:"position"`
	Status   string `json:"status"`
}

func HandleStatus(rdb *redis.Client, windowSize int) http.HandlerFunc {
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

		ctx := context.Background()
		windowEndStr, _ := rdb.Get(ctx, "queue:window-end").Result()
		windowEnd, _ := strconv.ParseInt(windowEndStr, 10, 64)

		var status string
		if pos <= windowEnd-int64(windowSize) {
			status = "expired"
		} else if pos > windowEnd {
			status = "waiting"
		} else {
			status = "active"
		}

		resp := StatusResponse{
			Position: pos,
			Status:   status,
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
		pos, err := rdb.Incr(ctx, "queue:current-position").Result()
		if err != nil {
			log.Println("Redis error during position increment:", err)
			http.Error(w, "failed to issue token", http.StatusInternalServerError)
			return
		}

		tkn, err := jwt.CreateToken(pos)
		if err != nil {
			log.Println("Failed to create JWT")
			http.Error(w, "failed to issue token", http.StatusInternalServerError)
			return
		}

		resp := JoinResponse{
			Token: tkn,
		}

		log.Printf("Issuing token (position %d)", pos)
		w.Header().Add("Content-Type", "application/json")
		json.NewEncoder(w).Encode(resp)
	}
}
