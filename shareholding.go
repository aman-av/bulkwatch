package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"
)

// ShareholdingPattern represents the shareholding pattern data
type ShareholdingPattern struct {
	ScripCode   string  `json:"scripCode"`
	CompanyName string  `json:"companyName"`
	Quarter     string  `json:"quarter"`
	Promoters   float64 `json:"promoters"`
	FII         float64 `json:"fii"`
	DII         float64 `json:"dii"`
	Public      float64 `json:"public"`
}

// DealWithContext represents a deal with shareholding context
type DealWithContext struct {
	Symbol       string  `json:"symbol"`
	ClientName   string  `json:"clientName"`
	DealType     string  `json:"dealType"`
	Quantity     float64 `json:"quantity"`
	Price        float64 `json:"price"`
	IsPromoter   bool    `json:"isPromoter"`
	IsFII        bool    `json:"isFII"`
	IsDII        bool    `json:"isDII"`
	VerifyURL    string  `json:"verifyUrl"`
}

var (
	symbolMappingCache map[string]string
	symbolMappingOnce  sync.Once
)

// loadSymbolMapping loads the NSE to BSE mapping from top_1000_marketcap.json
func loadSymbolMapping() map[string]string {
	symbolMappingOnce.Do(func() {
		symbolMappingCache = make(map[string]string)
		
		data, err := os.ReadFile("top_1000_marketcap.json")
		if err != nil {
			fmt.Printf("⚠️  Warning: Could not load top_1000_marketcap.json: %v\n", err)
			return
		}
		
		// Parse as array of generic maps to avoid type conflicts
		var mappings []map[string]interface{}
		if err := json.Unmarshal(data, &mappings); err != nil {
			fmt.Printf("⚠️  Warning: Could not parse top_1000_marketcap.json: %v\n", err)
			return
		}
		
		for _, mapping := range mappings {
			nseSymbol, _ := mapping["nseSymbol"].(string)
			bseScripCode, _ := mapping["bseScripCode"].(string)
			
			if nseSymbol != "" && bseScripCode != "" {
				symbolMappingCache[nseSymbol] = bseScripCode
			}
		}
		
		fmt.Printf("✅ Loaded %d NSE→BSE symbol mappings\n", len(symbolMappingCache))
	})
	
	return symbolMappingCache
}

// mapNSEToBSE maps NSE symbol to BSE script code using top_1000_marketcap.json
func mapNSEToBSE(nseSymbol string) string {
	mappings := loadSymbolMapping()
	return mappings[nseSymbol]
}

// classifyClient identifies if a client is Promoter/FII/DII
// Now uses investor registry for improved accuracy
func classifyClient(clientName string) (isPromoter, isFII, isDII bool) {
	return classifyClientWithRegistry(clientName)
}

// analyzeDealsWithBSE analyzes deals using actual BSE shareholding data
func analyzeDealsWithBSE(exchange Exchange, dealType DealType, limit int) error {
	fmt.Printf("\n🔍 Analyzing %s %s Deals with BSE Shareholding Data (First %d records)\n", exchange, dealType, limit)
	fmt.Println(strings.Repeat("=", 80))
	
	// Read the deal data
	filename := fmt.Sprintf("archive/%s/%s/%s_%s_%s_to_%s.json",
		exchange, dealType, exchange, dealType,
		time.Now().Format("2006-01-02"), time.Now().Format("2006-01-02"))
	
	data, err := os.ReadFile(filename)
	if err != nil {
		return fmt.Errorf("failed to read file: %w", err)
	}
	
	var deals []map[string]interface{}
	if exchange == NSE {
		var nseData NSEDealData
		if err := json.Unmarshal(data, &nseData); err != nil {
			return fmt.Errorf("failed to parse NSE data: %w", err)
		}
		deals = nseData.Data
	} else {
		var bseData BSEDealData
		if err := json.Unmarshal(data, &bseData); err != nil {
			return fmt.Errorf("failed to parse BSE data: %w", err)
		}
		deals = bseData.Table
	}
	
	// Limit to first N records
	if len(deals) > limit {
		deals = deals[:limit]
	}
	
	fmt.Printf("📊 Processing %d deals...\n", len(deals))
	
	// Group by symbol
	symbolDeals := make(map[string][]map[string]interface{})
	for _, deal := range deals {
		var symbol string
		if exchange == NSE {
			symbol = getString(deal, "BD_SYMBOL")
		} else {
			symbol = getString(deal, "scripname")
		}
		symbolDeals[symbol] = append(symbolDeals[symbol], deal)
	}
	
	// Fetch BSE shareholding data for each symbol
	fmt.Println("\n📥 Fetching BSE shareholding patterns...")
	shareholdingCache := make(map[string]*BSEShareholdingData)
	
	for symbol := range symbolDeals {
		scripCode := ""
		if exchange == NSE {
			scripCode = mapNSEToBSE(symbol)
		} else {
			if len(symbolDeals[symbol]) > 0 {
				scripCode = fmt.Sprintf("%v", symbolDeals[symbol][0]["SCRIP_CODE"])
			}
		}
		
		if scripCode != "" {
			fmt.Printf("  Fetching %s (ScripCode: %s)... ", symbol, scripCode)
			bseData, err := fetchBSEShareholding(scripCode, "129", "Mar-26")
			if err != nil {
				fmt.Printf("❌ Failed: %v\n", err)
				shareholdingCache[symbol] = nil
			} else {
				fmt.Printf("✓ Promoter=%.2f%% FII=%.2f%% DII=%.2f%%\n",
					bseData.PromoterTotal, bseData.FIITotal, bseData.DIITotal)
				shareholdingCache[symbol] = bseData
			}
		}
	}
	
	fmt.Println("\n📊 Deal Analysis with BSE Data:")
	fmt.Println(strings.Repeat("=", 80))
	
	// Analyze each symbol's deals
	for symbol, deals := range symbolDeals {
		bseData := shareholdingCache[symbol]
		
		fmt.Printf("\n📈 %s (%d deals)\n", symbol, len(deals))
		if bseData != nil {
			fmt.Printf("   BSE Holdings: Promoter=%.2f%% | FII=%.2f%% | DII=%.2f%% | Public=%.2f%%\n",
				bseData.PromoterTotal, bseData.FIITotal, bseData.DIITotal, bseData.PublicTotal)
		}
		fmt.Println(strings.Repeat("-", 80))
		
		for i, deal := range deals {
			var clientName string
			var quantity, price float64
			
			if exchange == NSE {
				clientName = getString(deal, "BD_CLIENT_NAME")
				quantity = getFloat64(deal, "BD_QTY_TRD")
				price = getFloat64(deal, "BD_TP_WATP")
			} else {
				clientName = getString(deal, "CLIENT_NAME")
				quantity = getFloat64(deal, "QUANTITY")
				price = getFloat64(deal, "PRICE")
			}
			
			// Classify using BSE data
			category, confidence, pct := classifyUsingBSEData(clientName, bseData)
			
			fmt.Printf("  %d. %s - %s\n", i+1, category, clientName)
			fmt.Printf("     Qty: %.0f | Price: %.2f | Value: %.2fM\n",
				quantity, price, (quantity*price)/1000000)
			if pct > 0 {
				fmt.Printf("     Classification: %s (%.2f%% holding)\n", confidence, pct)
			} else {
				fmt.Printf("     Classification: %s\n", confidence)
			}
			fmt.Println()
		}
	}
	
	return nil
}

// Helper functions
func getString(m map[string]interface{}, key string) string {
	if v, ok := m[key]; ok {
		if s, ok := v.(string); ok {
			return s
		}
	}
	return ""
}

func getFloat64(m map[string]interface{}, key string) float64 {
	if v, ok := m[key]; ok {
		switch val := v.(type) {
		case float64:
			return val
		case int:
			return float64(val)
		case int64:
			return float64(val)
		}
	}
	return 0.0
}

// Made with Bob
