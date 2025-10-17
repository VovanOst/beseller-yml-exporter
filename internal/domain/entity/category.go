package entity

// Category представляет категорию товаров
type Category struct {
	ID       string  // Уникальный идентификатор категории
	Name     string  // Название категории
	ParentID *string // ID родительской категории (nil для корневых категорий)
}

// IsRoot проверяет, является ли категория корневой
func (c *Category) IsRoot() bool {
	return c.ParentID == nil || *c.ParentID == ""
}

// Validate проверяет валидность категории
func (c *Category) Validate() error {
	if c.ID == "" {
		return ErrInvalidCategoryID
	}
	if c.Name == "" {
		return ErrInvalidCategoryName
	}
	return nil
}
