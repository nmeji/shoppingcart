# ShoppingCart

Simple Shopping Cart written in Golang.

## Go get it!

`go get github.com/nmeji/shoppingcart`

## Run

After go-getting, the binary build `$GOBIN/shoppingcart` should already be present.

Run it in the same folder with product-catalog.json

## Flexibility

### Product Catalog

You may edit product-catalog.json to reflect changes to available products

### Promotional Rules

At first, I thought of implementing the rules in something like pricing-rules.json. But due to the extra overhead I opted to write it in rules.go instead.

## Examples

You can either run the commands in $GOPATH/src/github.com/nmeji/shoppingcart or you may also copy the product-catalog.json file from that folder and put it in some other directory. The important thing is to make sure that $GOBIN variable is already set. If it is not yet set, I think it is set by default to $GOPATH/bin directory.

1. Run `shoppingcart` (This will search for product-catalog.json in the same folder)
2. You may type `help` to list all available commands
3. Just to provide you with some commands right away, here's a list of the important commands to remember:
4. Type `catalog` to show the different products and their prices.
5. In order to add items to cart, the usage is `add <id> [quantity]`
6. Use `cart` to show cart summary which provides you the total price and the total items you can expect to get after purchasing.
7. If you wish to apply a promo code, type `promo <code>` (eg. promo I<3AMAYSIM)

## Tests and Coverage

2 tests: amaysim/cart_test.go and amaysim/rules_test.go. In order to run both tests, run `go test -coverprofile=testprofile ./amaysim/`
