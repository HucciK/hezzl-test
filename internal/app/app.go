package app

import (
	"fmt"
	"hezzl/config"
	"hezzl/internal/repository/cache"
	"hezzl/internal/repository/itemRepository"
	"hezzl/internal/repository/logsRepository"
	"hezzl/internal/service/itemService"
	"hezzl/internal/service/logsService"
	"hezzl/internal/transport/broker"
	transport "hezzl/internal/transport/http"
	"hezzl/internal/transport/http/handlers/itemHandler"
	"hezzl/pkg/clickhouse"
	"hezzl/pkg/nats"
	"hezzl/pkg/postgres"
	"hezzl/pkg/redis"
	"log"
	"net/http"
)

type App struct {
	config config.Config
	logger *log.Logger
}

func New(cfg config.Config, log *log.Logger) *App {
	return &App{
		config: cfg,
		logger: log,
	}
}

func (a App) Run() error {
	pg, err := postgres.NewPostgres(a.config.PostgresConfig)
	if err != nil {
		return fmt.Errorf("error while trying to create postgres instance: %w", err)
	}

	ch, err := clickhouse.NewClickhouse(a.config.ClickhouseConfig)
	if err != nil {
		return fmt.Errorf("error while trying to create clickhouse instance: %w", err)
	}
	fmt.Println(a.config.ClickhouseConfig.ConnectionString())

	redisCli, err := redis.NewRedisClient(a.config.RedisConfig)
	if err != nil {
		return fmt.Errorf("error while trying to create redis instance: %w", err)
	}

	cache := cache.NewRedisCache(redisCli, a.config.RedisConfig.TTL)

	natsBroker, err := nats.NewNatsBroker(a.config.NatsConfig)
	if err != nil {
		return fmt.Errorf("error while trying to create nats instance: %w", err)
	}

	ItemPostgres := itemRepository.NewItemPostgres(pg)
	LogsClickhouse := logsRepository.NewLogsClickhouse(ch)

	ItemService := itemService.NewItemService(ItemPostgres, natsBroker, cache)
	LogsService := logsService.NewLogsService(LogsClickhouse)

	router := http.NewServeMux()

	_, err = broker.NewBroker("items_update", natsBroker, LogsService, a.logger)
	if err != nil {
		return fmt.Errorf("error while trying to create broker transport instance: %w", err)
	}

	ItemHandler := itemHandler.NewItemHandler(ItemService, a.logger)
	ItemHandler.RegisterItemHandlers(router)

	server := transport.NewServer(a.config.ServerConfig, router)
	if err := server.Run(); err != nil {
		return err
	}

	return nil
}
