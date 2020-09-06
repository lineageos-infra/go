package auth

import (
	"net/http"
)

func ApiKey(next http.Handler, keys map[string]bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		key := r.Header.Get("Api-Key")
		if key == "" {
			http.Error(w, "invalid api key", http.StatusUnauthorized)
			return
		}
		_, ok := keys[key]
		if !ok {
			http.Error(w, "invalid api key", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r)
	})
}
