# Event-driven market blog — design

Date: 2026-06-23
Status: Approved for planning

## Problem

The auto-blog generates posts from a fixed catalog of 10 evergreen topics on a
Mon/Wed/Fri schedule. Two problems:

1. **Not market-tracking.** Posts are timeless cohort/educational pieces with no
   tie to what actually happened in the market on the day they publish.
2. **Almost no readership.** 14 days of nginx logs show real article reads in the
   low dozens; most auto-posts get ~1 view. There is no embedded analytics, so we
   are measuring blind.

We want posts that are (a) triggered by real market events, (b) written to be
engaging, and (c) measurable — serving SEO discovery, the Telegram channel, and
credibility simultaneously.

## Goals

- Post **only when a real market event occurs** (event-triggered cadence).
- Anchor each post on a specific, named, dated event (ticker + number + date) so
  titles match what people search and are worth broadcasting.
- Keep the existing "no fabrication, data-anchored, no buy/sell rec" discipline.
- Add lightweight, bot-filtered readership measurement.
- Phase in real news headlines as the "why" layer without blocking phase 1.

## Non-goals

- No paid analytics or third-party tracker (privacy + no external dep).
- No intraday/real-time posting; one scan per trading day after close.
- No change to the Telegram broadcast mechanism or the quality gate's intent
  (it is reused).

## Architecture

New CLI subcommand `lmcli marketscan --out docs/blog`. It separates pure
detection (testable) from I/O (fetch) and prose (Claude):

```
lmcli marketscan --out docs/blog
  ├─ fetch:  VPSMultiple(VN30+HNX30) → quote (chg%, vol, bid/ask)
  │          EntradeHistory(each ticker, 30) → MA20 volume, prior closes
  │          YahooMultiple(^VIX,^GSPC,^N225,^HSI,GC=F,CL=F)
  │          CafefInsider(top movers)
  ├─ detect: DetectEvents(Snapshot) → []Event   ← PURE, no network
  ├─ gate:   maxScore < eventThreshold → no post (exit 0); else dedup check
  │          if no event AND days-since-last-post ≥ backfillGapDays → evergreen backfill
  └─ write:  top Event → event-mode prompt → Claude → docs/blog/<date>-event-<type>-<ticker>.md
```

`DetectEvents` takes a `Snapshot` struct (all fetched data) and returns scored
events. Fetching lives in builder code; detection is a pure function over structs
so it is unit-testable with fixtures and no network.

The existing 10-topic generator (`autoblog.go`) is retained and invoked as the
**backfill** path only.

### Data types (Go)

```go
type Snapshot struct {
    Date     time.Time
    Stocks   []StockSnap         // VN30 + HNX30
    Globals  []fetchers.YahooQuote
    Insiders []InsiderSnap
    FlowBuyPressure float64       // aggregate bid/(bid+ask) %
}

type StockSnap struct {
    Ticker     string
    Close      float64
    ChangePct  float64
    Volume     int64
    AvgVol20   float64
    RSI        float64
    MATrend    string            // uptrend|downtrend|mixed (reuse existing calc)
}

type Event struct {
    Type      string             // mover|volume|flow|insider|global
    Ticker    string             // "" for market-wide events
    Magnitude float64            // raw signal (chg%, vol multiple, etc.)
    Score     float64            // 0-100
    DataBlock string             // anchor numbers + cohort match, fed to Claude
}
```

## Event taxonomy and scoring (phase 1 — existing data only)

| Type      | Detection                                   | Score band |
|-----------|---------------------------------------------|------------|
| `mover`   | \|chg%\| ≥ 5% (HOSE limit ±7%)              | 50–85, +bonus if volume high |
| `volume`  | vol ≥ 2× AvgVol20, with price direction     | 40–70 by multiple |
| `flow`    | aggregate buy-pressure ≥ 60% or ≤ 40%       | 40–60 |
| `insider` | insider transaction value past threshold    | 45–65 |
| `global`  | VIX > 25/30, or S&P/HSI ±2%+                 | 45–75 ("spillover to VN" angle) |

Scoring functions are per-type and pure. Pick the single highest-scoring event.
Tie-break: prefer a named single-stock event (`mover`/`volume`/`insider`) over a
market-wide one (`flow`/`global`) — named tickers are stronger for SEO and titles.

Each Event carries a `DataBlock`: the real anchor numbers plus the historical
cohort match (reuse `matchCohort(matTrend, rsi)` already in `autoblog.go`). Claude
writes prose only; it never invents numbers.

### Threshold and dedup

- `eventThreshold = 55` (tunable const). `maxScore < eventThreshold` → no event post.
- Dedup reuses the recency mechanism added in the autoblog fix: content key
  `event-<type>-<ticker>`, cooldown `eventCooldownDays = 5`, so the same ticker +
  event type does not produce back-to-back posts.

### Backfill (enabled)

If a scan finds no event above threshold **and** `days-since-last-post ≥
backfillGapDays` (default **3**), publish one evergreen post via the existing
topic generator (which already has its own variant dedup). This keeps a minimum
cadence (~2–3 posts/week) for SEO freshness; event posts publish on top of that.

## Engaging writing (event-mode prompt)

A new prompt variant for event posts, applying proven finance-copywriting
principles while keeping the no-fabrication rules:

- **Title formula**: `[Ticker/Index] + [what happened + number] + [tension/curiosity]`,
  always containing ticker + number + date. E.g. *"VHM giảm sàn 7% hôm nay 24/6 —
  nhưng cohort 9 năm nói gì về vùng này?"*
- **Lede (2–3 sentences)**: BLUF — what happened — plus a curiosity gap the data
  resolves below.
- **Body as narrative**: ① what happened (real numbers) → ② why it may have
  happened (phase 1: sector/global/flow context; phase 2: real headlines) →
  ③ how similar historical cohorts played out next (the differentiator) →
  ④ what to watch.
- **Length**: 700–1100 words (shorter than evergreen; event posts read better tight).
- **Discipline kept**: no buy/sell recommendation, a "Verify reproducible" block
  (lmcli/pip), one-line disclaimer, never reveal as AI.

## Measurement

In-app view counter in `lotusai-site` (Go web server on the Pi):

- Middleware on `/blog/{slug}` increments a per-slug counter for non-bot User
  Agents (substring blocklist: bot, crawl, spider, slurp, bing, google, yandex,
  facebookexternal, semrush, ahrefs, petal, gpt, claude, python-requests, curl,
  go-http).
- Persist to `data/blog_views.json`, flushed periodically and on shutdown.
- Internal route `/blog/stats` returns counts (JSON).
- No third-party dependency; privacy-friendly; filters out the vulnerability
  scanners that currently dominate the logs.
- **Requires a one-time redeploy of the `lotusai-site` binary to the Pi.**

## Phase 2 — real headlines

Add `fetchers/news.go` pulling titles from **CafeF RSS** (more stable than HTML
scraping) → `sentiment.ClusterTitles` (already exists) groups them → select the
cluster relevant to the event's ticker → insert a "why" paragraph with citation.
Optionally promote a large news cluster to its own event type. Kept separate so
phase 1 ships without an external news dependency.

## Cron and scheduling

- Replace the `auto-blog` schedule (Mon/Wed/Fri 09:00 VN) with a **post-close**
  run: Mon–Fri at 09:00 UTC (16:00 VN, after the 15:00 close). Job runs
  `marketscan`; it posts only on event or backfill.
- Update `trigger-workflows.sh` on the Pi (fallback dispatcher) to match the new
  schedule and workflow.

## Testing

`marketscan_test.go`, all fixture-based (no network):

- Table-driven `DetectEvents`: each event type at boundary magnitudes.
- Threshold gate: below/at/above `eventThreshold`.
- Tie-break ordering: named single-stock beats market-wide at equal score.
- Dedup: content key shape and cooldown behavior.
- Backfill trigger: fires only when no event AND gap ≥ `backfillGapDays`.

## Decisions (confirmed)

- Universe: VN30 + HNX30 (~60 tickers).
- Backfill: **enabled**, `backfillGapDays = 3`.
- View counter requiring a Pi redeploy: approved.

## Open tunables (defaults, adjustable later)

`eventThreshold = 55`, `eventCooldownDays = 5`, `backfillGapDays = 3`,
volume-surge multiple `= 2×`, mover threshold `= 5%`.
