# Inventory Manager

Manages a typical inventory for a store front. It's CLI-based toy app to solve the interview problem, not based
on real-world models/use-cases.

## Models
```
User -> Customer
     -> Employee

Item -> SKU -> ProductGroup

Order : Item, User, Discount, Tax, Gross, Net

Inventory: Item, Order

```

## Usecases

* Replenish Inventory(Item)
* Place an order(User, Item)
* Summary of sales
* Summary of inventory
* Apply discount(Item/User)

## What can be better?

* Handling concurrency better: how do you handle multiple checkouts happening at different registers?
* More realistic models and usecases. Eg. "Silly things" like taxes are not calculated properly.
* Better error handling
* Caching/faster retrieval for sales and inventory summary
* More tests

## Installation

0. Install golang and set GOPATH.
1. Clone this repo into `src` folder of your `GOPATH`
2. Install glide for your OS: `curl https://glide.sh/get | sh` or `brew install glide`
3. Run `glide install` in your repo folder
4. Run `go run main.go`.

## Checks

1. Install gometalinter. `go get -u github.com/alecthomas/gometalinter`
2. Install linters: `gometalinter --install`
3. Run lints: `gometalinter ./...`

## Tests

1. Run `go test -cover $(glide nv)`