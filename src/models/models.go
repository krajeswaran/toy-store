package models

import (
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"sync"
	"time"
)

// Base fields are common across objects
type BaseFields struct {
	// Unique object identifier
	Id uuid.UUID
	// Created timestamp
	Created time.Time
	// Modified timestamp
	Modified time.Time
	// denotes status of the current object: like ENABLED, DISABLED, DELETED, PENDING etc
	// NOTE: not realistic in real world models
	Status string
}

type User struct {
	Name               string
	DiscountPercentage int
	BaseFields
}

type Customer struct {
	User
	DiscountPercentage int
	// Note: Many specific customer fields left to imagination
}

type Employee struct {
	User
	DiscountPercentage int
	// Note: Many specific employee fields left to imagination
}

type Item struct {
	DiscountPercentage int
	Name               string
	Description        string
	Price              decimal.Decimal
	BaseFields
	SKU
	ProductGroup
}

type SKU struct {
	SkuId              uuid.UUID
	Name               string
	Description        string
	DiscountPercentage int
}

// Unused currently, but can be used to group products at a higher level than SKU
type ProductGroup struct {
	ProductGroupId     uuid.UUID
	Name               string
	Description        string
	DiscountPercentage int
}

type OrderLineItem struct {
	Item     *Item
	Quantity int64
}

type Order struct {
	UserId      uuid.UUID
	LineItems   []OrderLineItem
	NetAmount   decimal.Decimal
	GrossAmount decimal.Decimal
	// Tag an order with particular notes. Eg. replenishment order vs purchase order
	Tag string
	BaseFields
}

type LedgerEntry struct {
	Order   *Order
	Item    *Item
	Credit  decimal.Decimal
	Debit   decimal.Decimal
	Balance decimal.Decimal
	BaseFields
}

type Inventory struct {
	Ledger []LedgerEntry
}

// Order Status
const (
	FailedOrderStatus    = "failed"
	PendingOrderStatus   = "pending"
	CompletedOrderStatus = "completed"
)

// Ledger entry Status
const (
	CreatedLedgerEntryStatus = "created"
	AbortedLedgerEnryStatus  = "aborted"
)

// User Status
const (
	EnabledUserStatus  = "enabled"
	DisabledUserStatus = "disabled"
)

// Item Status
const (
	AvailableItemStatus = "available"
	BlockedItemStatus   = "blocked"
)

// Let's create a master inventory and admin user when app starts once lazily(singletons)
var inventorySync sync.Once
var inventoryInstance *Inventory

func GetMasterInventory() *Inventory {
	inventorySync.Do(func() {
		inventoryInstance = &Inventory{
			Ledger: nil,
		}
	})
	return inventoryInstance
}

var adminSync sync.Once
var storeAdmin *User

func GetStoreAdmin() *User {
	adminSync.Do(func() {
		storeAdmin = &User{
			DiscountPercentage: 0,
			BaseFields: BaseFields{
				Id:       uuid.NewV4(),
				Created:  time.Now().UTC(),
				Modified: time.Now().UTC(),
				Status:   EnabledUserStatus,
			},
		}
	})
	return storeAdmin
}
