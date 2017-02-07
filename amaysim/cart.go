package amaysim

type ShoppingCart struct {
	PricingRules []Rule
	Added        map[string]*CartItem
	PromoCode    map[string]int
}

func (self ShoppingCart) Add(item CartItem) {
	if i, ok := self.Added[item.Code]; ok {
		i.Quantity += item.Quantity
	} else {
		self.Added[item.Code] = &item
	}
}

func (cart ShoppingCart) AddPromoCode(code string) {
	cart.PromoCode[code] = 1
}

func (cart ShoppingCart) UsesPromoCode(code string) bool {
	return cart.PromoCode[code] == 1
}

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

func (cart ShoppingCart) Remove(item CartItem) {
	if i, ok := cart.Added[item.Code]; ok {
		if q := item.Quantity - i.Quantity; q < 1 {
			delete(cart.Added, item.Code)
		} else {
			item.Quantity = q
		}
	}
}

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

func (cart ShoppingCart) ItemsAdded() map[string]*CartItem {
	return cart.Added
}
