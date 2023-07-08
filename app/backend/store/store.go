package store

import (
   "os"
   "fmt"
   "database/sql"
    _ "github.com/lib/pq"
)

type Store struct {
    Connection *sql.DB
}


func ConnectDB() Store {
    psqlconn := fmt.Sprintf(
        "host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
        os.Getenv("POSTGRES_HOST"),
        os.Getenv("POSTGRES_PORT"),
        os.Getenv("POSTGRES_USER"),
        os.Getenv("POSTGRES_PASSWORD"),
        os.Getenv("POSTGRESQL_DB"))
    fmt.Println(psqlconn)
    db, err := sql.Open("postgres", psqlconn)
    if err != nil {
        panic(err)
    }
    store := Store{
        Connection: db,
    }
    return store
}
