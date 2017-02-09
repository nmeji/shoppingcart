package amaysim

/*
Product represents the shop product given the product code, name, and price
*/
type Product struct {
	Code, Name string
	Price      float32
}

/*
CartItem represents the shop product inside the cart
*/
type CartItem struct {
	Quantity int
	SubTotal float32
	Product
}

/*
Cart interface contains the common functions of a typical shopping cart
*/
type Cart interface {
	Add(item CartItem)
	Remove(item CartItem)
	AddPromoCode(code string)
	UsesPromoCode(code string) bool
	PromoCodesApplied() []string
	Total() float32
	Items() []CartItem
	ItemsAdded() map[string]*CartItem
}

/*
Rule is something that can be applied to a cart
*/
type Rule interface {
	Apply(cart Cart) []CartItem
}

/*
Equals provides a way to compare two products
*/
func (a *Product) Equals(b *Product) bool {
	if a == b {
		return true
	}
	return a.Code == b.Code && a.Name == b.Name && a.Price == b.Price
}

/*
Equals provides a way to compare two cart items
*/
func (a *CartItem) Equals(b *CartItem) bool {
	if a == b {
		return true
	}
	return a.Quantity == b.Quantity && a.SubTotal == b.SubTotal && a.Product.Equals(&b.Product)
}
