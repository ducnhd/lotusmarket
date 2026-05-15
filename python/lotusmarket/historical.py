"""Cohort analysis on a panel of (ticker, date, features, fwd_returns) rows.

Bucket rows by RSI band, MA trend, MACD signal, Wyckoff stage, regime, and
joint(trend × RSI), then compute mean forward return + win rate per cohort.

DB-agnostic: caller queries their own table and converts to FeatureRow.
fwd_return values are interpreted as percent points (1.0 = 1%).
"""

from __future__ import annotations

from dataclasses import dataclass, field
from datetime import datetime
from typing import Dict, List, Optional, Sequence, Tuple


@dataclass
class FeatureRow:
    ticker: str
    date: datetime
    close: float
    rsi14: Optional[float] = None
    ma20: Optional[float] = None
    ma50: Optional[float] = None
    ma200: Optional[float] = None
    macd: Optional[float] = None
    macd_signal: Optional[float] = None
    wyckoff: Optional[int] = None  # 1-4
    regime: Optional[str] = None
    fwd5: Optional[float] = None  # percent points
    fwd20: Optional[float] = None
    fwd60: Optional[float] = None


@dataclass
class Stats:
    bucket: str
    n: int
    mean_5d: float
    mean_20d: float
    mean_60d: float
    win_5d: float  # fraction 0-1
    win_20d: float
    win_60d: float
    edge_60d: float


@dataclass
class Report:
    window: str
    clean_rows: int
    baseline: Stats
    groups: Dict[str, List[Stats]] = field(default_factory=dict)


def _clip(x: float, cap: float) -> Optional[float]:
    """Reject implausible fwd returns from unadjusted corporate actions."""
    if x > cap or x < -80:
        return None
    return x


def _summary(xs: List[float]) -> Tuple[float, float, int]:
    n = len(xs)
    if n == 0:
        return 0.0, 0.0, 0
    s = sum(xs)
    wins = sum(1 for x in xs if x > 0)
    return s / n, wins / n, n


def rsi_bucket(rsi: float) -> str:
    if rsi < 30:
        return "RSI<30 (oversold)"
    if rsi < 50:
        return "RSI 30-50"
    if rsi < 70:
        return "RSI 50-70"
    return "RSI>70 (overbought)"


def ma_trend(close: float, ma50: Optional[float], ma200: Optional[float]) -> str:
    if ma200 is None or ma50 is None:
        return ""
    above = close >= ma200
    ma50_above = ma50 >= ma200
    if above and ma50_above:
        return "uptrend"
    if not above and not ma50_above:
        return "downtrend"
    return "mixed"


def wyckoff_label(stage: int) -> str:
    return {
        1: "1-accumulation",
        2: "2-markup",
        3: "3-distribution",
        4: "4-decline",
    }.get(stage, str(stage))


def analyze(rows: Sequence[FeatureRow], window: str = "all") -> Report:
    """Run cohort analysis on a panel of rows. Returns a Report ready to render."""
    cohorts: Dict[
        str, Dict[str, List[Tuple[Optional[float], Optional[float], Optional[float]]]]
    ] = {
        "rsi": {},
        "trend": {},
        "macd": {},
        "wyckoff": {},
        "regime": {},
        "joint": {},
    }

    def add(group: str, label: str, row: FeatureRow) -> None:
        bucket = cohorts[group].setdefault(label, [])
        f5 = _clip(row.fwd5, 30) if row.fwd5 is not None else None
        f20 = _clip(row.fwd20, 60) if row.fwd20 is not None else None
        f60 = _clip(row.fwd60, 150) if row.fwd60 is not None else None
        bucket.append((f5, f20, f60))

    for r in rows:
        if r.rsi14 is not None:
            add("rsi", rsi_bucket(r.rsi14), r)
        trend = ma_trend(r.close, r.ma50, r.ma200)
        if trend:
            add("trend", trend, r)
        if r.macd is not None and r.macd_signal is not None:
            add("macd", "macd_bull" if r.macd > r.macd_signal else "macd_bear", r)
        if r.wyckoff is not None:
            add("wyckoff", wyckoff_label(r.wyckoff), r)
        if r.regime:
            add("regime", r.regime, r)
        if r.rsi14 is not None and trend:
            add("joint", f"{trend} × {rsi_bucket(r.rsi14)}", r)

    # Baseline
    all5: List[float] = []
    all20: List[float] = []
    all60: List[float] = []
    for r in rows:
        if r.fwd5 is not None:
            v = _clip(r.fwd5, 30)
            if v is not None:
                all5.append(v)
        if r.fwd20 is not None:
            v = _clip(r.fwd20, 60)
            if v is not None:
                all20.append(v)
        if r.fwd60 is not None:
            v = _clip(r.fwd60, 150)
            if v is not None:
                all60.append(v)
    m5, w5, _ = _summary(all5)
    m20, w20, _ = _summary(all20)
    m60, w60, _ = _summary(all60)
    baseline = Stats(
        bucket="BASELINE",
        n=len(all60),
        mean_5d=m5,
        mean_20d=m20,
        mean_60d=m60,
        win_5d=w5,
        win_20d=w20,
        win_60d=w60,
        edge_60d=0.0,
    )

    report = Report(window=window, clean_rows=len(rows), baseline=baseline)
    for group, buckets in cohorts.items():
        stats_list: List[Stats] = []
        for label in sorted(buckets.keys()):
            entries = buckets[label]
            xs5 = [e[0] for e in entries if e[0] is not None]
            xs20 = [e[1] for e in entries if e[1] is not None]
            xs60 = [e[2] for e in entries if e[2] is not None]
            mm5, ww5, n5 = _summary(xs5)
            mm20, ww20, _ = _summary(xs20)
            mm60, ww60, _ = _summary(xs60)
            stats_list.append(
                Stats(
                    bucket=label,
                    n=n5,
                    mean_5d=mm5,
                    mean_20d=mm20,
                    mean_60d=mm60,
                    win_5d=ww5,
                    win_20d=ww20,
                    win_60d=ww60,
                    edge_60d=mm60 - m60,
                )
            )
        report.groups[group] = stats_list
    return report


def to_markdown(report: Report) -> str:
    """Render a Report as a markdown leaderboard."""
    lines: List[str] = []
    lines.append("# Historical Cohort Analysis\n")
    lines.append(
        f"Generated: {datetime.now().strftime('%Y-%m-%d %H:%M')}  "
        f"\nWindow: {report.window}  "
        f"\nClean rows: {report.clean_rows} (excludes outliers)\n"
    )
    b = report.baseline
    lines.append("## Baseline\n")
    lines.append("| Horizon | Mean fwd return | Win rate | N |")
    lines.append("|---|---|---|---|")
    lines.append(f"| 5d | {b.mean_5d:+.2f}% | {b.win_5d * 100:.0f}% | {b.n} |")
    lines.append(f"| 20d | {b.mean_20d:+.2f}% | {b.win_20d * 100:.0f}% | {b.n} |")
    lines.append(f"| 60d | {b.mean_60d:+.2f}% | {b.win_60d * 100:.0f}% | {b.n} |\n")

    order = ["rsi", "trend", "macd", "wyckoff", "regime", "joint"]
    titles = {
        "rsi": "## RSI bucket",
        "trend": "## MA trend",
        "macd": "## MACD signal",
        "wyckoff": "## Wyckoff stage",
        "regime": "## Market regime",
        "joint": "## Joint: trend × RSI",
    }
    for g in order:
        stats = report.groups.get(g, [])
        if not stats:
            continue
        lines.append(titles[g] + "\n")
        lines.append(
            "| Bucket | N | Mean 5d | Win 5d | Mean 20d | Win 20d | **Mean 60d** | Win 60d | Edge 60d |"
        )
        lines.append("|---|---|---|---|---|---|---|---|---|")
        for s in stats:
            lines.append(
                f"| {s.bucket} | {s.n} | {s.mean_5d:+.2f}% | {s.win_5d * 100:.0f}% "
                f"| {s.mean_20d:+.2f}% | {s.win_20d * 100:.0f}% "
                f"| **{s.mean_60d:+.2f}%** | {s.win_60d * 100:.0f}% "
                f"| {s.edge_60d:+.2f}% |"
            )
        lines.append("")
    return "\n".join(lines)
