//go:build ignore
// +build ignore

// Standalone script to generate mappings with actual market cap data
// Run: go run run_marketcap.go

package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/xuri/excelize/v2"
)

// Import constants from parent package
// Since this is a standalone script, we'll define them here
const (
	NSEEquityListURL1 = "https://archives.nseindia.com/content/equities/EQUITY_L.csv"
	NSEEquityListURL2 = "https://www1.nseindia.com/content/equities/EQUITY_L.csv"
	NSEMarketCapURL   = "https://nsearchives.nseindia.com/web/sites/default/files/inline-files/MCAP31032021_2.xlsx"
	BSEScripListAPI   = "https://api.bseindia.com/BseIndiaAPI/api/ListofScripData/w?Group=&Scripcode=&industry=&segment=Equity&status=Active"
	BSEReferer        = "https://www.bseindia.com/"
)

type NSEEquity struct {
	Symbol      string
	CompanyName string
	ISIN        string
}

type BSEEquity struct {
	ScripCode   string
	CompanyName string
	ISIN        string
}

type SymbolMapping struct {
	NSESymbol    string  `json:"nseSymbol"`
	BSEScripCode string  `json:"bseScripCode"`
	CompanyName  string  `json:"companyName"`
	ISIN         string  `json:"isin"`
	MarketCap    float64 `json:"marketCap,omitempty"`
	MarketCapCr  string  `json:"marketCapCr,omitempty"`
	Rank         int     `json:"rank,omitempty"`
}

func fetchNSESymbols() ([]NSEEquity, error) {
	fmt.Println("📊 Fetching NSE symbols...")
	
	urls := []string{
		NSEEquityListURL1,
		NSEEquityListURL2,
	}
	
	client := &http.Client{Timeout: 30 * time.Second}
	
	var resp *http.Response
	var err error
	
	for _, url := range urls {
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("User-Agent", "Mozilla/5.0")
		req.Header.Set("Accept", "text/csv")
		
		resp, err = client.Do(req)
		if err == nil && resp.StatusCode == 200 {
			break
		}
		if resp != nil {
			resp.Body.Close()
		}
	}
	
	if err != nil || resp == nil {
		return nil, fmt.Errorf("failed to fetch NSE data")
	}
	defer resp.Body.Close()
	
	reader := csv.NewReader(resp.Body)
	records, _ := reader.ReadAll()
	
	var equities []NSEEquity
	for i, record := range records {
		if i == 0 || len(record) < 7 {
			continue
		}
		equities = append(equities, NSEEquity{
			Symbol:      strings.TrimSpace(record[0]),
			CompanyName: strings.TrimSpace(record[1]),
			ISIN:        strings.TrimSpace(record[6]),
		})
	}
	
	fmt.Printf("✓ Fetched %d NSE symbols\n", len(equities))
	return equities, nil
}

func fetchBSESymbols() ([]BSEEquity, error) {
	fmt.Println("📊 Fetching BSE symbols...")
	
	url := BSEScripListAPI
	
	client := &http.Client{Timeout: 30 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Referer", BSEReferer)
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	body, _ := io.ReadAll(resp.Body)
	
	var rawData []map[string]interface{}
	json.Unmarshal(body, &rawData)
	
	var equities []BSEEquity
	for _, item := range rawData {
		equities = append(equities, BSEEquity{
			ScripCode:   fmt.Sprintf("%v", item["SCRIP_CD"]),
			CompanyName: fmt.Sprintf("%v", item["SCRIP_NAME"]),
			ISIN:        fmt.Sprintf("%v", item["ISIN_NUMBER"]),
		})
	}
	
	fmt.Printf("✓ Fetched %d BSE symbols\n", len(equities))
	return equities, nil
}

func fetchMarketCap() (map[string]float64, error) {
	fmt.Println("💰 Fetching market cap from NSE Excel...")
	
	url := NSEMarketCapURL
	
	client := &http.Client{Timeout: 60 * time.Second}
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Set("User-Agent", "Mozilla/5.0")
	
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	
	tmpFile, _ := os.CreateTemp("", "mcap-*.xlsx")
	defer os.Remove(tmpFile.Name())
	
	io.Copy(tmpFile, resp.Body)
	tmpFile.Close()
	
	f, _ := excelize.OpenFile(tmpFile.Name())
	defer f.Close()
	
	sheets := f.GetSheetList()
	rows, _ := f.GetRows(sheets[0])
	
	marketCap := make(map[string]float64)
	
	for i, row := range rows {
		if i == 0 || len(row) < 4 {
			continue
		}
		
		symbol := strings.TrimSpace(row[1])
		mcapStr := strings.ReplaceAll(strings.TrimSpace(row[3]), ",", "")
		
		if mcapLakhs, err := strconv.ParseFloat(mcapStr, 64); err == nil && mcapLakhs > 0 {
			marketCap[symbol] = mcapLakhs / 100 // Convert to Crores
		}
	}
	
	fmt.Printf("✓ Fetched market cap for %d symbols\n", len(marketCap))
	return marketCap, nil
}

func main() {
	fmt.Println("🚀 NSE-BSE Mapping with Actual Market Cap")
	fmt.Println("==========================================\n")
	
	nse, _ := fetchNSESymbols()
	bse, _ := fetchBSESymbols()
	mcap, _ := fetchMarketCap()
	
	// Create ISIN map
	isinToBSE := make(map[string]BSEEquity)
	for _, b := range bse {
		if b.ISIN != "" && b.ISIN != "<nil>" {
			isinToBSE[b.ISIN] = b
		}
	}
	
	// Create mappings
	var mappings []SymbolMapping
	for _, n := range nse {
		if n.ISIN == "" || n.ISIN == "<nil>" {
			continue
		}
		
		if b, found := isinToBSE[n.ISIN]; found {
			m := SymbolMapping{
				NSESymbol:    n.Symbol,
				BSEScripCode: b.ScripCode,
				CompanyName:  n.CompanyName,
				ISIN:         n.ISIN,
			}
			
			if mc, ok := mcap[n.Symbol]; ok {
				m.MarketCap = mc
				if mc >= 100000 {
					m.MarketCapCr = fmt.Sprintf("₹%.2f Lakh Cr", mc/100000)
				} else {
					m.MarketCapCr = fmt.Sprintf("₹%.2f Cr", mc)
				}
			}
			
			mappings = append(mappings, m)
		}
	}
	
	// Sort by market cap
	sort.Slice(mappings, func(i, j int) bool {
		if mappings[i].MarketCap > 0 && mappings[j].MarketCap > 0 {
			return mappings[i].MarketCap > mappings[j].MarketCap
		}
		if mappings[i].MarketCap > 0 {
			return true
		}
		if mappings[j].MarketCap > 0 {
			return false
		}
		return mappings[i].NSESymbol < mappings[j].NSESymbol
	})
	
	// Add ranks
	for i := range mappings {
		mappings[i].Rank = i + 1
	}
	
	// Top 1000
	top1000 := mappings
	if len(mappings) > 1000 {
		top1000 = mappings[:1000]
	}
	
	// Save
	data, _ := json.MarshalIndent(top1000, "", "  ")
	os.WriteFile("top_1000_marketcap.json", data, 0644)
	
	// CSV
	f, _ := os.Create("top_1000_marketcap.csv")
	w := csv.NewWriter(f)
	w.Write([]string{"Rank", "NSE_Symbol", "BSE_ScripCode", "Company", "Market_Cap_Cr", "ISIN"})
	for _, m := range top1000 {
		w.Write([]string{
			fmt.Sprintf("%d", m.Rank),
			m.NSESymbol,
			m.BSEScripCode,
			m.CompanyName,
			m.MarketCapCr,
			m.ISIN,
		})
	}
	w.Flush()
	f.Close()
	
	withMC := 0
	for _, m := range top1000 {
		if m.MarketCap > 0 {
			withMC++
		}
	}
	
	fmt.Printf("\n✅ Complete!\n")
	fmt.Printf("   Total Mappings: %d\n", len(mappings))
	fmt.Printf("   Top 1000: %d\n", len(top1000))
	fmt.Printf("   With Market Cap: %d\n", withMC)
	fmt.Printf("\n📁 Files: top_1000_marketcap.json, top_1000_marketcap.csv\n")
}
