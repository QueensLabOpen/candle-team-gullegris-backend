package routes

import (
	"github.com/gorilla/mux"
	"github.com/unrolled/render"
	"net/http"
	"github.com/r3labs/sse"
	"github.com/QueensLabOpen/candle-team-gullegris-backend/lib/utils"
	"strconv"
	"time"
	"math/rand"
)

func NewRouter () *mux.Router {
	store := utils.NewStore()

	// Create a renderer
	rend := render.New()

	// Create a stream server
	server := sse.New()

	// Create a router
	r := mux.NewRouter()

	rand.Seed(time.Now().Unix())

	// Create a new game
	r.Handle("/create", corsHeaders(http.HandlerFunc(func (rw http.ResponseWriter, req *http.Request) {
		store.Games = append(store.Games, []int{1})
		server.CreateStream(strconv.FormatInt(int64(len(store.Games)), 10))

		server.Publish(strconv.FormatInt(int64(len(store.Games)), 10), &sse.Event{
			Data: []byte("wait"),
		})

		rend.JSON(rw, http.StatusOK, map[string]int{"gid": len(store.Games), "pid": store.Games[len(store.Games) - 1][len(store.Games[len(store.Games) - 1]) - 1]})
	}))).Methods("POST")

	// Join a game
	r.Handle("/join/{gid}", corsHeaders(http.HandlerFunc(func (rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		gid, _ := strconv.ParseInt(vars["gid"], 10, 64)
		if len(store.Games) < int(gid) {
			rend.JSON(rw, http.StatusOK, map[string]string{"error": "No such session"})
			return
		}

		pid := len(store.Games[gid - 1]) + 1
		store.Games[gid - 1] = append(store.Games[gid - 1], pid)
		rend.JSON(rw, http.StatusOK, map[string]int{"pid": pid})
	}))).Methods("POST")

	// Join stream
	r.Handle("/stream", corsHeaders(http.HandlerFunc(server.HTTPHandler)))

	// Start session
	r.Handle("/start/{gid}", corsHeaders(http.HandlerFunc(func (rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		gid, _ := strconv.ParseInt(vars["gid"], 10, 64)
		if len(store.Games) < int(gid) {
			rend.JSON(rw, http.StatusOK, map[string]string{"error": "No such session"})
			return
		}

		server.Publish(vars["gid"], &sse.Event{
			Data: []byte("start"),
		})

		rend.JSON(rw, http.StatusOK, []byte(""))
	}))).Methods("POST")

	// Trigger
	r.Handle("/trigger/{gid}", corsHeaders(http.HandlerFunc(func (rw http.ResponseWriter, req *http.Request) {
		vars := mux.Vars(req)
		gid, _ := strconv.ParseInt(vars["gid"], 10, 64)
		players := store.Games[gid - 1]

		server.Publish(vars["gid"], &sse.Event{
			Data: []byte(strconv.Itoa(players[rand.Intn(len(players))])),
		})

		rend.JSON(rw, http.StatusOK, []byte(""))
	}))).Methods("POST")

	return r
}

func corsHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) {
		rw.Header().Set("Access-Control-Allow-Origin", "*")
		rw.Header().Set("Access-Control-Allow-Headers", "Content-Type,X-Csrf-Token")
		rw.Header().Set("Access-Control-Allow-Methods", "PUT,POST,GET,OPTIONS,DELETE")
		rw.Header().Set("Access-Control-Expose-Headers", "X-Csrf-Token")
		rw.Header().Set("Access-Control-Allow-Credentials", "true")
		next.ServeHTTP(rw, req)
	})
}