package app

import (
	"context"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/TwiLightDM/diploma-user-service/internal/config"
	"github.com/TwiLightDM/diploma-user-service/internal/group-member-service"
	"github.com/TwiLightDM/diploma-user-service/internal/group-service"
	"github.com/TwiLightDM/diploma-user-service/internal/user-service"
	"github.com/TwiLightDM/diploma-user-service/package/databases"
	"github.com/TwiLightDM/diploma-user-service/package/utils"
	"github.com/TwiLightDM/diploma-user-service/proto/groupmemberservicepb"
	"github.com/TwiLightDM/diploma-user-service/proto/groupservicepb"
	"github.com/TwiLightDM/diploma-user-service/proto/userservicepb"
	"google.golang.org/grpc"
)

func Run(cfg *config.Config) error {
	postgres, err := databases.InitDB(
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Name,
	)
	if err != nil {
		return err
	}

	redis, err := databases.InitRedis(
		cfg.Redis.Host,
		cfg.Redis.Port,
		cfg.Redis.Password,
		0,
	)
	if err != nil {
		return err
	}

	validationService := utils.NewValidationService()
	encryptionService := utils.NewEncryptionService(cfg.SaltLength)
	jwtService := utils.NewJWTService(
		cfg.JWT.Secret,
		cfg.JWT.AccessDuration,
		cfg.JWT.RefreshDuration,
	)

	userRepo := user_service.NewUserRepository(postgres, redis)
	userService := user_service.NewUserService(
		userRepo,
		validationService,
		jwtService,
		encryptionService,
	)
	userHandler := user_service.NewUserHandler(userService)

	groupRepo := group_service.NewGroupRepository(postgres, redis)
	groupService := group_service.NewGroupService(groupRepo)
	groupHandler := group_service.NewGroupHandler(groupService)

	groupMemberRepo := group_member_service.NewGroupMemberRepository(postgres)
	groupMemberService := group_member_service.NewGroupMemberService(groupMemberRepo, userRepo)
	groupMemberHandler := group_member_service.NewGroupMemberHandler(groupMemberService)

	listener, err := net.Listen("tcp", cfg.GRPCPort)
	if err != nil {
		return err
	}

	log.Printf("Starting user-service on %s", listener.Addr().String())

	grpcServer := grpc.NewServer()
	userservicepb.RegisterUserServiceServer(grpcServer, userHandler)
	groupservicepb.RegisterGroupServiceServer(grpcServer, groupHandler)
	groupmemberservicepb.RegisterGroupMemberServiceServer(grpcServer, groupMemberHandler)

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

	sqlDB, err := postgres.DB()
	if err == nil {
		log.Println("Closing database connection...")
		_ = sqlDB.Close()
	}

	if err = redis.Close(); err != nil {
		log.Println("Error while disconnecting redis:", err)
	} else {
		log.Println("Closing redis connection...")
	}

	log.Println("User-service stopped gracefully")
	return nil
}
