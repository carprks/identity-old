package src

import "net/http"

func presetHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Strict-Transport-Security", "max-age=1000; includeSubdomains; preload")
		w.Header().Set("Content-Security-Policy", "upgrade-insecure-requests")
		w.Header().Set("Feature-Policy", "vibrate 'none'; geolocation 'none'")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}
