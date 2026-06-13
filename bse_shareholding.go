package main

import (
	"context"
	"fmt"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/chromedp/chromedp"
)

// BSEShareholdingData represents parsed shareholding data from BSE
type BSEShareholdingData struct {
	ScripCode       string                    `json:"scripCode"`
	CompanyName     string                    `json:"companyName"`
	Quarter         string                    `json:"quarter"`
	Categories      map[string]CategoryHolding `json:"categories"`
	PromoterTotal   float64                   `json:"promoterTotal"`
	FIITotal        float64                   `json:"fiiTotal"`
	DIITotal        float64                   `json:"diiTotal"`
	PublicTotal     float64                   `json:"publicTotal"`
}

// CategoryHolding represents holdings for a specific category
type CategoryHolding struct {
	CategoryName string  `json:"categoryName"`
	Count        int     `json:"count"`
	Percentage   float64 `json:"percentage"`
	Shares       int64   `json:"shares"`
}

// fetchBSEShareholding fetches shareholding data from both promoter and public pages
func fetchBSEShareholding(scripCode, qtrID, qtrName string) (*BSEShareholdingData, error) {
	// Create Chrome context
	opts := append(chromedp.DefaultExecAllocatorOptions[:],
		chromedp.Flag("headless", true),
		chromedp.Flag("disable-gpu", true),
		chromedp.Flag("no-sandbox", true),
	)
	
	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(), opts...)
	defer cancel()
	
	ctx, cancel := chromedp.NewContext(allocCtx)
	defer cancel()
	
	// Set timeout
	ctx, cancel = context.WithTimeout(ctx, 45*time.Second)
	defer cancel()
	
	data := &BSEShareholdingData{
		ScripCode:  scripCode,
		Quarter:    qtrName,
		Categories: make(map[string]CategoryHolding),
	}
	
	// Fetch promoter holdings
	promoterURL := BuildBSEShareholdingPromoterURL(scripCode, qtrID, qtrName)
	
	var promoterHTML string
	err := chromedp.Run(ctx,
		chromedp.Navigate(promoterURL),
		chromedp.Sleep(4*time.Second),
		chromedp.OuterHTML("html", &promoterHTML),
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to fetch promoter page: %w", err)
	}
	
	// Parse promoter data
	parseBSEShareholding(promoterHTML, data)
	
	// Fetch public shareholder holdings (FII, DII, MF, etc.)
	publicURL := BuildBSEShareholdingPublicURL(scripCode, qtrID, qtrName)
	
	var publicHTML string
	err = chromedp.Run(ctx,
		chromedp.Navigate(publicURL),
		chromedp.Sleep(4*time.Second),
		chromedp.OuterHTML("html", &publicHTML),
	)
	
	if err != nil {
		return nil, fmt.Errorf("failed to fetch public shareholder page: %w", err)
	}
	
	// Parse public shareholder data
	parseBSEShareholding(publicHTML, data)
	
	return data, nil
}

// parseBSEShareholding parses HTML content to extract shareholding data
func parseBSEShareholding(html string, data *BSEShareholdingData) {
	// Extract company name if not already set
	if data.CompanyName == "" {
		companyRe := regexp.MustCompile(`<h4[^>]*>([^<]+)</h4>`)
		if matches := companyRe.FindStringSubmatch(html); len(matches) > 1 {
			data.CompanyName = strings.TrimSpace(matches[1])
		}
	}
	
	// Parse table rows
	tableRe := regexp.MustCompile(`<table[^>]*>(.*?)</table>`)
	tables := tableRe.FindAllStringSubmatch(html, -1)
	
	for _, table := range tables {
		parseShareholdingTable(table[1], data)
	}
}

// parseShareholdingTable extracts shareholding data from table HTML
func parseShareholdingTable(tableHTML string, data *BSEShareholdingData) {
	rowRe := regexp.MustCompile(`<tr[^>]*>(.*?)</tr>`)
	rows := rowRe.FindAllStringSubmatch(tableHTML, -1)
	
	for _, row := range rows {
		cells := extractCells(row[1])
		
		if len(cells) < 4 {
			continue
		}
		
		// First cell is category name
		category := strings.TrimSpace(cells[0])
		
		// Skip header rows
		if category == "" || strings.Contains(strings.ToLower(category), "category") ||
			strings.Contains(strings.ToLower(category), "shareholder") {
			continue
		}
		
		// Extract count (usually 2nd cell)
		count := 0
		if len(cells) > 1 {
			count = parseInt(cells[1])
		}
		
		// Extract percentage (usually 4th or 5th cell)
		percentage := 0.0
		for i := 3; i < len(cells) && i < 6; i++ {
			if pct := parseFloat(cells[i]); pct > 0 {
				percentage = pct
				break
			}
		}
		
		// Extract shares (usually 3rd cell)
		shares := int64(0)
		if len(cells) > 2 {
			shares = parseInt64(cells[2])
		}
		
		// Store category data
		holding := CategoryHolding{
			CategoryName: category,
			Count:        count,
			Percentage:   percentage,
			Shares:       shares,
		}
		
		data.Categories[category] = holding
		
		// Aggregate into main categories
		categoryLower := strings.ToLower(category)
		
		switch {
		case strings.Contains(categoryLower, "promoter"):
			data.PromoterTotal += percentage
		case strings.Contains(categoryLower, "fii") || strings.Contains(categoryLower, "foreign"):
			data.FIITotal += percentage
		case strings.Contains(categoryLower, "mutual fund"):
			data.DIITotal += percentage
		case strings.Contains(categoryLower, "insurance"):
			data.DIITotal += percentage
		case strings.Contains(categoryLower, "dii") || strings.Contains(categoryLower, "domestic"):
			data.DIITotal += percentage
		case strings.Contains(categoryLower, "public"):
			data.PublicTotal += percentage
		}
	}
}

// extractCells extracts text content from table cells
func extractCells(rowHTML string) []string {
	cellRe := regexp.MustCompile(`<t[dh][^>]*>(.*?)</t[dh]>`)
	matches := cellRe.FindAllStringSubmatch(rowHTML, -1)
	
	cells := make([]string, 0, len(matches))
	for _, match := range matches {
		text := stripHTMLTags(match[1])
		text = strings.TrimSpace(text)
		cells = append(cells, text)
	}
	
	return cells
}

// stripHTMLTags removes HTML tags from string
func stripHTMLTags(s string) string {
	re := regexp.MustCompile(`<[^>]*>`)
	s = re.ReplaceAllString(s, "")
	// Decode HTML entities
	s = strings.ReplaceAll(s, "&nbsp;", " ")
	s = strings.ReplaceAll(s, "&", "&")
	return strings.TrimSpace(s)
}

// parseInt safely parses integer from string
func parseInt(s string) int {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	if val, err := strconv.Atoi(s); err == nil {
		return val
	}
	return 0
}

// parseInt64 safely parses int64 from string
func parseInt64(s string) int64 {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	if val, err := strconv.ParseInt(s, 10, 64); err == nil {
		return val
	}
	return 0
}

// parseFloat safely parses float from string
func parseFloat(s string) float64 {
	s = strings.ReplaceAll(s, ",", "")
	s = strings.TrimSpace(s)
	if val, err := strconv.ParseFloat(s, 64); err == nil {
		return val
	}
	return 0.0
}

// classifyUsingBSEData classifies client using actual BSE shareholding data
func classifyUsingBSEData(clientName string, bseData *BSEShareholdingData) (category string, confidence string, percentage float64) {
	if bseData == nil {
		// Fallback to heuristic
		isPromoter, isFII, isDII := classifyClient(clientName)
		if isPromoter {
			return "🏢 Promoter", "Heuristic", 0
		} else if isFII {
			return "🌍 FII", "Heuristic", 0
		} else if isDII {
			return "🏦 DII", "Heuristic", 0
		}
		return "Public", "Heuristic", 0
	}
	
	clientUpper := strings.ToUpper(clientName)
	
	// Check each category for name match
	for catName, holding := range bseData.Categories {
		catUpper := strings.ToUpper(catName)
		
		// Check if client name appears in category name or vice versa
		if strings.Contains(catUpper, clientUpper) || strings.Contains(clientUpper, catUpper) {
			catLower := strings.ToLower(catName)
			
			if strings.Contains(catLower, "promoter") {
				return "🏢 Promoter", "BSE Confirmed", holding.Percentage
			} else if strings.Contains(catLower, "fii") || strings.Contains(catLower, "foreign") {
				return "🌍 FII", "BSE Confirmed", holding.Percentage
			} else if strings.Contains(catLower, "mutual") || strings.Contains(catLower, "insurance") || strings.Contains(catLower, "dii") {
				return "🏦 DII", "BSE Confirmed", holding.Percentage
			}
		}
	}
	
	// Fallback to heuristic
	isPromoter, isFII, isDII := classifyClient(clientName)
	if isPromoter {
		return "🏢 Promoter", "Heuristic (not in BSE list)", 0
	} else if isFII {
		return "🌍 FII", "Heuristic (not in BSE list)", 0
	} else if isDII {
		return "🏦 DII", "Heuristic (not in BSE list)", 0
	}
	
	return "Public", "Not found in BSE data", 0
}

// Made with Bob
