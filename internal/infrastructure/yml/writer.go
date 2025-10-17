package yml

import (
	"encoding/xml"
	"fmt"
	"os"
	"time"

	"beseller-yml-exporter/internal/domain/entity"
)

// Logger интерфейс для логирования
type Logger interface {
	Debug(msg string, args ...interface{})
	Info(msg string, args ...interface{})
	Warn(msg string, args ...interface{})
	Error(msg string, args ...interface{})
}

// Writer реализует запись каталога в YML формат
type Writer struct {
	logger Logger
}

// NewWriter создаёт новый YML writer
func NewWriter(logger Logger) *Writer {
	return &Writer{
		logger: logger,
	}
}

// Write записывает каталог в YML файл
func (w *Writer) Write(
	outputPath string,
	shopName string,
	shopCompany string,
	shopURL string,
	currency string,
	categories []entity.Category,
	products []entity.Product,
) error {
	catalog := w.buildCatalog(shopName, shopCompany, shopURL, currency, categories, products)

	file, err := os.Create(outputPath)
	if err != nil {
		return fmt.Errorf("failed to create file: %w", err)
	}
	defer file.Close()

	// Запись XML заголовка
	if _, err := file.WriteString(xml.Header); err != nil {
		return fmt.Errorf("failed to write XML header: %w", err)
	}

	// Запись DOCTYPE
	if _, err := file.WriteString(`<!DOCTYPE yml_catalog SYSTEM "shops.dtd">` + "\n"); err != nil {
		return fmt.Errorf("failed to write DOCTYPE: %w", err)
	}

	// Кодирование XML
	encoder := xml.NewEncoder(file)
	encoder.Indent("", "  ")

	if err := encoder.Encode(catalog); err != nil {
		return fmt.Errorf("failed to encode XML: %w", err)
	}

	if err := encoder.Flush(); err != nil {
		return fmt.Errorf("failed to flush encoder: %w", err)
	}

	return nil
}

// buildCatalog создаёт структуру YML каталога
func (w *Writer) buildCatalog(
	shopName string,
	shopCompany string,
	shopURL string,
	currency string,
	categories []entity.Category,
	products []entity.Product,
) YMLCatalog {
	catalog := YMLCatalog{
		Date: time.Now().Format("2006-01-02 15:04"),
		Shop: Shop{
			Name:    shopName,
			Company: shopCompany,
			URL:     shopURL,
		},
	}

	// Валюты
	catalog.Shop.Currencies.Currency = []Currency{
		{
			ID:   currency,
			Rate: "1",
		},
	}

	// Категории
	catalog.Shop.Categories.Category = make([]Category, 0, len(categories))
	for _, cat := range categories {
		ymlCat := Category{
			ID:       cat.ID,
			ParentID: cat.ParentID,
			Name:     cat.Name,
		}
		catalog.Shop.Categories.Category = append(catalog.Shop.Categories.Category, ymlCat)
	}

	// Товары
	catalog.Shop.Offers.Offer = make([]Offer, 0, len(products))
	for _, prod := range products {
		available := "true"
		if !prod.Available {
			available = "false"
		}

		offer := Offer{
			ID:         prod.ID,
			Available:  available,
			URL:        prod.URL,
			Price:      prod.Price,
			CurrencyID: prod.Currency,
			CategoryID: prod.CategoryID,
			Name:       prod.Name,
		}

		// Картинки
		offer.Picture = prod.GetImageURLs()

		// Опциональные поля
		if prod.Vendor != nil && *prod.Vendor != "" {
			offer.Vendor = *prod.Vendor
		}
		if prod.Barcode != nil && *prod.Barcode != "" {
			offer.Barcode = *prod.Barcode
		}
		if prod.Description != nil && *prod.Description != "" {
			offer.Description = *prod.Description
		}

		catalog.Shop.Offers.Offer = append(catalog.Shop.Offers.Offer, offer)
	}

	return catalog
}
