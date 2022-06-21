package main

import (
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *application) routes() http.Handler {

	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.requireActivatedUser(app.healthcheckHandler))
	router.HandlerFunc(http.MethodGet, "/v1/rooms", app.requireActivatedUser(app.listRoomHandler))
	router.HandlerFunc(http.MethodPost, "/v1/rooms", app.requireActivatedUser(app.createRoomHandler))
	router.HandlerFunc(http.MethodGet, "/v1/rooms/:id", app.requireActivatedUser(app.showRoomHandler))
	router.HandlerFunc(http.MethodPatch, "/v1/rooms/:id", app.requireActivatedUser(app.updateRoomHandler))
	router.HandlerFunc(http.MethodDelete, "/v1/rooms/:id", app.requireActivatedUser(app.deleteRoomHandler))

	//router.HandlerFunc(http.MethodGet, "/v1/furniture", app.listFurnitureHandler) //TODO: implement handler
	//router.HandlerFunc(http.MethodPost, "/v1/furniture", app.createFurnitureHandler) //TODO: implement handler
	//router.HandlerFunc(http.MethodGet, "/v1/furniture/:id", app.showFurnitureHandler) //TODO: implement handler
	//router.HandlerFunc(http.MethodPatch, "/v1/furniture/:id", app.updateFurnitureHandler) //TODO: implement handler
	//router.HandlerFunc(http.MethodDelete, "/v1/furniture/:id", app.deleteFurnitureHandler) //TODO: implement handler

	router.HandlerFunc(http.MethodPost, "/v1/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPut, "/v1/users/activated", app.activateUserHandler)

	router.HandlerFunc(http.MethodPost, "/v1/tokens/authentication", app.createAuthenticationTokenHandler)

	return app.recoverPanic(app.rateLimit(app.authenticate(router)))
}
