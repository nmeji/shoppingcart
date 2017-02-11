package amaysim_test

import (
	"math"
	"testing"

	. "github.com/nmeji/shoppingcart/amaysim"
)

func AdjustPrecision(price float32) float32 {
	// round down to 2 decimal places
	return float32(math.Floor(float64(price*100)) / 100)
}

func TestBuy3For2(t *testing.T) {
	cart := SetUp()
	cart.Add(CartItem{3, 0.0, *unli1GB})
	cart.Add(CartItem{1, 0.0, *unli5GB})
	totalPrice := cart.Total()
	totalPrice = AdjustPrecision(totalPrice)
	expectedPrice := float32(94.70)
	if expectedPrice != totalPrice {
		t.Error("Assert failed: Buy3For2 Rule")
	}
}

func TestBulkDiscountOn5GB(t *testing.T) {
	cart := SetUp()
	cart.Add(CartItem{2, 0.0, *unli1GB})
	cart.Add(CartItem{4, 0.0, *unli5GB})
	totalPrice := cart.Total()
	totalPrice = AdjustPrecision(totalPrice)
	expectedPrice := float32(209.40)
	if expectedPrice != totalPrice {
		t.Errorf("Assert failed: BulkDiscountOn5GB Rule (actual: %.02f, expected: %.02f)", totalPrice, expectedPrice)
	}
}

func TestBuy1Get1(t *testing.T) {
	cart := SetUp()
	cart.Add(CartItem{1, 0.0, *unli1GB})
	cart.Add(CartItem{2, 0.0, *unli2GB})
	items := cart.Items()
	expectedItem, expectedQuantity := "1gb", 2
	found := false
	for _, item := range items {
		found = found || (expectedItem == item.Code && expectedQuantity == item.Quantity)
	}
	if !found {
		t.Error("Assert failed: Buy1Get1 Rule")
	}
}

func TestOverallDiscountPromo(t *testing.T) {
	// 1. test for unknown promo code (should have no effect)
	// 2. test for known promo code (eg. I<3AMAYSIM)
	cart := SetUp()
	cart.Add(CartItem{1, 0.0, *unli1GB})
	cart.Add(CartItem{1, 0.0, *data1GB})
	cart.AddPromoCode("code")
	totalPrice := AdjustPrecision(cart.Total())
	expectedPrice := float32(34.80) // no discount
	if expectedPrice != totalPrice {
		t.Error("Assert failed: OverallDiscountPromo Rule should not apply if promo code is not valid")
	}
	cart.AddPromoCode("I<3AMAYSIM")
	totalPrice = AdjustPrecision(cart.Total())
	expectedPrice = float32(31.32) // 10% discount
	if expectedPrice != totalPrice {
		t.Errorf("Assert failed: OverallDiscountPromo Rule (actual: %.02f, expected: %.02f)", totalPrice, expectedPrice)
	}
}
