package main

import (
	"fmt"
	"os"

	"mcp2go/rikei-search/ieee"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ieee-search.exe <query>")
		os.Exit(1)
	}
	res, err := ieee.SearchIEEE(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}
