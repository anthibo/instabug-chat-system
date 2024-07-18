package main

import (
	"net/http"

	"github.com/go-chi/chi"
)

func (app *AppDIContainer) routes() http.Handler {
	router := chi.NewRouter()

	router.Post("/applications/{token}/chats/{number}/messages", app.ApiCmdHandlers.CreateMessageCmdHandler)

	return router
}
