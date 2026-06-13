# Shareholding Pattern Analysis

## Overview

FlowWatch now includes shareholding pattern analysis that cross-references block/bulk deals with institutional holdings data from BSE. This helps identify whether deals are made by Promoters, FIIs (Foreign Institutional Investors), or DIIs (Domestic Institutional Investors).

## Features

### 1. **Automatic Classification**
The system automatically classifies clients into categories:

- **🏢 Promoter**: Companies with "Private Limited", "Trust", "Holdings", "Family" in their names
- **🌍 FII**: Foreign institutional investors (Goldman Sachs, Morgan Stanley, BNP Paribas, etc.)
- **🏦 DII**: Domestic institutional investors (Mutual Funds, Insurance companies, LIC, etc.)
- **Public**: Other investors

### 2. **NSE to BSE Mapping**
Maps NSE symbols to BSE script codes for verification:

| NSE Symbol | BSE Script Code | Company |
|------------|-----------------|---------|
| ADANIENT   | 512599         | Adani Enterprises |
| SWIGGY     | 543971         | Swiggy Limited |
| TATACAP    | 544574         | Tata Capital |
| FORCEMOT   | 500033         | Force Motors |
| M&M        | 500520         | Mahindra & Mahindra |
| SBILIFE    | 542163         | SBI Life Insurance |

### 3. **Verification URLs**
Each deal includes a direct link to BSE's shareholding pattern page where you can verify:
- Promoter holdings percentage
- FII holdings breakdown
- DII holdings breakdown
- Mutual fund holdings
- Insurance company holdings

## Usage

The analysis runs automatically when you execute the program:

```bash
go run .
```

### Sample Output

```
🔍 Analyzing NSE block Deals (First 20 records)
================================================================================
📊 Processing 20 deals...

📈 SWIGGY (2 deals)
--------------------------------------------------------------------------------
  1. Public - CYRUS SOLI POONAWALLA
     Qty: 1123500 | Price: 377.00 | Value: 423.56M
     🔗 Verify: https://www.bseindia.com/corporates/shppublicshareholder?scripcd=543971&qtrid=129.00&QtrName=Mar-26

  2. 🏢 Promoter - SERUM INSTITUTE OF INDIA PRIVATE LIMITED
     Qty: 1123500 | Price: 377.00 | Value: 423.56M
     🔗 Verify: https://www.bseindia.com/corporates/shppublicshareholder?scripcd=543971&qtrid=129.00&QtrName=Mar-26

📊 Summary:
   Promoter Deals: 10
   FII Deals:      3
   DII Deals:      0
   Public Deals:   7
```

## How to Verify Shareholding Patterns

1. **Click the verification URL** provided in the output
2. The BSE page shows:
   - **Promoters**: Percentage and number of shares held by company promoters
   - **FII**: Foreign institutional investor holdings
   - **DII**: Domestic institutional investor holdings
   - **Mutual Funds**: Detailed breakdown by fund house
   - **Insurance Companies**: Holdings by insurance firms
   - **Public**: Retail and other public shareholders

3. **Cross-reference** the client name from the block/bulk deal with the shareholding pattern to confirm their category

## Example Analysis

### Case Study: SWIGGY Block Deal

**Deal Details:**
- Client: SERUM INSTITUTE OF INDIA PRIVATE LIMITED
- Type: SELL
- Quantity: 1,123,500 shares
- Price: ₹377
- Value: ₹423.56M

**Classification:** 🏢 Promoter (based on "PRIVATE LIMITED" in name)

**Verification:**
Visit: https://www.bseindia.com/corporates/shppublicshareholder?scripcd=543971&qtrid=129.00&QtrName=Mar-26

On the BSE page, you can verify:
- If Serum Institute appears in the Promoter category
- Their total shareholding percentage
- Recent changes in their holdings

## Key Insights

### Promoter Deals
- Often indicate strategic moves (stake increase/decrease)
- May signal confidence or need for liquidity
- Important for corporate governance analysis

### FII Deals
- Indicate foreign investor sentiment
- Large FII deals can impact stock prices
- Track global fund flows into Indian markets

### DII Deals
- Show domestic institutional interest
- Mutual fund purchases indicate retail-friendly stocks
- Insurance company holdings suggest long-term stability

## Technical Details

### Data Sources
1. **Block/Bulk Deals**: NSE and BSE APIs
2. **Shareholding Patterns**: BSE Shareholding Pattern API
   - URL Format: `https://www.bseindia.com/corporates/shppublicshareholder?scripcd={SCRIP_CODE}&qtrid={QTR_ID}&QtrName={QTR_NAME}`
   - Quarter IDs: 129 (Mar-26), 128 (Dec-25), 127 (Sep-25), etc.

### Classification Logic

```go
// FII Keywords
"GOLDMAN SACHS", "MORGAN STANLEY", "BNP PARIBAS", "SOCIETE GENERALE",
"CITIGROUP", "BOFA", "CLSA", "NOMURA", "CREDIT SUISSE", "UBS"

// DII Keywords
"MUTUAL FUND", "INSURANCE", "LIC", "ICICI PRUDENTIAL",
"HDFC", "SBI LIFE", "ADITYA BIRLA", "NIPPON", "KOTAK"

// Promoter Keywords (Heuristic)
"PRIVATE LIMITED", "TRUST", "HOLDINGS", "FAMILY",
"INVESTMENT", "ENTERPRISES", "VENTURES"
```

## Limitations

1. **Sample Size**: Currently analyzes first 20 records to avoid context overload
2. **Heuristic Classification**: Promoter detection is based on name patterns, not definitive data
3. **Manual Verification**: BSE API requires browser session, so verification URLs are provided for manual checking
4. **Mapping Coverage**: NSE-to-BSE mapping covers major stocks but not exhaustive

## Future Enhancements

- [ ] Expand NSE-to-BSE symbol mapping database
- [ ] Add historical shareholding pattern tracking
- [ ] Implement automated BSE API authentication
- [ ] Add quarter-over-quarter change analysis
- [ ] Generate alerts for significant promoter/FII/DII movements
- [ ] Export analysis to CSV/Excel format

## Privacy & Data Usage

- **No Personal Data**: Only publicly available market data is used
- **No Storage**: Shareholding patterns are not stored, only verification URLs provided
- **Real-time**: Data is fetched fresh from NSE/BSE APIs
- **Transparency**: All data sources are clearly documented

## Contributing

To add more NSE-to-BSE mappings, edit [`shareholding.go`](shareholding.go):

```go
func mapNSEToBSE(nseSymbol string) string {
    mappings := map[string]string{
        "SYMBOL": "SCRIPCODE",
        // Add more mappings here
    }
    return mappings[nseSymbol]
}
```

## References

- [BSE Shareholding Pattern Page](https://www.bseindia.com/corporates/shppublicshareholder)
- [NSE Block Deals](https://www.nseindia.com/report-detail/eq_block)
- [BSE Bulk Deals](https://www.bseindia.com/markets/equity/EQReports/BulkDeals.aspx)

---

**Made with Bob** 🤖