# FlowWatch - Project Objective

## 🎯 Primary Goal

**Track institutional money flow to identify early investment opportunities**

Analyze bulk and block deals to detect when "smart money" (FIIs, DIIs, Promoters) is accumulating or distributing positions, enabling us to be early buyers or sellers to maximize profit.

---

## 💡 Core Strategy

### The Insight
When big institutional players create positions in a stock, it signals:
- They have extensive research capabilities
- They see long-term potential
- They have insider knowledge (legally obtained)
- Their accumulation = bullish signal for retail investors

### Our Edge
- **Early Detection:** Daily tracking of bulk/block deals
- **Smart Money Following:** Track what institutions are buying/selling
- **Data-Driven Decisions:** Objective analysis, not emotions
- **Automated Monitoring:** Run daily to stay updated
- **Incremental Processing:** Only process new data, not re-compute everything

---

## 📊 What We Track

### Data Sources
1. **NSE Bulk Deals** - Large trades on National Stock Exchange
2. **NSE Block Deals** - Massive institutional trades
3. **BSE Bulk Deals** - Large trades on Bombay Stock Exchange
4. **BSE Block Deals** - Massive institutional trades

### Key Metrics
- **Client Classification:** Promoter / FII / DII / Public
- **Transaction Type:** Buy (P) / Sell (S)
- **Deal Size:** Quantity, Price, Total Value
- **Shareholding %:** Current holding percentage
- **Trend:** Accumulation vs Distribution

---

## 🚀 Current System Status

### ✅ Completed Features
1. **Data Collection**
   - Fetch bulk/block deals from NSE & BSE
   - Store in organized archive structure
   - Handle API compression (gzip, brotli)
   - Fixed BSE date format (DD/MM/YYYY)

2. **Data Enrichment**
   - Classify investors (Promoter/FII/DII)
   - Fetch shareholding data from BSE
   - Match client names with shareholding categories
   - Calculate confidence levels

3. **Symbol Mapping**
   - Map NSE symbols to BSE scrip codes
   - Include market cap data
   - Top 1000 stocks by market cap

4. **Infrastructure**
   - Centralized API URLs in constants.go
   - Task tracking system (AGENT folder)
   - Debug configurations for VS Code
   - State tracking (state.json) for incremental processing

---

## 🎯 Next Steps to Complete

### Phase 1: Incremental Processing (Priority)
1. **State Management**
   - ✅ Created state.json to track last processed dates
   - [ ] Implement state reading/writing in code
   - [ ] Only fetch deals after last_processed date
   - [ ] Only enrich new deals, skip already enriched
   - [ ] Update state after successful processing

2. **Efficient Enrichment**
   - [ ] Check if deal already enriched before processing
   - [ ] Cache shareholding data to avoid re-fetching
   - [ ] Batch process new deals only
   - [ ] Update statistics in state.json

### Phase 2: Analysis & Insights
3. **Pattern Detection**
   - Identify repeated buying by same institutions
   - Track accumulation trends over time
   - Detect unusual activity (sudden large positions)

4. **Scoring System**
   - Rate stocks based on institutional activity
   - Weight by investor type (FII > DII > Public)
   - Consider deal size and frequency

5. **Trend Analysis**
   - Compare current deals with historical data
   - Identify stocks with increasing institutional interest
   - Track promoter buying/selling patterns

### Phase 3: Actionable Outputs
6. **Daily Report Generation**
   - Top stocks with institutional buying
   - Stocks with promoter accumulation
   - Unusual activity alerts
   - Summary statistics

7. **Alert System**
   - Notify when big players accumulate specific stocks
   - Alert on promoter buying (strong signal)
   - Flag unusual selling patterns

8. **Visualization**
   - Charts showing institutional flow
   - Trend graphs for specific stocks
   - Heatmaps of sector-wise activity

### Phase 4: Automation
9. **Daily Scheduler**
   - Run automatically every trading day
   - Fetch only new deals (incremental)
   - Enrich only new data
   - Generate reports
   - Send notifications

10. **Historical Analysis**
    - Backtest strategy effectiveness
    - Correlate with stock price movements
    - Measure prediction accuracy

---

## 📈 Success Metrics

### Short-term (1-3 months)
- Daily automated data collection
- Incremental processing working (no re-computation)
- Accurate investor classification (>90%)
- Basic pattern detection working
- Daily reports generated

### Medium-term (3-6 months)
- Proven correlation with stock movements
- Alert system operational
- Historical trend analysis
- Backtesting results positive

### Long-term (6-12 months)
- Consistent early detection of opportunities
- Measurable profit improvement
- Refined scoring algorithm
- Community/sharing platform

---

## 🔄 Daily Workflow (Target)

```
1. Morning (Pre-market)
   - Run FlowWatch (incremental mode)
   - Fetch only yesterday's new deals
   - Enrich only new data
   - Generate daily report
   - Review top opportunities

2. Market Hours
   - Monitor alerts for unusual activity
   - Cross-reference with price movements
   - Make informed decisions

3. Evening (Post-market)
   - Review day's institutional activity
   - Update watchlist
   - Plan next day's strategy
```

---

## 💾 State Management

### state.json Structure
```json
{
  "last_processed": {
    "NSE_bulk": "2026-06-13",
    "NSE_block": "2026-06-13",
    "BSE_bulk": "2026-06-13",
    "BSE_block": "2026-06-13"
  },
  "last_enriched": {
    "NSE_bulk": "2026-06-13",
    "NSE_block": "2026-06-13",
    "BSE_bulk": "2026-06-13",
    "BSE_block": "2026-06-13"
  },
  "statistics": {
    "total_deals_processed": 0,
    "total_deals_enriched": 0,
    "last_run": "2026-06-13T15:57:00Z"
  }
}
```

### Benefits
- **Efficiency:** Don't re-process old data
- **Speed:** Only fetch new deals
- **Cost:** Reduce API calls (shareholding data)
- **Reliability:** Track what's been processed
- **Resume:** Can restart from last checkpoint

---

## 💪 Our Competitive Advantage

1. **Speed:** Automated daily tracking vs manual research
2. **Efficiency:** Incremental processing, no re-computation
3. **Objectivity:** Data-driven vs emotional decisions
4. **Comprehensiveness:** Both NSE & BSE coverage
5. **Intelligence:** Enriched with shareholding data
6. **Actionability:** Clear buy/sell signals

---

## 📝 Key Principles

1. **Follow Smart Money:** Institutions have better information
2. **Early Detection:** Be ahead of retail crowd
3. **Data Over Emotions:** Trust the numbers
4. **Consistency:** Daily monitoring is crucial
5. **Efficiency:** Process only what's new
6. **Risk Management:** Use as one signal, not sole decision

---

**Last Updated:** 2026-06-13  
**Status:** Phase 1 - Foundation Complete, Implementing Incremental Processing  
**Next Milestone:** State Management & Incremental Enrichment