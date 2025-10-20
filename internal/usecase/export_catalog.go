package usecase

import (
	"context"
	"fmt"

	"beseller-yml-exporter/internal/domain/entity"
	"beseller-yml-exporter/internal/domain/repository"
	"beseller-yml-exporter/internal/usecase/dto"
)

// Logger определяет интерфейс для логирования
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// CatalogWriter определяет интерфейс для записи каталога в файл
type CatalogWriter interface {
	Write(
		outputPath string,
		shopName string,
		shopCompany string,
		shopURL string,
		currency string,
		categories []entity.Category,
		products []entity.Product,
	) error
}

// ExportCatalogUseCase реализует сценарий экспорта каталога в YML
type ExportCatalogUseCase struct {
	catalogRepo repository.CatalogRepository
	writer      CatalogWriter
	logger      Logger
}

// NewExportCatalogUseCase создаёт новый экземпляр use case
func NewExportCatalogUseCase(
	catalogRepo repository.CatalogRepository,
	writer CatalogWriter,
	logger Logger,
) *ExportCatalogUseCase {
	return &ExportCatalogUseCase{
		catalogRepo: catalogRepo,
		writer:      writer,
		logger:      logger,
	}
}

// Execute выполняет экспорт каталога
func (uc *ExportCatalogUseCase) Execute(ctx context.Context, req dto.ExportRequest) error {
	// Валидация запроса
	if err := req.Validate(); err != nil {
		return fmt.Errorf("invalid request: %w", err)
	}

	// 1. Получение категорий
	uc.logger.Info("Fetching categories...")
	categories, err := uc.catalogRepo.GetCategories(ctx)
	if err != nil {
		return fmt.Errorf("failed to fetch categories: %w", err)
	}
	uc.logger.Info(fmt.Sprintf("Found %d categories", len(categories)))

	// Валидация категорий
	validCategories := make([]entity.Category, 0, len(categories))
	for _, cat := range categories {
		if err := cat.Validate(); err != nil {
			uc.logger.Warn(fmt.Sprintf("Skipping invalid category %s: %v", cat.ID, err))
			continue
		}
		validCategories = append(validCategories, cat)
	}

	// 2. Получение товаров с нужным статусом
	uc.logger.Info(fmt.Sprintf("Fetching products with statusId=%d...", req.StatusID))
	products, err := uc.catalogRepo.GetProductsByStatus(ctx, req.StatusID)
	if err != nil {
		return fmt.Errorf("failed to fetch products: %w", err)
	}
	uc.logger.Info(fmt.Sprintf("Found %d products", len(products)))

	// Валидация и фильтрация товаров
	validProducts := make([]entity.Product, 0, len(products))
	for _, prod := range products {
		if err := prod.Validate(); err != nil {
			uc.logger.Warn(fmt.Sprintf("Skipping invalid product %s: %v", prod.ID, err))
			continue
		}
		// Дополнительная проверка статуса (на случай если API вернул лишнее)
		if !prod.IsNew() {
			uc.logger.Debug(fmt.Sprintf("Skipping product %s: statusId=%d", prod.ID, prod.StatusID))
			continue
		}
		validProducts = append(validProducts, prod)
	}

	if len(validProducts) == 0 {
		uc.logger.Warn("No valid products found for export")
	}

	// 3. Запись в YML файл
	uc.logger.Info("Generating YML file...")
	if err := uc.writer.Write(
		req.OutputPath,
		req.ShopName,
		req.ShopCompany,
		req.ShopURL,
		req.Currency,
		validCategories,
		validProducts,
	); err != nil {
		return fmt.Errorf("failed to write YML: %w", err)
	}

	uc.logger.Info(fmt.Sprintf("YML file created: %s", req.OutputPath))
	uc.logger.Info(fmt.Sprintf("Export completed (categories=%d, offers=%d)", len(validCategories), len(validProducts)))

	return nil
}
