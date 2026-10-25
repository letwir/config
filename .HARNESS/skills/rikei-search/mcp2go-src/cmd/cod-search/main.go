package main

import (
	"fmt"
	"os"

	"mcp2go/rikei-search/cod"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Println("Usage: cod-search.exe <formula>")
		os.Exit(1)
	}
	res, err := cod.SearchCOD(os.Args[1])
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}
