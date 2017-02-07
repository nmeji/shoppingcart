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

