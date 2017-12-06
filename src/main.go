package main

import (
	"controllers"
	"fmt"
	"gopkg.in/dixonwille/wmenu.v4"
)

func main() {
	fmt.Println("Welcome to the toy store!")

	// Replenish stock at beginning
	controllers.Cli.ReplenishStock()

	// start looping main menu
	for {
		menu := controllers.Cli.MainMenu()
		err := menu.Run()
		if err != nil {
			if wmenu.IsInvalidErr(err) {
				fmt.Println("Bad choice hombre... " + err.Error())
			} else {
				panic(fmt.Sprintf("error creating cli menu... %s", err))
			}
		}
	}
}
