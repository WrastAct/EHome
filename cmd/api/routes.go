package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {

	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	router.HandlerFunc(http.MethodGet, "/v1/editor", app.editRoomHandler)
	router.HandlerFunc(http.MethodPost, "/v1/rooms", app.createRoomHandler)
	router.HandlerFunc(http.MethodGet, "/v1/rooms/:id", app.showRoomHandler)

	return router
}
