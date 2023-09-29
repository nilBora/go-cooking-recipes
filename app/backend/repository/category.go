package repository

import (
    "database/sql"
    "log"
    "errors"
   "strconv"
)

type CategoryRepositoryInterface interface {
    Create(c Category) error
    GetList(params ListParams) ([]Category, error)
    GetOne(uuid string) (Category, error)
    Remove(uuid string) (error)
    Change(uuid string, c Category) (int64, error)
}

func NewCategoryRepository(conn *sql.DB) CategoryRepositoryInterface {
	return &CategoryRepository{
		Connection: conn,
	}
}

type CategoryRepository struct {
    Connection  *sql.DB
}


type Category struct {
    Uuid        string
    Name        string
	Description string
	Text        string
    Image       string
}

type ListParams struct {
    Limit int
    Offset int
    Order string
    Page int
    Size int
}

func (repo CategoryRepository) Create(c Category) error {
     sql := `INSERT INTO "categories"("uuid", "name", "description", "text") VALUES($1, $2, $3, $4)`

     _, err := repo.Connection.Exec(sql, c.Uuid, c.Name, c.Description, c.Text)
     if err != nil {
        return errors.New("Couldn't create category")
     }

     return nil
}

func (repo CategoryRepository) GetList(params ListParams) ([]Category, error) {
    sql := `SELECT uuid, name, description FROM "categories"`

    limit := params.Limit
    offset := params.Limit * (params.Page - 1)

    if params.Limit > 0 {
        sql = sql + ` LIMIT ` + strconv.Itoa(limit) + ` OFFSET ` + strconv.Itoa(offset)
    }

    rows, err := repo.Connection.Query(sql)

    categories := []Category{}

    if err != nil {
        return category, errors.New("Rows Not Found")
    }

    for rows.Next() {
        category := Category{}
        if err := rows.Scan(&category.Uuid, &category.Name, &category.Description); err != nil {
            log.Fatalf("could not scan row: %v", err)
        }
        categories = append(categories, category)
    }
    return categories, nil
}

func (repo CategoryRepository) GetOne(uuid string) (Category, error) {
    category := Category{}

    sql := `SELECT uuid, name, description FROM "categories" WHERE uuid = $1`
    row := repo.Connection.QueryRow(sql, uuid)

    err := row.Scan(&category.Uuid, &category.Name, &category.Description)

    if err != nil {
        return category, errors.New("Row Not Found")
    }

    return category, nil
}

func (repo CategoryRepository) Remove(uuid string) (error) {
    sql := `DELETE FROM "categories" WHERE uuid = $1`
    _, err := repo.Connection.Exec(sql, uuid)

    if err != nil {
        return errors.New("Can't delete row")
    }

    //count, err := res.RowsAffected()

    return nil
}

func (repo CategoryRepository) Change(uuid string, c Category) (int64, error) {
    sql := `UPDATE "categories" SET name = $1, description = $2, text = $3 WHERE uuid = $4`
    res, err := repo.Connection.Exec(sql, c.Name, c.Description, c.Text, uuid)
    if err != nil {
        return 0, errors.New("Can't update row")
    }
    count, err := res.RowsAffected()

    if err != nil {
        return 0, err
    }

    return count, nil
}
