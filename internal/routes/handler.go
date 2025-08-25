package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/HarrisonTCodes/session-queue/internal/jwt"
	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
)

func HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func HandleJoin(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	pos, _ := redisclient.Rdb.Incr(ctx, "queue:current-position").Result()

	tkn, err := jwt.CreateToken(pos)
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(`{"token":"` + tkn + `"}`))
}
