package usecases

import (
	"error"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"models"
	"sort"
	"time"
)

// InventoryUsecaseRepository contains business logic
type InventoryUsecaseRepository struct{}

// Replenish an item in inventory
func (i *InventoryUsecaseRepository) Replenish(item models.Item, count decimal.Decimal) (bool,
	error) {
	// check input
	if uuid.Equal(item.Id, uuid.Nil) {
		err := errors.NewError(errors.ReplenishError, "Empty item given")
		return false, err
	}

	// create a replenishment order
	order := models.Order{
		UserId:      models.GetStoreAdmin().Id,
		LineItems:   nil,
		NetAmount:   decimal.Zero,
		GrossAmount: decimal.Zero,
		Tag:         "replenishment",
		BaseFields: models.BaseFields{
			Id:       uuid.NewV4(),
			Created:  time.Now().UTC(),
			Modified: time.Now().UTC(),
			Status:   models.CompletedOrderStatus,
		},
	}

	// find or create ledger entry
	itemBalance := findItemBalanceInLedger(item)

	itemBalance = itemBalance.Add(count)

	entry := models.LedgerEntry{
		Order:   &order,
		Item:    &item,
		Credit:  count,
		Debit:   decimal.Zero,
		Balance: itemBalance,
		BaseFields: models.BaseFields{
			Id:       uuid.UUID{},
			Created:  time.Now().UTC(),
			Modified: time.Now().UTC(),
			Status:   models.CreatedLedgerEntryStatus,
		},
	}

	// add the ledger entry to inventory
	models.GetMasterInventory().Ledger = append(models.GetMasterInventory().Ledger, entry)

	// we are done
	return true, nil
}

// Place a purchase order for a user
func (i *InventoryUsecaseRepository) Purchase(lineItems *[]models.OrderLineItem,
	userId uuid.UUID, userDiscount int) (decimal.Decimal, error) {
	// check input
	if len(*lineItems) == 0 || uuid.Equal(userId, uuid.Nil) {
		err := errors.NewError(errors.OrderError, "Empty line items/user given")
		return decimal.Zero, err
	}

	// create a purchase order
	net, gross := CalcOrderAmounts(lineItems, userDiscount)
	order := models.Order{
		UserId:      userId,
		LineItems:   *lineItems,
		NetAmount:   net,
		GrossAmount: gross,
		Tag:         "purchase",
		BaseFields: models.BaseFields{
			Id:       uuid.NewV4(),
			Created:  time.Now().UTC(),
			Modified: time.Now().UTC(),
			Status:   models.CompletedOrderStatus,
		},
	}

	// process ledger entry for each line item
	for _, line := range *lineItems {
		itemQty := decimal.New(line.Quantity, 0)
		itemBalance := findItemBalanceInLedger(*line.Item)

		itemBalance = itemBalance.Sub(itemQty)
		if itemBalance.Cmp(decimal.Zero) < 0 {
			err := errors.NewError(errors.OrderError, "Inventory item balance will become negative")
			return decimal.Zero, err
		}

		entry := models.LedgerEntry{
			Order:   &order,
			Item:    line.Item,
			Credit:  decimal.Zero,
			Debit:   itemQty,
			Balance: itemBalance,
			BaseFields: models.BaseFields{
				Id:       uuid.UUID{},
				Created:  time.Now().UTC(),
				Modified: time.Now().UTC(),
				Status:   models.CreatedLedgerEntryStatus,
			},
		}

		// add the ledger entry to inventory
		models.GetMasterInventory().Ledger = append(models.GetMasterInventory().Ledger, entry)
	}

	// we are done
	return net, nil
}

func findItemBalanceInLedger(item models.Item) decimal.Decimal {
	ledger := models.GetMasterInventory().Ledger
	itemBalance := decimal.Zero
	if ledger != nil {
		sortLedger(ledger)
		for _, entry := range ledger {
			if uuid.Equal(entry.Item.Id, item.Id) {
				// found the item
				itemBalance = entry.Balance
				break
			}
		}
	}
	return itemBalance
}

func (i *InventoryUsecaseRepository) SaleSummary(from time.Time) (map[uuid.UUID]decimal.
	Decimal, decimal.Decimal) {
	summary := make(map[uuid.UUID]decimal.Decimal)
	totalSales := decimal.Zero
	orders := make(map[uuid.UUID]decimal.Decimal)

	ledger := models.GetMasterInventory().Ledger
	if ledger != nil {
		sortLedger(ledger)
		for _, entry := range ledger {
			if entry.Modified.After(from) {
				summary[entry.Item.Id] = summary[entry.Item.Id].Add(entry.Debit)
				orders[entry.Order.Id] = entry.Order.NetAmount
			}
		}
	}

	for _, value := range orders {
		totalSales = totalSales.Add(value)
	}

	return summary, totalSales
}

func (i *InventoryUsecaseRepository) InventorySummary(till time.Time) map[uuid.UUID]decimal.
	Decimal {
	summary := make(map[uuid.UUID]decimal.Decimal)

	ledger := models.GetMasterInventory().Ledger
	if ledger != nil {
		sortLedger(ledger)
		for _, entry := range ledger {
			if entry.Modified.Before(till) {
				if _, ok := summary[entry.Item.Id]; !ok {
					summary[entry.Item.Id] = entry.Balance
				}
			}
		}
	}

	return summary
}

// Reverse sort inventory based on timestamp
func sortLedger(ledger []models.LedgerEntry) {
	sort.Slice(ledger, func(i, j int) bool {
		return ledger[i].Modified.After(ledger[j].Modified)
	})
}

// Calculate net and gross amount for line items
// Note: public for testing purposes
func CalcOrderAmounts(lineItems *[]models.OrderLineItem, userDiscount int) (decimal.Decimal, decimal.Decimal) {
	netAmount := decimal.Zero
	grossAmount := decimal.Zero
	for _, line := range *lineItems {
		itemQty := decimal.New(line.Quantity, 0)
		if itemQty.Cmp(decimal.Zero) <= 0 {
			// skip negative/zero item qty
			continue
		}

		lineAmount := line.Item.Price.Mul(itemQty)
		grossAmount = grossAmount.Add(lineAmount)
		itemDiscount := 0
		if line.Item.DiscountPercentage == 0 {
			itemDiscount = line.Item.SKU.DiscountPercentage
		}
		discount := line.Item.Price.Mul(decimal.NewFromFloat(float64(itemDiscount) / 100))
		netAmount = netAmount.Add(lineAmount.Sub(discount))
	}

	userDiscountDec := netAmount.Mul(decimal.NewFromFloat(float64(userDiscount) / 100))
	netAmount = netAmount.Sub(userDiscountDec)
	return netAmount, grossAmount
}
