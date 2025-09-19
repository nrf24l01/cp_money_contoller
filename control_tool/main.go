package main

import (
	"log"
	"os"

	"github.com/nrf24l01/cp_money_controller/control_tool/tasks"
)

func main() {
	argsWithoutProg := os.Args[1:]

	if len(argsWithoutProg) == 3 && argsWithoutProg[0] == "create-user" {
		user, password := argsWithoutProg[1], argsWithoutProg[2]
		if err := tasks.CreateUser(user, password); err != nil {
			panic(err)
		}
		log.Printf("User %s created successfully", user)
	}
}