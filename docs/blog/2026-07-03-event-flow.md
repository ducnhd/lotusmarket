---
title: "Dòng tiền VN30 nghiêng bán hôm nay 03/07 — regime nào đang dẫn dắt?"
date: 2026-07-03
topic: flow
---

## Áp lực bán áp đảo VN30+HNX30 — lịch sử nói gì về những ngày như thế này?

Hôm nay dòng tiền nội trên rổ tổng hợp VN30+HNX30 ghi nhận buy-pressure chỉ đạt **31.3%** — tức là gần 7 trong 10 đồng tiền chảy qua hai rổ bluechip lớn nhất sàn đang đứng về phía bán. Con số này không phải "hơi yếu" mà nằm thẳng vào vùng **bán mạnh** theo phân loại regime của hệ thống. Câu hỏi thực sự là: trong 9 năm dữ liệu, những phiên tương tự đã kết thúc như thế nào sau 60 ngày giao dịch?

---

### ① Chuyện gì đang xảy ra — số liệu thô

Buy-pressure 31.3% là tín hiệu rõ ràng hơn nhiều so với mức nhiễu thông thường. Để đặt vào context: ngưỡng 50% là trung tính, dưới 40% đã là vùng áp lực bán có trọng lượng, còn 31.3% đưa thị trường vào trạng thái mà hệ thống phân loại là **CRISIS regime** — tức là dòng tiền đang co cụm, người mua rút lui đáng kể so với người bán trên cả hai sàn lớn.

Điều đáng chú ý là tín hiệu này đến từ **tổng hợp VN30 lẫn HNX30**, không phải một rổ đơn lẻ. Khi cả hai đồng thuận theo chiều bán mạnh, khả năng đây là noise kỹ thuật một phiên sẽ thấp hơn nhiều so với khi chỉ một rổ lệch.

---

### ② Vì sao regime CRISIS lại xuất hiện

Hệ thống không gán nhãn CRISIS dựa vào tin tức — nó đọc trực tiếp từ **hành vi dòng tiền thực tế** trên bảng khớp lệnh. Khi buy-pressure rơi xuống vùng này, thường đi kèm một trong các kịch bản: margin call lan rộng trong nhóm bluechip, tổ chức nước ngoài xả ròng liên tục, hoặc tâm lý phòng thủ xuất hiện trước một sự kiện vĩ mô lớn chưa được định giá vào thị trường.

Bất kể nguyên nhân cụ thể hôm nay là gì, **pattern hành vi** đã đủ để phân loại — và đó là thứ cohort lịch sử sẽ phản ánh.

---

### ③ Lịch sử 9 năm: cohort CRISIS đã đi đâu

Đây là phần data nói thẳng nhất. Trong 9 năm dữ liệu được phân tích, các phiên rơi vào **CRISIS regime** có forward return trung bình **+8.53% sau 60 ngày giao dịch**, với **tỷ lệ thắng 67%** — tức là 2 trong 3 lần thị trường hồi phục dương trong vòng 3 tháng tiếp theo.

So sánh thẳng với **STABLE regime** — khi dòng tiền bình thường, không có bán mạnh — forward 60d chỉ đạt +3.20%, win rate 48%, thậm chí **thấp hơn cả baseline +4.75%** của toàn bộ mẫu. Nói cách khác, mua vào những ngày "bình yên" không đem lại lợi thế thống kê so với mua ngẫu nhiên.

📊 Tóm gọn bằng số:

| Regime | Fwd 60d | Win Rate |
|--------|---------|----------|
| CRISIS (hôm nay) | **+8.53%** | **67%** |
| STABLE | +3.20% | 48% |
| Baseline | +4.75% | — |

Lịch sử cho thấy: **áp lực bán cực đoan có xu hướng xảy ra gần đáy hơn là gần đỉnh**. Điều này không có nghĩa đáy đã chạm — nhưng phân phối kết quả về phía dương rõ ràng hơn nhiều so với các phiên thị trường yên tĩnh.

---

### ④ Cần theo dõi gì tiếp theo

Ba thứ đáng quan sát trong những phiên tới:

**Thứ nhất**, buy-pressure có phục hồi về trên 40% không — đây là ngưỡng xác nhận áp lực bán đang được hấp thụ. Nếu buy-pressure tiếp tục duy trì dưới 35% nhiều phiên liên tiếp, cohort CRISIS có thể kéo dài thay vì đảo chiều nhanh.

**Thứ hai**, hành vi phân kỳ giữa VN30 và HNX30 — nếu một rổ bắt đầu tách ra (mua phục hồi ở HNX30 trong khi VN30 còn yếu, hoặc ngược lại), đó thường là dấu hiệu sớm dòng tiền đang định hướng lại.

**Thứ ba**, khối lượng khớp lệnh tuyệt đối — CRISIS với khối lượng thấp (bán vì không có người mua) khác với CRISIS khối lượng cao (bán có người mua hấp thụ tích cực). Bộ dữ liệu không phân biệt hai trường hợp này hôm nay, nhưng đây là biến số cần theo dõi thêm.

---

### Verify reproducible

Toàn bộ phân tích trên có thể tái tạo qua [lotusmarket](https://github.com/ducnhd/lotusmarket):

```python
pip install lotusmarket

from lotusmarket import flow, regime
df = flow.load("VN30+HNX30")
cohort = regime.backtest(df, metric="buy_pressure", forward_days=60)
print(cohort.groupby("regime")[["fwd_return","win_rate"]].mean())
```

---

*Bài viết sử dụng dữ liệu lịch sử và phân tích định lượng. Đây không phải lời khuyên đầu tư — mọi quyết định giao dịch thuộc về trách nhiệm cá nhân của nhà đầu tư.* 🔍
