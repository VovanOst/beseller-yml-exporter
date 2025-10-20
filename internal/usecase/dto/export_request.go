package dto

// ExportRequest содержит параметры для экспорта каталога в YML
type ExportRequest struct {
	OutputPath  string // Путь к выходному YML файлу
	ShopName    string // Название магазина
	ShopCompany string // Название компании
	ShopURL     string // URL магазина
	Currency    string // Валюта магазина (BYN, USD, RUB и т.д.)
	StatusID    int    // ID статуса товаров для экспорта (1 = новинка)
}

// Validate проверяет валидность запроса
func (r *ExportRequest) Validate() error {
	if r.OutputPath == "" {
		return ErrInvalidOutputPath
	}
	if r.ShopName == "" {
		return ErrInvalidShopName
	}
	if r.ShopCompany == "" {
		return ErrInvalidShopCompany
	}
	if r.ShopURL == "" {
		return ErrInvalidShopURL
	}
	if r.Currency == "" {
		return ErrInvalidCurrency
	}
	if r.StatusID <= 0 {
		return ErrInvalidStatusID
	}
	return nil
}
