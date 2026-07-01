---
title: "\"Mua đáy bán đỉnh\" — câu thần chú phổ biến nhất, và sai nhất với VN30"
date: 2026-07-01
topic: myth-buster
---

# "Mua đáy bán đỉnh" — Chiến lược nghe có lý nhưng data nói ngược lại

Retail VN bán đúng lúc nên giữ, mua đúng lúc nên tránh — và backtest 45,000 điểm dữ liệu trên VN30+HNX30 chứng minh điều đó bằng con số cụ thể.

---

Có một câu hỏi khá khó chịu: nếu "mua đáy bán đỉnh" là chiến lược đúng, tại sao đại đa số nhà đầu tư cá nhân lại thua thị trường?

Câu trả lời không nằm ở việc họ thiếu kỷ luật. Nó nằm ở chỗ **định nghĩa "đáy" và "đỉnh" hoàn toàn sai**.

---

## Câu chuyện quen thuộc ở mọi group chứng khoán

Cổ phiếu X rớt 15% trong hai tuần. RSI xuống dưới 30. Dòng comment bắt đầu xuất hiện: *"Oversold rồi, mua đáy thôi anh em."* Một tuần sau, cổ phiếu tăng thêm 20%. Dòng comment mới: *"RSI vượt 70 rồi, overbought, chốt lời thôi."*

Logic này nghe rất hợp lý. Nó cũng sai theo một cách rất có hệ thống.

---

## Số liệu thực sự nói gì

Backt test trên toàn bộ VN30 và HNX30, khoảng thời gian 9 năm, với 45,493 quan sát, baseline forward return 60 ngày là **+4.75%**.

Khi lọc theo RSI:

| Điều kiện | Forward return 60d | Win rate |
|---|---|---|
| RSI < 30 (oversold / "đáy") | +6.32% | 60% |
| RSI > 70 (overbought / "đỉnh") | +9.23% | 60% |

Dừng lại một chút.

Vùng "đỉnh" — nơi retail ùn ùn chốt lời — lại mang về return **cao hơn vùng "đáy"** tới 3 điểm phần trăm, với win rate tương đương. Nói cách khác, thống kê không ủng hộ việc bán khi RSI > 70. Nó ủng hộ điều ngược lại.

---

## Nhưng quan trọng hơn: MA trend là biến quyết định

Đây là phần mà hầu hết người dùng RSI bỏ qua hoàn toàn. Khi kết hợp RSI với MA trend (uptrend / downtrend), bức tranh trở nên rõ ràng hơn nhiều:

| Điều kiện kết hợp | Forward return 60d | Win rate |
|---|---|---|
| Uptrend × RSI > 70 | **+9.80%** | **63%** |
| Downtrend × RSI > 70 | **-1.40%** | **38%** |

Hai vùng đều có RSI > 70. Hai kết quả hoàn toàn đối lập.

Uptrend + RSI > 70 = momentum thực sự, letting winner run. Downtrend + RSI > 70 = dead cat bounce, đây là **vùng SELL duy nhất có căn cứ thống kê** trong toàn bộ dataset.

Vấn đề là retail VN không phân biệt hai trường hợp này. Họ thấy RSI > 70 là bán — và vô tình bán đúng lúc cần giữ (uptrend), đôi khi giữ đúng lúc cần bán (downtrend × RSI > 70 nhưng nhầm tưởng là "chưa đủ đỉnh").

---

## 30 ngày thực chiến: số liệu audit cá nhân 😬

Để kiểm tra thực tế, mình audit lại 18 khuyến nghị trong 30 ngày gần nhất.

Kết quả:

- **BUY/HOLD: đúng 14/15 (93%)**
- **SELL/CUT\_LOSS: đúng 0/3 (0%)**

Ba lần khuyến nghị bán, ba lần sai. Không phải ngẫu nhiên — đây là hệ quả trực tiếp của việc dùng RSI > 70 làm tín hiệu bán mà không lọc qua MA trend. Cả ba trường hợp đó đều là **uptrend × RSI > 70** — tức vùng mà data cho thấy forward return đạt +9.80% và win rate 63%.

Đây không phải xui xẻo. Đây là structural error, lặp lại có hệ thống.

---

## Tại sao myth này tồn tại dai dẳng?

Có lý do tâm lý rõ ràng. Bán khi giá cao cho cảm giác *thông minh*. Mua khi giá thấp cho cảm giác *can đảm*. Cả hai đều thỏa mãn narrative — nhưng thị trường không trả tiền cho narrative, nó trả tiền cho xác suất kỳ vọng.

RSI là oscillator — công cụ đo momentum ngắn hạn. Nó không được thiết kế để làm tín hiệu trend reversal độc lập. Dùng RSI một mình để "bắt đáy đỉnh" là dùng nhiệt kế để chẩn đoán bệnh tim: không sai về mặt dụng cụ, nhưng hoàn toàn sai về mặt ứng dụng.

---

## Điều data thực sự gợi ý 📊

*(Đây không phải lời khuyên đầu tư — xem disclaimer cuối bài)*

1. **RSI < 30 không tệ** — forward return +6.32%, win rate 60% so với baseline +4.75%. Nhưng N chỉ có 1,755 quan sát trong 9 năm, finding này cần thêm dữ liệu để xác nhận chắc chắn hơn.

2. **RSI > 70 trong uptrend**: data cho thấy cohort này thắng với +9.80% và win rate 63%. Bán vội ở đây là để lại tiền trên bàn.

3. **Vùng SELL có căn cứ duy nhất**: Downtrend × RSI > 70 — forward return -1.40%, win rate chỉ 38%. Đây là nơi đáng cân nhắc giảm vị thế, không phải khi thấy RSI vượt ngưỡng 70 bất kể context.

---

## Takeaway

> **1.** "Mua đáy bán đỉnh" theo RSI đơn thuần đi ngược lại xác suất kỳ vọng — RSI > 70 lịch sử cho return cao hơn RSI < 30 trên cùng dataset.

> **2.** MA trend là biến phân tách signal thực khỏi noise: cùng RSI > 70, uptrend cho +9.80% win 63%, downtrend cho -1.40% win 38%.

> **3.** Structural error trong SELL call không phải lỗi phán đoán ngẫu nhiên — nó là hệ quả tất yếu khi dùng RSI không có context trend.

---

## Verify reproducible

Toàn bộ backtest methodology có thể kiểm tra qua [Lotus AI](https://lotusai.servehttp.com) hoặc source tại [github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

```bash
pip install lotusmarket pandas numpy
```

```python
from lotusmarket import backtest, universe

df = universe.load("VN30+HNX30", years=9)
result = backtest.rsi_ma_joint(df, rsi_low=30, rsi_high=70, fwd_days=60)
print(result.summary())
```

---

*Bài viết phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Mọi quyết định mua bán là trách nhiệm của nhà đầu tư.*
