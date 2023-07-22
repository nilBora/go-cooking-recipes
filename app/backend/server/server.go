package server

import (
   "io"
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
   "encoding/json"
   recipe "go-cooking-recipes/v1/app/backend/repository"
   repository "go-cooking-recipes/v1/app/backend/repository"
   "github.com/google/uuid"
)

type JSON map[string]interface{}

type Server struct {
    Port           string
	PinSize        int
	MaxPinAttempts int
	MaxExpire      time.Duration
	WebRoot        string
	Version        string
	Repository     repository.Repository

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
	    //r.Use(Authentication)
		r.Get("/recipes", s.onListRecipe)
		r.Get("/recipes/{uuid}", s.onGetOneRecipe)
		r.Delete("/recipes/{uuid}", s.onDeleteOneRecipe)
        r.Post("/recipes", s.onCreateRecipe)
        r.Post("/recipes/{uuid}", s.onChangeOneRecipe)
	})

	router.Get("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
		render.PlainText(w, r, "User-agent: *\nDisallow: /api/\nDisallow: /show/\n")
	})

	return router
}

func (s Server) onListRecipe(w http.ResponseWriter, r *http.Request) {

    limit := r.URL.Query().Get("limit")
    offset := r.URL.Query().Get("offset")
    order := r.URL.Query().Get("order")

    listData, err := s.Repository.GetList()

    if err != nil {
         render.Status(r, http.StatusNotFound)
         render.JSON(w, r, JSON{"status": "error"})
         return
    }

    render.JSON(w, r, JSON{"status": "ok", "data": listData})
}

func (s Server) onGetOneRecipe(w http.ResponseWriter, r *http.Request) {
    uuid := chi.URLParam(r, "uuid")

    row, err := s.Repository.GetOne(uuid)

    if err != nil {
         render.Status(r, http.StatusNotFound)
         render.JSON(w, r, JSON{"status": "error"})
         return
    }

    render.JSON(w, r, JSON{"status": "ok", "data": row})
}

func (s Server) onDeleteOneRecipe(w http.ResponseWriter, r *http.Request) {
    uuid := chi.URLParam(r, "uuid")

    err := s.Repository.Remove(uuid)
    if err != nil {
         render.Status(r, http.StatusNotFound)
         render.JSON(w, r, JSON{"status": "error"})
         return
    }
    render.Status(r, http.StatusNoContent)
    render.JSON(w, r, JSON{"status": "deleted"})
}

func (s Server) onCreateRecipe(w http.ResponseWriter, r *http.Request) {

    var recipeData JSON

    b, err := io.ReadAll(r.Body)
    if err != nil {
        fmt.Printf("[ERROR] %s", err)
    }

    err = json.Unmarshal(b, &recipeData)

     if err != nil {
        fmt.Println("Error while decoding the data", err.Error())
     }
     uuid := uuid.New().String()

     rec := recipe.Recipe{
         Uuid: uuid,
         Name: recipeData["name"].(string),
         Description: recipeData["description"].(string),
         Text: recipeData["text"].(string),
         Image: recipeData["image"].(string),
         Labels: recipeData["labels"].(string),
     }

     err = s.Repository.Create(rec)
     if err != nil {
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, JSON{"status": "error", "message": err})
     }

     render.Status(r, http.StatusCreated)
     render.JSON(w, r, JSON{"status": "ok", "uuid": uuid})
}

func (s Server) onChangeOneRecipe(w http.ResponseWriter, r *http.Request) {
	 fmt.Fprint(w, "Change")
}
