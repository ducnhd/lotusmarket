---
title: "1 chỉ báo vô dụng, 2 chỉ báo có meaning — joint factor matters"
date: 2026-05-22
topic: data-insight
---

# Tại Sao Bạn Vẫn Thua Dù "Đọc Đúng" Chỉ Báo

Kết hợp đúng hai chỉ báo orthogonal tạo ra edge gần 5% trong 60 ngày — trong khi dùng từng cái riêng lẻ chỉ là tiếng ồn được ngụy trang thành tín hiệu.

---

Có một câu hỏi mình hay gặp trong các group cổ phiếu VN: *"RSI đã 72 rồi, bán chưa?"* Câu trả lời ngắn là: sai câu hỏi. Câu hỏi đúng là: *"RSI 72 trong bối cảnh nào?"*

## Chỉ Báo Đơn Lẻ Là Ảo Giác

Hãy bắt đầu với một thí nghiệm tư duy đơn giản. Nếu RSI > 70 thực sự là tín hiệu bán mạnh, thì cổ phiếu trong uptrend mạnh không bao giờ được chạm vùng đó — hoặc nếu chạm, phải sập ngay. Thực tế không diễn ra như vậy.

Cohort analysis trên dữ liệu lịch sử cho thấy điều này rất rõ ràng:

- **RSI > 70 đơn lẻ**: forward return 60 ngày cao hơn baseline **+4.48%**
- **Uptrend (MA200 + MA50) đơn lẻ**: chỉ **+1.20%**
- **MACD bullish đơn lẻ**: **+0.65%** — gần như không phân biệt được với noise

Số đầu tiên sẽ làm nhiều người ngạc nhiên. RSI > 70 không phải tín hiệu bán — ít nhất khi nhìn vào forward return thuần. Nhưng con số **+4.48%** này cũng chưa phải câu chuyện đầy đủ. Nó ẩn đi một sự phân kỳ cực kỳ quan trọng bên dưới.

## Khi Hai Chỉ Báo "Nói Chuyện" Với Nhau

Điều thú vị xảy ra khi kết hợp momentum (RSI) với trend context (MA50/MA200). Đây là hai thứ đo lường **orthogonal** — một cái đo tốc độ thay đổi giá gần đây, một cái đo hướng dài hạn. Khi cả hai đồng thuận, signal tách khỏi noise một cách rõ ràng:

| Combination | 60d Forward Return |
|---|---|
| Uptrend × RSI > 70 | **+5.05%** |
| Uptrend × RSI 50-70 | +1.45% |
| RSI > 70 đơn lẻ | +4.48% |
| Uptrend đơn lẻ | +1.20% |
| **Downtrend × RSI > 70** | **-6.15%** |

Con số cuối cùng là quan trọng nhất trong bảng này.

**Downtrend × RSI > 70 = -6.15%** — đây là signal SELL duy nhất xuất hiện trong toàn bộ cohort analysis. Không phải RSI > 70 đơn lẻ. Không phải "RSI overbought." Mà là RSI overbought *trong downtrend*.

Điều này giải thích tại sao cộng đồng retail VN hay bị lệch: họ thấy RSI > 70 và bán, bỏ lỡ toàn bộ uptrend continuation. Hoặc họ thấy RSI > 70 và mua thêm trong downtrend, rồi bị -6% trước khi kịp nhận ra mình nhầm.

## Câu Chuyện Thực Sự Của Số Liệu

Hãy hình dung hai kịch bản:

**Kịch bản A**: Một cổ phiếu đang trong uptrend rõ ràng (close > MA50 > MA200), momentum tăng mạnh (RSI vượt 70). Cohort lịch sử của nhóm này đạt **+5.05%** trong 60 ngày — gấp 4 lần so với chỉ nhìn MA trend đơn lẻ, gấp 1.1 lần RSI đơn lẻ. Trend + momentum đồng thuận = signal thực.

**Kịch bản B**: Cùng RSI > 70, nhưng cổ phiếu đang trong downtrend (close < MA50 hoặc MA50 < MA200). Cohort này có forward return **-6.15%**. Con số cùng chỉ báo, bối cảnh khác, kết quả ngược chiều hoàn toàn.

Cùng một con số RSI. Hai kết quả cách nhau hơn **11 percentage points**. Đây không phải sự trùng hợp — đây là lý do tại sao context không phải optional, nó là *bắt buộc*.

## Anti-Pattern Phổ Biến Nhất Trong Cộng Đồng VN 🎯

Mình không muốn phán xét ai, nhưng có ba myth lan truyền rộng đến mức cần nói thẳng:

**"RSI > 70 = bán"** — Như đã thấy ở trên, RSI > 70 trong uptrend có forward return **+5.05%**. Bán theo reflexe này là bán đúng cổ phiếu tốt nhất trong portfolio.

**"Phá MA200 = mua"** — MA trend đơn lẻ chỉ cho edge **+1.20%**. Không có momentum confirmation, đây là tín hiệu yếu. Breakout MA200 trong môi trường RSI thấp/neutral là khác hoàn toàn với breakout kèm momentum mạnh.

**"Volume tăng + giá tăng = mua"** — Volume là noisy variable. Trong cohort analysis, Wyckoff stage × volume class chỉ cho edge **±2%**, thấp hơn đáng kể so với combo momentum × trend (**±5%**). Volume cần được đặt trong Wyckoff framework, không dùng standalone.

Tất cả những anti-pattern này đều có chung một lỗi cấu trúc: **1 indicator + 0 context = noise được đặt tên đẹp**.

## Tại Sao Không Thêm Chỉ Báo Thứ 3?

Câu hỏi tự nhiên tiếp theo: nếu 2 tốt hơn 1, tại sao không dùng 3 hay 4?

Ba lý do:

**Sample size collapse.** Mỗi điều kiện thêm vào thu hẹp cohort. Khi còn 50-100 case lịch sử, con số trở nên không ổn định — một vài outlier có thể làm lệch toàn bộ kết quả.

**Transaction cost thực ăn vào edge.** Edge +5% trong 60 ngày sau khi trừ phí giao dịch VN (spread + thuế TNCN + phí môi giới) còn khoảng 3-4%. Thêm điều kiện thứ 3 thu hẹp cơ hội giao dịch, không nhất thiết tăng edge net.

**Overfitting.** Ba điều kiện fit đẹp trên in-sample data, chạy out-of-sample thì edge biến mất. Đây là classic curse of dimensionality trong quant trading.

Ba best combos từ cohort data, theo thứ tự edge:
1. **MA trend × RSI bucket** → edge ±5%
2. **Market regime × MA trend** → edge ±3-4%
3. **Wyckoff stage × volume class** → edge ±2%

Stick với combo đầu tiên cho đến khi có đủ dữ liệu để validate combo thứ hai.

## Framework Thực Hành

Mỗi lần nhìn vào một cổ phiếu, hãy tự hỏi theo thứ tự:

1. **Trend context trước**: Close đang ở đâu so với MA50 và MA200? Uptrend / downtrend / sideways?
2. **Momentum sau**: RSI đang ở bucket nào? (<30 / 30-50 / 50-70 / >70)
3. **Hai con số này đồng thuận hay mâu thuẫn?**

Nếu đồng thuận (ví dụ: Uptrend + RSI > 70), data cho thấy cohort lịch sử thắng baseline **+5.05%** trong 60 ngày. Nếu mâu thuẫn (Downtrend + RSI > 70), lịch sử cho thấy cohort đó thua baseline **-6.15%**.

Đây không phải prediction. Đây là base rate từ lịch sử — dùng để calibrate xác suất, không phải đảm bảo kết quả.

---

## 3 Takeaway

> **1.** RSI > 70 trong uptrend = +5.05% edge. RSI > 70 trong downtrend = -6.15%. Cùng con số, context khác nhau, kết quả ngược chiều.
>
> **2.** Hai chỉ báo orthogonal (momentum × trend) tạo signal thực. Một chỉ báo đơn lẻ — dù là RSI, MA, hay MACD — phần lớn là noise được đặt tên đẹp.
>
> **3.** Thêm điều kiện thứ 3 trở đi làm giảm sample size, tăng phí, và tăng rủi ro overfitting. Hai là đủ.

---

## Verify Reproducible

Toàn bộ cohort analysis này có thể reproduce được qua [Lotus Market](https://github.com/ducnhd/lotusmarket) hoặc xem thêm tại [lotusai.servehttp.com](https://lotusai.servehttp.com).

```bash
pip install lotusmarket pandas numpy

```

```python
from lotusmarket import CohortAnalysis

cohort = CohortAnalysis(universe="HOSE", lookback_days=252)
cohort.run(factors=["ma_trend", "rsi_bucket"], forward_days=60)
cohort.plot_heatmap()
```

---

*Bài viết này là phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Mọi quyết định mua bán là trách nhiệm của nhà đầu tư.*
