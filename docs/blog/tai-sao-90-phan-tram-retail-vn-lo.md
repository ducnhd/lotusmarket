---
title: "Tại sao 90% retail Việt Nam lỗ — và 3 quy tắc đảo ngược (có data thật)"
date: 2026-05-22
tags: [vietnam, stocks, retail, behavior, quant]
---

# Tại sao 90% retail Việt Nam lỗ — và 3 quy tắc đảo ngược

> TL;DR — Statistic Trung tâm Lưu ký Chứng khoán (SSC) hằng năm cho thấy ~85-90% tài khoản retail VN lỗ ròng sau 5 năm giao dịch. Mình chạy backtest 16 năm trên 30 mã VN30 + cohort analysis 47,000 data points để tìm ra **3 quy tắc** ngược với mọi lời khuyên đang lan truyền — và đây là data thật, code MIT open source tại [lotusmarket](https://github.com/ducnhd/lotusmarket).

---

## Lý do gốc rễ: retail làm ngược chiều market

Quan sát 30 ngày plan recommendation trên một bot trading cá nhân (vnstock-bot, source data thật):

| Action | Số lệnh | Đúng (sau 30d) | Tỷ lệ thắng |
|---|---|---|---|
| BUY_SUGGEST | 10 | 9 | **90%** |
| HOLD | 5 | 5 | **100%** |
| SELL / CUT_LOSS | 3 | 0 | **0%** |

Pattern rõ ràng: **mua/giữ đúng, bán/cắt lỗ luôn sai**. Retail bán đáy, giữ đỉnh — ngược với quy luật "buy low sell high" họ tự nói với nhau.

Tại sao? **Tâm lý bear chiếm ưu thế khi giá rớt 15-30%.** Đây là vùng RSI < 40, oversold — và data của mình nói đây là **vùng mua, không phải bán**.

---

## Quy tắc đảo ngược #1: RSI > 70 là tín hiệu MUA (trong uptrend)

Cohort analysis trên ~45,000 row clean VN30+HNX30:

| Joint cohort | Mean fwd return 60d | Win rate | N |
|---|---|---|---|
| **Uptrend × RSI > 70** | **+9.80%** | **63%** | 2,814 |
| Uptrend × RSI 50-70 | +6.21% | 55% | 11,912 |
| Uptrend × RSI 30-50 | +4.00% | 52% | 7,142 |
| Mixed × RSI > 70 | +6.89% | 49% | 657 |
| **Downtrend × RSI > 70** | **-1.40%** | **38%** | 139 |
| Downtrend × RSI 30-50 | +3.48% | 53% | 7,582 |

Baseline 60d = +4.75%. Tức là **uptrend × RSI>70 cho edge +5.05% so với trung bình thị trường** — gần gấp đôi.

**Vùng SELL duy nhất** trong toàn bộ data là `downtrend × RSI>70`: tức stock đang xu hướng giảm (giá < MA200, MA50 < MA200) nhưng có pump ngắn hạn lên overbought. Đây là vùng bear trap — bán ở đây thắng 62% (mất ít hơn baseline).

**Bài học:** đừng nhìn RSI một mình. Phải nhìn cả `MA200 trend` trước. Pattern matching trên Tradingview chỉ ra RSI 78 hay 82 không có giá trị nếu không biết ticker đang uptrend hay downtrend toàn cảnh.

---

## Quy tắc đảo ngược #2: CRISIS không phải rủi ro. Là cơ hội.

Phân loại regime deterministic (no AI) — chỉ dùng VIX + VN-Index change + foreign flow + global contagion:

| Regime | Mean fwd return 60d | Win rate | N | Hành động sai phổ biến |
|---|---|---|---|---|
| **CRISIS** | **+8.53%** | **67%** | 2,582 | Bán đáy panic — thắng 33% |
| EUPHORIA | +6.12% | 59% | 10,279 | Đu đỉnh — thắng 59%, OK |
| VOLATILE | +6.34% | 56% | 7,482 | Tránh — bỏ lỡ +6.34% |
| STABLE | +3.20% | 48% | 23,868 | Trade nhiều — phí ăn hết |

**CRISIS có forward return cao nhất, đồng thời win rate cao nhất.** Vì lúc đó panic seller đã bán hết, định giá lại về fair value, mã tốt bị bán cùng mã xấu.

**Caveat quan trọng**: CRISIS có 2 sub-type — đừng làm ngược cả 2:

- **CRISIS_PANIC** (contagion từ ngoài: VIX spike Mỹ, news tier 1, COVID, war) → mua. Hồi phục 3-12 tháng.
- **CRISIS_FUNDAMENTAL** (vấn đề cấu trúc VN: bank run, kiểm tra hình sự lãnh đạo lớn, BĐS sập) → tránh. Có thể không hồi phục.

Vạn Thịnh Phát tháng 10/2022 là **FUNDAMENTAL** — đúng để tránh. COVID tháng 3/2020 là **PANIC** — đúng để mua. Nếu retail làm ngược 2 cái này, lỗ kép.

---

## Quy tắc đảo ngược #3: Active trading thua buy-hold dài hạn

Backtest 16 năm (2010-01 → 2026-04), 30 mã VN30+HNX30, **phí TPS retail thực tế** (0.15% commission + 0.1% bán + 0.2% slippage + lot 100):

| Strategy | CAGR | MaxDD | Tần suất rebalance |
|---|---|---|---|
| **Buy-hold equal-weight VN30** | **+15.37%** | -36% | 0 (mua đầu, giữ đến cuối) |
| trend_200ma | +21.6% (10y) | -51% | Khi MA200 cross |
| momentum_12_1 | +21.7% (10y) | -44% | Tháng |
| low_vol top-5 | +26.2% (10y) | -42% | Tháng |
| dual_momentum | +15.7% (10y) | -55% | Tháng |

Trên window 10 năm (2017-2026), active strategy trông đẹp. Trên 16 năm (thêm 2010-2016), **buy-hold đánh bại tất cả**.

Lý do toán học:
- Mỗi rebalance lose ~0.35% bị phí + slippage ăn
- Momentum/low_vol rebalance ~12 lần/năm = **4.2%/năm** phí ăn vào alpha
- Alpha thật của các strategy chỉ ~3-6%/năm → bị phí ăn gần hết

**Hệ quả thực tế**: nếu bạn không phải full-time quant, **đừng trade**. Mua đều VN30 ETF (FUEVFVND, E1VFVN30) mỗi tháng, giữ 10+ năm. Năm 16 sau, bạn đánh bại 95% retail trader đang trade hằng tuần.

---

## Vậy "3 quy tắc đảo ngược" cụ thể là gì?

Tổng hợp lại — đây là 3 quy tắc retail VN nên áp dụng ngược lại với mọi khuyên cũ:

### Quy tắc 1: KHÔNG bán đáy khi market panic
**Tại sao retail làm sai**: tâm lý "save what's left" khi -30%.
**Data nói**: CRISIS 60d forward = +8.53%, win 67%. Mua thêm hoặc giữ.
**Action**: setup quỹ "panic buying" 10-20% portfolio, chờ regime CRISIS_PANIC, dồn vào VN30 top quality.

### Quy tắc 2: KHÔNG chốt lời sớm khi mã đang uptrend
**Tại sao retail làm sai**: RSI > 70 thấy "overbought, sắp giảm".
**Data nói**: Uptrend × RSI>70 60d = +9.80%, win 63%. Letting winners run.
**Action**: chỉ exit khi MA50 cross down MA200 (trend đổi), không phải khi RSI cao.

### Quy tắc 3: KHÔNG trade khi không có edge rõ
**Tại sao retail làm sai**: cảm giác "phải làm gì đó" khi mở app thấy danh mục.
**Data nói**: STABLE regime 60d = +3.20%, win 48%. Trade nhiều = phí ăn hết.
**Action**: chỉ rebalance khi có signal mạnh (regime change, breakout MA200, volume surge). Mặc định: ngồi yên.

---

## Cách verify

Tất cả số liệu trong bài này reproducible. Trong 5 phút bạn có thể chạy lại:

```bash
# Python (3 dòng)
pip install lotusmarket==0.5.0
python -c "from lotusmarket.historical import analyze; ..."

# Hoặc Go CLI
go install github.com/ducnhd/lotusmarket/go/cmd/lmcli@v0.5.0
lmcli screen --rsi=60-100  # filter VN30 uptrend × RSI>60
lmcli rate ACB             # 6-dim star ratings
```

Daily reports cập nhật mỗi ngày 15:30 VN tại:

- **Dashboard**: <https://ducnhd.github.io/lotusmarket/>
- **Telegram**: <https://t.me/vnlotusmarket>
- **RSS feed**: <https://ducnhd.github.io/lotusmarket/feed.xml>

Tất cả MIT license, không paid feed, không hot picks, không "thầy".

---

## 🎯 Câu hỏi quan trọng hơn 3 quy tắc

3 quy tắc trên là **tactical** — quan trọng nhưng không đủ. Câu hỏi strategic là:

**"Tại sao tôi đầu tư cổ phiếu?"**

Nếu là vì:
- **"Muốn giàu nhanh"** → 90% xác suất bạn sẽ là 1 trong 85-90% retail lỗ. Đầu tư cổ phiếu không phải lottery. Hãy đặt mục tiêu thực tế (10-15%/năm) hoặc làm việc khác.
- **"Bảo vệ vốn từ lạm phát"** → ok. Buy-hold VN30 ETF. Không cần đọc tiếp.
- **"Học cách hiểu thị trường"** → tốt. Đi học. Backtest. Đọc data. Nhưng KHÔNG bằng tiền thật cho đến khi có edge chứng minh được.
- **"Có thu nhập thay lương"** → cần 2.5-8 tỷ vốn tùy mức sống (mình đã chứng minh bằng Monte Carlo tại [bài 1](16-nam-backtest-vn-stocks.md)). 500 triệu không đủ.

**Câu trả lời "tại sao" quyết định câu trả lời "làm gì"** — không phải ngược lại.

---

**Author**: [Dat Duc Nguyen Huu](https://github.com/ducnhd) — backend dev, VN retail investor 8 năm.

**Discuss**: [Telegram channel](https://t.me/vnlotusmarket) (daily auto-posts) · [GitHub Issues](https://github.com/ducnhd/lotusmarket/issues)

*Past performance does not guarantee future results. Not investment advice.*
