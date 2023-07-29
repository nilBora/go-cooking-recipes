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

func Cors(next http.Handler) http.Handler {
    fn := func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "OPTIONS" {
            w.Header().Set("Access-Control-Allow-Origin", "*")
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Api-Token, HX-Request, HX-Current-URL")
            return
        }
        next.ServeHTTP(w, r)
    }
    return http.HandlerFunc(fn)
}
