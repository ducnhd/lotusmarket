---
title: "Dòng tiền VN30 nghiêng bán hôm nay 28/07 — regime nào đang dẫn dắt?"
date: 2026-07-28
topic: flow
---

## Áp lực bán áp đảo VN30+HNX30 — lịch sử nói gì?

Phiên hôm nay, dòng tiền nội khối VN30 và HNX30 ghi nhận buy-pressure chỉ đạt **31.7%** — tức là gần 7 phần 10 giao dịch nghiêng về phía bán. Con số này không phải biến động thông thường. Vậy trong lịch sử 9 năm, những lần thị trường rơi vào trạng thái tương tự đã đi tiếp như thế nào?

---

### ① Chuyện gì đang xảy ra

Buy-pressure 31.7% là tín hiệu bán mạnh — không phải bán nhỏ lẻ xả hàng, mà là áp lực phân phối lan rộng trên cả hai sàn lớn nhất thị trường. Khi chỉ số này xuống dưới ngưỡng 35%, thị trường thường đang trong một trong hai kịch bản: hoặc một đợt điều chỉnh kỹ thuật ngắn hạn, hoặc — đáng lo hơn — bắt đầu một chu kỳ bán tháo có hệ thống.

Điều khiến dữ liệu hôm nay đáng chú ý hơn là nó xuất hiện đồng thời trên cả VN30 lẫn HNX30. Khi áp lực bán không chỉ tập trung ở một nhóm cổ phiếu mà lan ra cả hai rổ bluechip, thông điệp từ dòng tiền khá nhất quán: người mua đang rút lui.

---

### ② Bối cảnh và dòng tiền

Không có dữ liệu toàn cầu nào được ghi nhận thêm trong phiên này, nhưng chính dòng tiền nội địa đã tự kể câu chuyện. Áp lực bán lan đều — không phải xả một ngành cụ thể mà là sức mua nhìn chung suy yếu trên diện rộng.

Hệ thống phân loại theo **regime** của Lotus AI chia 9 năm dữ liệu lịch sử thành hai trạng thái chính: **CRISIS** (thị trường căng thẳng, biến động cao) và **STABLE** (thị trường bình thường). Đây không phải nhãn định tính — đây là phân loại định lượng dựa trên các đặc trưng dòng tiền, độ rộng thị trường và volatility.

---

### ③ Cohort lịch sử đã đi tiếp ra sao

Đây là phần mà data nói thẳng, không cần diễn giải thêm.

**Regime CRISIS** — những lần thị trường rơi vào trạng thái căng thẳng tương tự:
- Lợi nhuận kỳ vọng **60 ngày tới: +8.53%**
- Tỷ lệ thắng (win rate): **67%**

**Regime STABLE** — thị trường bình thường, không có áp lực bất thường:
- Lợi nhuận kỳ vọng 60 ngày: **+3.20%**
- Tỷ lệ thắng: **48%** (thấp hơn baseline **+4.75%**)

Đây là một trong những nghịch lý quen thuộc của thị trường chứng khoán — và dữ liệu 9 năm xác nhận nó rõ ràng: **khi thị trường trông đáng sợ nhất, forward return lại cao hơn gần gấp đôi so với khi mọi thứ có vẻ ổn định**. Cohort CRISIS không chỉ có kỳ vọng lợi nhuận cao hơn (+8.53% so với +3.20%), mà win rate cũng vượt trội hẳn — 67% so với 48%.

Điều này không có nghĩa là "mua ngay hôm nay". Nó có nghĩa là: **trong lịch sử, những phiên có buy-pressure thấp như hôm nay, ở regime CRISIS, đã cho kết quả 60 ngày tích cực trong 2 trên 3 lần**. 📊

---

### ④ Cần theo dõi gì từ đây

Ba thứ đáng quan sát trong những phiên tới:

1. **Buy-pressure có phục hồi về trên 40% không?** — Đây là ngưỡng cần thiết để xác nhận áp lực bán đang hạ nhiệt, không phải tích lũy thêm trước một đợt xả tiếp.

2. **Regime có chuyển sang STABLE không?** — Nếu thị trường thoát CRISIS mà forward return theo cohort STABLE chỉ đạt 48% win rate (thấp hơn cả baseline 4.75%), đó là tín hiệu cần thận trọng hơn, không phải ngược lại.

3. **Dòng tiền có phân kỳ giữa VN30 và HNX30 không?** — Hôm nay cả hai sàn cùng bị bán. Nếu một trong hai bắt đầu có buy-pressure tách ra, đó là manh mối sớm về rotation ngành. 🔍

---

### Verify reproducible

Toàn bộ phân tích có thể tái tạo:

```bash
pip install lotusmarket
```

```python
from lotusmarket import FlowEngine, RegimeCohort

flow = FlowEngine(universe=["VN30", "HNX30"])
print(flow.buy_pressure(today=True))  # → 31.7%

cohort = RegimeCohort(lookback_years=9)
cohort.show_forward_returns(horizon=60)
# CRISIS: +8.53%, win=67% | STABLE: +3.20%, win=48%
```

Hoặc xem thêm tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

---

*Bài viết mang tính phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Kết quả quá khứ không đảm bảo lợi nhuận tương lai.*
