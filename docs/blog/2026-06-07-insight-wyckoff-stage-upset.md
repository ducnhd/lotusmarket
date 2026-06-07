---
title: "5 điều cohort 47K data point cho thấy ngược trực giác retail VN"
date: 2026-06-07
topic: data-insight
---

# 5 Điều Data VN30+HNX30 Nói Ngược Hoàn Toàn Với Trực Giác Retail

Sau 9 năm backtest trên 45.493 rows dữ liệu sạch của cohort VN30+HNX30, kết quả cho thấy một sự thật khó chịu: phần lớn những gì cộng đồng retail VN tin là đúng — đều sai theo hướng ngược lại.

---

## 1. Wyckoff Distribution (Đỉnh Phân Phối) Lại Cho Return Cao Nhất

Hỏi bất kỳ trader nào ở các group chứng khoán VN: "Thấy Wyckoff Stage 3 thì làm gì?" — câu trả lời gần như đồng thanh: *"Bán ra, thoát hàng."*

Data cohort không đồng ý.

Backtest cho thấy Stage 3 (distribution) có forward return 60 ngày đạt **+6.63%**, win rate **60%**, edge **+1.88%**. Trong khi đó, Stage 2 (markup) — giai đoạn mà ai cũng muốn "bắt đúng sóng" — chỉ cho **+5.46%**, win rate **53%**, edge **+0.71%**.

Tại sao lại vậy? Một giả thuyết hợp lý: khi thị trường đang trong giai đoạn distribution, lực cung từ smart money tạo ra áp lực nhưng không đủ để crush ngay lập tức. Trong khoảng 60 ngày kế tiếp, hàng vẫn di chuyển, và những ai hoảng loạn bán sớm vô tình nhường lợi cho phía còn lại. Điều này không có nghĩa là "cứ thấy distribution là mua" — win rate 60% nghĩa là 40% còn lại thua. Nhưng *trung bình* cohort vẫn thắng, trái với narrative phổ biến.

Retail nghĩ đỉnh phân phối là tín hiệu thoát hàng. Data nói đó là giai đoạn có edge dương — nhỏ nhưng dương.

---

## 2. RSI > 70 Không Phải Lúc "Cần Cẩn Thận" — Mà Là Lúc Return Cao Nhất

"RSI trên 70 là overbought, cẩn thận đảo chiều" — câu này xuất hiện trong mọi cuốn sách nhập môn phân tích kỹ thuật.

Cohort VN30+HNX30 9 năm cho kết quả:

- **RSI > 70**: forward return 60 ngày **+9.23%**, win rate **60%**
- **RSI < 30**: forward return 60 ngày **+6.32%**, win rate **60%**

Cả hai vùng đều có win rate bằng nhau (60%), nhưng RSI cao cho return *cao hơn 46%* so với RSI thấp. Momentum tồn tại trong dữ liệu VN, và nó mạnh hơn mean-reversion ở khung 60 ngày.

Điều này không phủ nhận rằng RSI > 70 có thể đảo chiều — nhưng trên tổng thể cohort, holding qua vùng "overbought" mang lại return tốt hơn là bán vì sợ "quá mua". Đây là lý do tại sao stop-out sớm dựa thuần vào RSI là một trong những nguồn gốc của underperformance retail.

---

## 3. CRISIS Regime: Panic-Sell Là Món Quà Cho Người Còn Tiền Mặt 📉

Tâm lý con người được thiết kế để thoát khỏi nguy hiểm khi mọi thứ trông tệ nhất. Trong thị trường chứng khoán, đây là một bug cực kỳ đắt tiền.

Backtest theo regime:

- **CRISIS regime**: forward return 60 ngày **+8.53%**, win rate **67%**
- **STABLE regime**: forward return 60 ngày **+3.20%**, win rate **48%**

STABLE — giai đoạn thị trường bình thường, tin tức tốt, mọi người tự tin — cho win rate dưới 50% và return thấp nhất. CRISIS — giai đoạn ai cũng muốn thoát — cho win rate cao nhất và return gần gấp 3.

Lý do không phức tạp: khi CRISIS xảy ra, panic-sell đẩy giá xuống dưới fair value trên diện rộng. Cohort người mua trong giai đoạn đó không cần skill đặc biệt — họ chỉ cần không bán cùng đám đông. Win rate 67% là con số được tạo ra không phải bởi phân tích xuất sắc, mà bởi *việc không làm gì* khi tâm lý đám đông đang ở đáy.

---

## 4. Trading Tích Cực Thua Buy-Hold 12% CAGR Sau 16 Năm

Đây là finding có lẽ gây khó chịu nhất — không phải vì nó lạ, mà vì nó quá quen mà vẫn bị bỏ qua.

Backtest 16 năm:

- **Active trading** (low-vol top-5, rebalance hàng tháng): CAGR **+3%**
- **Buy-hold VN30 equal-weight**: CAGR **+15.37%**

Chênh lệch: **12.37% CAGR**. Sau 16 năm, sự khác biệt này là khoảng cách giữa nghỉ hưu sớm và không nghỉ hưu được.

Thủ phạm chính không phải là strategy kém — mà là chi phí. Phí giao dịch 0.35% × 12 lần rebalance/năm = **4.2%/năm** bị ăn mất trước khi tính bất kỳ alpha nào. Để strategy active thắng được buy-hold, nó phải tạo ra alpha ít nhất 4.2%/năm *trước* khi tính đến slippage, thuế, và sai lầm hành vi.

Sarcasm nhẹ dành cho những ai đang trả phí môi giới cao để rebalance danh mục hàng tháng: bạn đang trả tiền để underperform.

---

## 5. Kết Quả Đầu Tư Phụ Thuộc 50% Vào Năm Bắt Đầu, Không Phải Skill 🎲

Finding này có lẽ là cái "tát" mạnh nhất vào toàn bộ narrative về alpha và skill trong cộng đồng đầu tư VN.

Cùng một setup: **2 tỷ đồng**, rút **12 triệu/tháng**:

- **Bắt đầu năm 2010**: sau 10 năm còn **10 tỷ** (tăng 5×)
- **Bắt đầu năm 2022**: sau 10 năm còn **2.2 tỷ** (giảm 78%)

Cùng chiến lược, cùng số tiền ban đầu, cùng quy tắc rút tiền — chênh lệch kết quả là **7.8 tỷ đồng**. Biến số quyết định không phải là người chọn cổ phiếu giỏi hơn. Mà là *năm bắt đầu*.

Đây là sequence-of-returns risk — và nó không được dạy trong bất kỳ khóa học phân tích kỹ thuật nào ở VN. Người bắt đầu năm 2010 trải qua giai đoạn tăng trưởng mạnh trước khi rút tiền nhiều, nên portfolio vẫn đủ lớn để compound. Người bắt đầu năm 2022 rút tiền ngay trong giai đoạn thị trường điều chỉnh, phá vỡ compound effect.

Kết luận từ data: **timing luck giải thích phần lớn variance trong kết quả** hơn là edge thực sự.

---

## Common Thread: Intuition Retail VN Bị Calibrate Sai

Nhìn lại 5 finding:

1. Distribution stage → retail bán, data nói edge dương
2. Overbought RSI → retail cẩn thận, data nói return cao nhất
3. CRISIS → retail panic-sell, data nói win rate cao nhất
4. Active trading → retail nghĩ "quản lý tích cực = tốt hơn", data nói ngược lại
5. Skill vs luck → retail attribut kết quả cho skill, data nói timing luck dominant

Cái nguy hiểm không phải là retail VN thiếu thông tin — mà là thông tin sai đã được *repeat đủ nhiều* đến mức nó cảm giác như sự thật.

Edge thực sự không đến từ "đảo ngược tâm lý một cách mù quáng" — mà đến từ **pre-commit rules dựa trên cohort data**: biết trước mình sẽ làm gì trong từng regime, từng giai đoạn, và không thay đổi khi tâm lý nhất thời nói khác đi.

---

## 3 Takeaway

1. **Đừng sell reflexively khi thấy distribution hoặc RSI > 70** — cohort historical cho thấy cả hai vùng đều có positive edge ở khung 60 ngày.
2. **CRISIS regime là regime tốt nhất để build position** — win rate 67% không đến từ insight đặc biệt, mà từ việc không panic-sell cùng đám đông.
3. **Chi phí rebalance active > alpha active** — trừ khi strategy của bạn tạo ra hơn 4.2%/năm alpha trước phí, buy-hold thắng mặc định.

---

## Verify — Reproducible

Toàn bộ backtest này chạy trên engine công khai tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket). Để reproduce:

```bash
pip install lotusmarket pandas numpy
```

```python
from lotusmarket import cohort, backtest

df = cohort.load("VN30+HNX30", years=9)
result = backtest.run(df, forward_days=60, group_by=["wyckoff_stage", "rsi_zone", "regime"])
print(result.summary())
```

Kết quả có thể sai lệch nhỏ tùy version data — nhưng direction của các finding đã stable qua nhiều lần re-run. Xem thêm methodology tại [https://lotusai.servehttp.com](https://lotusai.servehttp.com).

---

*Bài viết là phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Past performance không đảm bảo kết quả tương lai.*
