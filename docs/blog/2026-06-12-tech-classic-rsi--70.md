---
title: "RSI > 70 có thực sự là tín hiệu BÁN? Backtest 9 năm VN30 nói ngược"
date: 2026-06-12
topic: tech-classic
---

# RSI > 70 Không Phải Tín Hiệu Bán — Data 45,000 Dòng Nói Gì Khác Hẳn

Cộng đồng chứng khoán Việt Nam đã sai về RSI overbought suốt nhiều năm — và con số +9.23% return sau 60 ngày là bằng chứng khó chối cãi nhất.

---

## Cái Bẫy Tâm Lý Tên "Quá Mua"

Hỏi bất kỳ trader bán lẻ nào ở các group Facebook chứng khoán: "RSI vượt 70 thì làm gì?" — câu trả lời gần như đồng thanh: *"Chốt lời ngay, cổ phiếu đang overbought rồi."*

Nghe có vẻ hợp lý. RSI > 70 = quá mua = sắp đảo chiều = bán. Logic tròn trịa, dễ nhớ, dễ dạy cho người mới. Vấn đề duy nhất: nó sai.

Không phải sai theo kiểu "đôi khi không hiệu quả". Sai theo kiểu backtest trên **45,000 dòng dữ liệu thực** từ VN30 và HNX30 giai đoạn 2017–2026 cho thấy RSI > 70 là **cohort có return cao nhất** trong toàn bộ phân tích — trung bình **+9.23%** trong 60 ngày giao dịch tiếp theo, tỷ lệ thắng 60%, edge so với baseline lên đến **+4.48%**.

Baseline của toàn thị trường trong cùng kỳ? +4.75%, win rate 53%.

Nói cách khác: thay vì bán khi RSI vượt 70, data lịch sử cho thấy đây là vùng *outperform* mạnh nhất.

---

## Bốn Vùng RSI, Một Kết Quả Bất Ngờ

Phân tích cohort chia toàn bộ 45,000 quan sát thành bốn nhóm RSI, đo forward return 60 ngày:

| RSI Bucket | N | Return 60d | Win Rate | Edge vs Baseline |
|---|---|---|---|---|
| < 30 (oversold) | 1,755 | +6.32% | 60% | +1.57% |
| 30–50 | 18,933 | +3.67% | 52% | -1.08% |
| 50–70 | 20,681 | +4.82% | 53% | +0.07% |
| **> 70 (overbought)** | **4,015** | **+9.23%** | **60%** | **+4.48%** |

Nhìn vào đây một lúc. Vùng "nguy hiểm nhất" theo lý thuyết truyền thống lại có return cao nhất. Vùng RSI 30–50 — nơi nhiều người cho là "an toàn, chờ mua" — thực ra là vùng **dưới baseline**, edge âm -1.08%.

Vùng oversold (RSI < 30) cũng tốt, +6.32% và edge +1.57% — nhưng vẫn thua overbought gần 3 điểm phần trăm.

---

## Tại Sao? Context Là Tất Cả

Đây là chỗ phân tích trở nên thú vị thực sự.

Khi tách thêm biến MA trend — cụ thể là điều kiện uptrend (close > MA200 *và* MA50 > MA200) so với downtrend — bức tranh rõ hơn nhiều:

**Uptrend × RSI > 70** (N=2,814): forward return **+9.80%**, win rate **63%**, edge **+5.05%**

**Downtrend × RSI > 70** (N=139): forward return **-1.40%**, win rate **38%**, edge **-6.15%**

Đây là sự phân kỳ lớn nhất trong toàn bộ dataset. Cùng một tín hiệu RSI, hai bối cảnh trend khác nhau cho ra hai kết quả gần như đối lập hoàn toàn.

Downtrend × RSI > 70 là **cohort duy nhất có return âm** trong toàn bộ phân tích. Win rate chỉ 38% — tức là cầm cổ phiếu 60 ngày, 6/10 lần thua. Nếu có một vùng xứng đáng được gọi là "tín hiệu cẩn trọng", đây mới là nó.

Còn uptrend × RSI > 70? Lịch sử cho thấy đây là một trong những setup có xác suất thuận lợi nhất — không phải vì RSI cao, mà vì RSI cao *trong bối cảnh xu hướng đang mạnh*. Momentum có xu hướng tiếp diễn, không đảo chiều ngay lập tức. 📊

---

## Cái Sai Của Myth Và Tại Sao Nó Dai Dẳng

Myth "RSI > 70 = bán" không phải vô căn cứ hoàn toàn. Nó đúng trong một điều kiện cụ thể: **khi xu hướng đang yếu hoặc đảo chiều**. Downtrend × RSI > 70 với edge -6.15% là bằng chứng rõ ràng.

Vấn đề là myth này được dạy và áp dụng *không phân biệt context* — bất kể cổ phiếu đang trong uptrend mạnh hay sideways hay downtrend. Kết quả là rất nhiều người bán sớm những cổ phiếu đang trong momentum tốt nhất, rồi đứng nhìn chúng tiếp tục tăng thêm 10–15% trong hai tháng tiếp theo.

Indicator không sai. Cách dùng indicator mới là vấn đề.

---

## Một Góc Nhìn Khác: Sample Size Quan Trọng

Cần lưu ý: downtrend × RSI > 70 chỉ có N=139 — nhỏ hơn đáng kể so với uptrend × RSI > 70 với N=2,814. Finding về vùng downtrend là **sơ bộ và cần thận trọng hơn** khi suy rộng. Xu hướng rõ ràng, nhưng confidence interval trên 139 quan sát sẽ rộng hơn nhiều.

Ngược lại, uptrend × RSI > 70 với gần 3,000 quan sát qua 9 năm là một dataset đủ lớn để các con số có ý nghĩa thống kê thực sự. 🔍

---

## Ba Takeaway Từ Data

**1. RSI > 70 trong uptrend là cohort thắng, không phải tín hiệu thoát.**
Lịch sử VN30+HNX30 (2017–2026) cho thấy setup này return +9.80% trung bình sau 60 ngày, win rate 63%.

**2. Vùng duy nhất data ủng hộ cẩn trọng là Downtrend × RSI > 70.**
Return -1.40%, win rate 38%, edge -6.15% — nhưng sample size 139 nên cần quan sát thêm.

**3. Context trend quan trọng hơn giá trị RSI đơn lẻ.**
Cùng RSI > 70, uptrend cho +9.80% còn downtrend cho -1.40%. Dùng indicator mà thiếu context là bỏ đi một nửa thông tin.

---

## Verify Reproducible

Dataset và methodology có thể kiểm tra lại qua [Lotus Market](https://github.com/ducnhd/lotusmarket):

```bash
pip install lotusmarket pandas numpy
```

```python
from lotusmarket import load_cohort, rsi_bucket_analysis

df = load_cohort(universe=["VN30", "HNX30"], start="2017-01-01", end="2026-01-01")
results = rsi_bucket_analysis(df, forward_days=60, ma_filter=True)
print(results.summary())
```

Hoặc xem thêm công cụ phân tích tại [lotusai.servehttp.com](https://lotusai.servehttp.com).

---

*Disclaimer: Toàn bộ nội dung trên là phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Past performance không đảm bảo kết quả tương lai.*
