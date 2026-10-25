package pubchem

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type PubChemPropertyItem struct {
	CID             int     `json:"CID"`
	MolecularFormula string  `json:"MolecularFormula"`
	MolecularWeight  interface{} `json:"MolecularWeight"`
	IUPACName        string  `json:"IUPACName"`
	CanonicalSMILES  string  `json:"CanonicalSMILES"`
	InChIKey         string  `json:"InChIKey"`
}

type PubChemPropertyTable struct {
	Properties []PubChemPropertyItem `json:"Properties"`
}

type PubChemResponse struct {
	PropertyTable PubChemPropertyTable `json:"PropertyTable"`
}

func SearchPubChem(query string) (string, error) {
	reqURL := fmt.Sprintf("https://pubchem.ncbi.nlm.nih.gov/rest/pug/compound/name/%s/property/IUPACName,MolecularFormula,MolecularWeight,CanonicalSMILES,InChIKey/JSON", url.PathEscape(query))

	client := &http.Client{Timeout: 15 * time.Second}
	resp, err := client.Get(reqURL)
	if err != nil {
		return "", fmt.Errorf("PubChem API Request Failed: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		return fmt.Sprintf("🔍 PubChem: 化合物 '%s' の直接ヒットはありませんでしたわ。\n(Status: %d, Message: %s)", query, resp.StatusCode, string(body)), nil
	}

	var pcResp PubChemResponse
	if err := json.Unmarshal(body, &pcResp); err != nil {
		return "", fmt.Errorf("Failed to parse PubChem JSON: %w", err)
	}

	props := pcResp.PropertyTable.Properties
	if len(props) == 0 {
		return fmt.Sprintf("🔍 PubChem: '%s' に該当するプロパティは見つかりませんでしたわ。", query), nil
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("### ⚗️ PubChem (物理化学物性) 検索結果 (Query: %s)\n\n", query))
	for i, item := range props {
		sb.WriteString(fmt.Sprintf("#### %d. CID: `%d` (%s)\n", i+1, item.CID, query))
		if item.IUPACName != "" {
			sb.WriteString(fmt.Sprintf("- **IUPAC名**: `%s`\n", item.IUPACName))
		}
		sb.WriteString(fmt.Sprintf("- **化学式 (Formula)**: `%s`\n", item.MolecularFormula))
		sb.WriteString(fmt.Sprintf("- **分子量 (MW)**: `%v g/mol`\n", item.MolecularWeight))
		sb.WriteString(fmt.Sprintf("- **SMILES**: `%s`\n", item.CanonicalSMILES))
		sb.WriteString(fmt.Sprintf("- **InChIKey**: `%s`\n", item.InChIKey))
		sb.WriteString(fmt.Sprintf("- **詳細ページ**: [https://pubchem.ncbi.nlm.nih.gov/compound/%d](https://pubchem.ncbi.nlm.nih.gov/compound/%d)\n", item.CID, item.CID))
		sb.WriteString(fmt.Sprintf("- **分光・NMR・IRデータ (SDBS代替)**: [PubChem Spectras %d](https://pubchem.ncbi.nlm.nih.gov/compound/%d#section=Spectra)\n\n", item.CID, item.CID))
	}

	return sb.String(), nil
}
