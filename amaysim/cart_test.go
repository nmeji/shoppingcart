package amaysim_test

import (
	"reflect"
	"testing"

	. "github.com/nmeji/shoppingcart/amaysim"
)

func CreateCartItem(qty int, subtotal float32, code string, name string, price float32) *CartItem {
	return &CartItem{
		Quantity: qty,
		SubTotal: subtotal,
		Product: Product{
			Code:  code,
			Name:  name,
			Price: price,
		},
	}
}

func TestAdd(t *testing.T) {
	cart := ShoppingCart{Added: make(map[string]*CartItem)}
	initialInput := CreateCartItem(1, 2.3, "ult_small", "Unlimited 1GB", 24.90)
	sameItemInput := CreateCartItem(2, 2.3, "ult_small", "Unlimited 1GB", 24.90)
	expectedItemWithAdjustedQty := CreateCartItem(3, 2.3, "ult_small", "Unlimited 1GB", 24.90)

	// after calling cart.Add(input), we expect that it should be present in cart.ItemsAdded()
	// and it should be equal to the item we added
	cart.Add(*initialInput)
	if item, ok := cart.ItemsAdded()["ult_small"]; !ok || !item.Equals(initialInput) {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Item should be present in ShoppingCart.ItemsAdded()")
	}

	// adding of the same item should increase the quantity based on the new item to be added
	// item total quantity = present item quantity + new item quantity
	cart.Add(*sameItemInput)
	if item, ok := cart.ItemsAdded()["ult_small"]; !ok || !item.Equals(expectedItemWithAdjustedQty) {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Item quantity should adjust based on the new item to be added")
	}
}

func TestRemove(t *testing.T) {
	cart := ShoppingCart{Added: map[string]*CartItem{
		"ult_large": CreateCartItem(5, 0.0, "ult_large", "Unlimited 5GB", 44.90),
	},
	}
	itemToRemove := CreateCartItem(2, 0.0, "ult_large", "Unlimited 5GB", 44.90)

	// calling cart.remove(item) that does not reduce the quantity to 0 should still exist in the cart
	cart.Remove(*itemToRemove)
	expectedQuantityLeft := 3
	if item, ok := cart.ItemsAdded()["ult_large"]; !ok {
		t.Error("Assert failed: ShoppingCart.Remove(CartItem) - Item is removed from the cart even though quantity > 0")
	} else if item.Quantity != expectedQuantityLeft {
		t.Error("Assert failed: ShoppingCart.Remove(CartItem) - Item quantity did not match with expected quantity")
	}

	// callinig cart.remove(item) that reduces item quantity to zero or less should also remove the item from the cart
	itemToRemove.Quantity = 4 // after removing this much, item quantity = -1
	cart.Remove(*itemToRemove)
	if _, ok := cart.ItemsAdded()["ult_large"]; ok {
		t.Error("Assert failed: ShoppingCart.Remove(CartItem) - After quantity drops to zero or less, it should be removed from the cart")
	}
}

func TestAddPromoCode(t *testing.T) {
	// 1. test for unknown promo code (should have no effect)
	// 2. test for known promo code (eg. I<3AMAYSIM)
	cart := ShoppingCart{
		PricingRules: AllRules, // contains rule for promo code "I<3AMAYSIM"
		Added:        make(map[string]*CartItem),
		PromoCode:    make(map[string]int),
	}
	cart.Add(*CreateCartItem(1, 0.0, "ult_small", "Unlimited 1GB", 24.90))
	cart.Add(*CreateCartItem(1, 0.0, "1gb", "1 GB Data-Pack", 9.90))
	cart.AddPromoCode("QWERTY")
	total := cart.Total()
	expectedPrice := float32(34.80) // no discount
	if total != expectedPrice {
		t.Error("Assert failed: ShoppingCart.AddPromoCode(string) - Expected price did not match after unknown promo code is applied")
	}
	cart.AddPromoCode("I<3AMAYSIM")
	total = cart.Total()
	expectedDiscountedPrice := float32(31.32) // 10% discount
	if total != expectedDiscountedPrice {
		t.Error("Assert failed: ShoppingCart.AddPromoCode(string) - Expected discounted price did not match after promo code is applied")
	}
}

func TestUsesPromoCode(t *testing.T) {
	cart := ShoppingCart{PromoCode: map[string]int{
		"I<3AMAYSIM": 1,
	},
	}
	// 1. test for missing promo code
	if cart.UsesPromoCode("QWERTY") {
		t.Error("Assert failed: ShoppingCart.UsesPromoCode(string) - Expected false for promo codes not added")
	}
	// 2. test for added promo code
	if !cart.UsesPromoCode("I<3AMAYSIM") {
		t.Error("Assert failed: ShoppingCart.UsesPromoCode(string) - Expected true for promo codes already added")
	}
}

func TestPromoCodesApplied(t *testing.T) {
	cart := ShoppingCart{PromoCode: map[string]int{
		"I<3AMAYSIM": 1,
		"QWERTY":     1,
	},
	}
	// should match with all added promo codes
	if !reflect.DeepEqual([]string{"I<3AMAYSIM", "QWERTY"}, cart.PromoCodesApplied()) {
		t.Error("Assert failed: ShoppingCart.PromoCodesApplied() - Expected match")
	}
}

func MatchTotal(input []*CartItem, expected float32) bool {
	cart := ShoppingCart{Added: make(map[string]*CartItem)}
	for _, item := range input {
		cart.Add(*item)
	}
	return expected == cart.Total()
}

func TestTotal(t *testing.T) {
	// 1. should compute zero if no items are added yet
	// 2. should have correct computation of summation of subtotals
	tests := []struct {
		Input    []*CartItem
		Expected float32
	}{
		{nil, 0},
		{[]*CartItem{
			CreateCartItem(1, 0.0, "a", "A", 24.90),
			CreateCartItem(1, 0.0, "b", "B", 9.90),
		},
			float32(34.80),
		},
	}
	for _, test := range tests {
		if !MatchTotal(test.Input, test.Expected) {
			t.Error("Assert failed: ShoppingCart.Total() - Expected match")
		}
	}
}

func TestItems(t *testing.T) {
	cart := ShoppingCart{Added: make(map[string]*CartItem)}
	// should contain items added
	input := CreateCartItem(1, 0, "a", "A", 1.0)
	cart.Add(*input)
	allItems := cart.Items()
	found := false
	for _, item := range allItems {
		found = found || input.Equals(&item)
	}
	if !found {
		t.Error("Assert failed: ShoppingCart.Items() - Expected item exists after being added")
	}
}

func TestItemsAdded(t *testing.T) {
	cart := ShoppingCart{Added: make(map[string]*CartItem)}
	// should contain items added
	input := CreateCartItem(1, 0, "a", "A", 1.0)
	cart.Add(*input)
	if item, ok := cart.ItemsAdded()["a"]; !ok || !input.Equals(item) {
		t.Error("Assert failed: ShoppingCart.ItemsAdded() - Expected item exists after being added")
	}
}
