package controllers

import (
	"bufio"
	"error"
	"fmt"
	"github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
	"gopkg.in/dixonwille/wmenu.v4"
	"math/rand"
	"models"
	"os"
	"strconv"
	"strings"
	"time"
	"usecases"
)

type CliController struct {
}

var uRepo = new(usecases.InventoryUsecaseRepository)
var Cli = new(CliController)
var fakeModels = new(models.Mocks)

// main menu action handler
func MainMenuAction(opts []wmenu.Opt) error {
	for _, opt := range opts {
		switch opt.ID {
		case 0:
			Cli.ReplenishStock()
		case 1:
			// clear stale line items
			fakeModels.ClearMockModelCache()
			Cli.UserMenu()
		case 2:
			Cli.SalesSummary()
		case 3:
			Cli.InventoryStatus()
		case 4:
			fmt.Println("Bye!")
			os.Exit(0)
		default:
			fmt.Println("Bad choice hombre...")
		}

	}
	return nil
}

func (c *CliController) UserMenu() {
	for {
		menu := wmenu.NewMenu("Choose a user who is placing order > ")
		menu.Action(UserMenuAction)
		for _, customer := range fakeModels.Customers {
			msg := fmt.Sprintf("%s Discount: %d%%", customer.Name, customer.DiscountPercentage)
			menu.Option(msg, customer.Id, false, nil)
		}

		for _, employee := range fakeModels.Employees {
			msg := fmt.Sprintf("%s (Employee) Discount: %d%%", employee.Name, employee.DiscountPercentage)
			menu.Option(msg, employee.Id, false, nil)
		}

		err := menu.Run()
		if err != nil {
			e, ok := err.(errors.ApplicationError)
			if ok && e.ErrorType == errors.ErrorMap[errors.PurchaseDoneBreak].ErrorType {
				// we are done with this menu
				return
			} else if wmenu.IsInvalidErr(err) {
				fmt.Println("Bad choice hombre... " + err.Error())
			} else {
				panic(fmt.Sprintf("error in user cli menu... %s", err))
			}
		}
	}
}

func UserMenuAction(opts []wmenu.Opt) error {
	for _, opt := range opts {
		optId := opt.Value.(uuid.UUID)

		if strings.Contains(opt.Text, "Employee") {
			for _, e := range fakeModels.Employees {
				if uuid.Equal(e.Id, optId) {
					fakeModels.PurchaseUserId = optId
					fakeModels.PurchaseUserDiscount = e.DiscountPercentage
					break
				}
			}
		} else {
			for _, c := range fakeModels.Customers {
				if uuid.Equal(c.Id, optId) {
					fakeModels.PurchaseUserId = c.Id
					fakeModels.PurchaseUserDiscount = c.DiscountPercentage
					break
				}
			}
		}
	}

	// now go to purchase menu
	Cli.PurchaseMenu()

	return errors.NewError(errors.PurchaseDoneBreak, "Done with purchase order")
}

// purchase menu action handler
func PurchaseMenuAction(opts []wmenu.Opt) error {
	for _, opt := range opts {
		optId := opt.Value.(uuid.UUID)
		if uuid.Equal(optId, uuid.Nil) {
			return errors.NewError(errors.PurchaseDoneBreak, "Done with purchase order")
		}
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("Enter quantity: ")
		qtyS, _ := reader.ReadString('\n')
		qtyS = strings.TrimSpace(qtyS)
		qty, _ := strconv.ParseInt(qtyS, 10, 64)

		Cli.AddToPurchaseOrder(optId, qty)
	}
	return nil
}

func (c *CliController) AddToPurchaseOrder(itemId uuid.UUID, qty int64) {
	for _, item := range fakeModels.Items {
		if uuid.Equal(item.Id, itemId) {
			lineItem := new(models.OrderLineItem)
			lineItem.Item = &item
			lineItem.Quantity = qty
			fakeModels.LineItems = append(fakeModels.LineItems, *lineItem)
			break
		}
	}
}

func (c *CliController) InventoryStatus() {
	inventory := uRepo.InventorySummary(time.Now().UTC())
	printInventory(inventory)
}

func (c *CliController) SalesSummary() {
	inventory, total := uRepo.SaleSummary(time.Now().AddDate(0, 0, -1).UTC())

	fmt.Println("*** Itemwise sales today so far ***")
	printInventory(inventory)
	fmt.Println("*** Total sales today so far ***")
	fmt.Println(total.StringFixedCash(5))
}

func printInventory(inventory map[uuid.UUID]decimal.Decimal) {
	for k, v := range inventory {
		var item models.Item
		for i := 0; i < len(fakeModels.Items); i++ {
			if uuid.Equal(k, fakeModels.Items[i].Id) {
				item = fakeModels.Items[i]
				fmt.Printf("%s(Id %s, SKU: %s, Price %s) : %s\n",
					item.Name, item.Id, item.SkuId, item.Price.StringFixedCash(5), v.String())
			}
		}
	}
}

func (c *CliController) PurchaseMenu() {
	// loop purchase menu until user breaks out
	for {
		menu := wmenu.NewMenu("Choose items to purchase > ")
		menu.Action(PurchaseMenuAction)
		menu.Option("Done with purchase", uuid.Nil, false, nil)
		for _, item := range fakeModels.Items {
			menu.Option(item.Name, item.Id, false, nil)
		}
		err := menu.Run()
		if err != nil {
			e, ok := err.(errors.ApplicationError)
			if ok && e.ErrorType == errors.ErrorMap[errors.PurchaseDoneBreak].ErrorType {
				c.PlaceOrder()
				// we are done with this menu
				return
			} else if wmenu.IsInvalidErr(err) {
				fmt.Println("Bad choice hombre... " + err.Error())
			} else {
				panic(fmt.Sprintf("error in purchase cli menu... %s", err))
			}
		}
	}
}

func (c *CliController) PlaceOrder() {
	amt, err := uRepo.Purchase(&fakeModels.LineItems, fakeModels.PurchaseUserId, fakeModels.PurchaseUserDiscount)

	if err != nil {
		fmt.Println("Purchase failed!, retry again later. Reason: " + err.Error())
	} else {
		fmt.Println("Thanks for placing order! You need to pay " + amt.StringFixedCash(5))
	}
}

// main menu Cli
func (c *CliController) MainMenu() *wmenu.Menu {
	menu := wmenu.NewMenu("What do you want to do > ")
	menu.Action(MainMenuAction)
	menu.Option("Replenish stock again", nil, true, nil)
	menu.Option("Purchase", nil, false, nil)
	menu.Option("Today's sales summary", nil, false, nil)
	menu.Option("Inventory status", nil, false, nil)
	menu.Option("Exit", nil, false, nil)

	return menu
}

// create stock with dummy items
func (c *CliController) ReplenishStock() {
	fakeModels.InitInventory()
	fakeModels.InitUsers()

	for _, item := range fakeModels.Items {
		// add to inventory
		rand.NewSource(time.Now().Unix())
		ok, err := uRepo.Replenish(item, decimal.New(rand.Int63n(10), 0))
		if !ok || err != nil {
			fmt.Println("Replenish failed, try again later...")
			break
		}
	}

	fmt.Println("Inventory replenished.")
}
