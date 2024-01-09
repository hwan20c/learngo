package main

import (
	"fmt"
	"strings"
)

func lenAndUpper(name string) (length int, upppercase string) {
	defer fmt.Println("I'm done.")
	length = len(name)
	upppercase = strings.ToUpper(name)
	return
}

func main() {
	totalLength, up := lenAndUpper("nico")
	fmt.Println(totalLength, up)
}
