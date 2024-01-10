package main

import "fmt"

func main() {
	nico := map[string]string{"name": "nico", "age": "12"}
	for key, value := range nico {
		fmt.Println(key, value)
	} // key나 value를 ignore 할 수 있다.
}
