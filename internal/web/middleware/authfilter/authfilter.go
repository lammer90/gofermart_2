package authfilter

import (
	"net/http"

	"github.com/gorilla/sessions"
	"github.com/lammer90/gofermart/internal/logger"
	"github.com/lammer90/gofermart/internal/services/authservice"
)

func New(authenticationService authservice.AuthenticationService, cookieStore *sessions.CookieStore) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		logger.Log.Info("auth middleware enabled")

		fn := func(w http.ResponseWriter, r *http.Request) {
			var found bool
			var err error
			var login string

			if !skipAuth(r.URL.String()) {
				for _, cookie := range r.Cookies() {
					if cookie.Name == "Authorization" {
						found = true
						login, err = authenticationService.CheckAuthentication(cookie.Value)
					}
				}
				if err != nil {
					logger.Log.Error("Error during auth user", err)
					w.WriteHeader(http.StatusUnauthorized)
					return
				}
				if !found {
					logger.Log.ErrorMsg("Not Authorized")
					w.WriteHeader(http.StatusUnauthorized)
					return
				}

				session, _ := cookieStore.Get(r, "Authorization")
				session.Values["login"] = login
			}

			next.ServeHTTP(w, r)
		}

		return http.HandlerFunc(fn)
	}
}

func skipAuth(url string) bool {
	return url == "/api/user/register" || url == "/api/user/login"
}
