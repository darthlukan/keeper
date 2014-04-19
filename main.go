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
	goweb.Map("/user/{username}", usernameHandler)
	goweb.Map("/user/{username}/contacts", contactsHandler)
	goweb.Map("/user/{username}/contact/{contactname}", contactHandler)
	goweb.Map("/user/{username}/preferences", preferencesHandler)
	// goweb.Map("/user/{username}/preference/{prefname}", preferenceHandler)
	goweb.Map("/user/{username}/notifications", notificationsHandler)
	// goweb.Map("/user/{username}/notification/{id}", notificationHandler)
	goweb.Map("/user/{username}/location", locationHandler)
	goweb.Map("/user/{username}/location/{lat}/{lon}", locationHandler)

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

// Route Hanlders
func rootHandler(ctx context.Context) error {
	// TODO: Do something meaningful
	m := objx.MSI("status", 200, "message", "Success!")
	return goweb.API.RespondWithData(ctx, m)
}

func usernameHandler(ctx context.Context) error {
	username := ctx.PathValue("username")
	m := objx.MSI("status", 200, "username", username)
	return goweb.API.RespondWithData(ctx, m)
}

func contactsHandler(ctx context.Context) error {
	username := ctx.PathValue("username")
	m := objx.MSI("status", 200, "username", username, "contacts", []string{"foo", "bar", "baz"})
	return goweb.API.RespondWithData(ctx, m)
}

func contactHandler(ctx context.Context) error {
	username := ctx.PathValue("username")
	contactname := ctx.PathValue("contactname")
	m := objx.MSI("status", 200, "username", username, "contact", contactname)
	return goweb.API.RespondWithData(ctx, m)
}

func preferencesHandler(ctx context.Context) error {
	username := ctx.PathValue("username")
	m := objx.MSI("status", 200, "username", username, "preferences", []string{"foo", "bar", "baz"})
	return goweb.API.RespondWithData(ctx, m)
}

func notificationsHandler(ctx context.Context) error {
	username := ctx.PathValue("username")
	notifications := []string{"foo", "bar", "baz"}
	m := objx.MSI("status", 200, "username", username, "notifications", notifications)
	return goweb.API.RespondWithData(ctx, m)
}

func locationHandler(ctx context.Context) error {
	username := ctx.PathValue("username")
	if ctx.HttpRequest().Method == "POST" {
		lat := ctx.PathValue("lat")
		lon := ctx.PathValue("lon")

		type message struct {
			Content  string
			Lat      string
			Lon      string
			Username string
		}
		msg := new(message)
		msg.Content = "Location Set!"
		msg.Username = username
		msg.Lat = lat
		msg.Lon = lon

		return goweb.API.RespondWithData(ctx, msg)
	}

	m := objx.MSI("status", 200, "username", username, "location", []string{"lat", "lon"})
	return goweb.API.RespondWithData(ctx, m)
}

// Logging
type LoggedHandler struct {
	baseHandler http.Handler
}

func (handler *LoggedHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	go log.Printf("%s Request for: %s, From: %s", r.Method, r.RequestURI, r.RemoteAddr)
	handler.baseHandler.ServeHTTP(w, r)
}
