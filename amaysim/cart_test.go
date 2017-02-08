package amaysim_test

import (
	"testing"

	. "github.com/nmeji/shoppingcart/amaysim"
)

func AssertCartItem(item *CartItem, t *testing.T) {
	if item.Quantity != 4 {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Quantity did not match")
	}
	if item.Code != "ult_small" {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Product.Code did not match")
	}
	if item.Name != "Unlimited 1GB" {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Product.Name did not match")
	}
	if item.Price != 24.90 {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Product.Price did not match")
	}
}

func TestAdd(t *testing.T) {
	cart := ShoppingCart{Added: make(map[string]*CartItem)}
	cart.Add(CartItem{
		Quantity: 4,
		Product: Product{
			Code:  "ult_small",
			Name:  "Unlimited 1GB",
			Price: 24.90,
		},
	})
	if item, ok := cart.ItemsAdded()["ult_small"]; !ok {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Item is not present in ShoppingCart.ItemsAdded()")
	} else {
		AssertCartItem(item, t)
	}
	if cartSize := len(cart.Items()); cartSize != 1 {
		t.Error("Assert failed: ShoppingCart.Add(CartItem) - Item is not present in ShoppingCart.Items")
	} else {
		AssertCartItem(&cart.Items()[0], t)
	}
}
