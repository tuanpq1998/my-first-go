package middlewares

import "net/http"

func AllowCors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allowed-Origins", "*")
		w.Header().Set("Access-Control-Allowed-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allowed-Headers", "*")
		w.Header().Set("Access-Control-Exposed-Headers", "Link")
		w.Header().Set("Access-Control-Allow-Credentials", "false")
		w.Header().Set("Access-Control-Max-Age", "300")

		next.ServeHTTP(w, r)
	})
}
