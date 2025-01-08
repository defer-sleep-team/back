package main

import (
	"context"
	"flag"
	"log"
	"net/http"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/gorilla/sessions"

	"github.com/defer-sleep-team/Aether_backend/sso/internal/login"
	"github.com/defer-sleep-team/Aether_backend/sso/internal/logout"
	"github.com/defer-sleep-team/Aether_backend/sso/internal/validate"
)

var store = sessions.NewCookieStore([]byte(""))

func main() {
	key := flag.String("key", "aetherblockchain", "the key used to encrypt sessions")
	flag.Parse()

	store = sessions.NewCookieStore([]byte(*key))

	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			session, _ := store.Get(r, "auth")
			r = r.WithContext(context.WithValue(r.Context(), "session", session))
			next.ServeHTTP(w, r)
		})
	})
	r.Route("/api_auth_aether_server_sso/", func(r chi.Router) {
		r.Post("/login", login.LoginHandler(store))
		r.Delete("/logout", logout.LogoutHandler(store))
		r.Get("/validate", validate.ValidateHandler(store))
	})
	log.Fatal(http.ListenAndServe(":8530", r))
}
