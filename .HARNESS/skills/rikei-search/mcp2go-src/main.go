package main

import (
	"flag"
	"fmt"
	"os"

	"mcp2go/rikei-search/ads"
	"mcp2go/rikei-search/aflow"
	"mcp2go/rikei-search/cod"
	"mcp2go/rikei-search/ieee"
	"mcp2go/rikei-search/mp"
	"mcp2go/rikei-search/pubchem"
)

func main() {
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(0)
	}

	cmd := os.Args[1]

	switch cmd {
	case "mp", "materials-project":
		runMP(os.Args[2:])
	case "aflow":
		runAFLOW(os.Args[2:])
	case "pubchem", "chem":
		runPubChem(os.Args[2:])
	case "ieee":
		runIEEE(os.Args[2:])
	case "ads":
		runADS(os.Args[2:])
	case "cod":
		runCOD(os.Args[2:])
	case "search", "--search", "-s":
		runAllSearch(os.Args[2:])
	case "--help", "-h", "help":
		printUsage()
	default:
		// Default to all search if query string is passed directly
		runAllSearch(os.Args[1:])
	}
}

func printUsage() {
	fmt.Println(`=======================================================
 🔬 Rikei-Search: 理系（自然科学・工学）統合API CLI
=======================================================
使用方法:
  rikei.exe <command> [query]

利用可能なサブコマンド:
  search <query>          全自然科学・工学データベース横断検索 (推奨)
  mp <formula>            Materials Project (結晶構造・バンドギャップ・生成熱)
  aflow <species/query>   AFLOW/AFLUX (結晶・材料熱力学DB)
  pubchem <name>          PubChem (分子量・SMILES・物理化学定数・分光リンク)
  ieee <query>            IEEE Xplore (工学・電子・情報・物理論文)
  ads <query>             NASA ADS (物理学・天文学・宇宙工学論文)
  cod <formula>           Crystallography Open Database (CIF結晶データ)

環境変数 (オプション):
  MP_API_KEY              Materials Project API Key
  IEEE_API_KEY            IEEE Xplore API Key
  ADS_API_TOKEN           NASA ADS API Bearer Token

使用例:
  rikei.exe search "LiFePO4"
  rikei.exe mp "Fe2O3"
  rikei.exe pubchem "aspirin"
  rikei.exe ieee "quantum computing"
  rikei.exe ads "black hole"`)
}

func runMP(args []string) {
	fs := flag.NewFlagSet("mp", flag.ExitOnError)
	fs.Parse(args)
	query := fs.Arg(0)
	if query == "" {
		fmt.Println("Error: 化学式または元素リストを指定してください。(例: rikei.exe mp LiFePO4)")
		os.Exit(1)
	}
	res, err := mp.SearchMP(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}

func runAFLOW(args []string) {
	fs := flag.NewFlagSet("aflow", flag.ExitOnError)
	fs.Parse(args)
	query := fs.Arg(0)
	if query == "" {
		fmt.Println("Error: 元素リストを指定してください。(例: rikei.exe aflow Ti,O)")
		os.Exit(1)
	}
	res, err := aflow.SearchAFLOW(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}

func runPubChem(args []string) {
	fs := flag.NewFlagSet("pubchem", flag.ExitOnError)
	fs.Parse(args)
	query := fs.Arg(0)
	if query == "" {
		fmt.Println("Error: 化合物名/SMILESを指定してください。(例: rikei.exe pubchem aspirin)")
		os.Exit(1)
	}
	res, err := pubchem.SearchPubChem(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}

func runIEEE(args []string) {
	fs := flag.NewFlagSet("ieee", flag.ExitOnError)
	fs.Parse(args)
	query := fs.Arg(0)
	if query == "" {
		fmt.Println("Error: 検索キーワードを指定してください。(例: rikei.exe ieee \"semiconductor\")")
		os.Exit(1)
	}
	res, err := ieee.SearchIEEE(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}

func runADS(args []string) {
	fs := flag.NewFlagSet("ads", flag.ExitOnError)
	fs.Parse(args)
	query := fs.Arg(0)
	if query == "" {
		fmt.Println("Error: 検索キーワードを指定してください。(例: rikei.exe ads \"astronomy\")")
		os.Exit(1)
	}
	res, err := ads.SearchADS(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}

func runCOD(args []string) {
	fs := flag.NewFlagSet("cod", flag.ExitOnError)
	fs.Parse(args)
	query := fs.Arg(0)
	if query == "" {
		fmt.Println("Error: 化学式を指定してください。(例: rikei.exe cod TiO2)")
		os.Exit(1)
	}
	res, err := cod.SearchCOD(query)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(res)
}

func runAllSearch(args []string) {
	fs := flag.NewFlagSet("search", flag.ExitOnError)
	fs.Parse(args)
	query := fs.Arg(0)
	if query == "" {
		fmt.Println("Error: 検索クエリを指定してください。")
		os.Exit(1)
	}

	fmt.Printf("# 🔬 自然科学・工学統合検索レポート: '%s'\n\n", query)

	// PubChem (Chemistry & Physchem)
	if res, err := pubchem.SearchPubChem(query); err == nil {
		fmt.Println(res)
	}

	// Materials Project (Materials)
	if res, err := mp.SearchMP(query); err == nil {
		fmt.Println(res)
	}

	// AFLOW (Materials)
	if res, err := aflow.SearchAFLOW(query); err == nil {
		fmt.Println(res)
	}

	// COD (Crystallography)
	if res, err := cod.SearchCOD(query); err == nil {
		fmt.Println(res)
	}

	// IEEE Xplore (Engineering)
	if res, err := ieee.SearchIEEE(query); err == nil {
		fmt.Println(res)
	}

	// NASA ADS (Physics & Astronomy)
	if res, err := ads.SearchADS(query); err == nil {
		fmt.Println(res)
	}
}
