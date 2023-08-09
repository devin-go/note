package main

import (
	"io/ioutil"
	"net/http"
	"time"
)

func main() {
	// Create a new HTTP server with a caching middleware
	cacheHandler := CacheHandler(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Write the response body
		data,_:=ioutil.ReadFile(`E:\devin\github\note\golang\服务保障.md`)
		//w.Write([]byte("Hello, World!sdfsdddddddddddddddddddddddddasdfsdfsdfasdfsadfsedfwerwefsdcsdfsdfsdfsdfsdfsdf"))
		time.Sleep(time.Second)
		w.Write(data)
	}))
	http.ListenAndServe(":8080", cacheHandler)
}

// CacheHandler is a middleware that sets cache headers in the HTTP response.
func CacheHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Set cache headers in the response
		w.Header().Set("Cache-Control", "public, max-age=3600")
		w.Header().Set("Expires", time.Now().Add(time.Hour).Format(http.TimeFormat))
		w.Header().Set("Last-Modified", time.Now().Format(http.TimeFormat))

		// Call the next handler in the chain
		next.ServeHTTP(w, r)
	})
}