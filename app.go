package main

import (
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"gopkg.in/validator.v2"
	"log"
	"net/http"
)

type App struct {
	Router *mux.Router
}

type ShortenLinkRequest struct {
	Url                 string `json:"url" validate:"nonzero"`
	ExpirationInMinutes int64  `json:"expirationInMinutes" validate:"min=0"`
}

type ShortenLinkResponse struct {
	ShortLink string `json:"shortLink"`
}

func (app *App) Initialise() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	app.Router = mux.NewRouter()
	app.initialiseRoutes()
}

func (app *App) initialiseRoutes() {
	app.Router.HandleFunc("/api/shorten", app.createShortLink).Methods("POST")
	app.Router.HandleFunc("/api/info", app.getShortLinkInfo).Methods("GET")
	app.Router.HandleFunc("/{shortLink:[a-zA-Z0-9]{1,11}}", app.redirect).Methods("GET")
}

func (app *App) createShortLink(writer http.ResponseWriter, request *http.Request) {
	var shortenLinkRequest ShortenLinkRequest
	if err := json.NewDecoder(request.Body).Decode(&shortenLinkRequest); err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Errorf("parse parameters failed: %v", request))
		return
	}
	if err := validator.Validate(shortenLinkRequest); err != nil {
		respondWithError(writer, http.StatusBadRequest, fmt.Errorf("validate parameters failed: %v", request))
		return
	}
	defer request.Body.Close()
	fmt.Printf("%v\n", shortenLinkRequest)
}

func (app *App) getShortLinkInfo(writer http.ResponseWriter, request *http.Request) {
	vals := request.URL.Query()
	info := vals.Get("shortLink")
	fmt.Printf("shortLink: %s\n", info)
}

func (app *App) redirect(writer http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	fmt.Printf("shortLink: %s\n", vars["shortLink"])
}

func respondWithError(writer http.ResponseWriter, code int, err error) {
	switch e := err.(type) {
	case Error:
		log.Printf("HTTP %d - %s\n", e.Status(), e)
		respondWithJson(writer, e.Status(), e.Error())
	default:
		respondWithJson(writer, http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))
	}
}

func respondWithJson(writer http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	writer.WriteHeader(code)
	writer.Header().Set("Content-Type", "application/json")
	writer.WriteHeader(code)
	writer.Write(response)
}

func (app *App) Run(address string) {
	log.Fatal(http.ListenAndServe(address, app.Router))
}