package app

import (
	"context"
	"fmt"
	"log"

	"github.com/digisata/todo-service/config"
	"github.com/digisata/todo-service/internal/handler"
	"github.com/digisata/todo-service/internal/repository"
	"github.com/digisata/todo-service/internal/usecase"
	"github.com/digisata/todo-service/pkg/grpcserver"
	"github.com/digisata/todo-service/pkg/interceptor"
	"github.com/digisata/todo-service/pkg/postgres"
	activityPB "github.com/digisata/todo-service/stubs/activity"
	taskPB "github.com/digisata/todo-service/stubs/task"
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
	taskRepository := repository.NewTask(pg)
	taskService := usecase.NewTask(taskRepository)
	taskHandler := handler.NewTask(taskService)

	activityRepository := repository.NewActivity(pg)
	activityService := usecase.NewActivity(activityRepository)
	activityCategoryHandler := handler.NewActivity(activityService)

	// Setup grpc server
	im := interceptor.NewInterceptorManager(sugar)
	grpcServer, err := grpcserver.NewGrpcServer(cfg.GrpcServer, sugar, im)
	if err != nil {
		panic(err)
	}
	defer grpcServer.Stop(ctx)

	taskPB.RegisterTaskServiceServer(grpcServer, taskHandler)
	activityPB.RegisterActivityServiceServer(grpcServer, activityCategoryHandler)
	grpc_health_v1.RegisterHealthServer(grpcServer.Server, health.NewServer())

	err = grpcServer.Run()
	if err != nil {
		panic(err)
	}

	<-ctx.Done()
}
