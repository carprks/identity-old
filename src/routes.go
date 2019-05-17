package src

import (
	"fmt"
	"github.com/carprks/identity/src/healthcheck"
	"github.com/carprks/identity/src/identity"
	"github.com/carprks/identity/src/probe"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"os"
	"time"
)

// Routes get the routes for the service
func Routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.Timeout(60 * time.Second))
	router.Use(presetHeaders)

	// health check
	router.Get(fmt.Sprintf("%s/healthcheck", os.Getenv("SITE_PREFIX")), healthcheck.HTTP)

	// Probe
	router.Get("/probe", probe.HTTP)
	router.Get(fmt.Sprintf("%s/probe", os.Getenv("SITE_PREFIX")), probe.HTTP)

	// Create
	router.Post(fmt.Sprintf("%s/", os.Getenv("SITE_PREFIX")), identity.Create)

	// Retrieve
	router.Get(fmt.Sprintf("%s/", os.Getenv("SITE_PREFIX")), identity.RetrieveAllIdentities)

	// User
	router.Route(fmt.Sprintf("%s/{identityID}", os.Getenv("SITE_PREFIX")), func(r chi.Router) {
		r.Get("/", identity.Retrieve)
		r.Put("/", identity.Update)
		r.Delete("/", identity.Delete)
	})

	return router
}
