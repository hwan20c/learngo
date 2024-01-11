package main

import (
	"fmt"
	"name/accounts"
)

func main() {
	account := accounts.NewAccount("nico")
	fmt.Println(account)
}
