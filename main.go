package main

import (
	"fmt"
	"github.com/carprks/identity/src/healthcheck"
	"github.com/carprks/identity/src/identity"
	"github.com/carprks/identity/src/probe"
	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"time"
)

func presetHeaders(next http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()

		w.Header().Set("Strict-Transport-Security", "max-age=1000; includeSubdomains; preload")
		w.Header().Set("Content-Security-Policy", "upgrade-insecure-requests")
		w.Header().Set("Feature-Policy", "vibrate 'none'; geolocation 'none'")
		next.ServeHTTP(w, r.WithContext(ctx))
	}
	return http.HandlerFunc(fn)
}

func _main(args []string) int {
	// development
	if len(args) >= 1 {
		if args[0] == "localDev" {
			err := godotenv.Load()
			if err != nil {
				fmt.Println(fmt.Errorf("godotenv err: %v", err))
			}
			fmt.Println("Running LocalDev")
		}
	}

	port := "80"
	if len(os.Getenv("PORT")) >= 2 {
		port = os.Getenv("PORT")
	}

	// Router
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

	// Start Server
	fmt.Println(fmt.Sprintf("Starting Server on Port :%s", port))

	if err := http.ListenAndServe(fmt.Sprintf(":%s", port), router); err != nil {
		fmt.Println(fmt.Sprintf("HTTP: %v", err))
		return 1
	}

	fmt.Println("Died but nicely")
	return 0
}

func main() {
	os.Exit(_main(os.Args[1:]))
}
