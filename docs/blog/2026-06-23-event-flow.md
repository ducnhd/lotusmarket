---
title: "Dòng tiền VN30 nghiêng bán hôm nay 23/06 — regime nào đang dẫn dắt?"
date: 2026-06-23
topic: flow
---

## Áp lực bán lan rộng, VN30+HNX30 chỉ còn 20% buy-pressure

Hôm nay dòng tiền nội tổng của rổ VN30 và HNX30 ghi nhận buy-pressure chỉ đạt **20.0%** — ngưỡng mà thị trường thường gọi thẳng là "bán mạnh". Câu hỏi mà data 9 năm có thể trả lời: khi thị trường rơi vào trạng thái này, lịch sử đã đi tiếp ra sao?

---

### ① Chuyện gì đang xảy ra

Buy-pressure 20% có nghĩa là trong tổng khối lượng khớp lệnh của hai rổ lớn nhất sàn, chỉ 1 phần 5 đến từ phía mua chủ động. Bốn phần còn lại là người bán đang nắm thế chủ động. Đây không phải thị trường đang "điều chỉnh nhẹ" hay "rung lắc kỹ thuật" — con số 20% nằm sâu trong vùng áp lực bán cực đoan, nơi mà tâm lý đám đông đang nghiêng hẳn về phía thoát hàng.

Điều đáng chú ý là tín hiệu này đến từ **tổng hợp cả VN30 lẫn HNX30**, tức là áp lực không chỉ tập trung ở một sàn hay một nhóm ngành. Khi cả hai rổ đồng loạt ghi nhận buy-pressure thấp như vậy, thường phản ánh tâm lý thận trọng lan rộng trên diện rộng, không phải câu chuyện riêng của vài mã bluechip.

---

### ② Vì sao trạng thái này xuất hiện

Dòng tiền nội tổng phản ánh hành vi thực tế của nhà đầu tư trong nước — từ tổ chức đến cá nhân — tại thời điểm khớp lệnh. Khi buy-pressure co về vùng 20%, một trong hai kịch bản thường xảy ra: hoặc là lực cầu rút lui chờ đợi tín hiệu rõ hơn, hoặc là lực bán chủ động gia tăng do các yếu tố kích hoạt từ bên ngoài — macro toàn cầu, tỷ giá, hay dòng vốn ngoại đảo chiều.

Dữ liệu regime phân loại trạng thái hôm nay thuộc **CRISIS** — không phải STABLE. Đây là sự phân biệt quan trọng: cùng một mức buy-pressure thấp nhưng xảy ra trong môi trường CRISIS sẽ có hàm ý khác hẳn so với giai đoạn thị trường bình thường.

---

### ③ Cohort lịch sử 9 năm nói gì

Đây là phần mà data trả lời thẳng, không cần phỏng đoán.

Trong 9 năm dữ liệu lịch sử, khi thị trường ở **regime CRISIS** với cấu hình tương tự:

- **Forward return 60 ngày: +8.53%**
- **Win rate: 67%**

So sánh với baseline regime **STABLE**:
- Forward return 60 ngày: +3.20%
- Win rate: 48% (thấp hơn cả baseline +4.75%)

Nói thẳng từ data: cohort CRISIS lịch sử không những có return kỳ vọng cao hơn gần **2.7 lần** so với cohort STABLE, mà win rate cũng vượt trội rõ ràng — 67% so với 48%.

Điều này nghe có vẻ phản trực giác. Thị trường đang bị bán mạnh, áp lực cực đoan, tâm lý xấu — nhưng chính xác ở những thời điểm như thế này, cohort lịch sử lại ghi nhận kết quả thuận lợi hơn trong 60 ngày tiếp theo. Cơ chế phổ biến được quan sát: áp lực bán cực đoan thường đi kèm với định giá bị nén, và khi lực bán cạn kiệt, thị trường phục hồi với biên độ rộng hơn bình thường. 🔍

Tuy nhiên, win rate 67% cũng đồng nghĩa 33% các lần trong lịch sử, thị trường **không** hồi phục trong khung 60 ngày. Một phần ba không phải con số nhỏ. Đây là lý do data không nên được đọc như một bảo đảm.

---

### ④ Cần theo dõi gì tiếp theo

Ba biến số cần quan sát trong các phiên tới:

**Buy-pressure có phục hồi về vùng 35–40% không?** Đây là ngưỡng mà lực cầu bắt đầu cân bằng lại với lực bán. Nếu buy-pressure tiếp tục duy trì dưới 25% trong nhiều phiên liên tiếp, regime CRISIS có thể deepening — cohort lịch sử cho loại diễn biến đó sẽ cần đánh giá lại.

**Regime có chuyển từ CRISIS sang STABLE không?** Nếu có, lưu ý rằng cohort STABLE lại có kết quả kém hơn baseline (+3.20% so với +4.75%). Sự chuyển dịch regime không tự động là tín hiệu tích cực.

**Độ đồng thuận giữa VN30 và HNX30.** Khi hai rổ phân kỳ — một sàn bắt đầu tích lũy trong khi sàn kia vẫn bị bán — đó thường là dấu hiệu sớm của sự phân hóa dòng tiền có chủ đích. 📊

---

### Verify reproducible

Toàn bộ tính toán cohort và buy-pressure có thể tái tạo qua [lotusmarket](https://github.com/ducnhd/lotusmarket):

```bash
pip install lotusmarket
```

```python
from lotusmarket import FlowAnalyzer
fa = FlowAnalyzer(universe=["VN30", "HNX30"])
fa.cohort_backtest(regime="CRISIS", forward_days=60)
```

---

*Nội dung trên chỉ phản ánh phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Mọi quyết định giao dịch thuộc trách nhiệm của nhà đầu tư.*
