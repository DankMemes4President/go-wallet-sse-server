package main

import (
	"flag"
	"github.com/go-redis/redis/v8"
	"go-wallet-sse-server/config"
	"go.uber.org/zap"
	"log"
)

func main() {

	// initialize logger
	mainLogger, _ := zap.NewDevelopment()
	logger := mainLogger.Sugar()
	stdLogger := zap.NewStdLog(mainLogger)
	err := run(logger, stdLogger, mainLogger)
	if err != nil {
		logger.Fatal(err.Error())
	}
}

func run(logger *zap.SugaredLogger, stdLogger *log.Logger, mainLogger *zap.Logger) error {
	var cfg config.Config
	flag.StringVar(&cfg.BaseURL, "base-url", "http://localhost:4444", "base URL for the sse server")
	flag.IntVar(&cfg.HttpPort, "http-port", 4445, "port to listen on for HTTP requests")
	flag.StringVar(&cfg.Jwt.SecretKey, "jwt-secret-key", "rbztegymvi2bxjdh2tftkvd7b44z5akg", "secret key for JWT authentication")

	flag.Parse()
	rdb := redis.NewClient(&redis.Options{
		Addr: "wallet-redis:6379",
	})
	// Initialize PubSub
	pubSub := rdb.Subscribe(rdb.Context(), "broadcast")
	app := &config.Application{
		Config:     cfg,
		Logger:     logger,
		MainLogger: mainLogger,
		StdLogger:  stdLogger,
		Rdb:        rdb,
		PubSub:     pubSub,
	}
	return serveHTTP(app)
}
