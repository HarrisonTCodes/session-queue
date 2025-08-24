package routes

import (
	"fmt"
	"net/http"
)

func HandleStatus(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello world")
}
