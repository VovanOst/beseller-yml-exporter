package graphql

const (
	// QueryCategories - запрос для получения всех категорий
	QueryFilterCategories = `
		query FilterCategory {
			filterCategory {
				id
				name
			}
		}
	`

	// QueryProducts - запрос для получения всех товаров
	QueryFilterProduct = `
		query FilterProduct {
			filterProduct {
				id
				name
				statusId
				category {
                          id
                          name
                          page {
                            url
                           }
                          parentCategory {
                                      id
                                      name
                                      page {
                                        url
                                      }
                                      parentCategory {
                                           id
                                           name
                                           page {
                                              url
                                           }
                                        parentCategory {
                                            id
                                            name
                                            page {
                                                url
                                            }
                                        }
                                      }
                          }
                }
                 
				price
				priceToShow {     # объект Price, содержащий currency
                             name
                             value
                             suffix
                }		
                images {
					image
				}
                page {
                       url           # slug товара
                       links {
                           parentUrl         # путь каждой категории-родителя
                       }
                } 
				itemCode
				vendorCode				
			}
		}
	`
)
