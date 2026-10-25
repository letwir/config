package main

import (
	"fmt"
	"os"

	"mcp2go/rikei-search/ads"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: ads-search.exe <query>")
		os.Exit(1)
	}
	res, err := ads.SearchADS(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}
