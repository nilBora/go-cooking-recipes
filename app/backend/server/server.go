package server

import (
   "context"
   "time"
   "log"
   "net/http"
   "github.com/pkg/errors"
   "github.com/didip/tollbooth/v7"
   "github.com/didip/tollbooth_chi"
   "github.com/go-chi/chi/v5"
   "github.com/go-chi/chi/v5/middleware"
   "github.com/go-chi/render"
   "github.com/jtrw/go-rest"
   "fmt"
)

type Server struct {
    Port           string
	PinSize        int
	MaxPinAttempts int
	MaxExpire      time.Duration
	WebRoot        string
	Version        string
}

func (s Server) Run(ctx context.Context) error {
	log.Printf("[INFO] activate rest server")
    log.Printf("[INFO] Port: %s", s.Port)

	httpServer := &http.Server{
		Addr:              ":"+s.Port,
		Handler:           s.routes(),
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      30 * time.Second,
		IdleTimeout:       30 * time.Second,
	}

	go func() {
		<-ctx.Done()
		if httpServer != nil {
			if clsErr := httpServer.Close(); clsErr != nil {
				log.Printf("[ERROR] failed to close proxy http server, %v", clsErr)
			}
		}
	}()

	err := httpServer.ListenAndServe()
	log.Printf("[WARN] http server terminated, %s", err)

	if err != http.ErrServerClosed {
		return errors.Wrap(err, "server failed")
	}
	return nil
}

func (s Server) routes() chi.Router {
	router := chi.NewRouter()

	router.Use(middleware.RequestID, middleware.RealIP)
	router.Use(middleware.Throttle(1000), middleware.Timeout(60*time.Second))
	//router.Use(rest.AppInfo("recipe", "Jrtw", s.Version), rest.Ping)
	router.Use(rest.Ping)
	router.Use(tollbooth_chi.LimitHandler(tollbooth.NewLimiter(10, nil)))
    router.Use(middleware.Logger)

	router.Route("/api/v1", func(r chi.Router) {
		r.Get("/recipe", s.onListRecipe)
        r.Post("/recipe", s.onCreateRecipe)
	})

	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "User-agent: *\nDisallow: /api/\nDisallow: /show/\n")
	})

	return router
}

func (s Server) onListRecipe(w http.ResponseWriter, r *http.Request) {
	 fmt.Fprint(w, "List")
}

func (s Server) onCreateRecipe(w http.ResponseWriter, r *http.Request) {
	 fmt.Fprint(w, "Create")
}