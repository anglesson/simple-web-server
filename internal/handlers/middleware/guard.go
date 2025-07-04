package middleware

import "net/http"

func AuthGuard(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if _, err := authorizer(r); err == nil {
			http.Redirect(w, r, "/dashboard", http.StatusSeeOther)
		}
		next.ServeHTTP(w, r)
	})
}
