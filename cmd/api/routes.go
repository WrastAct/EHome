package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() *httprouter.Router {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/rooms", app.editRoomHandler) //TODO: implement handler
	router.HandlerFunc(http.MethodPost, "/v1/rooms", app.createRoomHandler)
	router.HandlerFunc(http.MethodGet, "/v1/rooms/:id", app.showRoomHandler)
	router.HandlerFunc(http.MethodPut, "/v1/rooms/:id", app.updateRoomHandler)
	router.HandlerFunc(http.MethodDelete, "/v1/rooms/:id", app.deleteRoomHandler)
	//router.HandlerFunc(http.MethodGet, "/v1/furniture", app.editFurnitureHandler) //TODO: implement handler
	//router.HandlerFunc(http.MethodPost, "/v1/furniture", app.createFurnitureHandler) //TODO: implement handler
	//router.HandlerFunc(http.MethodGet, "/v1/furniture/:id", app.showFurnitureHandler) //TODO: implement handler
	//router.HandlerFunc(http.MethodDelete, "/v1/furniture/:id", app.deleteFurnitureHandler) //TODO: implement handler
	//router.HandlerFunc(http.MethodPut, "/v1/furniture/:id", app.updateFurnitureHandler) //TODO: implement handler

	return router
}
