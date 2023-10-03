package repository

import (
    "database/sql"
    "log"
    "errors"
   "strconv"
)

type TagRepositoryInterface interface {
    Create(t Tag) error
    GetList(params ListParams) ([]Tag, error)
    GetOne(uuid string) (Tag, error)
    Remove(uuid string) (error)
    Change(uuid string, t Tag) (int64, error)
}

func NewtTagRepository(conn *sql.DB) TagRepositoryInterface {
	return &TagRepository{
		Connection: conn,
	}
}

type TagRepository struct {
    Connection  *sql.DB
}


type Tag struct {
    Uuid string
    Name string
}

func (repo TagRepository) Create(t Tag) error {
     sql := `INSERT INTO "tags"("uuid", "name") VALUES($1, $2)`

     _, err := repo.Connection.Exec(sql, t.Uuid, t.Name)
     if err != nil {
        return errors.New("Couldn't create tag")
     }

     return nil
}

func (repo TagRepository) GetList(params ListParams) ([]Tag, error) {
    sql := `SELECT uuid, name FROM "tags"`

    limit := params.Limit
    offset := params.Limit * (params.Page - 1)

    if params.Limit > 0 {
        sql = sql + ` LIMIT ` + strconv.Itoa(limit) + ` OFFSET ` + strconv.Itoa(offset)
    }

    rows, err := repo.Connection.Query(sql)

    tags := []Tag{}

    if err != nil {
        return tags, errors.New("Rows Not Found")
    }

    for rows.Next() {
        tag := Tag{}
        if err := rows.Scan(&tag.Uuid, &tag.Name); err != nil {
            log.Fatalf("could not scan row: %v", err)
        }
        tags = append(tags, tag)
    }
    return tags, nil
}

func (repo TagRepository) GetOne(uuid string) (Tag, error) {
    tag := Tag{}

    sql := `SELECT uuid, name FROM "tags" WHERE uuid = $1`
    row := repo.Connection.QueryRow(sql, uuid)

    err := row.Scan(&tag.Uuid, &tag.Name)

    if err != nil {
        return tag, errors.New("Row Not Found")
    }

    return tag, nil
}

func (repo TagRepository) Remove(uuid string) (error) {
    sql := `DELETE FROM "tags" WHERE uuid = $1`
    _, err := repo.Connection.Exec(sql, uuid)

    if err != nil {
        return errors.New("Can't delete row")
    }

    //count, err := res.RowsAffected()

    return nil
}

func (repo TagRepository) Change(uuid string, t Tag) (int64, error) {
    sql := `UPDATE "tags" SET name = $1 WHERE uuid = $1`
    res, err := repo.Connection.Exec(sql, t.Name, uuid)
    if err != nil {
        return 0, errors.New("Can't update row")
    }
    count, err := res.RowsAffected()

    if err != nil {
        return 0, err
    }

    return count, nil
}
