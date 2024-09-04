package middleware

import (
	"net/http"

	"github.com/Cirqach/GoStream/internal/auth"
)

// TODO: add CORS, mayby this is the reason why redirect wont work
func Auth() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := r.Cookie("token")
			if err != nil {
				http.Redirect(w, r, "/login", http.StatusMovedPermanently)
				return
			}
			if auth.VerifyToken(token.Value) {
				next.ServeHTTP(w, r)
				return
			}

			http.Redirect(w, r, "/login", http.StatusMovedPermanently)
		})
	}
}
