package handler

import (
	"net/http"

	"github.com/defer-sleep-team/Aether_backend/proxy/pkg/service"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

    "github.com/prometheus/client_golang/prometheus"
    "github.com/prometheus/client_golang/prometheus/promhttp"
)

type Handler struct {
	//server - config the http server
	Server *service.Server
}

const (
	IP        = "45.156.21.196"
	UsersPort = "8003"
	SSOPort   = "8530"
)
var (
    httpRequestsTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "http_requests_total",
            Help: "Total number of HTTP requests.",
        },
        []string{"path", "method"},
    )
)
func init() {
    prometheus.MustRegister(httpRequestsTotal)
}
func SetCORSOriginMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "https://aether-net.ru")
		w.Header().Set("Access-Control-Allow-Credentials", "true")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type")
		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}
		httpRequestsTotal.WithLabelValues(r.URL.Path, r.Method).Inc()
		next.ServeHTTP(w, r)
	})
}

type LoginRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}
type RegRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Constructor of a handler
func NewHandler(server *service.Server) *Handler {
	return &Handler{Server: server}
}

func (h *Handler) InitRoutes() *chi.Mux {
	/////////////////////////////////////////////////////////////////////////////////////////////
	//init new router
	r := chi.NewRouter()
	// redirect /auth/ to /auth
	r.Use(middleware.RedirectSlashes)
	//serve all the api-routes

	r.Use(SetCORSOriginMiddleware)

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Post("/login", (h.Login))

	r.Post("/reg", (h.Register))

	r.Get("/auth/yandex", (h.LoginYandex))

	r.Delete("/logout", (h.Logout))

	r.Get("/metrics", func(w http.ResponseWriter, r *http.Request) {
		promhttp.Handler().ServeHTTP(w, r)
	})
	r.Route("/api", func(r chi.Router) {

		r.Route("/put", func(r chi.Router) {
			r.Put("/like/{id}", FetchUserMiddleware(h.LikeHandler))
			r.Put("/unlike/{id}", FetchUserMiddleware(h.UnlikeHandler))
			r.Put("/changeuser", FetchUserMiddleware(h.ChangeUserHandler))
			r.Put("/subscribe/{id}", FetchUserMiddleware(h.SubscribeHandler))
			r.Put("/unsubscribe/{id}", FetchUserMiddleware(h.UnSubscribeHandler))
		})
		r.Route("/get", func(r chi.Router) {
			r.Get("/subscriptions", FetchUserMiddleware(h.GetSubscriptions))
			r.Get("/user/{userID}", NotStrictFetchUserMiddleware(h.GetUser))
			r.Get("/post/{postID}", (h.GetPost))
			r.Get("/comments/{postID}", (h.GetComments))
			r.Get("/trends/{n}", NotStrictFetchUserMiddleware(h.GetTrends))
			r.Get("/posts_of/{userID}/{n}", (h.GetPostsOf))
			r.Get("/recommendations/{n}", FetchUserMiddleware(h.GetRecs))
		})
		r.Route("/create", func(r chi.Router) {
			r.Post("/post", FetchUserMiddleware(h.CreatePostHandler))
			r.Post("/comment/{postID}", FetchUserMiddleware(h.AddCommentHandler))
		})
		r.Route("/delete", func(r chi.Router) {
			r.Delete("/sudo_post/{id}", FetchUserMiddleware(h.DeletePostHandlerSudo))
			r.Delete("/post/{id}", FetchUserMiddleware(h.DeletePostHandler))
			r.Delete("/comment/{id}", FetchUserMiddleware(h.DeleteCommentHandler))
		})

	})
	return r

}
