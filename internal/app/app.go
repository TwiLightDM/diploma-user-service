package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/TwiLightDM/diploma-user-service/internal/config"
	"github.com/TwiLightDM/diploma-user-service/internal/user-service"
	"github.com/TwiLightDM/diploma-user-service/package/databases"
	"github.com/TwiLightDM/diploma-user-service/package/utils"
	"github.com/TwiLightDM/diploma-user-service/proto/userservicepb"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config) error {
	db, err := databases.InitDB(
		cfg.DB.Host,
		cfg.DB.Port,
		cfg.DB.User,
		cfg.DB.Password,
		cfg.DB.Name,
	)
	if err != nil {
		return err
	}

	encryptionService := utils.NewEncryptionService(cfg.SaltLength)
	jwtService := utils.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessDuration,
		cfg.JWT.RefreshDuration,
	)

	userRepo := user_service.NewUserRepository(db)
	userService := user_service.NewUserService(
		userRepo,
		jwtService,
		encryptionService,
	)
	userHandler := user_service.NewUserHandler(userService)

	listener, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		return err
	}

	log.Printf("Starting user-service on %s", listener.Addr().String())

	grpcServer := grpc.NewServer()
	userservicepb.RegisterUserServiceServer(grpcServer, userHandler)

	ctx, stop := signal.NotifyContext(
		context.Background(),
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()

	go func() {
		if err = grpcServer.Serve(listener); err != nil {
			log.Printf("gRPC server stopped: %v", err)
		}
	}()

	<-ctx.Done()

	log.Println("Shutting down gRPC server...")

	grpcServer.GracefulStop()

	sqlDB, err := db.DB()
	if err == nil {
		log.Println("Closing database connection...")
		_ = sqlDB.Close()
	}

	log.Println("User-service stopped gracefully")
	return nil
}
