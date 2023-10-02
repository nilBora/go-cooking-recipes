package repository

import (
    "database/sql"
    "os"
)

type Repository struct {
    Connection  *sql.DB
}

type ListParams struct {
    Limit int
    Offset int
    Order string
    Page int
    Size int
}

func ConnectDB() Repository {
    db, err := sql.Open("postgres", os.Getenv("POSTGRES_DSN"))
    if err != nil {
        panic(err)
    }
    rep := Repository{
        Connection: db,
    }
    return rep
}
