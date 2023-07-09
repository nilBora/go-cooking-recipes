package repository

import (
    "database/sql"
   //"fmt"
   //"go-cooking-recipes/v1/app/backend/store"
)


type Recipe struct {
    Connection  *sql.DB
    Uuid        string
    Name        string
	Description string
	Text        string
    Image       string
	Labels      string
}

func (r Recipe) Create() {
     sql := `INSERT INTO "recipes"("uuid", "name", "description", "text") VALUES($1, $2, $3, $4)`
     _, err := r.Connection.Exec(sql, r.Uuid, r.Name, r.Description, r.Text)
     if err != nil {
        panic(err)
     }
}

func (r Recipe) getList() {
    sql := `SELECT * FROM "recipes"`
     rows, err := r.Connection.Query(sql)
     if err != nil {
        panic(err)
     }
     for rows.Next() {
     }
}
