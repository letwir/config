package main

import (
	"fmt"
	"os"

	"mcp2go/rikei-search/aflow"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: aflow-search.exe <species_or_query>")
		os.Exit(1)
	}
	res, err := aflow.SearchAFLOW(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}
