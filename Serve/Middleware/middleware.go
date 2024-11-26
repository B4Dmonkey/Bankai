package middleware 

import (
	"log"
	"net/http"
	"time"
)

type Middleware func(http.HandlerFunc) http.HandlerFunc

func Use(h http.HandlerFunc, m ...Middleware) http.HandlerFunc {
	if len(m) < 1 {
		return h
	}

	wrapped := h

	// * loop in reverse to preserve middleware order
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}

	return wrapped
}

func Logger(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start))
	}
}

// func auth(next http.HandlerFunc) http.HandlerFunc {
// 	return func(w http.ResponseWriter, r *http.Request) {
// 		userName, password, ok := r.BasicAuth()
// 		if ok {
// 			log.Println("Authenticating user", userName, password)
// 			if false {
// 				next.ServeHTTP(w, r)
// 			}
// 		}
// 		w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
// 		http.Error(w, "Unauthorized", http.StatusUnauthorized)
// 	}
// }