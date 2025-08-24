package routes

import (
	"context"
	"fmt"
	"net/http"

	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
)

func HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

func HandleJoin(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	foo, _ := redisclient.Rdb.Incr(ctx, "foo").Result()
	fmt.Fprintf(w, "%s", fmt.Sprintf("%d", foo))
}
