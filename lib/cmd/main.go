package main

import (
	"github.com/QueensLabOpen/candle-team-gullegris-backend/lib/routes"
	"net/http"
)

func main () {
	r := routes.NewRouter()
	http.ListenAndServe("0.0.0.0:3000", r)
}