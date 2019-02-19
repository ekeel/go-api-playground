package main

import (
	"encoding/json"
	"fmt"
)

func main() {
	ret, _ := json.Marshal([]string{"1111", "2222", "3333", "4444"})

	fmt.Print(string(ret))
}
