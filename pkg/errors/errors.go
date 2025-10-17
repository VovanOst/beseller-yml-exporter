package errors

import "errors"

var (
	// Repository errors
	ErrCategoryNotFound = errors.New("category not found")
	ErrProductNotFound  = errors.New("product not found")
	ErrNetworkTimeout   = errors.New("network timeout")
	ErrInvalidResponse  = errors.New("invalid response from API")

	// Validation errors
	ErrEmptyData = errors.New("empty data received")
)
