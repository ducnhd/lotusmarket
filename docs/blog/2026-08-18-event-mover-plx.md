---
title: "PLX tăng 6.9% hôm nay 18/08 — cohort lịch sử nói gì?"
date: 2026-08-18
topic: mover
---

## PLX tăng gần 7% trong một phiên — volume nói gì khi giá đã chạy xa?

Petrolimex bất ngờ bứt phá +6.90% trong phiên hôm nay, đóng cửa tại 37.950 VND — mức tăng hiếm gặp với một cổ phiếu vốn hóa lớn trong ngành phân phối xăng dầu. Điều khiến người quan sát phải dừng lại: volume chỉ đạt 535.090 cổ phiếu, tương đương 0.2× MA20 — tức là thanh khoản chưa bằng một phần năm trung bình 20 phiên. Câu hỏi đặt ra là: một cú tăng mạnh trên nền volume "im lặng" như vậy thường dẫn đến đâu?

---

### ① Chuyện gì vừa xảy ra

Số liệu trực tiếp từ phiên hôm nay không để lại nhiều tranh cãi. PLX tăng +6.90%, đây là mức dao động một ngày thuộc nhóm đáng chú ý với bất kỳ cổ phiếu bluechip nào. RSI(14) đang ở 65.2 — chưa chạm vùng overbought kỹ thuật (thường tính từ 70), nhưng cũng đã rời khá xa vùng trung tính 50.

Tuy nhiên, điểm dị thường nằm ở volume. 535.090 đơn vị giao dịch, bằng 0.2 lần MA20 volume — nghĩa là phần lớn thị trường **không tham gia** vào cú tăng này. Giá di chuyển mạnh nhưng lực mua thực sự đo được qua khớp lệnh lại rất mỏng. Đây là tổ hợp cần đọc kỹ hơn là phấn khích.

MA trend được ghi nhận là **mixed** — các đường trung bình ngắn và dài hạn chưa xếp hàng cùng chiều, tức là về mặt xu hướng, PLX chưa có tín hiệu breakout sạch.

---

### ② Vì sao có thể có cú tăng này

Petrolimex là doanh nghiệp đầu mối phân phối xăng dầu lớn nhất Việt Nam, do đó giá cổ phiếu thường phản ứng với hai nhân tố chính: **biến động giá dầu thế giới** và **chính sách điều hành giá xăng trong nước**.

Trong bối cảnh giá dầu thô quốc tế có những phiên hồi phục gần đây sau giai đoạn rung lắc, thị trường có thể đang định giá lại biên lợi nhuận đầu mối của PLX. Ngoài ra, chu kỳ điều chỉnh giá xăng trong nước — nếu đi kèm tín hiệu thuận lợi từ cơ quan quản lý — thường tạo ra phản ứng nhanh trên cổ phiếu này.

Tuy nhiên, data hôm nay **không cung cấp thêm thông tin về catalyst cụ thể**. Điều mình có thể khẳng định từ data: giá tăng mạnh nhưng lực cầu đo được (volume) ở mức rất thấp. Trong nhiều trường hợp lịch sử, dạng tăng này xuất phát từ áp lực bán giảm đột ngột hơn là dòng tiền mới ồ ạt vào.

---

### ③ Cohort lịch sử tương tự đã đi như thế nào

Đây là phần quan trọng nhất — và data trả lời thẳng.

Khi chạy cohort các phiên PLX có tổ hợp tương tự (tăng mạnh >5%, RSI trong vùng 60-70, volume dưới 0.3× MA20, MA trend mixed), **hệ thống không tìm được cohort khớp pattern đủ mạnh**. Baseline cohort edge được ghi nhận là **~0%** — tức là về mặt thống kê, lịch sử không ủng hộ cũng không phản bác xu hướng tiếp diễn.

Nói thẳng: đây là **vùng neutral theo cohort**. Không có edge lịch sử rõ ràng để kỳ vọng PLX tiếp tục tăng, cũng không có edge rõ để kỳ vọng đảo chiều ngay. Thị trường đang ở trạng thái mà data chưa có câu trả lời định lượng.

Điều này quan trọng hơn mọi nhận định cảm tính: khi edge ~0%, rủi ro đến từ việc **hành động dựa trên narrative** thay vì xác suất thực từ lịch sử.

---

### ④ Cần theo dõi gì từ đây

Ba biến số đáng quan sát trong các phiên tới:

**Volume phục hồi hay không.** Nếu các phiên sau volume quay về trên 0.5× MA20 kèm giá giữ được vùng 37.000–38.000, đó là tín hiệu dòng tiền thực bắt đầu xác nhận cú tăng hôm nay. Ngược lại, volume tiếp tục mỏng mà giá không giữ được, cú tăng này có thể chỉ là kỹ thuật ngắn hạn.

**RSI vượt 70 hay không.** RSI 65.2 còn cách vùng overbought khoảng 5 điểm. Nếu giá tiếp tục tăng mà RSI leo qua 70 trên volume thấp, đó là tổ hợp lịch sử thường kèm biến động ngược chiều.

**MA trend chuyển sang aligned.** Khi các đường MA ngắn-dài xếp cùng chiều tăng, tín hiệu xu hướng mới có trọng lượng hơn. Hiện tại mixed MA là lý do để thận trọng hơn là hành động vội.

---

### Verify reproducible 🔍

Toàn bộ cohort analysis trên có thể tự kiểm tra:

```bash
pip install lotusmarket
```

```python
from lotusmarket import cohort
result = cohort.run("PLX", rsi_range=(60,70), volume_ratio_max=0.3, ma_trend="mixed")
print(result.edge_summary())
```

Hoặc dùng CLI:

```bash
lmcli cohort PLX --rsi 60-70 --vol-ratio 0.3 --ma mixed
```

Tài liệu và source tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

---

*Bài viết mang tính phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Mọi quyết định giao dịch thuộc trách nhiệm của nhà đầu tư.*
