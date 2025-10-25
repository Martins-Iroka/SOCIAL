package main

import (
	"net/http"

	"github.com/Martins-Iroka/social/internal/store"
)

func (app *application) getUserFeedHandler(w http.ResponseWriter, r *http.Request) {

	fq := PaginatedFeedQueryAPi{
		Limit:  20,
		Offset: 0,
		Sort:   "desc",
	}

	fq, err := fq.Parse(r)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	if err := Validate.Struct(fq); err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	feedQuery := &store.PaginatedFeedQuery{
		Limit:  fq.Limit,
		Offset: fq.Offset,
		Sort:   fq.Sort,
		Tags:   fq.Tags,
		Search: fq.Search,
		Since:  fq.Since,
		Until:  fq.Until,
	}

	ctx := r.Context()

	feed, err := app.store.Post.GetUserFeed(ctx, int64(1), feedQuery)
	if err != nil {
		app.internalServerError(w, r, err)
		return
	}

	if err := jsonResponse(w, http.StatusOK, feed); err != nil {
		app.internalServerError(w, r, err)
	}
}
