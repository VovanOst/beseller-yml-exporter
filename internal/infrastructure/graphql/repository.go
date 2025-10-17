package graphql

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"beseller-yml-exporter/internal/domain/entity"
	"beseller-yml-exporter/internal/domain/repository"
)

// PriceDTO представляет объект Price из GraphQL API
type PriceDTO struct {
	Name   string  `json:"name"`   // аббревиатура валюты, например "BYN"
	Value  float64 `json:"value"`  // значение цены
	Suffix string  `json:"suffix"` // отображаемый суффикс, например "руб."
}

type PageInfoDTO struct {
	URL string `json:"url"`
}

type CategoryDTO struct {
	ID             int          `json:"id"`
	Name           string       `json:"name"`
	Page           PageInfoDTO  `json:"page"`
	ParentCategory *CategoryDTO `json:"parentCategory"` // рекурсивная ссылка
}

// ImageDTO представляет изображение из PageImage
type ImageDTO struct {
	Image string `json:"image"` // имя файла изображения
}

// PageLinkDTO представляет ссылку на родительскую категорию
type PageLinkDTO struct {
	ParentURL string `json:"parentUrl"`
}

// PageDTO представляет страницу товара с slug и родительскими ссылками
type PageDTO struct {
	URL   string        `json:"url"`
	Links []PageLinkDTO `json:"links"`
}

// ProductDTO представляет товар из GraphQL API
type ProductDTO struct {
	ID          int         `json:"id"`
	Name        string      `json:"name"`
	StatusID    int         `json:"statusId"`
	Category    CategoryDTO `json:"category"`
	Price       float64     `json:"price"`
	PriceToShow []PriceDTO  `json:"priceToShow"` // массив объектов Price
	ItemCode    *string     `json:"itemCode"`
	VendorCode  *string     `json:"vendorCode"`
	Images      []ImageDTO  `json:"images"`
	Page        PageDTO     `json:"page"`
}

type CategoriesResponse struct {
	FilterCategory []CategoryDTO `json:"filterCategory"`
}

// ProductsResponse представляет ответ на запрос товаров
type ProductsResponse struct {
	FilterProduct []ProductDTO `json:"filterProduct"`
}

// CatalogRepository реализует repository.CatalogRepository через GraphQL
type CatalogRepository struct {
	client  *Client
	logger  Logger
	shopURL string
}

func NewCatalogRepository(c *Client, log Logger, shopURL string) repository.CatalogRepository {
	return &CatalogRepository{client: c, logger: log, shopURL: shopURL}
}

// getCurrencyFromPrice извлекает валюту из объекта Price
func (r *CatalogRepository) getCurrencyFromPrice(prices []PriceDTO) string {
	if len(prices) > 0 {
		// Берём аббревиатуру валюты из первого объекта Price
		return prices[0].Name
	}
	// Fallback валюта, если PriceToShow пустой
	return "BYN"
}

func buildCategoryPath(cat *CategoryDTO) []string {
	if cat == nil {
		return nil
	}

	// Сначала получаем путь родителя (от корня к текущей)
	segments := buildCategoryPath(cat.ParentCategory)

	// Добавляем текущую категорию
	if seg := strings.Trim(cat.Page.URL, "/"); seg != "" {
		segments = append(segments, seg)
	}

	return segments
}

func (r *CatalogRepository) GetCategories(ctx context.Context) ([]entity.Category, error) {
	var resp CategoriesResponse

	if err := r.client.Query(ctx, QueryFilterCategories, nil, &resp); err != nil {
		return nil, fmt.Errorf("failed to query categories: %w", err)
	}

	categories := make([]entity.Category, 0, len(resp.FilterCategory))
	for _, dto := range resp.FilterCategory {
		categories = append(categories, entity.Category{
			ID:   strconv.Itoa(dto.ID),
			Name: dto.Name,
		})
	}

	r.logger.Debug(fmt.Sprintf("Fetched %d categories", len(categories)))
	return categories, nil
}

func (r *CatalogRepository) GetProductsByStatus(ctx context.Context, statusID int) ([]entity.Product, error) {
	var resp ProductsResponse
	vars := map[string]interface{}{
		"first":  100,
		"offset": 0,
		"filter": map[string]interface{}{"statusId": statusID},
	}

	if err := r.client.Query(ctx, QueryFilterProduct, vars, &resp); err != nil {
		return nil, fmt.Errorf("failed to query products: %w", err)
	}

	r.logger.Debug(fmt.Sprintf("Received %d products from API, filtering by statusId=%d", len(resp.FilterProduct), statusID))

	products := make([]entity.Product, 0, len(resp.FilterProduct))
	for _, dto := range resp.FilterProduct {
		// Маппинг изображений с полным URL
		imgs := make([]entity.Image, 0, len(dto.Images))
		for _, img := range dto.Images {
			if img.Image != "" {
				// Формируем полный URL изображения
				fullImageURL := fmt.Sprintf("%s/pics/items/%s",
					strings.TrimRight(r.shopURL, "/"), img.Image)
				imgs = append(imgs, entity.Image{URL: fullImageURL})
			}
		}

		// Собираем путь из иерархии категорий
		categorySegments := buildCategoryPath(&dto.Category)

		// Добавляем slug товара
		allSegments := append(categorySegments, strings.Trim(dto.Page.URL, "/"))

		// Формируем полный URL
		path := "/" + strings.Join(allSegments, "/") + "/"
		fullPageURL := strings.TrimRight(r.shopURL, "/") + path

		// Извлекаем валюту из priceToShow
		currency := r.getCurrencyFromPrice(dto.PriceToShow)

		prod := entity.Product{
			ID:         strconv.Itoa(dto.ID),
			Name:       dto.Name,
			StatusID:   dto.StatusID,
			CategoryID: strconv.Itoa(dto.Category.ID),
			Price:      dto.Price,
			Currency:   currency, // валюта из priceToShow[0].name
			URL:        fullPageURL,
			Images:     imgs,
		}

		// Опциональные поля
		if dto.VendorCode != nil {
			prod.Vendor = dto.VendorCode
		}
		if dto.ItemCode != nil {
			prod.Barcode = dto.ItemCode
		}

		products = append(products, prod)
	}

	r.logger.Debug(fmt.Sprintf("Mapped %d products with currency from API", len(products)))
	return products, nil
}
