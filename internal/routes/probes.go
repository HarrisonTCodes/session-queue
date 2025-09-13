package routes

import (
	"fmt"
	"net/http"
)

func HandleLivez(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}
