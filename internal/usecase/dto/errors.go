package dto

import "errors"

var (
	ErrInvalidOutputPath  = errors.New("output path is required")
	ErrInvalidShopName    = errors.New("shop name is required")
	ErrInvalidShopCompany = errors.New("shop company is required")
	ErrInvalidShopURL     = errors.New("shop URL is required")
	ErrInvalidCurrency    = errors.New("currency is required")
	ErrInvalidStatusID    = errors.New("status ID must be positive")
)
