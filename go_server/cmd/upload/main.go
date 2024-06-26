package main

import (
	"context"
	"database/sql"
	"fmt"
	"os"

	"github.com/IBM/sarama"
	"github.com/Shresth72/server/internals/data"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"github.com/rs/zerolog"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
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
		oauth struct {
			googleClientID     string
			googleClientSecret string
			redirectURL        string
		}
		kafka struct {
			brokers []string
			topic   string
      version string
		}
	}
	logger      *zerolog.Logger
	models      *data.Models
	redis       *redis.Client
	oauth       *oauth2.Config
	producer    sarama.SyncProducer
	logproducer sarama.AsyncProducer
	consumer    sarama.Consumer
}

func main() {
	var app application

	// Logger
	var logger = zerolog.New(os.Stdout).With().Timestamp().Logger()
	app.logger = &logger

	err := godotenv.Load()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to load env variables")
		return
	}

	// Config
	app.cfg.port = GetEnvInt("PORT", 8000)
	app.cfg.db.conn = os.Getenv("DATABASE_URL")
	app.cfg.oauth.googleClientID = os.Getenv("GOOGLE_CLIENT_ID")
	app.cfg.oauth.googleClientSecret = os.Getenv("GOOGLE_CLIENT_SECRET")
	app.cfg.oauth.redirectURL = os.Getenv("REDIRECT_URL")
	app.cfg.kafka.brokers = []string{os.Getenv("KAFKA_BROKER")}
	app.cfg.kafka.topic = os.Getenv("KAFKA_TOPIC")
  app.cfg.kafka.version = sarama.DefaultVersion.String()

	// DB
	db, err := connectDB(&app)
	if err != nil {
		logger.Panic().Err(err).Msg("failed to initialize db connection")
	}
	defer db.Close()

	app.models = data.NewModels(db)

	// Redis (store Refresh Tokens)
	// However, Generating New Token using RefreshToken not implemented
	app.redis, err = connectRedis()
	if err != nil {
		logger.Panic().Err(err).Msg("failed to initialize redis connection")
	}
	defer app.redis.Close()

	// OAuth2
	app.oauth = &oauth2.Config{
		ClientID:     app.cfg.oauth.googleClientID,
		ClientSecret: app.cfg.oauth.googleClientSecret,
		RedirectURL:  app.cfg.oauth.redirectURL,
		Endpoint:     google.Endpoint,
		Scopes: []string{"https://www.googleapis.com/auth/userinfo.profile",
			"https://www.googleapis.com/auth/userinfo.email"},
	}

  // Kafka
  version, err := sarama.ParseKafkaVersion(*&app.cfg.kafka.version)
  if err != nil {
		logger.Panic().Err(err).Msg("wrong kafka version")
  }

  app.producer = app.newDataCollector(app.cfg.kafka.brokers, version)
  defer app.producer.Close()
  // app.logproducer = app.newLogProducer(app.cfg.kafka.brokers, version)
  // defer app.logproducer.Close()

	// Server
	logger.Info().Msg(fmt.Sprintf("Server starting at port: %d", app.cfg.port))
	logger.Panic().Err(app.Serve())
}

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

func connectRedis() (*redis.Client, error) {
	conn := redis.NewClient(&redis.Options{
		Addr: os.Getenv("REDIS_URL"),
	})

	_, err := conn.Ping(context.Background()).Result()
	if err != nil {
		return nil, err
	}

	return conn, nil
}
