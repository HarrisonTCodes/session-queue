package routes

import (
	"fmt"
	"net/http"

	"github.com/redis/go-redis/v9"
)

func HandleLivez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func HandleReadyz(rdb *redis.Client) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		err := rdb.Ping(ctx).Err()
		if err != nil {
			http.Error(w, "Redis not ready", http.StatusServiceUnavailable)
			return
		}

		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, "ok")
	}
}
