package server

import (
	"net/http"
)


func Authentication(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
	    token := r.Header.Get("Api-Token")
	    if token != "111" {
	        w.Write([]byte("Unauthorized"));
	        w.WriteHeader(http.StatusUnauthorized)
        	return

	    }
	    next.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
