# Inventory Manager

Manages a typical inventory using a ledger. It's CLI-based toy app to solve the interview problem, not based
on real-world models/use-cases.

## Models
High level view of the models used:

```uml
User <- Customer
     <- Employee

Item <- SKU <- ProductGroup

Order -<> Items, User, Net, Gross 

Ledger -<> Order, Item, Credit, Debit, Balance

```

## Usecases

* Replenish an item in Inventory
* Place an order by an User for a list of Items
* Summary of sales so far today
* Summary of inventory
* Discounts are at both user level(mock users created with different types of discount) and at item/SKU/Product Group levels

## What can be better?

* More tests
* Handling concurrency better: how do you handle multiple checkouts happening at different registers?
* More realistic models and usecases. Eg. "Silly things" like taxes are not considered in this toy model
* Caching/faster retrieval for sales and inventory summary
* Better error handling

## Installation

0. Install golang(version > 1.8.1) and set GOPATH. For eg. `$HOME/gopath`.
1. Clone this repo into a folder and add that folder to your `GOPATH`. Eg. `export GOPATH=$HOME/gopath:$HOME/src/toy-store`
2. `cd toy-store/src` and run `go run main.go`.

## Tests

1. Run `go test -cover $(glide nv)`
