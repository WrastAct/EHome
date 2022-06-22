package main

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/WrastAct/EHome/internal/data"
	"github.com/WrastAct/EHome/internal/validator"
)

func (app *application) createFurnitureHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Name        string     `json:"name"`
		Price       float64    `json:"price"`
		Description string     `json:"description"`
		Width       int64      `json:"width"`
		Height      int64      `json:"height"`
		Image       string     `json:"image"` // Path to the image
		Shape       data.Shape `json:"shape"` // To improve collision detection
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	furniture := &data.Furniture{
		Name:        input.Name,
		Price:       input.Price,
		Description: input.Description,
		Width:       input.Width,
		Height:      input.Height,
		Image:       input.Image,
		Shape:       input.Shape,
	}

	v := validator.New()

	if data.ValidateFurniture(v, furniture); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Furniture.Insert(furniture)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	headers := make(http.Header)
	headers.Set("Location", fmt.Sprintf("/v1/furniture/%d", furniture.ID))

	err = app.writeJSON(w, http.StatusCreated, envelope{"furniture": furniture}, headers)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) showFurnitureHandler(w http.ResponseWriter, r *http.Request) {

	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	furniture, err := app.models.Furniture.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"furniture": furniture}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) updateFurnitureHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	furniture, err := app.models.Furniture.Get(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	var input struct {
		Name        *string     `json:"name"`
		Price       *float64    `json:"price"`
		Description *string     `json:"description"`
		Width       *int64      `json:"width"`
		Height      *int64      `json:"height"`
		Image       *string     `json:"image"` // Path to the image
		Shape       *data.Shape `json:"shape"` // To improve collision detection
	}

	err = app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if input.Name != nil {
		furniture.Name = *input.Name
	}

	if input.Price != nil {
		furniture.Price = *input.Price
	}

	if input.Description != nil {
		furniture.Description = *input.Description
	}

	if input.Width != nil {
		furniture.Width = *input.Width
	}

	if input.Height != nil {
		furniture.Height = *input.Height
	}

	if input.Image != nil {
		furniture.Image = *input.Image
	}

	if input.Shape != nil {
		furniture.Shape = *input.Shape
	}

	v := validator.New()

	if data.ValidateFurniture(v, furniture); !v.Valid() {
		app.failedValidationResponse(w, r, v.Errors)
		return
	}

	err = app.models.Furniture.Update(furniture)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"furniture": furniture}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) deleteFurnitureHandler(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}

	err = app.models.Furniture.Delete(id)
	if err != nil {
		switch {
		case errors.Is(err, data.ErrRecordNotFound):
			app.notFoundResponse(w, r)
		default:
			app.serverErrorResponse(w, r, err)
		}
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"message": "furniture successfully deleted"}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}

func (app *application) listFurnitureHandler(w http.ResponseWriter, r *http.Request) {
	furniture, err := app.models.Furniture.GetAll()
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}

	err = app.writeJSON(w, http.StatusOK, envelope{"furniture": furniture}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
