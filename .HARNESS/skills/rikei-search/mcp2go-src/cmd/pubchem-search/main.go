package main

import (
	"fmt"
	"os"

	"mcp2go/rikei-search/pubchem"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: pubchem-search.exe <compound_name>")
		os.Exit(1)
	}
	res, err := pubchem.SearchPubChem(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}
