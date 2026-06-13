package main

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
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

	// Fetch shareholding data for unique symbols in parallel
	fmt.Printf("  Fetching shareholding data for %d unique symbols (parallel with 1000 workers)...\n", len(symbolDeals))
	shareholdingCache := make(map[string]*BSEShareholdingData)
	var cacheMutex sync.RWMutex

	// Create work queue
	type symbolWork struct {
		symbol    string
		scripCode string
		index     int
		total     int
	}

	workQueue := make(chan symbolWork, len(symbolDeals))

	// Populate work queue
	symbolIndex := 0
	for symbol := range symbolDeals {
		symbolIndex++
		scripCode := ""
		if exchange == NSE {
			scripCode = mapNSEToBSE(symbol)
		} else {
			// Get from first deal of this symbol
			firstDealIdx := symbolDeals[symbol][0]
			scripCode = fmt.Sprintf("%v", deals[firstDealIdx]["SCRIP_CODE"])
		}

		if scripCode != "" {
			workQueue <- symbolWork{
				symbol:    symbol,
				scripCode: scripCode,
				index:     symbolIndex,
				total:     len(symbolDeals),
			}
		}
	}
	close(workQueue)

	// Worker pool - optimal size based on I/O-bound operations
	// For HTTP requests, 50-100 workers is typically optimal
	// Too many workers can overwhelm the server or cause rate limiting
	maxWorkers := 100
	if len(symbolDeals) < 100 {
		maxWorkers = len(symbolDeals) // Don't create more workers than tasks
	}

	var wg sync.WaitGroup

	fmt.Printf("  Using %d parallel workers\n", maxWorkers)

	// Progress tracking
	var processedCount int
	var countMutex sync.Mutex

	// Start workers
	for i := 0; i < maxWorkers; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()

			for work := range workQueue {
				bseData, err := fetchBSEShareholding(work.scripCode, "129", "Mar-26")

				cacheMutex.Lock()
				if err != nil {
					shareholdingCache[work.symbol] = nil
				} else {
					shareholdingCache[work.symbol] = bseData
				}
				cacheMutex.Unlock()

				// Update progress
				countMutex.Lock()
				processedCount++
				if processedCount%10 == 0 || processedCount == work.total {
					fmt.Printf("    Progress: %d/%d symbols processed (%.1f%%)\n",
						processedCount, work.total, float64(processedCount)/float64(work.total)*100)
				}
				countMutex.Unlock()
			}
		}(i)
	}

	// Wait for all workers to complete
	wg.Wait()
	fmt.Printf("  ✅ Completed fetching shareholding data for all symbols\n")

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
