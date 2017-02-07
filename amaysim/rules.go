package amaysim

type Buy3For2 struct{}
type Buy1Get1 struct{}
type BulkDiscountOn5gb struct{}
type OverallDiscountPromo struct{}

func (rule Buy3For2) Apply(cart Cart) []CartItem {
	added := cart.ItemsAdded()
	if item, ok := added["ult_small"]; ok {
		if i := item.Quantity / 3; i > 0 {
			item.SubTotal -= float32(i) * item.Price
		}
	}
	return nil
}

func (rule Buy1Get1) Apply(cart Cart) []CartItem {
	added := cart.ItemsAdded()
	if item, ok := added["ult_medium"]; ok {
		freebie, err := catalog.Lookup("1gb")
		if err != nil {
			panic(err)
		}
		return []CartItem{
			CartItem{item.Quantity, 0.0, *freebie},
		}
	}
	return nil
}

func (rule BulkDiscountOn5gb) Apply(cart Cart) []CartItem {
	added := cart.ItemsAdded()
	if item, ok := added["ult_large"]; ok {
		if item.Quantity > 3 {
			item.SubTotal = 39.90 * float32(item.Quantity)
		}
	}
	return nil
}

func (rule OverallDiscountPromo) Apply(cart Cart) []CartItem {
	if cart.UsesPromoCode("I<3AMAYSIM") {
		added := cart.ItemsAdded()
		for _, item := range added {
			item.SubTotal *= 0.9
		}
	}
	return nil
}

var AllRules []Rule

func init() {
	AllRules = make([]Rule, 4)
	var rule1, rule2, rule3, rule4 Rule
	rule1 = Buy3For2{}
	rule2 = Buy1Get1{}
	rule3 = BulkDiscountOn5gb{}
	rule4 = OverallDiscountPromo{}
	AllRules[0] = rule1
	AllRules[1] = rule2
	AllRules[2] = rule3
	AllRules[3] = rule4
}
