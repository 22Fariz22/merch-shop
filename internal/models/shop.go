package models

type Shop struct {
	Catalog map[string]int `gorm:"-"`
}

func NewShop() *Shop {
	return &Shop{
		Catalog: map[string]int{
			"t-shirt":    80,
			"cup":        20,
			"book":       50,
			"pen":        10,
			"powerbank":  200,
			"hoody":      300,
			"umbrella":   200,
			"socks":      10,
			"wallet":     50,
			"pink-hoody": 500,
		},
	}
}

func (s *Shop) GetPrice(itemName string) (int, bool) {
	price, exists := s.Catalog[itemName]
	return price, exists
}

var DefaultShop *Shop

func init() {
	DefaultShop = NewShop()
}
