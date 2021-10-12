package configs

import (
	"fmt"
	"github.com/VSKrivoshein/test/internal/app/api/grpc_api"
	api "github.com/VSKrivoshein/test/internal/app/api/grpc_api/proto"
	"github.com/VSKrivoshein/test/internal/app/service"
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"net"
	"os"
	"os/signal"
	"syscall"
)

func StartGrpc(service service.User) error {
	grpcPort := fmt.Sprintf(":%s", os.Getenv("GRPC_PORT"))
	log.Infof("Starting grpc server on port %s", grpcPort)
	listen, err := net.Listen("tcp", grpcPort)
	if err != nil {
		log.Fatalf("Lister start fatal: %v", err)
	}

	s := grpc.NewServer(
		grpc.ChainUnaryInterceptor(grpc_api.Err, grpc_api.Logger),
	)

	srv := grpc_api.NewGRPCServer(service)
	api.RegisterUserServer(s, srv)

	return s.Serve(listen)
}

func GracefulShutdown() error {
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT)
	return fmt.Errorf("%s", <-signals)
}
