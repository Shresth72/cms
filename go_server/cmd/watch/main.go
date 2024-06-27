package main

import (
	"database/sql"
	"fmt"
	"os"

	"github.com/Shresth72/server/internals/data"
	"github.com/Shresth72/server/internals/utils"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
)

const (
	maxChunkSize = int64(5 << 20) // 5Mb
	maxReqBytes  = 1 << 20        // 1Mb
	maxRetries   = 3
)

type application struct {
	cfg struct {
		port int
		db   struct {
			conn string
		}
  }
	logger      *zerolog.Logger
	models      *data.Models
	oauth       *oauth2.Config
}

func main() {
  var app application

  // Logger 
  logger := zerolog.New(os.Stdout).With().Timestamp().Logger()
	app.logger = &logger
  
  err := godotenv.Load()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to load env variables")
		return
	}

	// Config
	app.cfg.port = utils.GetEnvInt("PORT", 8001)
	app.cfg.db.conn = os.Getenv("DATABASE_URL")

  db, err := connectDB(&app)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to initialize db connection")
	}
  defer db.Close()

	app.models = data.NewModels(db)

	// Server
	logger.Info().Msg(fmt.Sprintf("Server starting at port: %d", app.cfg.port))
  logger.Panic().Err(app.Serve())
}

// Connect to services
func connectDB(app *application) (*sql.DB, error) {
	db, err := sql.Open("postgres", app.cfg.db.conn)
	if err != nil {
		return nil, err
	}

	if err = db.Ping(); err != nil {
		return nil, err
	}

	return db, nil
}
