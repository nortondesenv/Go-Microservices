package server

import (
	"context"
	"net"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-playground/validator/v10"
	grpcrecovery "github.com/grpc-ecosystem/go-grpc-middleware/recovery"
	grpc_ctxtags "github.com/grpc-ecosystem/go-grpc-middleware/tags"
	grpc_opentracing "github.com/grpc-ecosystem/go-grpc-middleware/tracing/opentracing"
	grpc_prometheus "github.com/grpc-ecosystem/go-grpc-prometheus"
	"github.com/opentracing/opentracing-go"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"

	"github.com/nortondesenv/Go-Microservice/config"
	product "github.com/nortondesenv/Go-Microservice/internal/product/delivery/grpc"
	"github.com/nortondesenv/Go-Microservice/internal/product/repository"
	"github.com/nortondesenv/Go-Microservice/internal/product/usecase"
	"github.com/nortondesenv/Go-Microservice/pkg/logger"
	productsService "github.com/nortondesenv/Go-Microservice/proto/product"
)

// server
type server struct {
	log    logger.Logger
	cfg    *config.Config
	tracer opentracing.Tracer
}

func NewServer(log logger.Logger, cfg *config.Config, tracer opentracing.Tracer) *server {
	return &server{log: log, cfg: cfg, tracer: tracer}
}

func (s *server) Run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	validate := validator.New()

	productMongoRepo := repository.NewProductMongoRepo()
	productUC := usecase.NewProductUC(productMongoRepo, s.log)

	l, err := net.Listen("tcp", s.cfg.Server.Port)
	if err != nil {
		return errors.Wrap(err, "net.Listen")
	}
	defer l.Close()

	grpcServer := grpc.NewServer(grpc.KeepaliveParams(keepalive.ServerParameters{
		MaxConnectionIdle: s.cfg.Server.MaxConnectionIdle * time.Minute,
		Timeout:           s.cfg.Server.Timeout * time.Second,
		MaxConnectionAge:  s.cfg.Server.MaxConnectionAge * time.Minute,
		Time:              s.cfg.Server.Timeout * time.Minute,
	}),
		grpc.ChainUnaryInterceptor(
			grpc_ctxtags.UnaryServerInterceptor(),
			grpc_opentracing.UnaryServerInterceptor(),
			grpc_prometheus.UnaryServerInterceptor,
			grpcrecovery.UnaryServerInterceptor(),
		),
	)

	productService := product.NewProductService(s.log, productUC, validate)
	productsService.RegisterProductsServiceServer(grpcServer, productService)
	grpc_prometheus.Register(grpcServer)

	go func() {
		s.log.Infof("GRPC Server is listening on port: %v", s.cfg.Server.Port)
		s.log.Fatal(grpcServer.Serve(l))
	}()

	if s.cfg.Server.Development {
		reflection.Register(grpcServer)
	}

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)

	select {
	case v := <-quit:
		s.log.Errorf("signal.Notify: %v", v)
	case done := <-ctx.Done():
		s.log.Errorf("ctx.Done: %v", done)
	}

	// if err := s.echo.Server.Shutdown(ctx); err != nil {
	// 	return errors.Wrap(err, "echo.Server.Shutdown")
	// }

	grpcServer.GracefulStop()
	s.log.Info("Server Exited Properly")

	return nil
}
