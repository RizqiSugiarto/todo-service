package config

import (
	"fmt"
	"log"

	"github.com/digisata/invitation-service/pkg/grpcserver"
	"github.com/digisata/invitation-service/pkg/postgres"
	"github.com/digisata/invitation-service/pkg/rabbitclient"
	"github.com/spf13/viper"
)

type Config struct {
	AppEnv                   string              `mapstructure:"app_env"`
	WeddingInvitationBaseUrl string              `mapstructure:"wedding_invitation_base_url"`
	EncryptionKey            string              `mapstructure:"encryption_key"`
	Postgres                 postgres.Config     `mapstructure:"postgres"`
	GrpcServer               grpcserver.Config   `mapstructure:"grpc_server"`
	RabbitMQ                 rabbitclient.Config `mapstructure:"rabbitmq"`
}

func Load() (*Config, error) {
	var cfg Config

	viper.SetConfigFile("config.yaml")
	viper.AutomaticEnv()

	err := viper.ReadInConfig()
	if err != nil {
		return nil, fmt.Errorf("can't find the config file: %v", err)
	}

	err = viper.Unmarshal(&cfg)
	if err != nil {
		return nil, fmt.Errorf("environment can't be loaded: %v", err)
	}

	log.Printf("The App is running in %s environment", cfg.AppEnv)

	return &cfg, nil
}
