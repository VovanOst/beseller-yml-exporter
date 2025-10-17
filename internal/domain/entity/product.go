package entity

import "errors"

var (
	// Category errors
	ErrInvalidCategoryID   = errors.New("category ID is required")
	ErrInvalidCategoryName = errors.New("category name is required")

	// Product errors
	ErrInvalidProductID       = errors.New("product ID is required")
	ErrInvalidProductName     = errors.New("product name is required")
	ErrInvalidProductPrice    = errors.New("product price must be positive")
	ErrInvalidProductCategory = errors.New("product must belong to a category")
	ErrInvalidCurrency        = errors.New("product currency is required")
)

// Image представляет изображение товара
type Image struct {
	URL string // URL изображения
}

// Product представляет товар
type Product struct {
	ID          string  // Уникальный идентификатор товара
	Name        string  // Название товара
	StatusID    int     // Статус товара (1 = новинка)
	CategoryID  string  // ID категории
	Price       float64 // Цена товара
	Currency    string  // Валюта (BYN, USD, RUB и т.д.)
	URL         string  // URL страницы товара
	Images      []Image // Изображения товара
	Vendor      *string // Производитель/бренд
	Barcode     *string // Штрих-код
	Description *string // Описание товара
	Available   bool    // Доступен ли товар для заказа
}

// IsNew проверяет, является ли товар новинкой
func (p *Product) IsNew() bool {
	return p.StatusID == 1
}

// HasImages проверяет, есть ли у товара изображения
func (p *Product) HasImages() bool {
	return len(p.Images) > 0
}

// GetImageURLs возвращает список URL изображений
func (p *Product) GetImageURLs() []string {
	urls := make([]string, 0, len(p.Images))
	for _, img := range p.Images {
		if img.URL != "" {
			urls = append(urls, img.URL)
		}
	}
	return urls
}

// Validate проверяет валидность товара
func (p *Product) Validate() error {
	if p.ID == "" {
		return ErrInvalidProductID
	}
	if p.Name == "" {
		return ErrInvalidProductName
	}
	if p.Price < 0 {
		return ErrInvalidProductPrice
	}
	if p.CategoryID == "" {
		return ErrInvalidProductCategory
	}
	if p.Currency == "" {
		return ErrInvalidCurrency
	}
	return nil
}
