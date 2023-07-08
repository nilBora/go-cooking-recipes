package recipe

import (
   //"fmt"
   "go-cooking-recipes/v1/app/backend/store"
)

type Recipe struct {
    Store       store.Store
    Uuid        string
    Name        string
	Description string
	Text        string
    Image       string
	Labels      string
}

func (r Recipe) Create() {
     sql := `INSERT INTO "recipes"("uuid", "name", "description", "text") VALUES($1, $2, $3, $4)`
     _, err := r.Store.Connection.Exec(sql, r.Uuid, r.Name, r.Description, r.Text)
     if err != nil {
        panic(err)
     }
}
