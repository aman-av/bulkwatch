package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"
)

// EnrichedDealsFile represents the enriched output file
type EnrichedDealsFile struct {
	Data    []map[string]interface{} `json:"data,omitempty"`  // For NSE
	Table   []map[string]interface{} `json:"Table,omitempty"` // For BSE
	Summary DealsSummary             `json:"summary"`
}

// DealsSummary provides statistics
type DealsSummary struct {
	PromoterDeals int `json:"promoterDeals"`
	FIIDeals      int `json:"fiiDeals"`
	DIIDeals      int `json:"diiDeals"`
	PublicDeals   int `json:"publicDeals"`
}

// enrichDealsFile enriches a deals file with classification data
func enrichDealsFile(exchange Exchange, dealType DealType) error {
	fmt.Printf("\n📝 Enriching %s %s deals...\n", exchange, dealType)
	
	// Read original file
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
	
	fmt.Printf("  Processing %d deals...\n", len(deals))
	
	// Group by symbol to fetch shareholding data once per symbol
	symbolDeals := make(map[string][]int) // symbol -> deal indices
	for i, deal := range deals {
		var symbol string
		if exchange == NSE {
			symbol = getString(deal, "BD_SYMBOL")
		} else {
			symbol = getString(deal, "scripname")
		}
		symbolDeals[symbol] = append(symbolDeals[symbol], i)
	}
	
	// Fetch shareholding data for unique symbols
	fmt.Printf("  Fetching shareholding data for %d unique symbols...\n", len(symbolDeals))
	shareholdingCache := make(map[string]*BSEShareholdingData)
	
	symbolCount := 0
	for symbol := range symbolDeals {
		symbolCount++
		scripCode := ""
		if exchange == NSE {
			scripCode = mapNSEToBSE(symbol)
		} else {
			// Get from first deal of this symbol
			firstDealIdx := symbolDeals[symbol][0]
			scripCode = fmt.Sprintf("%v", deals[firstDealIdx]["SCRIP_CODE"])
		}
		
		if scripCode != "" {
			fmt.Printf("    [%d/%d] %s (ScripCode: %s)... ", symbolCount, len(symbolDeals), symbol, scripCode)
			bseData, err := fetchBSEShareholding(scripCode, "129", "Mar-26")
			if err != nil {
				fmt.Printf("❌\n")
				shareholdingCache[symbol] = nil
			} else {
				fmt.Printf("✓\n")
				shareholdingCache[symbol] = bseData
			}
		}
	}
	
	// Enrich each deal by adding fields directly
	summary := DealsSummary{}
	
	for _, deal := range deals {
		var symbol, clientName string
		if exchange == NSE {
			symbol = getString(deal, "BD_SYMBOL")
			clientName = getString(deal, "BD_CLIENT_NAME")
		} else {
			symbol = getString(deal, "scripname")
			clientName = getString(deal, "CLIENT_NAME")
		}
		
		bseData := shareholdingCache[symbol]
		clientType, source, pct := classifyUsingBSEData(clientName, bseData)
		
		// Remove emoji from clientType for JSON
		cleanType := clientType
		switch clientType {
		case "🏢 Promoter":
			cleanType = "Promoter"
			summary.PromoterDeals++
		case "🌍 FII":
			cleanType = "FII"
			summary.FIIDeals++
		case "🏦 DII":
			cleanType = "DII"
			summary.DIIDeals++
		default:
			cleanType = "Public"
			summary.PublicDeals++
		}
		
		// Add new fields directly to the deal
		deal["clientType"] = cleanType
		deal["classificationSource"] = source
		if pct > 0 {
			deal["holdingPercentage"] = pct
		}
	}
	
	// Create enriched file with same structure as original
	enrichedFile := EnrichedDealsFile{
		Summary: summary,
	}
	
	if exchange == NSE {
		enrichedFile.Data = deals
	} else {
		enrichedFile.Table = deals
	}
	
	// Save enriched file
	enrichedFilename := fmt.Sprintf("archive/%s/%s/%s_%s_enriched_%s.json",
		exchange, dealType, exchange, dealType, time.Now().Format("2006-01-02"))
	
	enrichedData, err := json.MarshalIndent(enrichedFile, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal enriched data: %w", err)
	}
	
	if err := os.WriteFile(enrichedFilename, enrichedData, 0644); err != nil {
		return fmt.Errorf("failed to write enriched file: %w", err)
	}
	
	fmt.Printf("\n  ✅ Saved enriched file: %s\n", enrichedFilename)
	fmt.Printf("  📊 Summary: Promoter=%d | FII=%d | DII=%d | Public=%d\n",
		summary.PromoterDeals, summary.FIIDeals, summary.DIIDeals, summary.PublicDeals)
	
	return nil
}

// enrichAllDeals enriches all deal files
func enrichAllDeals() error {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("📊 ENRICHING DEALS WITH SHAREHOLDING CLASSIFICATION")
	fmt.Println(strings.Repeat("=", 80))
	
	configs := []struct {
		Exchange Exchange
		DealType DealType
	}{
		{NSE, BlockDeal},
		{NSE, BulkDeal},
		{BSE, BlockDeal},
		{BSE, BulkDeal},
	}
	
	for _, config := range configs {
		if err := enrichDealsFile(config.Exchange, config.DealType); err != nil {
			fmt.Printf("❌ Error enriching %s %s: %v\n", config.Exchange, config.DealType, err)
		}
	}
	
	fmt.Println("\n✅ Enrichment complete!")
	return nil
}

// Made with Bob
