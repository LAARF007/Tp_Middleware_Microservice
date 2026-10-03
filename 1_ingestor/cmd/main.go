package main

import (
	"context"
	"errors"
	"middleware/ingestor/internal/controllers/articles"
	"middleware/ingestor/internal/controllers/sources"
	"middleware/ingestor/internal/helpers"
	_ "middleware/ingestor/internal/models"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/sirupsen/logrus"
)

func main() {
	port := helpers.GetEnv("INGESTOR_PORT", "8080")

	logrus.SetLevel(logrus.DebugLevel)

	// SIGTERM + SIGINT are caught, which allows us to close everything
	ctx, stopListeningForSignals := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stopListeningForSignals()

	server := startServer(port)

	<-ctx.Done()
	logrus.Info("stopping : finishing the requests in progress, then closing the database")
	stopServer(server)
	helpers.CloseDB()
}

// startServer
// Starts the HTTP server in the background and hands it back, so that main can stop it.
func startServer(port string) *http.Server {
	r := chi.NewRouter()

	r.Route("/sources", func(r chi.Router) { // route /sources
		r.Get("/", sources.GetSources)        // GET /sources
		r.Post("/", sources.PostSource)       // POST /sources
		r.Route("/{id}", func(r chi.Router) { // route /sources/{id}
			r.Use(sources.Context)        // Use Context method to get source ID
			r.Get("/", sources.GetSource) // GET /sources/{id}
		})
	})

	r.Route("/articles", func(r chi.Router) { // route /articles
		r.Get("/", articles.GetArticles)      // GET /articles?source_id=...&limit=...
		r.Route("/{id}", func(r chi.Router) { // route /articles/{id}
			r.Use(articles.Context)         // Use Context method to get article ID
			r.Get("/", articles.GetArticle) // GET /articles/{id}
		})
	})

	server := &http.Server{Addr: ":" + port, Handler: r}

	// This is a go routine (go func(){})
	// It's a program running in parallel, allowing us to continue and not be blocked at ListenAndServer
	go func() {
		// ListenAndServe always ends with an error; the one Shutdown causes is expected
		err := server.ListenAndServe()
		if err != nil && !errors.Is(err, http.ErrServerClosed) {
			logrus.Fatalf("web server stopped : %s", err.Error())
		}
	}()

	logrus.Infof("Web server started. Now listening on *:%s", port)
	return server
}

// stopServer
// Stops server properly
func stopServer(server *http.Server) {
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	err := server.Shutdown(ctx)
	if err != nil {
		logrus.Errorf("error while stopping the web server : %s", err.Error())
	}
}

// In Go, the init function is implicitly called when the package is called, at the very start
// We use it as our program initialization
func init() {
	logrus.Debugf("init db if needed")
	databasePath := helpers.GetEnv("INGESTOR_DB_PATH", "ingestor.db")

	err := helpers.OpenDB(databasePath)
	if err != nil {
		logrus.Fatalf("error while opening database : %s", err.Error())
	}

	// TODO : add your db models here
	schemes := []string{
		`CREATE TABLE IF NOT EXISTS sources (
			id              VARCHAR(255) PRIMARY KEY NOT NULL,
			name            VARCHAR(255) NOT NULL,
			url             TEXT NOT NULL
		);`,
		`CREATE TABLE IF NOT EXISTS articles (
			id            VARCHAR(255) NOT NULL PRIMARY KEY,
			canonical_url TEXT NOT NULL,
			source_id     VARCHAR(255) NOT NULL,
    		url           TEXT NOT NULL,
			first_seen_at DATETIME NOT NULL,
			FOREIGN KEY (source_id) REFERENCES sources(id)
		);`,
	}

	for _, scheme := range schemes {
		if _, err := helpers.DB.Exec(scheme); err != nil {
			logrus.Fatalln("Could not generate table ! Error was : " + err.Error())
		}
	}
}
