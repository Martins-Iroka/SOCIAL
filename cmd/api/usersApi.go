package main

import (
	"net/http"
	"strconv"

	"github.com/Martins-Iroka/social/internal/store"
	"github.com/go-chi/chi/v5"
)

func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := app.store.Users.GetUserByID(ctx, userID)
	if err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFoundResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := jsonResponse(w, http.StatusOK, user); err != nil {
		app.internalServerError(w, r, err)
		return
	}
}
