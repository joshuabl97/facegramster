package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/go-chi/chi/v5"
	_ "github.com/jackc/pgx/v4/stdlib"
	"github.com/rs/zerolog"

	"github.com/joshuabl97/facegramster/controllers"
	"github.com/joshuabl97/facegramster/models"
	"github.com/joshuabl97/facegramster/ui"
	"github.com/joshuabl97/facegramster/ui/templates"
)

func main() {
	debug := flag.Bool("debug", false, "sets log level to debug")
	timeZone := flag.String("timezone", "Etc/Greenwich", "An official TZ identifier")
	portNum := flag.String("port_number", "3000", "The port number the server runs on")

	flag.Parse()

	// Create logger
	l := zerolog.New(os.Stderr).With().Timestamp().Logger()
	// Setting timezone
	loc, err := time.LoadLocation(*timeZone)
	if err != nil {
		l.Error().Msg("Couldn't determine timezone, using local machine time")
	} else if err == nil {
		time.Local = loc
	}

	// Make the logs look pretty
	l = l.Output(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339})

	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	if *debug {
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
	}
	l.Debug().Msg("This message appears only when log level set to Debug")

	// open the db
	cfg := models.DefaultPostgresConfig()
	db, err := models.Open(&l, cfg)
	if err != nil {
		l.Fatal().Err(err).Msg("failed to open db...")
	}
	defer db.Close()

	ui := ui.New(&l)

	r := chi.NewRouter()

	r.Use(requestLogger(l))

	r.Get("/", controllers.StaticHandler(ui.Must(ui.ParseFS(templates.FS,
		"home.html.tmpl", "default-wrapper.html.tmpl"))))

	r.Get("/contact", controllers.StaticHandler(ui.Must(ui.ParseFS(templates.FS,
		"contact.html.tmpl", "default-wrapper.html.tmpl"))))

	r.Get("/faq", controllers.FAQ(ui.Must(ui.ParseFS(templates.FS,
		"faq.html.tmpl", "default-wrapper.html.tmpl"))))

	usersC := controllers.Users{
		UserService: &models.UserService{
			DB: db,
			Lg: &l,
		},
	}
	usersC.Templates.New = ui.Must(ui.ParseFS(templates.FS, "signup.html.tmpl", "default-wrapper.html.tmpl"))
	usersC.Templates.SignIn = ui.Must(ui.ParseFS(templates.FS, "signin.html.tmpl", "default-wrapper.html.tmpl"))

	r.Get("/signup", usersC.New)
	r.Get("/signin", usersC.SignIn)
	r.Post("/signin", usersC.ProcessSignIn)
	r.Post("/users", usersC.Create)

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Page not found", http.StatusNotFound)
		// TODO: find a 404 template that looks swagalicious to embed
	})

	// Creates a new http server
	errorLog := &zerologLogger{l}
	s := http.Server{
		Addr:         ":" + *portNum,           // configure the bind address
		Handler:      r,                        // set the default handler
		IdleTimeout:  120 * time.Second,        // max duration to wait for the next request when keep-alives are enabled
		ReadTimeout:  5 * time.Second,          // max duration for reading the request
		WriteTimeout: 10 * time.Second,         // max duration before returning the request
		ErrorLog:     log.New(errorLog, "", 0), // set the logger for the server
	}

	// this go function starts the server
	// when the function is done running, that means we need to shutdown the server
	// we can do this by killing the program, but if there are requests being processed
	// we want to give them time to complete
	l.Info().Msg(fmt.Sprintf("Starting server on port %v....", *portNum))
	go func() {
		err := s.ListenAndServe()
		if err != nil {
			l.Fatal().Err(err)
		}
	}()

	// Sending kill and interrupt signals to os.Signal channel
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)

	// Does not invoke 'graceful shutdown' unless the signalChannel is closed
	<-sigChan

	l.Info().Msg("Received terminate, graceful shutdown")

	// This timeoutContext allows the server 30 seconds to complete all requests (if any) before shutting down
	timeoutCtx, _ := context.WithTimeout(context.Background(), 30*time.Second)
	err = s.Shutdown(timeoutCtx)
	if err != nil {
		l.Fatal().Err(err).Msg("Shutdown exceeded timeout")
		os.Exit(1)
	}
}

type zerologLogger struct {
	logger zerolog.Logger
}

// implement the io.Writer interface using my fancy custom logger for the server error log
func (l *zerologLogger) Write(p []byte) (n int, err error) {
	l.logger.Error().Msg(string(p))
	return len(p), nil
}

// Middleware to log request method and URI using ZeroLog
func requestLogger(logger zerolog.Logger) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			logger.Info().Str("Method", r.Method).Str("URI", r.RequestURI).Msg("Request details")
			next.ServeHTTP(w, r)
		})
	}
}
