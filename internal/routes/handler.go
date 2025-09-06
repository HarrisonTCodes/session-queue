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

type Status string

const (
	StatusExpired Status = "expired"
	StatusWaiting Status = "waiting"
	StatusActive  Status = "active"
)

type StatusResponse struct {
	Position int64  `json:"position"`
	Status   Status `json:"status"`
}

func HandleStatus(rdb *redis.Client, jwtSecret []byte, windowSize int, activeWindowCount int) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tknHeader := r.Header.Get("Authorization")
		if tknHeader == "" {
			http.Error(w, "missing Authorization header", http.StatusBadRequest)
			return
		}

		tknString := strings.TrimPrefix(tknHeader, "Bearer")
		tknString = strings.TrimSpace(tknString)

		pos, err := jwt.ValidateToken(tknString, jwtSecret)
		if err != nil {
			http.Error(w, "invalid token", http.StatusBadRequest)
			return
		}

		ctx := context.Background()
		windowEndStr, _ := rdb.Get(ctx, "queue:window-end").Result()
		windowEnd, _ := strconv.ParseInt(windowEndStr, 10, 64)

		var status Status
		if pos <= windowEnd-int64(windowSize*activeWindowCount) {
			status = StatusExpired
		} else if pos > windowEnd {
			status = StatusWaiting
		} else {
			status = StatusActive
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

func HandleJoin(rdb *redis.Client, jwtSecret []byte) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := context.Background()
		pos, err := rdb.Incr(ctx, "queue:current-position").Result()
		if err != nil {
			log.Println("Redis error during position increment:", err)
			http.Error(w, "failed to issue token", http.StatusInternalServerError)
			return
		}

		tkn, err := jwt.CreateToken(pos, jwtSecret)
		if err != nil {
			log.Println("Failed to create JWT:", err)
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
