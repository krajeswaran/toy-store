package usecases

import (
	"fmt"
	"math/rand"
	"models"
	"reflect"
	"testing"
	"time"

	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

var testRepo = new(InventoryUsecaseRepository)

func TestInventoryUsecaseRepository_Replenish(t *testing.T) {
	type args struct {
		item  models.Item
		count decimal.Decimal
	}
	tests := []struct {
		name                       string
		InventoryUsecaseRepository *InventoryUsecaseRepository
		args                       args
		want                       bool
		wantErr                    bool
	}{
		{
			name: "Test Nil Item",
			InventoryUsecaseRepository: testRepo,
			args: args{
				item: models.Item{
					DiscountPercentage: 0,
					Name:               "Test Item",
					Description:        "",
					Price:              decimal.Decimal{},
					BaseFields: models.BaseFields{
						Id:       uuid.Nil,
						Created:  time.Time{},
						Modified: time.Time{},
						Status:   "",
					},
					SKU: models.SKU{
						SkuId:              uuid.UUID{},
						Name:               "",
						Description:        "",
						DiscountPercentage: 0,
					},
				},
				count: decimal.New(5, 0),
			},
			want:    false,
			wantErr: true,
		},
		{
			name: "Test Add Item",
			InventoryUsecaseRepository: testRepo,
			args: args{
				item: models.Item{
					DiscountPercentage: 0,
					Name:               "Test Item",
					Description:        "",
					Price:              decimal.New(10, 0),
					BaseFields: models.BaseFields{
						Id:       uuid.NewV4(),
						Created:  time.Now().UTC(),
						Modified: time.Now().UTC(),
						Status:   models.AvailableItemStatus,
					},
					SKU: models.SKU{
						SkuId:              uuid.NewV4(),
						Name:               "Test SKU",
						Description:        "",
						DiscountPercentage: 10,
					},
				},
				count: decimal.New(5, 0),
			},
			want:    true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.InventoryUsecaseRepository.Replenish(tt.args.item, tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("InventoryUsecaseRepository.Replenish() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InventoryUsecaseRepository.Replenish() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestInventoryUsecaseRepository_Purchase(t *testing.T) {
	// setup
	fM := new(models.Mocks)
	fM.InitInventory()
	fM.InitUsers()
	for _, item := range fM.Items {
		// add to inventory
		ok, err := testRepo.Replenish(item, decimal.New(rand.Int63n(10), 0))
		if !ok || err != nil {
			fmt.Println("Replenish failed, try again later...")
			break
		}
	}

	// init some test data
	singleLineItem := &[]models.OrderLineItem{
		{
			Item:     fM.GetMockedItem(1),
			Quantity: 2,
		},
	}
	singleItemPriceWant, _ := CalcOrderAmounts(singleLineItem, fM.GetMockedUser(1).DiscountPercentage)

	multiLineItems := &[]models.OrderLineItem{
		{
			Item:     fM.GetMockedItem(0),
			Quantity: 2,
		},
		{
			Item:     fM.GetMockedItem(1),
			Quantity: 1,
		},
		{
			Item:     fM.GetMockedItem(2),
			Quantity: 3,
		},
	}
	multiItemPriceWant, _ := CalcOrderAmounts(multiLineItems, fM.GetMockedUser(1).DiscountPercentage)

	type args struct {
		lineItems    *[]models.OrderLineItem
		userId       uuid.UUID
		userDiscount int
	}
	tests := []struct {
		name                       string
		InventoryUsecaseRepository *InventoryUsecaseRepository
		args                       args
		want                       decimal.Decimal
		wantErr                    bool
	}{
		{
			name: "Test Empty Order",
			InventoryUsecaseRepository: testRepo,
			args: args{
				lineItems:    &[]models.OrderLineItem{},
				userId:       uuid.UUID{},
				userDiscount: 0,
			},
			want:    decimal.Zero,
			wantErr: true,
		},
		{
			name: "Test single item purchase order",
			InventoryUsecaseRepository: testRepo,
			args: args{
				lineItems:    singleLineItem,
				userId:       fM.GetMockedUser(1).Id,
				userDiscount: fM.GetMockedUser(1).DiscountPercentage,
			},
			want:    singleItemPriceWant,
			wantErr: false,
		},
		{
			name: "Test multi item purchase order",
			InventoryUsecaseRepository: testRepo,
			args: args{
				lineItems:    multiLineItems,
				userId:       fM.GetMockedUser(1).Id,
				userDiscount: fM.GetMockedUser(1).DiscountPercentage,
			},
			want:    multiItemPriceWant,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := tt.InventoryUsecaseRepository.Purchase(tt.args.lineItems, tt.args.userId, tt.args.userDiscount)
			if (err != nil) != tt.wantErr {
				t.Errorf("InventoryUsecaseRepository.Purchase() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("InventoryUsecaseRepository.Purchase() = %v, want %v", got, tt.want)
			}
		})
	}

	// teardown
	fM.ClearMockModelCache()
}
