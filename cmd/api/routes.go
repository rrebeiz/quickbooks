package main

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := chi.NewRouter()

	router.Use(middleware.Recoverer)
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))
	router.Group(func(router chi.Router) {
		router.Use(app.authTokenMiddleware)
		router.Get("/v1/users/auth", app.authenticateToken)
		router.Get("/v1/users/{id}", app.getUserHandler)
		router.Patch("/v1/users/{id}", app.updateUserHandler)
		router.Post("/v1/books", app.createBookHandler)
		router.Patch("/v1/books/{id}", app.updateBookHandler)
		router.Delete("/v1/books/{id}", app.deleteBookHandler)
	})
	router.Group(func(router chi.Router) {
		router.Use(app.adminMiddleware)
		router.Get("/v1/users", app.getAllUsersHandler)
		router.Get("/v1/users/authenticated", app.getAllAuthenticatedUsersHandler)
		router.Delete("/v1/users/{id}", app.deleteUserHandler)
		router.Delete("/v1/users/logout/{id}", app.adminLogoutHandler)
	})

	router.Get("/healthcheck", app.healthCheckHandler)

	router.Post("/v1/users/login", app.loginHandler)
	router.Get("/v1/users/logout", app.logoutHandler)
	router.Post("/v1/users", app.createUserHandler)
	// Book routes
	router.Get("/v1/books", app.getAllBooksHandler)
	router.Get("/v1/books/{id}", app.getBookByIDHandler)
	router.Get("/v1/books/slug", app.getBookBySlugHandler)

	return router
}
