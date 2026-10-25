package main

import (
	"fmt"
	"os"

	"mcp2go/rikei-search/mp"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: mp-search.exe <formula_or_elements>")
		os.Exit(1)
	}
	res, err := mp.SearchMP(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}
