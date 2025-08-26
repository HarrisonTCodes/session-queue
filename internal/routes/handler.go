package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/HarrisonTCodes/session-queue/internal/jwt"
	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
)

func HandleStatus(w http.ResponseWriter, r *http.Request) {
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

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"position":"` + strconv.FormatInt(pos, 10) + `"}`))
}

func HandleJoin(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	pos, _ := redisclient.Rdb.Incr(ctx, "queue:current-position").Result()

	tkn, err := jwt.CreateToken(pos)
	if err != nil {
		log.Fatal(err)
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write([]byte(`{"token":"` + tkn + `"}`))
}
