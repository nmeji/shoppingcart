package amaysim_test

import (
	"reflect"
	"strings"
	"testing"

	. "github.com/nmeji/shoppingcart/amaysim"
)

const TEST_PRODUCTS = `[
  {"Code": "ult_small", "Name": "Unlimited 1GB", "Price": 24.90},
  {"Code": "ult_medium", "Name": "Unlimited 2GB", "Price": 29.90},
  {"Code": "ult_large", "Name": "Unlimited 5GB", "Price": 44.90},
  {"Code": "1gb", "Name": "1 GB Data-pack", "Price": 9.90}
  ]`

var unli1GB, unli2GB, unli5GB, data1GB *Product

func init() {
	products, err := InitProducts(strings.NewReader(TEST_PRODUCTS))
	if err != nil {
		panic(err)
	}
	unli1GB, _ = products.Lookup("ult_small")
	unli2GB, _ = products.Lookup("ult_medium")
	unli5GB, _ = products.Lookup("ult_large")
	data1GB, _ = products.Lookup("1gb")
}

func SetUp() Cart {
	shoppingcart := ShoppingCart{
		PricingRules: AllRules,
		Added:        make(map[string]*CartItem),
		PromoCode:    make(map[string]int),
	}
	var cart Cart
	cart = shoppingcart
	return cart
}

func TestAdd(t *testing.T) {
	cart := SetUp()
	initialInput := CartItem{1, 0.0, *unli1GB}

	// after calling cart.Add(input), we expect that it should be present in cart.ItemsAdded()
	// and it should be equal to the item we added
	cart.Add(initialInput)
	if item, ok := cart.ItemsAdded()["ult_small"]; !ok || !item.Equals(&initialInput) {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Item should be present in ShoppingCart.ItemsAdded()")
	}

	// adding of the same item should increase the quantity based on the new item to be added
	// item total quantity = present item quantity + new item quantity
	cart.Add(initialInput)
	expectedItemWithAdjustedQty := CartItem{2, 0.0, *unli1GB}
	if item, ok := cart.ItemsAdded()["ult_small"]; !ok || !item.Equals(&expectedItemWithAdjustedQty) {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Item quantity should adjust based on the new item to be added")
	}
}

func TestRemove(t *testing.T) {
	cart := SetUp()
	cart.Add(CartItem{5, 0.0, *unli5GB})

	// calling cart.remove(item) that does not reduce the quantity to 0 should still exist in the cart
	cart.Remove(CartItem{2, 0.0, *unli5GB})
	expectedQuantityLeft := 3
	if item, ok := cart.ItemsAdded()["ult_large"]; !ok {
		t.Error("Assert failed: ShoppingCart.Remove(CartItem) - Item is removed from the cart even though quantity > 0")
	} else if item.Quantity != expectedQuantityLeft {
		t.Error("Assert failed: ShoppingCart.Remove(CartItem) - Item quantity did not match with expected quantity")
	}

	// callinig cart.remove(item) that reduces item quantity to zero or less should also remove the item from the cart
	cart.Remove(CartItem{4, 0.0, *unli5GB}) // after removing this much, item quantity = -1
	if _, ok := cart.ItemsAdded()["ult_large"]; ok {
		t.Error("Assert failed: ShoppingCart.Remove(CartItem) - After quantity drops to zero or less, it should be removed from the cart")
	}
}

func TestAddPromoCode(t *testing.T) {
	cart := SetUp()
	cart.AddPromoCode("code")
	// expected to return true
	if !cart.UsesPromoCode("code") {
		t.Error("Assert failed: ShoppingCart.AddPromoCode(string) expects UsesPromoCode(string) to return true")
	}
	if promoCodes := cart.PromoCodesApplied(); len(promoCodes) == 0 {
		t.Error("Assert failed: ShoppingCart.AddPromoCode(string) expects non-empty PromoCodesApplied()")
	} else {
		found := false
		for _, promoCode := range promoCodes {
			found = found || promoCode == "code"
		}
		if !found {
			t.Error("Assert failed: ShoppingCart.AddPromoCode(string) expects promo code to exist in PromoCodesApplied()")
		}
	}
}

func TestUsesPromoCode(t *testing.T) {
	cart := SetUp()
	cart.AddPromoCode("I<3AMAYSIM")
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
	cart := SetUp()
	cart.AddPromoCode("I<3AMAYSIM")
	cart.AddPromoCode("QWERTY")
	// should match with all added promo codes
	if !reflect.DeepEqual([]string{"I<3AMAYSIM", "QWERTY"}, cart.PromoCodesApplied()) {
		t.Error("Assert failed: ShoppingCart.PromoCodesApplied() - Expected match")
	}
}

func MatchTotal(input []*CartItem, expected float32) bool {
	cart := SetUp()
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
			&CartItem{1, 0.0, *unli1GB},
			&CartItem{1, 0.0, *data1GB},
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
	cart := SetUp()
	// should contain items added
	input := CartItem{1, 0.0, *unli2GB}
	cart.Add(input)
	allItems := cart.Items()
	found := false
	for _, item := range allItems {
		found = found || item.Equals(&input)
	}
	if !found {
		t.Error("Assert failed: ShoppingCart.Items() - Expected item exists after being added")
	}
}

func TestItemsAdded(t *testing.T) {
	cart := SetUp()
	// should contain items added
	input := CartItem{1, 0.0, *unli5GB}
	cart.Add(input)
	if item, ok := cart.ItemsAdded()[unli5GB.Code]; !ok || !item.Equals(&input) {
		t.Error("Assert failed: ShoppingCart.ItemsAdded() - Expected item exists after being added")
	}
}
