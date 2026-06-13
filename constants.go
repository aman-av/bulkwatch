package main

// API URLs and Endpoints
const (
	// NSE URLs
	NSEBulkBlockDealsAPI = "https://www.nseindia.com/api/historicalOR/bulk-block-short-deals"
	NSEEquityListURL1    = "https://archives.nseindia.com/content/equities/EQUITY_L.csv"
	NSEEquityListURL2    = "https://www1.nseindia.com/content/equities/EQUITY_L.csv"
	NSEMarketCapURL      = "https://nsearchives.nseindia.com/web/sites/default/files/inline-files/MCAP31032021_2.xlsx"
	
	// BSE URLs
	BSEBulkBlockDealsAPI       = "https://api.bseindia.com/BseIndiaAPI/api/BulkDealData_ng/w"
	BSEScripListAPI            = "https://api.bseindia.com/BseIndiaAPI/api/ListofScripData/w"
	BSEShareholdingPromoterURL = "https://www.bseindia.com/corporates/shppromoterngroup"
	BSEShareholdingPublicURL   = "https://www.bseindia.com/corporates/shppublicshareholder"
	BSEOrigin                  = "https://www.bseindia.com"
	BSEReferer                 = "https://www.bseindia.com/"
)

// URL Builder Functions

// BuildNSEBulkBlockURL builds NSE bulk/block deals API URL
func BuildNSEBulkBlockURL(optionType, fromDate, toDate string) string {
	return NSEBulkBlockDealsAPI + "?optionType=" + optionType + "&from=" + fromDate + "&to=" + toDate
}

// BuildBSEBulkBlockURL builds BSE bulk/block deals API URL
func BuildBSEBulkBlockURL(dealType, fromDate, toDate string) string {
	return BSEBulkBlockDealsAPI + "?DealType=" + dealType + "&sc_code=&FDate=" + fromDate + "&TDate=" + toDate
}

// BuildBSEScripListURL builds BSE scrip list API URL
func BuildBSEScripListURL() string {
	return BSEScripListAPI + "?Group=&Scripcode=&industry=&segment=Equity&status=Active"
}

// BuildBSEShareholdingPromoterURL builds BSE promoter shareholding URL
func BuildBSEShareholdingPromoterURL(scripCode, qtrID, qtrName string) string {
	return BSEShareholdingPromoterURL + "?scripcd=" + scripCode + "&qtrid=" + qtrID + "&QtrName=" + qtrName
}

// BuildBSEShareholdingPublicURL builds BSE public shareholding URL
func BuildBSEShareholdingPublicURL(scripCode, qtrID, qtrName string) string {
	return BSEShareholdingPublicURL + "?scripcd=" + scripCode + "&qtrid=" + qtrID + "&QtrName=" + qtrName
}

// Made with Bob
