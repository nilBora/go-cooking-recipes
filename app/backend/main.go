package main

import (
    "os"
   "context"
   "log"
   "time"
   "github.com/jessevdk/go-flags"
   "fmt"
   server "go-cooking-recipes/v1/app/backend/server"
   "github.com/joho/godotenv"
   "database/sql"
   _ "github.com/lib/pq"
)

type Server struct {
	PinSize        int
	MaxPinAttempts int
	WebRoot        string
	Version        string
	Port           string
}

type Options struct {
    Port string `short:"p" long:"port" env:"SERVER_PORT" default:"8080" description:"Port web server"`
    PinSize int `long:"pinszie" env:"PIN_SIZE" default:"5" description:"pin size"`
    MaxExpire time.Duration `long:"expire" env:"MAX_EXPIRE" default:"24h" description:"max lifetime"`
    MaxPinAttempts int `long:"pinattempts" env:"PIN_ATTEMPTS" default:"3" description:"max attempts to enter pin"`
    WebRoot string `long:"web" env:"WEB" default:"/" description:"web ui location"`
}

var revision string

func main() {
    if err := godotenv.Load(); err != nil {
        panic("No .env file found")
    }

    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("recipe %s\n", revision)

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

    rows, err := db.Query(`SELECT id, uuid FROM recipes`)
    if err != nil {
        panic(err)
    }
    defer rows.Close()
    for rows.Next() {
        var id int
        var uuid string

        err = rows.Scan(&id, &uuid)
        if err != nil {
            panic(err)
        }

        fmt.Println(id, uuid)
    }

    srv := server.Server{
        Port:           opts.Port,
        PinSize:        opts.PinSize,
        MaxExpire:      opts.MaxExpire,
        MaxPinAttempts: opts.MaxPinAttempts,
        WebRoot:        opts.WebRoot,
        Version:        revision,
    }
    if err := srv.Run(context.Background()); err != nil {
        log.Printf("[ERROR] failed, %+v", err)
    }
}
