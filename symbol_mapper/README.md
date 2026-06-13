# 🎯 NSE-BSE Symbol Mapper

Get top 1000 companies by market cap with NSE↔BSE mapping.

## 🚀 Quick Start

```bash
cd symbol_mapper
go run run_marketcap.go
```

## 📁 Output Files

- **`top_1000_marketcap.json`** - JSON format with all details
- **`top_1000_marketcap.csv`** - CSV format for Excel/spreadsheets

## 📊 Data Format

### JSON
```json
{
  "nseSymbol": "RELIANCE",
  "bseScripCode": "500325",
  "companyName": "Reliance Industries Limited",
  "isin": "INE002A01018",
  "marketCap": 1269853.61,
  "marketCapCr": "₹12.70 Lakh Cr",
  "rank": 1
}
```

### CSV
```
Rank, NSE_Symbol, BSE_ScripCode, Company, Market_Cap_Cr, ISIN
1, RELIANCE, 500325, Reliance Industries Limited, ₹12.70 Lakh Cr, INE002A01018
```

## 📈 Top 10 Companies

| Rank | NSE | BSE | Company | Market Cap |
|------|-----|-----|---------|------------|
| 1 | RELIANCE | 500325 | Reliance Industries | ₹12.70 Lakh Cr |
| 2 | TCS | 532540 | Tata Consultancy Services | ₹11.76 Lakh Cr |
| 3 | HDFCBANK | 500180 | HDFC Bank | ₹8.23 Lakh Cr |
| 4 | INFY | 500209 | Infosys | ₹5.83 Lakh Cr |
| 5 | HINDUNILVR | 500696 | Hindustan Unilever | ₹5.71 Lakh Cr |
| 6 | ICICIBANK | 532174 | ICICI Bank | ₹4.03 Lakh Cr |
| 7 | KOTAKBANK | 500247 | Kotak Mahindra Bank | ₹3.47 Lakh Cr |
| 8 | SBIN | 500112 | State Bank of India | ₹3.25 Lakh Cr |
| 9 | BAJFINANCE | 500034 | Bajaj Finance | ₹3.10 Lakh Cr |
| 10 | BHARTIARTL | 532454 | Bharti Airtel | ₹2.84 Lakh Cr |

## 🔍 Data Sources

1. **NSE Equity List** - All NSE symbols with ISIN
2. **BSE API** - Active equity securities  
3. **NSE Market Cap Excel** - Official market capitalization data
4. **ISIN Matching** - Cross-exchange mapping

## 📊 Statistics

- Total Mappings: **2,236 companies**
- Top 1000: **All with verified market cap**
- Match Rate: **94.5%**
- Sorted by: **Actual market cap (descending)**

## 💡 Usage Example

```go
import "encoding/json"

type Mapping struct {
    NSESymbol    string  `json:"nseSymbol"`
    BSEScripCode string  `json:"bseScripCode"`
    CompanyName  string  `json:"companyName"`
    MarketCap    float64 `json:"marketCap"`
    Rank         int     `json:"rank"`
}

// Load mappings
data, _ := os.ReadFile("top_1000_marketcap.json")
var mappings []Mapping
json.Unmarshal(data, &mappings)

// Create lookup map
nseToBSE := make(map[string]string)
for _, m := range mappings {
    nseToBSE[m.NSESymbol] = m.BSEScripCode
}

// Use it
bseCode := nseToBSE["RELIANCE"]  // "500325"
```

## 🎯 Use Cases

- Trading applications
- Portfolio management
- Market analysis
- Data integration
- Research & analytics

---

**Data Sources:** NSE & BSE Official APIs  
**Last Updated:** 2026-06-09
