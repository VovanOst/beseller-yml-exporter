package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"time"

	//"beseller-yml-exporter/internal/domain/repository"
	"beseller-yml-exporter/internal/infrastructure/config"
	"beseller-yml-exporter/internal/infrastructure/graphql"
	"beseller-yml-exporter/internal/infrastructure/yml"
	"beseller-yml-exporter/internal/logger"
	"beseller-yml-exporter/internal/usecase"
	"beseller-yml-exporter/internal/usecase/dto"
)

func main() {
	// Парсинг флагов командной строки
	cfg := parseFlags()

	// Инициализация логгера
	log := logger.New(cfg.LogLevel)
	log.Info("Starting BeSeller YML Exporter")

	// Инициализация инфраструктуры
	ctx := context.Background()

	// GraphQL клиент и репозиторий
	log.Info("Connecting to GraphQL endpoint")
	gqlClient := graphql.NewClient(cfg.GraphQLEndpoint, cfg.HTTPTimeout, log)
	catalogRepo := graphql.NewCatalogRepository(gqlClient, log, cfg.ShopURL)

	// YML writer
	ymlWriter := yml.NewWriter(log)

	// Инициализация use case
	exportUC := usecase.NewExportCatalogUseCase(catalogRepo, ymlWriter, log)

	// Подготовка запроса на экспорт
	req := dto.ExportRequest{
		OutputPath:  cfg.OutputPath,
		ShopName:    cfg.ShopName,
		ShopCompany: cfg.ShopCompany,
		ShopURL:     cfg.ShopURL,
		Currency:    cfg.Currency,
		StatusID:    cfg.StatusID,
	}

	// Выполнение экспорта
	if err := exportUC.Execute(ctx, req); err != nil {
		log.Error("Export failed", "error", err)
		os.Exit(1)
	}

	log.Info("Export completed successfully")
}

func parseFlags() *config.Config {
	cfg := &config.Config{}

	// Сначала загружаем из .env
	envCfg := config.LoadFromEnv()

	// Затем парсим флаги (они имеют приоритет над .env)
	flag.StringVar(&cfg.GraphQLEndpoint, "endpoint", envCfg.GraphQLEndpoint, "GraphQL endpoint URL with token")
	flag.StringVar(&cfg.OutputPath, "out", envCfg.OutputPath, "Output YML file path")
	flag.StringVar(&cfg.ShopName, "shop-name", envCfg.ShopName, "Shop name")
	flag.StringVar(&cfg.ShopCompany, "shop-company", envCfg.ShopCompany, "Company name")
	flag.StringVar(&cfg.ShopURL, "shop-url", envCfg.ShopURL, "Shop URL")
	flag.StringVar(&cfg.Currency, "currency", envCfg.Currency, "Currency code (e.g., BYN, USD, RUB)")
	flag.IntVar(&cfg.StatusID, "status-id", envCfg.StatusID, "Product status ID to filter (1 for new)")
	flag.DurationVar(&cfg.HTTPTimeout, "timeout", envCfg.HTTPTimeout, "HTTP request timeout")
	flag.StringVar(&cfg.LogLevel, "log-level", envCfg.LogLevel, "Log level (debug, info, warn, error)")

	flag.Parse()

	// Валидация обязательных параметров
	if cfg.GraphQLEndpoint == "" {
		fmt.Fprintln(os.Stderr, "Error: GraphQL endpoint is required (use --endpoint flag or GRAPHQL_ENDPOINT env variable)")
		os.Exit(2)
	}

	if cfg.OutputPath == "" {
		cfg.OutputPath = "export.yml"
	}

	if cfg.HTTPTimeout == 0 {
		cfg.HTTPTimeout = 30 * time.Second
	}

	if cfg.LogLevel == "" {
		cfg.LogLevel = "info"
	}

	return cfg
}
