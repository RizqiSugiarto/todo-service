package app

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/digisata/invitation-service/config"
	"github.com/digisata/invitation-service/internal/handler"
	"github.com/digisata/invitation-service/internal/repository"
	"github.com/digisata/invitation-service/internal/usecase"
	"github.com/digisata/invitation-service/pkg/grpcserver"
	"github.com/digisata/invitation-service/pkg/interceptor"
	"github.com/digisata/invitation-service/pkg/postgres"
	"github.com/digisata/invitation-service/pkg/rabbitclient"
	invitationPB "github.com/digisata/invitation-service/stubs/invitation"
	invitationCategoryPB "github.com/digisata/invitation-service/stubs/invitation-category"
	invitationLabelPB "github.com/digisata/invitation-service/stubs/invitation-label"
	"go.uber.org/zap"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
)

func Run(cfg *config.Config) {
	ctx := context.Background()

	// Setup DB
	dbUsername := cfg.Postgres.DbUser
	dbPassword := cfg.Postgres.DbPass
	dbHost := cfg.Postgres.DbHost
	dbPort := cfg.Postgres.DbPort
	dbName := cfg.Postgres.DbName

	logger, _ := zap.NewProduction()
	defer logger.Sync()

	sugar := logger.Sugar()

	postgresUrl := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", dbUsername, dbPassword, dbHost, dbPort, dbName)

	pg, err := postgres.New(postgresUrl)
	if err != nil {
		log.Fatalf("app - run - maria.New: %v", err.Error())
	}
	defer pg.Close()

	err = RunMigrate(pg.Db)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err.Error())
	}

	// Dependencies injection
	invitationRepository := repository.NewInvitation(pg)
	invitationService := usecase.NewInvitation(cfg, invitationRepository)
	invitationHandler := handler.NewInvitation(invitationService)

	invitationLabelRepository := repository.NewInvitationLabel(pg)
	invitationLabelService := usecase.NewInvitationLabel(invitationLabelRepository)
	invitationLabelHandler := handler.NewInvitationLabel(invitationLabelService)

	invitationCategoryRepository := repository.NewInvitationCategory(pg)
	invitationCategoryService := usecase.NewInvitationCategory(invitationCategoryRepository)
	invitationCategoryHandler := handler.NewInvitationCategory(invitationCategoryService)

	// Setup RabbitMQ
	rabbitMQ, err := rabbitclient.New(cfg.RabbitMQ)
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer rabbitMQ.Conn.Close()

	err = rabbitMQ.NewChannel()
	if err != nil {
		log.Fatalf(err.Error())
	}
	defer rabbitMQ.Ch.Close()

	invitationEventHandler := handler.NewInvitationEvent(invitationService, rabbitMQ, cfg.RabbitMQ.ConsumerCount)

	// Declare queues
	for _, queueName := range strings.Split(cfg.RabbitMQ.Queues, ", ") {
		err := rabbitMQ.DeclareQueue(queueName)
		if err != nil {
			log.Fatalf(err.Error())
		}
	}

	// Start consume messages
	err = invitationEventHandler.StartConsumeInvitationOpenEvent()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = invitationEventHandler.StartConsumeInvitationComingEvent()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = invitationEventHandler.StartConsumeInvitationSendMoneyEvent()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = invitationEventHandler.StartConsumeInvitationSendGiftEvent()
	if err != nil {
		log.Fatalf(err.Error())
	}

	err = invitationEventHandler.StartConsumeInvitationCheckInEvent()
	if err != nil {
		log.Fatalf(err.Error())
	}

	// Setup grpc server
	im := interceptor.NewInterceptorManager(sugar)
	grpcServer, err := grpcserver.NewGrpcServer(cfg.GrpcServer, sugar, im)
	if err != nil {
		panic(err)
	}
	defer grpcServer.Stop(ctx)

	invitationPB.RegisterInvitationServiceServer(grpcServer, invitationHandler)
	invitationLabelPB.RegisterInvitationLabelServiceServer(grpcServer, invitationLabelHandler)
	invitationCategoryPB.RegisterInvitationCategoryServiceServer(grpcServer, invitationCategoryHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer.Server, health.NewServer())

	err = grpcServer.Run()
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
}
