package repository

import (
	"context"
	"beseller-yml-exporter/internal/domain/entity"
)

// CatalogRepository определяет интерфейс для работы с каталогом товаров
type CatalogRepository interface {
	// GetCategories возвращает все категории из каталога
	GetCategories(ctx context.Context) ([]entity.Category, error)

	// GetProductsByStatus возвращает товары с указанным статусом
	// statusID: 1 - новинка, 2 - хит продаж, и т.д.
	GetProductsByStatus(ctx context.Context, statusID int) ([]entity.Product, error)
}
