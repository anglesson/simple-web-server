package middleware

import (
	"net/http"

	middleware2 "github.com/anglesson/simple-web-server/internal/authentication/middleware"
)

func TrialMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		user := middleware2.GetCurrentUserID(r)
		if user == "" {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		// Skip trial check for these paths
		excludedPaths := map[string]bool{
			"/settings": true,
			"/logout":   true,
		}

		if excludedPaths[r.URL.Path] {
			next.ServeHTTP(w, r)
			return
		}

		//if !user.IsInTrialPeriod() && !user.IsSubscribed() {
		//	http.Redirect(w, r, "/settings", http.StatusSeeOther)
		//	return
		//}

		next.ServeHTTP(w, r)
	})
}
