package yml

import "encoding/xml"

// YMLCatalog представляет корневой элемент YML
type YMLCatalog struct {
	XMLName xml.Name `xml:"yml_catalog"`
	Date    string   `xml:"date,attr"`
	Shop    Shop     `xml:"shop"`
}

// Shop представляет элемент shop
type Shop struct {
	Name       string     `xml:"name"`
	Company    string     `xml:"company"`
	URL        string     `xml:"url"`
	Currencies Currencies `xml:"currencies"`
	Categories Categories `xml:"categories"`
	Offers     Offers     `xml:"offers"`
}

// Currencies представляет список валют
type Currencies struct {
	Currency []Currency `xml:"currency"`
}

// Currency представляет валюту
type Currency struct {
	ID   string `xml:"id,attr"`
	Rate string `xml:"rate,attr"`
}

// Categories представляет список категорий
type Categories struct {
	Category []Category `xml:"category"`
}

// Category представляет категорию
type Category struct {
	ID       string  `xml:"id,attr"`
	ParentID *string `xml:"parentId,attr,omitempty"`
	Name     string  `xml:",chardata"`
}

// Offers представляет список предложений
type Offers struct {
	Offer []Offer `xml:"offer"`
}

// Offer представляет товарное предложение
type Offer struct {
	ID          string   `xml:"id,attr"`
	Available   string   `xml:"available,attr"`
	URL         string   `xml:"url,omitempty"`
	Price       float64  `xml:"price"`
	CurrencyID  string   `xml:"currencyId"`
	CategoryID  string   `xml:"categoryId"`
	Picture     []string `xml:"picture,omitempty"`
	Name        string   `xml:"name"`
	Vendor      string   `xml:"vendor,omitempty"`
	Barcode     string   `xml:"barcode,omitempty"`
	Description string   `xml:"description,omitempty"`
}
