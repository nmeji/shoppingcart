package amaysim

type Product struct {
	Code, Name string
	Price      float32
}

type CartItem struct {
	Quantity int
	SubTotal float32
	Product
}

type Cart interface {
	Add(item CartItem)
	Remove(item CartItem)
	AddPromoCode(code string)
	UsesPromoCode(code string) bool
	Total() float32
	Items() []CartItem
	ItemsAdded() map[string]*CartItem
}

type Rule interface {
	Apply(cart Cart) []CartItem
}
