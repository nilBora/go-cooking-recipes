package main

import (
   "context"
   "log"
   "time"
   "github.com/jessevdk/go-flags"
   "fmt"
   server "go-cooking-recipes/v1/app/backend/server"
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
    var opts Options
    parser := flags.NewParser(&opts, flags.Default)
    _, err := parser.Parse()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("recipe %s\n", revision)

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
