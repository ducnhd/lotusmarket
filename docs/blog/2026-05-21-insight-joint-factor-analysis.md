---
title: "1 chỉ báo vô dụng, 2 chỉ báo có meaning — joint factor matters"
date: 2026-05-21
topic: data-insight
---

# Tại Sao Bạn Vẫn Thua Dù "Đọc Đúng" Tín Hiệu?

Một chỉ báo kỹ thuật đứng một mình không phải signal — đó là noise được đóng gói đẹp.

---

Hãy thử một thí nghiệm tư duy. Bạn mở chart, RSI chạm 72, tay đã chuẩn bị bấm bán. Người kia nhìn cùng mã đó, thấy MA50 đang nằm trên MA200, giá bứt phá khỏi vùng tích lũy 3 tuần — và họ mua. Sáu mươi ngày sau, ai đúng?

Câu trả lời không nằm ở việc ai "đọc biểu đồ tốt hơn." Nó nằm ở dữ liệu cohort. Và cohort nói rất rõ ràng.

---

## Một Chỉ Báo Đứng Một Mình: Bạn Đang Chơi Xúc Xắc

Phân tích cohort trên dữ liệu thực cho thấy ba chỉ báo phổ biến nhất trong cộng đồng VN đều có *edge dương* khi nhìn riêng lẻ — nghe có vẻ tốt. Nhưng đây là vấn đề:

| Chỉ báo đơn lẻ | Edge 60 ngày so với baseline |
|---|---|
| RSI > 70 | +4.48% |
| Uptrend (MA200 + MA50) | +1.20% |
| MACD bullish | +0.65% |

MACD bullish cho edge **+0.65%** sau 60 ngày. Phí giao dịch trên sàn VN đã ngốn hết phần đó rồi, chưa kể spread và thuế. RSI > 70 trông ổn hơn với +4.48%, nhưng đây là trung bình toàn bộ thị trường — bao gồm cả lúc thị trường đang uptrend lẫn downtrend. Khi bạn tách ra, câu chuyện thay đổi hoàn toàn.

Điểm mấu chốt: **một chỉ báo đo một chiều của giá.** RSI đo momentum. MA đo trend. MACD đo divergence momentum. Không cái nào trong số đó biết "bức tranh toàn cảnh" là gì — đó là việc của người kết hợp chúng lại.

---

## Khi Hai Chỉ Báo Vuông Góc Gặp Nhau: Edge Nhân Lên

"Vuông góc" ở đây theo nghĩa toán học: hai chỉ báo đo *hai chiều khác nhau*, không overlap thông tin. MA trend cho biết thị trường đang đi đâu về cấu trúc dài hạn. RSI cho biết momentum ngắn hạn đang ở mức nào. Kết hợp chúng tạo ra một bức tranh 2D — và bức tranh đó thay đổi mọi thứ.

Nhìn vào số liệu joint factor:

**Uptrend × RSI > 70: +5.05%**

So sánh:
- MA trend đơn lẻ: +1.20%
- RSI > 70 đơn lẻ: +4.48%
- Kết hợp cả hai: +5.05%

Uptrend × RSI > 70 **gấp 4 lần** MA trend đứng một mình và **vượt RSI đơn lẻ thêm 12%.** Không phải cộng dồn — mà là tương tác. Hai chỉ báo xác nhận lẫn nhau, loại bỏ trường hợp nhiễu, và cohort còn lại là những setup chất lượng hơn hẳn.

Nhưng combination đáng chú ý nhất không phải là combo tốt nhất — mà là combo **tệ nhất:**

**Downtrend × RSI > 70: -6.15%** 🚨

Đây là signal SELL duy nhất xuất hiện trong toàn bộ phân tích cohort. RSI > 70 trong downtrend không phải "overbought rồi sẽ điều chỉnh nhẹ" — đó là **dead cat bounce** trong xu hướng giảm, và lịch sử cho thấy cohort này thua baseline tới 6.15% sau 60 ngày. Người mua ở đây không phải "bắt đáy sớm" — họ đang đứng dưới dao đang rơi.

---

## Anti-Pattern Của Cộng Đồng Retail VN

Đây là phần sẽ khiến một số người khó chịu, nhưng data không quan tâm đến cảm xúc.

**"RSI > 70 = bán"** — Myth phổ biến nhất. Trong uptrend, RSI > 70 cho edge dương +5.05%. Bán ở đây nghĩa là bỏ lại 5% trên bàn. RSI cao trong uptrend là dấu hiệu momentum mạnh, không phải tín hiệu đảo chiều.

**"Phá MA200 = mua"** — MA200 đứng một mình chỉ cho +1.20% edge. Không đủ để cover phí và rủi ro. Câu hỏi phải hỏi là: momentum lúc đó ở đâu? RSI đang ở bucket nào? Nếu RSI 50-70 trong uptrend, cohort cho +1.45% — vẫn khiêm tốn. Nếu RSI > 70 trong uptrend, nhảy lên +5.05%.

**"Volume tăng + giá tăng = mua"** — Volume là chỉ báo noisy bậc nhất trên thị trường VN với thanh khoản mỏng và nhiều phiên khớp lệnh cuối ngày bất thường. Cohort Wyckoff stage × volume class chỉ cho edge ±2% — thấp nhất trong ba combo tốt nhất. Volume có giá trị, nhưng đứng một mình nó là nhiễu.

Pattern chung của cả ba anti-pattern: **chỉ dùng một chiều thông tin để ra quyết định hai chiều (mua/bán).** Đó là lý do thua.

---

## Ba Combo Tốt Nhất Từ Cohort Analysis

Data từ phân tích cho thấy ba combination có edge cao nhất và đáng tin cậy nhất:

**1. MA trend × RSI bucket → edge ±5%**
Combo mạnh nhất. Kết hợp cấu trúc xu hướng dài hạn (MA200 vs MA50) với trạng thái momentum ngắn hạn (RSI bucket). Đây là hai chiều thực sự vuông góc — không overlap, thông tin bổ sung cho nhau.

**2. Regime × MA trend → edge ±3-4%**
Regime (CRISIS / STABLE / EUPHORIA) đóng vai trò macro filter. Cùng một MA uptrend trong giai đoạn STABLE cho kết quả khác hoàn toàn so với trong CRISIS. Macro context làm thay đổi xác suất cơ bản của toàn bộ cohort.

**3. Wyckoff stage × volume class → edge ±2%**
Edge thấp hơn, nhưng có giá trị như confirmation layer khi hai combo kia đã xác nhận. Đứng một mình: noise. Kết hợp như tầng thứ ba: tăng conviction.

---

## Tại Sao Không Thêm Điều Kiện Thứ Ba, Thứ Tư?

Nghe có vẻ logic: thêm filter → setup sạch hơn → edge cao hơn. Thực tế ngược lại.

Khi thêm mỗi điều kiện, sample size giảm theo cấp số nhân. Với 3+ điều kiện, cohort còn lại thường quá nhỏ để kết luận có ý nghĩa thống kê — kết quả tốt trên backtest là do **may mắn của sample nhỏ, không phải edge thực.** Đây gọi là overfitting.

Chưa kể: mỗi điều kiện thêm vào đồng nghĩa với ít cơ hội giao dịch hơn, phí cố định ăn vào edge thực tế nhiều hơn theo tỷ lệ, và thời gian nắm giữ chờ đủ điều kiện kéo dài hơn.

**Hai điều kiện vuông góc là sweet spot.** Đây không phải quan điểm cá nhân — cohort data khẳng định điều này.

---

## Verify Tự Mình

Toàn bộ phân tích này có thể tái tạo với dữ liệu công khai. Không cần tin vào bài viết — chạy và kiểm tra:

```bash
pip install lotusmarket pandas numpy
```

```python
from lotusmarket import cohort_analysis

# Load data và chạy joint factor analysis
result = cohort_analysis(
    factor_a="ma_trend",      # uptrend / downtrend
    factor_b="rsi_bucket",    # <30 / 30-50 / 50-70 / >70
    forward_window=60
)
print(result.pivot_table(values="fwd_return", index="factor_a", columns="factor_b"))
```

Chi tiết về methodology và source code tại [github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket). Dashboard trực quan tại [lotusai.servehttp.com](https://lotusai.servehttp.com).

---

## Ba Điều Cần Nhớ

1. **Một chỉ báo = noise.** RSI alone, MA alone, MACD alone — tất cả đều không đủ để cover phí giao dịch một cách đáng tin cậy.
2. **Hai chỉ báo vuông góc = signal.** MA trend × RSI bucket cho edge ±5% — đủ có ý nghĩa thực tế. Downtrend × RSI > 70 là warning rõ ràng nhất mà cohort tạo ra.
3. **Ba+ điều kiện = overfitting.** Thêm filter không tăng edge — nó giảm sample size và tăng rủi ro pathological backtest.

---

*Bài viết mang tính phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Mọi quyết định giao dịch là trách nhiệm cá nhân của người thực hiện.*
