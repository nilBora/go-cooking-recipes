package handler

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    repository "go-cooking-recipes/v1/app/backend/repository"
)

type Handler struct {
    Repository repository.Repository
}

func NewHandler() Handler {
    return Handler{}
}


func (h Handler) onListRecipe(w http.ResponseWriter, r *http.Request) {
    limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
    order := r.URL.Query().Get("order")
    page, _ := strconv.Atoi(r.URL.Query().Get("page"))

    param := recipe.ListParams{
        Limit: limit,
        Order: order,
        Page: page,
     }
    recipeRepository := recipe.NewRecipeRepository(h.Repository.Connection)
    listData, err := recipeRepository.GetList(param)

    if err != nil {
         render.Status(r, http.StatusNotFound)
         render.JSON(w, r, JSON{"status": "error"})
         return
    }

    render.JSON(w, r, JSON{"status": "ok", "data": listData})
}
