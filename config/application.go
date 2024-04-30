package config

import (
	"github.com/go-redis/redis/v8"
	"go.uber.org/zap"
	"log"
)

// for ease of development we will switch to a closure pattern of dependency injection
// The reason I made this decision is because of two reasons
//		- Using method based dependency injection will not work because my handlers will be split in many packages.
//		I do not want to push in all of my handlers in one file only
//
//		- Using an external library introduces 'magic' which I do not want.

type Config struct {
	BaseURL  string
	HttpPort int
	Jwt      struct {
		SecretKey string
	}
}

type Application struct {
	Config     Config
	Logger     *zap.SugaredLogger
	MainLogger *zap.Logger
	StdLogger  *log.Logger
	Rdb        *redis.Client
	PubSub     *redis.PubSub
}
