package main

import (
    "github.com/stretchr/goweb"
    "github.com/stretchr/goweb/context"
    "github.com/stretchr/objx"
    "log"
    "net"
    "net/http"
    "os"
    "path"
    "time"
)

var (
    projectRoot string
    goPath      = os.Getenv("GOPATH")
)

func main() {
    log.Println("Starting server...")

    // We want to be able to execute from anywhere
    if goPath == "" {
        projectRoot = "."
    } else {
        projectRoot = path.Join(goPath, "src", "github.com", "darthlukan", "keeper")
    }

    // Routes
    goweb.Map("/", rootHandler)

    // Server ENV
    address := ":3000"
    if port := os.Getenv("PORT"); port != "" {
        address = ":" + port
    }

    server := &http.Server{
        Addr:           address,
        Handler:        &LoggedHandler{goweb.DefaultHttpHandler()},
        ReadTimeout:    10 * time.Second,
        WriteTimeout:   10 * time.Second,
        MaxHeaderBytes: 1 << 20,
    }

    listener, listenErr := net.Listen("tcp", address)
    if listenErr != nil {
        log.Panicf("Could not listen for TCP on %s: %s", address, listenErr)
    }

    log.Println("Server loaded, check localhost" + address)
    // Lobbeth thy holy hand grenade
    server.Serve(listener)
}

// Logging
type LoggedHandler struct {
    baseHandler http.Handler
}

func (handler *LoggedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    go log.Printf("%s Request for: %s, From: %s", r.Method, r.RequestURI, r.RemoteAddr)
    handler.baseHandler.ServeHTTP(w, r)
}

// Route Hanlders
func rootHandler(ctx context.Context) error {
    // TODO: Do something meaningful
    m := objx.MSI("status", 200, "message", "Success!")
    return goweb.API.RespondWithData(ctx, m)
}
