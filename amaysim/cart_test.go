package amaysim_test

import (
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
	shoppincart := ShoppingCart{Added: map[string]*CartItem{
		"ult_large": CreateCartItem(5, 0.0, "ult_large", "Unlimited 5GB", 44.90),
	},
	}
	itemToRemove := CreateCartItem(2, 0.0, "ult_large", "Unlimited 5GB", 44.90)

	// calling cart.remove(item) that does not reduce the quantity to 0 should still exist in the cart
	shoppincart.Remove(*itemToRemove)
	expectedQuantityLeft := 3
	if item, ok := shoppincart.ItemsAdded()["ult_large"]; !ok {
		t.Error("Assert failed: ShoppingCart.Remove(CartItem) - Item is removed from the cart even though quantity > 0")
	} else if item.Quantity != expectedQuantityLeft {
		t.Error("Assert failed: ShoppingCart.Remove(CartItem) - Item quantity did not match with expected quantity")
	}

	// callinig cart.remove(item) that reduces item quantity to zero or less should also remove the item from the cart
	itemToRemove.Quantity = 4 // after removing this much, item quantity = -1
	shoppincart.Remove(*itemToRemove)
	if _, ok := shoppincart.ItemsAdded()["ult_large"]; ok {
		t.Error("Assert failed: ShoppingCart.Remove(CartItem) - After quantity drops to zero or less, it should be removed from the cart")
	}
}

func TestAddPromoCode(t *testing.T) {

}
