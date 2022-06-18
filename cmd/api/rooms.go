package main

import (
	"fmt"
	"net/http"
	"time"

	"github.com/WrastAct/EHome/internal/data"
)

func (app *application) editRoomHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "open room editor")
}

func (app *application) createRoomHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "create a new room")
}

func (app *application) showRoomHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	furniture1 := data.Furniture{
		ID:          1,
		Name:        "Chair",
		Price:       423.75,
		Description: "heh",
		X:           10,
		Y:           25,
		Width:       25,
		Height:      25,
		Image:       "../img",
		Shape:       data.Circle,
	}

	furniture2 := data.Furniture{
		ID:     2,
		Name:   "Table",
		Price:  103.24,
		X:      10,
		Y:      25,
		Width:  25,
		Height: 25,
		Image:  "../img",
		Shape:  data.Circle,
	}

	room := data.Room{
		ID:            id,
		Data:          time.Now(),
		Description:   "Hehw",
		Title:         "Custom room",
		Width:         500,
		Height:        300,
		FurnitureList: []data.Furniture{furniture1, furniture2},
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"room": room}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
