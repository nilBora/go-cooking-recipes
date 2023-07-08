package main

import (
   "context"
   "log"
   "time"
   "github.com/jessevdk/go-flags"
   "fmt"
   server "go-cooking-recipes/v1/app/backend/server"
   store "go-cooking-recipes/v1/app/backend/store"
   "github.com/joho/godotenv"
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
    s := store.ConnectDB()

    srv := server.Server{
        Port:           opts.Port,
        PinSize:        opts.PinSize,
        MaxExpire:      opts.MaxExpire,
        MaxPinAttempts: opts.MaxPinAttempts,
        WebRoot:        opts.WebRoot,
        Version:        revision,
        Store:          s,
    }
    if err := srv.Run(context.Background()); err != nil {
        log.Printf("[ERROR] failed, %+v", err)
    }
}
