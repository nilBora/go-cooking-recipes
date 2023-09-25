package repository

import (
    "database/sql"
    "log"
    "errors"
  // "fmt"
   "strconv"
)

type RecipeRepositoryInterface interface {
    Create(r Recipe) error
    GetList(params ListParams) ([]Recipe, error)
    GetOne(uuid string) (Recipe, error)
    Remove(uuid string) (error)
    Change(uuid string, r Recipe) (int64, error)
}

func NewRecipeRepository(conn *sql.DB) RecipeRepositoryInterface {
	return &RecipeRepository{
		Connection: conn,
	}
}

type RecipeRepository struct {
    Connection  *sql.DB
}


type Recipe struct {
    Uuid        string
    Name        string
	Description string
	Text        string
    Image       string
	Labels      []string
}

type ListParams struct {
    Limit int
    Offset int
    Order string
    Page int
    Size int
}

// func NewRecipeRepository(conn *sql.DB) *RecipeRepository {
// 	return &RecipeRepository{
// 		Connection: conn,
// 	}
// }

func (repo RecipeRepository) Create(r Recipe) error {
     sql := `INSERT INTO "recipes"("uuid", "name", "description", "text") VALUES($1, $2, $3, $4)`
     _, err := repo.Connection.Exec(sql, r.Uuid, r.Name, r.Description, r.Text)
     if err != nil {
        return errors.New("Couldn't create recipe")
     }

     return nil
}

func (repo RecipeRepository) GetList(params ListParams) ([]Recipe, error) {
    sql := `SELECT uuid, name, description FROM "recipes"`

    limit := params.Limit
    offset := params.Limit * (params.Page - 1)

    if params.Limit > 0 {
        sql = sql + ` LIMIT ` + strconv.Itoa(limit) + ` OFFSET ` + strconv.Itoa(offset)
    }

    rows, err := repo.Connection.Query(sql)

    recipes := []Recipe{}

    if err != nil {
        return recipes, errors.New("Rows Not Found")
    }

    for rows.Next() {
        recipe := Recipe{}
        if err := rows.Scan(&recipe.Uuid, &recipe.Name, &recipe.Description); err != nil {
            log.Fatalf("could not scan row: %v", err)
        }
        recipe.Labels = []string{"test", "test2"}
        recipes = append(recipes, recipe)
    }
    return recipes, nil
}

func (repo RecipeRepository) GetOne(uuid string) (Recipe, error) {
    recipe := Recipe{}

    sql := `SELECT uuid, name, description FROM "recipes" WHERE uuid = $1`
    row := repo.Connection.QueryRow(sql, uuid)

    err := row.Scan(&recipe.Uuid, &recipe.Name, &recipe.Description)

    if err != nil {
        return recipe, errors.New("Row Not Found")
    }

    return recipe, nil
}

func (repo RecipeRepository) Remove(uuid string) (error) {
    sql := `DELETE FROM "recipes" WHERE uuid = $1`
    _, err := repo.Connection.Exec(sql, uuid)

    if err != nil {
        return errors.New("Can't delete row")
    }

    //count, err := res.RowsAffected()

    return nil
}

func (repo RecipeRepository) Change(uuid string, r Recipe) (int64, error) {
    sql := `UPDATE "recipes" SET name = $1, description = $2, text = $3 WHERE uuid = $4`
    res, err := repo.Connection.Exec(sql, r.Name, r.Description, r.Text, uuid)
    if err != nil {
        return 0, errors.New("Can't update row")
    }
    count, err := res.RowsAffected()

    if err != nil {
        return 0, err
    }

    return count, nil
}
