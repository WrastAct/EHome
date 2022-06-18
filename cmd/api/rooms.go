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
		http.NotFound(w, r)
		return
	}

	furniture1 := data.Furniture{
		ID:     1,
		Name:   "Chair",
		X:      10,
		Y:      25,
		Width:  25,
		Height: 25,
		Image:  "../img",
		Shape:  data.Circle,
	}

	furniture2 := data.Furniture{
		ID:     2,
		Name:   "Table",
		X:      10,
		Y:      25,
		Width:  25,
		Height: 25,
		Image:  "../img",
		Shape:  data.Circle,
	}

	room := data.Room{
		ID:            id,
		CreatedAt:     time.Now(),
		Title:         "Custom room",
		Width:         500,
		Height:        300,
		FurnitureList: []data.Furniture{furniture1, furniture2},
	}

	err = app.writeJSON(w, http.StatusOK, room, nil)
	if err != nil {
		app.logger.Println(err)
		http.Error(w, "The server encountered a problem and could not process your request", http.StatusInternalServerError)
	}

}
