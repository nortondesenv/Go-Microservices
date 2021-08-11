package main

import (
	"context"
	"log"

	"github.com/nortondesenv/Go-Microservice/internal/server"
	"github.com/nortondesenv/Go-Microservice/pkg/jaeger"
	"github.com/nortondesenv/Go-Microservice/pkg/kafka"
	"github.com/nortondesenv/Go-Microservice/pkg/logger"
	"github.com/nortondesenv/Go-Microservice/pkg/mongodb"
	"github.com/nortondesenv/Go-Microservice/pkg/redis"
	"github.com/opentracing/opentracing-go"

	"github.com/nortondesenv/Go-Microservice/config"
)

// @title Products microservice
// @version 1.0
// @description Products REST API
// @termsOfService http://swagger.io/terms/
// @contact.name Norton Victor
// @contact.url https://github.com/nortondesenv
// @contact.email nortonvictor1999@gmail.com
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:5007
// @BasePath /api/v1
func main() {
	log.Println("Starting service...")
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg, err := config.ParseConfig()
	if err != nil {
		log.Fatal(err)
	}

	appLogger := logger.NewApiLogger(cfg)
	appLogger.InitLogger()
	appLogger.Info("Starting user server")
	appLogger.Infof(
		"AppVersion: %s, LogLevel: %s, DevelopmentMode: %s",
		cfg.AppVersion,
		cfg.Logger.Level,
		cfg.Server.Development,
	)
	appLogger.Infof("Success parsed config: %#v", cfg.AppVersion)

	mongoDBConn, err := mongodb.NewMongoDBConn(ctx, cfg)
	if err != nil {
		appLogger.Fatal("cannot connect mongodb", err)
	}
	defer func() {
		if err := mongoDBConn.Disconnect(ctx); err != nil {
			appLogger.Fatal("mongoDBConn.Disconnect", err)
		}
	}()
	appLogger.Infof("MongoDB connected: %v", mongoDBConn.NumberSessionsInProgress())

	conn, err := kafka.NewKafkaConn(cfg)
	if err != nil {
		appLogger.Fatal("NewKafkaConn", err)
	}
	defer conn.Close()
	brokers, err := conn.Brokers()
	if err != nil {
		appLogger.Fatal("conn.Brokers", err)
	}
	appLogger.Infof("Kafka connected: %v", brokers)

	tracer, closer, err := jaeger.InitJaeger(cfg)
	if err != nil {
		appLogger.Fatal("cannot create tracer", err)
	}
	appLogger.Info("Jaeger connected")

	opentracing.SetGlobalTracer(tracer)
	defer closer.Close()
	appLogger.Info("Opentracing connected")

	redisClient := redis.NewRedisClient(cfg)
	appLogger.Info("Redis connected")

	s := server.NewServer(appLogger, cfg, tracer, mongoDBConn, redisClient)
	appLogger.Fatal(s.Run())
}
