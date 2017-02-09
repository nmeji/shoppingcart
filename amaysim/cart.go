package amaysim

/*
ShoppingCart is a Cart implementation
*/
type ShoppingCart struct {
	PricingRules []Rule
	Added        map[string]*CartItem
	PromoCode    map[string]int
}

/*
Add item(s) to the ShoppingCart
*/
func (cart ShoppingCart) Add(item CartItem) {
	if i, ok := cart.Added[item.Code]; ok {
		i.Quantity += item.Quantity
	} else {
		cart.Added[item.Code] = &item
	}
}

/*
AddPromoCode allows to set Promo Codes
*/
func (cart ShoppingCart) AddPromoCode(code string) {
	cart.PromoCode[code] = 1
}

/*
UsesPromoCode returns true if a specific code matches the set promo code(s)
*/
func (cart ShoppingCart) UsesPromoCode(code string) bool {
	return cart.PromoCode[code] == 1
}

/*
PromoCodesApplied returns a slice of all promo codes added
*/
func (cart ShoppingCart) PromoCodesApplied() []string {
	if size := len(cart.PromoCode); size > 0 {
		promoCodes := make([]string, 0, size)
		for code, _ := range cart.PromoCode {
			promoCodes = append(promoCodes, code)
		}
		return promoCodes
	}
	return nil
}

/*
Remove item(s) from the cart
*/
func (cart ShoppingCart) Remove(item CartItem) {
	if i, ok := cart.Added[item.Code]; ok {
		if q := i.Quantity - item.Quantity; q < 1 {
			delete(cart.Added, item.Code)
		} else {
			i.Quantity = q
		}
	}
}

/*
Total computes all the sub total and applies the pricing rules if there is any and returns the sum of all sub totals
*/
func (cart ShoppingCart) Total() float32 {
	for _, item := range cart.Added {
		item.SubTotal = float32(item.Quantity) * item.Price
	}
	var c Cart
	for _, rule := range cart.PricingRules {
		c = cart
		rule.Apply(c)
	}
	TotalPrice := float32(0)
	for _, item := range cart.Added {
		TotalPrice += item.SubTotal
	}
	return TotalPrice
}

/*
Items returns all the added items including the promo items
*/
func (cart ShoppingCart) Items() []CartItem {
	copyItems := make(map[string]*CartItem)
	for k, v := range cart.Added {
		copyItems[k] = v
	}
	var c Cart
	for _, rule := range cart.PricingRules {
		c = cart
		freebies := rule.Apply(c)
		for _, item := range freebies {
			if i, ok := copyItems[item.Code]; ok {
				i.Quantity += item.Quantity
			} else {
				copyItems[item.Code] = &item
			}
		}
	}
	items := make([]CartItem, 0, len(copyItems))
	for _, item := range copyItems {
		items = append(items, *item)
	}
	return items
}

/*
ItemsAdded returns the items added by the client
*/
func (cart ShoppingCart) ItemsAdded() map[string]*CartItem {
	return cart.Added
}
