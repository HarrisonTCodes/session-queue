package routes

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/HarrisonTCodes/session-queue/internal/redisclient"
	"github.com/golang-jwt/jwt/v5"
)

func HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}

var hmacSecret = []byte("secret")

func HandleJoin(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	position, _ := redisclient.Rdb.Incr(ctx, "queue:current-position").Result()

	claims := jwt.MapClaims{
		"sub": position,
		"iat": time.Now().Unix(),
		"exp": time.Now().Add(time.Hour * 3).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(hmacSecret)
	if err != nil {
		log.Fatal(err)
	}

	w.Write([]byte(`{"token":"` + tokenString + `"}`))
}
