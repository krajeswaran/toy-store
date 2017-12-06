// Fake models. SAD!
package models

import (
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

type Mocks struct {
	Items     []Item
	LineItems []OrderLineItem
	Customers []Customer
	Employees []Employee

	PurchaseUserId       uuid.UUID
	PurchaseUserDiscount int
}

func (m *Mocks) InitInventory() {
	if len(m.Items) == 0 {
		// fake fill Items
		itemNames := []string{"Dora", "Teddy", "Superman", "Spiderman", "Batman"}
		superHeroSkuId := uuid.NewV4()
		superHeroDiscount := rand.Intn(50)

		// init Items first time around
		for i := 0; i < len(itemNames); i++ {
			item := new(Item)
			item.Name = itemNames[i]
			item.Id = uuid.NewV4()
			item.Price = decimal.New(rand.Int63n(100), 0)
			item.SKU = *new(SKU)
			if strings.Contains(item.Name, "man") {
				// assign a superhero sku with a discount on sku
				item.SKU.Name = "SuperHeroToy"
				item.SkuId = superHeroSkuId
				item.SKU.DiscountPercentage = superHeroDiscount
			} else {
				// no discounts for others
				item.SkuId = uuid.NewV4()
				item.SKU.Name = item.Name + strconv.Itoa(i)
			}
			item.Status = AvailableItemStatus
			item.Created = time.Now().UTC()
			item.Modified = time.Now().UTC()
			m.Items = append(m.Items, *item)
		}
	}
}

// public for testability
func (m *Mocks) InitUsers() {
	if len(m.Customers) == 0 {
		// fake fill Customers
		customerNames := []string{"Alpha", "Bravo"}

		// init Items first time around
		for i := 0; i < len(customerNames); i++ {
			customer := new(Customer)
			customer.Name = customerNames[i]
			customer.Id = uuid.NewV4()
			customer.DiscountPercentage = rand.Intn(10)
			customer.Status = EnabledUserStatus
			customer.Created = time.Now().UTC()
			customer.Modified = time.Now().UTC()
			m.Customers = append(m.Customers, *customer)
		}
	}

	if len(m.Employees) == 0 {
		// fake fill Employees
		employeeNames := []string{"Anna", "Boris"}

		// init Items first time around
		for i := 0; i < len(employeeNames); i++ {
			employee := new(Employee)
			employee.Name = employeeNames[i]
			employee.Id = uuid.NewV4()
			// Employees have better discount
			employee.DiscountPercentage = rand.Intn(30)
			employee.Status = EnabledUserStatus
			employee.Created = time.Now().UTC()
			employee.Modified = time.Now().UTC()
			m.Employees = append(m.Employees, *employee)
		}
	}
}

// for testability
func (m *Mocks) GetMockedUser(i int) *Customer {
	return &m.Customers[i]
}

// for testability
func (m *Mocks) GetMockedItem(i int) *Item {
	return &m.Items[i]
}

// for testability
func (m *Mocks) ClearMockModelCache() {
	m.LineItems = nil
	m.PurchaseUserDiscount = 0
	m.PurchaseUserId = uuid.Nil
}
