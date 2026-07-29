---
title: "VRE tăng 6.8% hôm nay 29/07 — cohort lịch sử nói gì?"
date: 2026-07-29
topic: mover
---

## VRE bật +6.81% từ vùng RSI 26: tín hiệu đảo chiều hay bẫy hồi kỹ thuật?

Phiên hôm nay, VRE (Vincom Retail) bứt phá mạnh với mức tăng **+6.81%**, đóng cửa tại **22.750 đồng**. Điều khiến tín hiệu này đáng chú ý hơn là bối cảnh nền — RSI(14) chỉ ở **26.7**, tức vùng quá bán sâu. Câu hỏi trung tâm: liệu đây là cú nảy kỹ thuật đơn lẻ, hay điểm khởi đầu của một pha phục hồi có chiều sâu hơn?

---

### ① Chuyện gì xảy ra hôm nay

Con số **+6.81%** trong một phiên không phải chuyện bình thường với một blue-chip bất động sản bán lẻ như VRE. Nhưng ẩn dưới mức tăng đó là hai dữ liệu cần đọc kỹ hơn.

Thứ nhất, **volume chỉ đạt 831.480 cổ phiếu — bằng 0.2× MA20 volume**. Tức là thanh khoản hôm nay chỉ bằng một phần năm mức trung bình 20 phiên. Một cây nến xanh đậm với khối lượng cạn như vậy thường phản ánh lực cầu chủ động yếu — giá leo lên do lực bán co lại, chứ chưa hẳn do dòng tiền mới ào vào.

Thứ hai, **RSI(14) = 26.7**, tức nằm dưới ngưỡng quá bán 30. Thông thường RSI ở vùng này báo hiệu cổ phiếu đã bị bán quá mức trong ngắn hạn, áp lực bán đã "vắt kiệt" một phần. Kết hợp thêm **MA trend: mixed** — xu hướng các đường trung bình chưa thống nhất chiều — cho thấy bức tranh kỹ thuật vẫn còn mâu thuẫn nội tại.

---

### ② Vì sao tín hiệu này xuất hiện

VRE vận hành chuỗi trung tâm thương mại Vincom, mô hình kinh doanh phụ thuộc sức tiêu dùng nội địa và tỉ lệ lấp đầy mặt bằng. Trong bối cảnh lãi suất đã hạ nhiệt và kỳ vọng phục hồi tiêu dùng nửa cuối năm, nhóm bất động sản bán lẻ thường được định giá lại khi thị trường tìm kiếm các cổ phiếu defensive yield.

Trước phiên hôm nay, chuỗi giảm điểm tích lũy đủ để kéo RSI xuống vùng 26 — một vùng mà lịch sử thị trường chứng khoán Việt Nam ít khi duy trì lâu mà không có ít nhất một nhịp hồi. Volume thấp bất thường trong phiên tăng có thể giải thích theo hai hướng: hoặc nhà đầu tư tổ chức chưa vào nhưng cũng không bán thêm (lực cung yếu đẩy giá lên), hoặc đây là giai đoạn tích lũy âm thầm trước khi volume bùng nổ xác nhận.

---

### ③ Cohort lịch sử cho thấy gì

Đây là phần dữ liệu cứng nhất.

Lotus AI đã chạy backtest trên toàn bộ lịch sử giao dịch với **cohort RSI < 30 oversold (N = 1.755 trường hợp)**. Kết quả:

| Chỉ số | Giá trị |
|---|---|
| Forward return 60 ngày | **+6.32%** |
| Win rate | **60%** |
| Edge | **+1.57%** |

Nói cách khác, trong 1.755 lần cổ phiếu rơi vào vùng RSI dưới 30 tương tự, **60% số trường hợp sinh lời dương sau 60 ngày**, với mức sinh lời trung bình **+6.32%**. Edge dương **+1.57%** so với baseline cho thấy cohort này không phải ngẫu nhiên — tín hiệu RSI quá bán có giá trị thống kê xác thực.

Tuy nhiên, 40% còn lại của cohort tiếp tục giảm hoặc không hồi đủ. Volume thấp hôm nay là yếu tố phân kỳ cần theo dõi — các trường hợp trong cohort thắng thường có volume xác nhận trong 3–5 phiên tiếp theo.

---

### ④ Cần theo dõi gì trong các phiên tới

Bốn điểm cụ thể cần quan sát:

**Volume:** Nếu VRE duy trì đà tăng mà volume vẫn dưới 0.5× MA20, tín hiệu kém tin cậy hơn. Một phiên volume đột biến vượt MA20 sẽ là xác nhận dòng tiền thực sự vào.

**RSI:** Cần RSI thoát khỏi vùng dưới 30 và duy trì trên 35–40 để xu hướng quá bán được xem là đã giải tỏa.

**MA trend:** Hiện tại "mixed" — cần quan sát MA ngắn hạn (5, 10 ngày) có cắt lên MA dài hạn (20, 50 ngày) không. Đây là điều kiện để trend kỹ thuật chuyển từ mâu thuẫn sang rõ hướng.

**Ngưỡng giá:** 22.750 đồng là mức đóng cửa hôm nay. Nếu các phiên tiếp theo không giữ được trên mức này, cú bật có thể chỉ là dead cat bounce trong downtrend chưa kết thúc.

---

### Verify reproducible 🔍

Toàn bộ cohort analysis trên có thể tái tạo qua [lotusmarket](https://github.com/ducnhd/lotusmarket):

```bash
pip install lotusmarket
```

```python
from lotusmarket import backtest
result = backtest.cohort("VRE", filter="RSI_14<30", forward_days=60)
print(result.summary())
```

Hoặc dùng CLI:

```bash
lmcli cohort --ticker VRE --signal rsi_oversold --fwd 60
```

---

*Bài viết mang tính phân tích dữ liệu, không phải lời khuyên đầu tư. Mọi quyết định giao dịch thuộc trách nhiệm của nhà đầu tư.*
