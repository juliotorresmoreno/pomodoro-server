package router

import (
	"net/http"
	"os"
	"time"

	"github.com/pingcap/tidb/terror"

	"github.com/juliotorresmoreno/pomodoro-server/controllers"
	"github.com/juliotorresmoreno/pomodoro-server/db"
	"github.com/juliotorresmoreno/pomodoro-server/models"
	"github.com/juliotorresmoreno/pomodoro-server/ws"
	"github.com/justinas/alice"
	"github.com/rs/cors"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/hlog"

	"github.com/gorilla/mux"
)

type handlerFunc func(http.ResponseWriter, *http.Request, models.Session)

func NewRouter() http.Handler {
	router := mux.NewRouter()
	hub := ws.NewHub()
	timer := controllers.NewTimer(hub)
	auth := controllers.NewAuth(hub, timer.TaskManager)

	router.HandleFunc("/auth/login", auth.Login).Methods("POST")
	router.HandleFunc("/auth/register", auth.Register).Methods("POST")
	router.HandleFunc("/auth/session", auth.Session).Methods("GET")
	router.HandleFunc("/timer/new", newRouterProtect(timer.NewPomodoro)).Methods("PUT")

	router.HandleFunc("/ws", newRouterProtect(func(w http.ResponseWriter, r *http.Request, session models.Session) {
		hub.ServeWs(w, r, session)
	}))

	router.PathPrefix("/").Handler(http.FileServer(http.Dir("webroot")))

	return normalize(router)
}

func newRouterProtect(routerFunc handlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.URL.Query().Get("token")
		if token == "" {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		conn, err := db.NewConnection()
		if err != nil {
			terror.Log(err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		defer conn.Close()
		session := models.Session{}
		conn.Find(&session, "token = ?", token)
		routerFunc(w, r, session)
	}
}

func normalize(router http.Handler) http.Handler {
	c := alice.New()
	log := zerolog.New(os.Stdout).With().
		Timestamp().
		Str("role", "Pomodoro").
		Logger()
	// Install the logger handler with default output on the console
	c = c.Append(hlog.NewHandler(log))

	// Install some provided extra handler to set some request's context fields.
	// Thanks to those handler, all our logs will come with some pre-populated fields.
	c = c.Append(hlog.AccessHandler(func(r *http.Request, status, size int, duration time.Duration) {
		hlog.FromRequest(r).Info().
			Str("method", r.Method).
			Str("url", r.URL.String()).
			Int("status", status).
			Int("size", size).
			Dur("duration", duration).
			Msg("")
	}))
	c = c.Append(hlog.RemoteAddrHandler("ip"))
	c = c.Append(hlog.UserAgentHandler("user_agent"))
	c = c.Append(hlog.RefererHandler("referer"))
	c = c.Append(hlog.RequestIDHandler("req_id", "Request-Id"))
	c = c.Append(cors.New(cors.Options{OptionsPassthrough: true}).Handler)

	return c.Then(router)
}
