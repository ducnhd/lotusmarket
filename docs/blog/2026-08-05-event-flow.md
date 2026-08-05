---
title: "Dòng tiền VN30 nghiêng bán hôm nay 05/08 — regime nào đang dẫn dắt?"
date: 2026-08-05
topic: flow
---

## Thị trường rơi vào vùng bán mạnh — lịch sử nói gì?

Phiên hôm nay, áp lực bán trên rổ VN30+HNX30 đẩy buy-pressure xuống chỉ còn **32.0%** — ngưỡng mà hệ thống phân loại thẳng là *bán mạnh*. Câu hỏi thực sự không phải là "thị trường có đang xấu không" — data đã trả lời rồi — mà là: **những lần tương tự trong 9 năm qua, giá đi đâu tiếp theo?**

---

### ① Chuyện gì đang xảy ra

Buy-pressure 32.0% có nghĩa là chưa đến một phần ba dòng tiền nội trên hai rổ chủ lực đang đứng phía mua. Phần còn lại — hơn 68% — đang chạy về phía bán hoặc đứng ngoài. Đây không phải trạng thái "giằng co nhẹ"; đây là trạng thái mà bên bán đang kiểm soát hoàn toàn nhịp khớp lệnh.

Khi dòng tiền nội tổng nghiêng lệch đến mức này, hệ thống Lotus AI phân loại phiên vào **regime CRISIS** — không phải vì dùng từ cho có kịch tính, mà vì đây là nhóm quan sát trong cơ sở dữ liệu có hành vi giá phía trước khác hẳn so với regime bình thường.

---

### ② Vì sao dòng tiền co lại

Không có một nguyên nhân đơn lẻ nào giải thích đủ, nhưng cấu trúc của chính số liệu đã nói lên phần lớn câu chuyện. Buy-pressure 32% trên *tổng* VN30+HNX30 — tức là cả hai sàn HOSE lẫn HNX cùng lúc — cho thấy đây không phải hiện tượng cục bộ một nhóm ngành. Khi cả hai rổ lớn đồng loạt rơi vào trạng thái này, lực bán mang tính lan rộng, không phân biệt vốn hóa lớn hay trung bình.

Trong những giai đoạn như thế này, hành vi điển hình là nhà đầu tư ngắn hạn thoát margin, quỹ rebalance danh mục, còn dòng tiền mới chưa dám vào vì chưa thấy đáy xác nhận. Kết quả: thanh khoản mỏng, lệnh bán nhỏ cũng tạo ra biên độ giảm lớn hơn bình thường.

---

### ③ 9 năm dữ liệu nói gì

Đây là phần quan trọng nhất — và cũng là lý do mình viết bài này thay vì chỉ nhìn bảng giá.

Từ cohort phân tích theo regime trải dài 9 năm:

| Regime | Fwd 60 ngày | Win rate |
|--------|-------------|----------|
| **CRISIS** | **+8.53%** | **67%** |
| STABLE | +3.20% | 48% |
| Baseline toàn bộ | +4.75% | — |

Kết quả này phản trực giác hoàn toàn. Regime **CRISIS** — tức là những lần buy-pressure rơi xuống vùng bán mạnh như hôm nay — lại cho forward return 60 ngày trung bình cao nhất: **+8.53%**, với tỷ lệ thắng **67%**. Trong khi đó, regime STABLE với thị trường "bình thường, ít biến động" lại chỉ mang về +3.20% với win rate dưới baseline.

Cơ chế đằng sau không có gì bí ẩn: khi bán mạnh đã xảy ra và buy-pressure chạm đáy, phần lớn lực bán đã *giải phóng xong*. Bên còn lại trên thị trường là người giữ kiên nhẫn hơn — và khi sentiment đảo chiều dù chỉ một phần, không cần nhiều lực mua để đẩy giá.

Điều cần nhấn mạnh: 67% win rate không có nghĩa là 33% còn lại không đau. Những lần thị trường rơi vào CRISIS mà *không* phục hồi trong 60 ngày thường đi kèm với sự kiện vĩ mô kéo dài — không phải câu chuyện kỹ thuật đơn thuần.

---

### ④ Cần theo dõi gì trong thời gian tới

Ba tín hiệu đáng xem sát:

**Buy-pressure có phục hồi về trên 40% không?** Đây là ngưỡng mà dòng tiền nội bắt đầu cân bằng trở lại. Nếu trong 3–5 phiên tới buy-pressure vẫn dưới 35%, regime CRISIS có thể kéo dài — lịch sử cho thấy phục hồi chậm hơn trong trường hợp đó.

**Sự phân kỳ giữa VN30 và HNX30** — nếu một trong hai rổ bắt đầu tách ra (buy-pressure tăng độc lập), đó là dấu hiệu dòng tiền đang tìm điểm vào chọn lọc thay vì thoát đồng loạt.

**Khối lượng tổng phiên** — regime CRISIS kết thúc thường đi kèm với một phiên có khối lượng đột biến và buy-pressure đảo chiều mạnh trong cùng ngày. Thiếu một trong hai, tín hiệu chưa xác nhận. ⚡

---

### Verify — tự kiểm tra số liệu

```bash
pip install lotusmarket
```

```python
from lotusmarket import regime, cohort

df = regime.load("VN30+HNX30", lookback="9y")
result = cohort.forward_return(df, regime_filter="CRISIS", fwd_days=60)
print(result[["mean_return", "win_rate"]])
```

Hoặc dùng CLI:

```bash
lmcli cohort --ticker VN30+HNX30 --regime CRISIS --fwd 60
```

Tài liệu và mã nguồn tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket). 📊

---

*Bài viết chỉ mang tính phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Mọi quyết định giao dịch thuộc trách nhiệm của nhà đầu tư.*
