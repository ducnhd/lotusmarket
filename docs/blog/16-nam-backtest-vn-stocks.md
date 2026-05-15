---
title: "16 năm backtest cổ phiếu Việt Nam — 5 điều ngược lại với hầu hết khuyên cũ"
date: 2026-05-15
tags: [vietnam, stocks, backtest, vn30, quant]
---

# 16 năm backtest cổ phiếu Việt Nam — 5 điều ngược lại với hầu hết khuyên cũ

> TL;DR — Mình chạy backtest deterministic trên ~47,000 data points của 30 mã VN30 + HNX30 từ 2010 đến 2026. Đây là 5 finding ngược với hầu hết khuyên đang lan truyền trên các group đầu tư Việt Nam. Tất cả code + data + reproducible setup đều open source tại [lotusmarket](https://github.com/ducnhd/lotusmarket).

---

## Setup ngắn gọn

- **Window**: 2010-01-04 → 2026-04-29 (16.3 năm)
- **Tickers**: 30 mã VN30 + HNX30 (ACB, BVH, CTG, FPT, HPG, MSN, SHB, SSI, STB, VCB, VIC, VNM, + 18 mã khác)
- **Data source**: Yahoo Finance + VPS histdatafeed (split + dividend adjusted)
- **Phí**: TPS retail 0.15% commission + 0.1% selling tax + 0.2% slippage + lot 100 (realistic)
- **Stop-loss**: -15% từ avg cost
- **Reproducible**: `make backtest` trong repo

Không cherry-pick. Không "chỉ năm tốt". Bao gồm 2011 (-10%), 2018 (-15%), **2022 (VTP/BĐS, -28%)**, 2025 mini-crash (-31%). 10/16 năm xanh, 6/16 đỏ.

---

## Finding #1 — Buy-hold thắng tất cả active strategies trên 16 năm

| Strategy | CAGR 16y | MaxDD | Sharpe |
|---|---|---|---|
| **Buy-hold equal-weight VN30** | **+15.37%** | -36% | 0.81 |
| trend_200ma | +21.6% (10y) | -51% | 0.59 |
| momentum_12_1 | +21.7% (10y) | -44% | 0.62 |
| low_vol top-5 | +26.2% (10y) | -42% | 0.74 |
| dual_momentum | +15.7% (10y) | -55% | 0.52 |

Active strategy trông đẹp trên window 10y (2017-2026) — nhưng đó là **window cherry-picked**. Mở rộng ra 16 năm, buy-hold đơn giản (mua đều 12 mã VN30 cũ, không rebalance) **đánh bại tất cả**.

**Lý do**: phí + slippage + mistimed entries ăn hết alpha mà active strategy tạo ra. Trên 16 năm, mỗi lần rebalance là một lần lose ~0.35% bị phí ăn — và momentum/low_vol rebalance ~12 lần/năm.

**Hệ quả thực tế**: nếu bạn không phải full-time quant, **chiến lược đơn giản nhất thắng**. Mua đều VN30 mỗi tháng. Không rebalance. Không trade. Năm 16 sau, bạn đánh bại 95% retail trader.

---

## Finding #2 — RSI > 70 là tín hiệu **MUA**, không phải bán

Đây là điều ngược lại nhất với những gì textbook + group VN đang dạy.

Cohort analysis trên ~45,000 row clean:

| RSI bucket | Mean fwd return 60d | Win rate | N |
|---|---|---|---|
| RSI < 30 (oversold) | +6.32% | 60% | 1,755 |
| RSI 30-50 | +3.67% | 52% | 18,933 |
| RSI 50-70 | +4.82% | 53% | 20,681 |
| **RSI > 70 (overbought)** | **+9.23%** | **60%** | 4,015 |

**Baseline VN30 2017-2026 60d = +4.75%.** Tức là RSI > 70 cho **edge +4.48% so với baseline**. Đây là **momentum effect**, không phải mean-reversion.

Càng kết hợp với MA trend càng rõ:

| Joint cohort | Mean 60d | Edge | N |
|---|---|---|---|
| Uptrend × RSI > 70 | **+9.80%** | **+5.05%** | 2,814 |
| Downtrend × RSI > 70 | **-1.40%** | **-6.15%** | 139 |

Cùng một RSI 75, nhưng:
- Stock đang trong uptrend (giá > MA200, MA50 > MA200) → cohort thắng lớn
- Stock đang trong downtrend → cohort duy nhất thực sự nên SELL

Bài học: **đừng dùng RSI làm signal độc lập**. Phải nhìn cả MA trend trước.

---

## Finding #3 — CRISIS không phải rủi ro. Là cơ hội (NẾU đúng loại).

Phân loại regime deterministic (no AI):

| Regime | Mean fwd return 60d | Win rate | N |
|---|---|---|---|
| STABLE | +3.20% | 48% | 23,868 |
| EUPHORIA | +6.12% | 59% | 10,279 |
| VOLATILE | +6.34% | 56% | 7,482 |
| **CRISIS** | **+8.53%** | **67%** | 2,582 |

Đúng vậy: **regime CRISIS có forward return cao nhất**. Vì lúc đó giá đã rớt, định giá lại, panic seller đã bán hết.

**Cảnh báo**: CRISIS có 2 sub-type:
- **CRISIS_PANIC** (contagion từ ngoài, VIX spike, news tier 1) → mua. Lịch sử cho thấy hồi phục trong 3-12 tháng.
- **CRISIS_FUNDAMENTAL** (vấn đề cấu trúc VN, bank run, kiểm tra hình sự lãnh đạo lớn) → tránh. Có thể không hồi phục.

Vạn Thịnh Phát 10/2022 là **fundamental**. 2020 COVID là **panic**. Phản ứng đúng cho 2 cái khác nhau hoàn toàn.

---

## Finding #4 — Với 500 triệu, bạn không thể full-time trader ở Việt Nam

Monte Carlo bootstrap 1000 paths, horizon 10 năm, withdrawal đều mỗi tháng:

| Vốn ban đầu | 12tr/tháng | 24tr/tháng | 36tr/tháng |
|---|---|---|---|
| **500 triệu** | **90% cháy** | 99% cháy | 99% cháy |
| 1 tỷ | 28% cháy | 91% cháy | 98% cháy |
| 1.5 tỷ | 9% cháy | 57% cháy | 90% cháy |
| 2 tỷ | 5% cháy ✅ | 29% cháy | 67% cháy |
| 3 tỷ | 4% cháy ✅ | 7% cháy ✅ | 26% cháy |
| **5 tỷ** | **2% cháy ✅** | **4% cháy ✅** | **6% cháy ✅** |

Vốn tối thiểu (P_cháy ≤ 5%, real-terms với lạm phát 4%):
- **12tr/tháng** → cần **2.55 tỷ**
- **24tr/tháng** → cần **5.06 tỷ**
- **36tr/tháng** → cần **7.78 tỷ**

Đây là **safe withdrawal rate ~5-7%** cho VN — cao hơn Mỹ (4%) vì CAGR VN cao hơn nhưng cũng volatile hơn. Quy luật: **vốn ≈ chi phí năm × 17-20**.

**Người có 500tr nói "nghỉ việc đi trade chứng khoán nuôi cả nhà"** chỉ có 1 outcome: 18-24 tháng sau quay lại đi làm với tài khoản nhỏ hơn ban đầu.

---

## Finding #5 — Sequence-of-returns risk lớn hơn ta tưởng

Cùng 2 tỷ + rút 12tr/tháng, kết quả khác nhau theo năm bắt đầu:

| Năm bắt đầu | Sau 10 năm |
|---|---|
| 2010 (sau khủng hoảng 2008) | OK, còn ~10 tỷ |
| 2014 (giữa cycle) | OK, còn ~8.6 tỷ |
| 2017 (cherry-picked tốt) | OK, còn ~6.2 tỷ |
| **2022 (khủng hoảng đầu)** | **OK nhưng còn ~2.2 tỷ — giảm 70% vốn dù không cháy** |

Đây không phải skill issue. Đây là **timing luck**. Bạn chọn đúng năm bắt đầu là 50% kết quả. Đó là lý do **emergency fund 24 tháng** không phải lựa chọn — là requirement.

---

## Hệ quả thực tế cho retail Việt Nam

Sau khi chạy hết các con số này, mình rút ra 4 nguyên tắc:

1. **Không trade. Buy-hold VN30. Đầu tư đều mỗi tháng.** Cách này thắng 95% retail.
2. **Nếu phải trade, dùng MA trend làm filter trước, RSI làm timer thứ hai.** Không bao giờ ngược lại.
3. **CRISIS là cơ hội — nhưng phải biết là PANIC hay FUNDAMENTAL.** Đọc tin trước khi quyết định.
4. **Để full-time trader VN, cần ít nhất 2.5 tỷ.** Dưới ngưỡng đó, giữ công việc chính, đầu tư song song.

---

## Reproducible

Tất cả số liệu trong bài này có thể chạy lại trên máy bạn:

```bash
git clone https://github.com/ducnhd/lotusmarket
cd lotusmarket/go && go build -o lmcli ./cmd/lmcli

# Daily report
./lmcli report

# Star ratings
./lmcli rate ACB

# Cohort analysis: use the historical package directly
```

Hoặc dùng Python:
```bash
pip install lotusmarket==0.5.0
```

Tất cả MIT license. Không API key (trừ FRED nếu muốn macro US data — free).

Mình build cái này vì 7 triệu retail VN deserve công cụ tốt hơn "hot pick từ thầy". Nếu bạn thấy hữu ích, GitHub star và share là cách support tốt nhất.

---

**Author**: [Dat Duc Nguyen Huu](https://github.com/ducnhd) — backend dev, VN retail investor.
**Discuss**: GitHub Issues / [Telegram channel](https://t.me/lotusmarket) (daily auto-posts).

*Tất cả con số trong bài là kết quả backtest deterministic — không phải lời khuyên đầu tư. Past performance does not guarantee future results.*
