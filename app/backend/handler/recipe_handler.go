package handler

import (
    "io"
    "fmt"
    "net/http"
    "database/sql"
    "strconv"
    "github.com/go-chi/chi/v5"
    "github.com/go-chi/render"
    recipe "go-cooking-recipes/v1/app/backend/repository"
    "encoding/json"
    "github.com/google/uuid"
)

type JSON map[string]interface{}

type Handler struct {
    Connection *sql.DB
    RecipeRepository recipe.RecipeRepositoryInterface
}

func NewHandler(conn *sql.DB) Handler {
    return Handler{Connection: conn}
}


func (h Handler) OnListRecipe(w http.ResponseWriter, r *http.Request) {
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    order := r.URL.Query().Get("order")
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))

    param := recipe.ListParams{
        Limit: limit,
        Order: order,
        Page: page,
     }
    recipeRepository := recipe.NewRecipeRepository(h.Connection)
    listData, err := recipeRepository.GetList(param)

    if err != nil {
         render.Status(r, http.StatusNotFound)
         render.JSON(w, r, JSON{"status": "error"})
         return
    }

    render.JSON(w, r, JSON{"status": "ok", "data": listData})
}

func (h Handler) OnGetOneRecipe(w http.ResponseWriter, r *http.Request) {
    uuid := chi.URLParam(r, "uuid")

    recipeRepository := recipe.NewRecipeRepository(h.Connection)
    row, err := recipeRepository.GetOne(uuid)

    if err != nil {
         render.Status(r, http.StatusNotFound)
         render.JSON(w, r, JSON{"status": "error"})
         return
    }

    render.JSON(w, r, JSON{"status": "ok", "data": row})
}

func (h Handler) OnDeleteOneRecipe(w http.ResponseWriter, r *http.Request) {
    uuid := chi.URLParam(r, "uuid")

    recipeRepository := recipe.NewRecipeRepository(h.Connection)
    err := recipeRepository.Remove(uuid)
    if err != nil {
         render.Status(r, http.StatusNotFound)
         render.JSON(w, r, JSON{"status": "error"})
         return
    }
    render.Status(r, http.StatusNoContent)
    render.JSON(w, r, JSON{"status": "deleted"})
}

func (h Handler) OnCreateRecipe(w http.ResponseWriter, r *http.Request) {

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

     labels := recipeData["labels"].(map[string]interface{})
     var labelsArray []string
     for _, v := range labels {
        labelsArray = append(labelsArray, v.(string))
     }

     rec := recipe.Recipe{
         Uuid: uuid,
         Name: recipeData["name"].(string),
         Description: recipeData["description"].(string),
         Text: recipeData["text"].(string),
         Image: recipeData["image"].(string),
         Labels: labelsArray,
     }
     recipeRepository := recipe.NewRecipeRepository(h.Connection)
     err = recipeRepository.Create(rec)
     if err != nil {
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, JSON{"status": "error", "message": err})
        return
     }

     render.Status(r, http.StatusCreated)
     render.JSON(w, r, JSON{"status": "ok", "uuid": uuid})
}

func (h Handler) OnChangeOneRecipe(w http.ResponseWriter, r *http.Request) {
    uuid := chi.URLParam(r, "uuid")

    var recipeData JSON

    b, err := io.ReadAll(r.Body)
    if err != nil {
        fmt.Printf("[ERROR] %s", err)
    }

    err = json.Unmarshal(b, &recipeData)

     if err != nil {
        fmt.Println("Error while decoding the data", err.Error())
     }

    labels := recipeData["labels"].(map[string]interface{})
    var labelsArray []string
    for _, v := range labels {
        labelsArray = append(labelsArray, v.(string))
    }

     rec := recipe.Recipe{
         Name: recipeData["name"].(string),
         Description: recipeData["description"].(string),
         Text: recipeData["text"].(string),
         Image: recipeData["image"].(string),
         Labels: labelsArray,
     }
     recipeRepository := recipe.NewRecipeRepository(h.Connection)
     _, err = recipeRepository.Change(uuid, rec)
     if err != nil {
        render.Status(r, http.StatusBadRequest)
        render.JSON(w, r, JSON{"status": "error", "message": err})
        return
     }

     render.Status(r, http.StatusOK)
     render.JSON(w, r, JSON{"status": "ok", "uuid": uuid})
}

