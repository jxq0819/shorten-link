package main

import (
	"log"
	"net/http"
	"time"
)

type Middleware struct {
}

func (middleware Middleware) LoggingHandler(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		time1 := time.Now()
		next.ServeHTTP(writer, request)
		time2 := time.Now()
		log.Printf("[%s] %q %v\n", request.Method, request.URL.String(), time2.Sub(time1))
	}

	return http.HandlerFunc(fn)
}

func (middleware Middleware) RecoverHandler(next http.Handler) http.Handler {
	fn := func(writer http.ResponseWriter, request *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				log.Printf("recover from panic: %+v", err)
				http.Error(writer, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
			}
		}()
		next.ServeHTTP(writer, request)
	}

	return http.HandlerFunc(fn)
}
