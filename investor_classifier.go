package main

import (
	"encoding/json"
	"os"
	"strings"
	"sync"
)

// InvestorRegistry holds loaded FII and DII lists
type InvestorRegistry struct {
	FIIList map[string]bool
	DIIList map[string]bool
	mu      sync.RWMutex
}

// Global registry instance
var investorRegistry *InvestorRegistry
var registryOnce sync.Once

// DIIRegistry structure for JSON
type DIIRegistry struct {
	LastUpdated string `json:"last_updated"`
	Source      string `json:"source"`
	Count       int    `json:"count"`
	Categories  struct {
		MutualFunds  []string `json:"mutual_funds"`
		Insurance    []string `json:"insurance"`
		Banks        []string `json:"banks"`
		NBFCs        []string `json:"nbfcs"`
		PensionFunds []string `json:"pension_funds"`
		Government   []string `json:"government"`
	} `json:"categories"`
}

// FIIRegistryJSON structure for JSON
type FIIRegistryJSON struct {
	LastUpdated string `json:"last_updated"`
	Source      string `json:"source"`
	Count       int    `json:"count"`
	Investors   []struct {
		Name               string `json:"name"`
		Type               string `json:"type"`
		RegistrationNumber string `json:"registration_number,omitempty"`
	} `json:"investors"`
}

// GetInvestorRegistry returns the singleton registry instance
func GetInvestorRegistry() *InvestorRegistry {
	registryOnce.Do(func() {
		investorRegistry = &InvestorRegistry{
			FIIList: make(map[string]bool),
			DIIList: make(map[string]bool),
		}
		investorRegistry.loadRegistries()
	})
	return investorRegistry
}

// loadRegistries loads FII and DII lists from JSON files
func (r *InvestorRegistry) loadRegistries() {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Load DII list
	if data, err := os.ReadFile("investor_registry/dii_list.json"); err == nil {
		var diiReg DIIRegistry
		if err := json.Unmarshal(data, &diiReg); err == nil {
			// Add all DIIs to map
			for _, name := range diiReg.Categories.MutualFunds {
				r.DIIList[strings.ToUpper(name)] = true
			}
			for _, name := range diiReg.Categories.Insurance {
				r.DIIList[strings.ToUpper(name)] = true
			}
			for _, name := range diiReg.Categories.Banks {
				r.DIIList[strings.ToUpper(name)] = true
			}
			for _, name := range diiReg.Categories.NBFCs {
				r.DIIList[strings.ToUpper(name)] = true
			}
			for _, name := range diiReg.Categories.PensionFunds {
				r.DIIList[strings.ToUpper(name)] = true
			}
			for _, name := range diiReg.Categories.Government {
				r.DIIList[strings.ToUpper(name)] = true
			}
		}
	}

	// Load FII list - try basic list first, then full list if available
	fiiFiles := []string{
		"investor_registry/fii_list_basic.json",
		"investor_registry/fii_list.json",
	}
	
	for _, fiiFile := range fiiFiles {
		if data, err := os.ReadFile(fiiFile); err == nil {
			// Try new format first (with metadata)
			var fiiReg FIIRegistryJSON
			if err := json.Unmarshal(data, &fiiReg); err == nil {
				for _, investor := range fiiReg.Investors {
					r.FIIList[strings.ToUpper(investor.Name)] = true
				}
				break // Successfully loaded, stop trying other files
			}
			
			// Try simple format (array of objects with just "name")
			var simpleFII []struct {
				Name string `json:"name"`
			}
			if err := json.Unmarshal(data, &simpleFII); err == nil {
				for _, investor := range simpleFII {
					r.FIIList[strings.ToUpper(investor.Name)] = true
				}
				break // Successfully loaded, stop trying other files
			}
		}
	}
}

// IsFII checks if a client name is in the FII registry
func (r *InvestorRegistry) IsFII(clientName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	upper := strings.ToUpper(clientName)
	
	// Exact match
	if r.FIIList[upper] {
		return true
	}
	
	// Partial match - check if any FII name is contained in client name
	for fiiName := range r.FIIList {
		if strings.Contains(upper, fiiName) || strings.Contains(fiiName, upper) {
			return true
		}
	}
	
	return false
}

// IsDII checks if a client name is in the DII registry
func (r *InvestorRegistry) IsDII(clientName string) bool {
	r.mu.RLock()
	defer r.mu.RUnlock()
	
	upper := strings.ToUpper(clientName)
	
	// Exact match
	if r.DIIList[upper] {
		return true
	}
	
	// Partial match - check if any DII name is contained in client name
	for diiName := range r.DIIList {
		if strings.Contains(upper, diiName) || strings.Contains(diiName, upper) {
			return true
		}
	}
	
	return false
}

// classifyClientWithRegistry uses the investor registry for classification
func classifyClientWithRegistry(clientName string) (isPromoter, isFII, isDII bool) {
	registry := GetInvestorRegistry()
	upper := strings.ToUpper(clientName)
	
	// Check FII registry first
	if registry.IsFII(clientName) {
		return false, true, false
	}
	
	// Check DII registry
	if registry.IsDII(clientName) {
		return false, false, true
	}
	
	// Fallback to keyword matching for FII
	fiiKeywords := []string{
		"GOLDMAN SACHS", "MORGAN STANLEY", "BNP PARIBAS", "SOCIETE GENERALE",
		"CITIGROUP", "BOFA", "CLSA", "NOMURA", "CREDIT SUISSE", "UBS",
		"BARCLAYS", "DEUTSCHE", "HSBC", "JP MORGAN", "FOREIGN", "OFFSHORE",
	}
	for _, keyword := range fiiKeywords {
		if strings.Contains(upper, keyword) {
			return false, true, false
		}
	}
	
	// Fallback to keyword matching for DII
	diiKeywords := []string{
		"MUTUAL FUND", "INSURANCE", "LIC", "ICICI PRUDENTIAL",
		"HDFC", "SBI LIFE", "ADITYA BIRLA", "NIPPON", "KOTAK",
		"AXIS MUTUAL", "UTI", "DSP", "TATA MUTUAL", "MIRAE ASSET",
		"BANK", "NBFC", "PENSION", "PROVIDENT FUND", "EPFO",
	}
	for _, keyword := range diiKeywords {
		if strings.Contains(upper, keyword) {
			return false, false, true
		}
	}
	
	// Promoter patterns (heuristic)
	promoterKeywords := []string{
		"PRIVATE LIMITED", "TRUST", "HOLDINGS", "FAMILY",
		"INVESTMENT", "ENTERPRISES", "VENTURES",
	}
	for _, keyword := range promoterKeywords {
		if strings.Contains(upper, keyword) {
			return true, false, false
		}
	}
	
	return false, false, false
}

// Made with Bob
