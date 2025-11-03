package main

import (
	"net/http"
	"strconv"

	"github.com/Martins-Iroka/social/internal/store"
	"github.com/go-chi/chi/v5"
)

type FollowUser struct {
	UserID int64 `json:"user_id"`
}

type userKey string

const userContextKey userKey = "user"

// GetUser godoc
//
//	@summary		Fetches a user
//	@description	Fetches a user profile by ID
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			userID	path		int	true	"User ID"
//	@success		200		{object}	store.User
//	@failure		400		{object}	error
//	@failure		404		{object}	error
//	@failure		500		{object}	error
//	@security		ApiKeyAuth
//	@router			/users/{userID}	[get]
func (app *application) getUserHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := strconv.ParseInt(chi.URLParam(r, "userID"), 10, 64)

	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	ctx := r.Context()

	user, err := app.getUser(ctx, userID)
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
	}
}

// Follow User godoc
//
//	@summary		Follows a user
//	@description	Follows a user by ID
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			userID	path		int		true	"User ID"
//	@success		204		{string}	string	"User followed"
//	@failure		400		{object}	error
//	@failure		404		{object}	error	"User not found"
//	@failure		500		{object}	error
//	@security		ApiKeyAuth
//	@router			/users/{userID}/follow	[put]
func (app *application) followUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.User.FollowUser(r.Context(), user.ID, payload.UserID); err != nil {
		switch err {
		case store.ErrorUserFollowConflict:
			app.conflictResponse(w, r, err)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

// Unfollow User godoc
//
//	@summary		Unfollows a user
//	@description	Unfollows a user by ID
//	@tags			users
//	@accept			json
//	@produce		json
//	@param			userID	path		int		true	"User ID"
//	@success		204		{string}	string	"User unfollowed"
//	@failure		400		{object}	error
//	@failure		404		{object}	error	"User not found"
//	@failure		500		{object}	error
//	@security		ApiKeyAuth
//	@router			/users/{id}/unfollow	[put]
func (app *application) unfollowUserHandler(w http.ResponseWriter, r *http.Request) {
	user := getUserFromContext(r)

	var payload FollowUser
	if err := readJSON(w, r, &payload); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := app.store.User.UnFollowUser(r.Context(), user.ID, payload.UserID); err != nil {
		switch err {
		case store.ErrorNotFound:
			app.conflictResponse(w, r, store.ErrorUserUnFollowConflict)
			return
		default:
			app.internalServerError(w, r, err)
			return
		}
	}

	if err := jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}

func getUserFromContext(r *http.Request) *store.User {
	user, _ := r.Context().Value(userContextKey).(*store.User)
	return user
}

// ActivateUser godoc
//
//	@summary		Activates/Register a user
//	@description	Activates/Register a user by invitation token
//	@tags			users
//	@produce		json
//	@param			token	path		string	true	"Invitation token"
//	@success		204		{string}	string	"User activated"
//	@failure		404		{object}	error
//	@failure		500		{object}	error
//	@security		ApiKeyAuth
//	@router			/users/activate/{token} [put]
func (app *application) activateUserHandler(w http.ResponseWriter, r *http.Request) {
	token := chi.URLParam(r, "token")

	if err := app.store.User.ActivateUser(r.Context(), token); err != nil {
		switch err {
		case store.ErrorNotFound:
			app.notFoundResponse(w, r, err)
		default:
			app.internalServerError(w, r, err)
		}
		return
	}

	if err := jsonResponse(w, http.StatusNoContent, nil); err != nil {
		app.internalServerError(w, r, err)
	}
}
