package repository

import (
    //"database/sql"
    "log"
    "errors"
   //"fmt"
)

type Recipe struct {
    Uuid        string
    Name        string
	Description string
	Text        string
    Image       string
	Labels      string
}

func (repo Repository) Create(r Recipe) {
     sql := `INSERT INTO "recipes"("uuid", "name", "description", "text") VALUES($1, $2, $3, $4)`
     _, err := repo.Connection.Exec(sql, r.Uuid, r.Name, r.Description, r.Text)
     if err != nil {
        panic(err)
     }
}

func (repo Repository) GetList() []Recipe {
    sql := `SELECT uuid, name, description FROM "recipes"`
    rows, err := repo.Connection.Query(sql)
    if err != nil {
        panic(err)
    }
    recipes := []Recipe{}
    for rows.Next() {
        recipe := Recipe{}
        if err := rows.Scan(&recipe.Uuid, &recipe.Name, &recipe.Description); err != nil {
            log.Fatalf("could not scan row: %v", err)
        }

        recipes = append(recipes, recipe)
    }
    return recipes
}

func (repo Repository) GetOne(uuid string) (Recipe, error) {
    recipe := Recipe{}

    sql := `SELECT uuid, name, description FROM "recipes" WHERE uuid = $1`
    row := repo.Connection.QueryRow(sql, uuid)

    err := row.Scan(&recipe.Uuid, &recipe.Name, &recipe.Description)

    if err != nil {
        return recipe, errors.New("Row Not Found")
    }

    return recipe, nil
}

func (repo Repository) Remove(uuid string) (error) {
    sql := `DELETE FROM "recipes" WHERE uuid = $1`
    _, err := repo.Connection.Exec(sql, uuid)

    if err != nil {
        return errors.New("Can't delete row")
    }

    //count, err := res.RowsAffected()

    return nil
}
