---
title: "SSI tăng 7.0% hôm nay 21/08 — cohort lịch sử nói gì?"
date: 2026-08-21
topic: mover
---

## SSI bật +6.96% — khối lượng mỏng nói lên điều gì?

Phiên hôm nay, SSI đóng cửa tại 20.750 VND, tăng gần 7% chỉ trong một phiên. Nhưng điều đáng chú ý hơn mức tăng đó là khối lượng giao dịch — và chính con số này mới là câu hỏi cần trả lời.

---

### ① Chuyện gì vừa xảy ra

SSI khép phiên với mức tăng +6.96%, một con số đủ để lọt vào nhóm cổ phiếu nổi bật nhất thị trường trong ngày. RSI(14) leo lên 65.2 — chưa vào vùng overbought (>70), nhưng đã tiến khá gần.

Thế nhưng volume lại kể một câu chuyện khác: **3.870.680 cổ phiếu khớp lệnh**, chỉ bằng **0.2× MA20 volume**. Nghĩa là khối lượng hôm nay chưa bằng một phần năm mức trung bình 20 phiên gần nhất. Một mức tăng gần 7% mà dòng tiền tham gia lại thưa đến vậy — đây là điểm bất thường đáng để dừng lại.

Về mặt kỹ thuật, MA trend hiện ở trạng thái **mixed** — tức là các đường trung bình ngắn hạn và dài hạn chưa đồng thuận rõ ràng về xu hướng. Không phải uptrend xác nhận, cũng không phải downtrend — SSI đang đứng ở vùng tranh chấp.

---

### ② Vì sao có thể có nhịp bật này

Nhóm chứng khoán vốn nhạy cảm với thanh khoản thị trường và kỳ vọng lãi suất. Khi thị trường chung có những phiên hồi phục tâm lý hoặc có thông tin hỗ trợ từ phía cơ quan quản lý — dù là kỳ vọng nới room, câu chuyện margin, hay đơn giản là dòng tiền đầu cơ ngắn hạn tìm đến beta cao — các mã môi giới lớn như SSI thường phản ứng nhanh và mạnh hơn mặt bằng chung.

Tuy nhiên, khi volume chỉ đạt 0.2× MA20, điều đó gợi ý rằng nhịp tăng này **không có sự tham gia rộng của dòng tiền lớn**. Có thể là lực cầu bắt đáy tập trung ở một vài phiên khớp lệnh cụ thể trong ngày, hoặc áp lực bán mỏng khiến giá bật lên mà không cần lực mua đặc biệt lớn — kiểu tăng trên nền trống.

RSI ở 65.2 cho thấy đà tăng ngắn hạn đang có lực, nhưng chưa đến mức cảnh báo quá mua ngay lập tức. Vùng 65–70 thường là nơi thị trường phân kỳ: một phần muốn chốt lời, một phần vẫn còn kỳ vọng tiếp diễn.

---

### ③ Cohort lịch sử đã đi tiếp ra sao

Đây là phần thú vị — và cũng là phần thẳng thắn nhất.

Khi chạy phân tích cohort trên dữ liệu lịch sử SSI với các điều kiện tương tự (tăng mạnh ~7% trong một phiên, volume thấp hơn đáng kể so với MA20, RSI tiệm cận 65, MA trend mixed), **không có cohort nào khớp pattern đủ mạnh để rút ra edge rõ ràng**.

Baseline cohort edge được ghi nhận ở mức **~0%** — nghĩa là lịch sử không cho thấy thiên lệch xác suất đáng kể theo chiều nào sau tín hiệu này. Không phải tín hiệu bullish mạnh, cũng không phải tín hiệu đảo chiều có độ tin cậy cao. Đây là **vùng neutral theo định nghĩa thống kê**.

Điều này không có nghĩa là giá sẽ đứng yên — mà có nghĩa là data lịch sử chưa đủ để xác nhận hướng tiếp theo với độ tự tin cao. Nói cách khác: SSI đang ở điểm mà nhiều kịch bản đều có thể xảy ra với xác suất tương đương nhau.

---

### ④ Cần theo dõi gì từ đây

Ba điểm cần quan sát trong các phiên tới:

**Volume xác nhận.** Nếu SSI tiếp tục tăng nhưng volume vẫn dưới 0.5× MA20, nhịp tăng này thiếu nền tảng dòng tiền. Ngược lại, nếu volume bùng lên vượt MA20 kèm theo giá giữ được trên 20.750, đó mới là tín hiệu đáng chú ý hơn.

**RSI vượt 70 hay không.** Vùng 65–70 là trung gian. RSI vượt 70 mà volume vẫn thấp thường là cảnh báo phân kỳ. RSI thoái về dưới 60 trong các phiên tới mà không có volume bán mạnh lại là dấu hiệu tích lũy.

**MA trend có đồng thuận không.** Hiện tại mixed — cần theo dõi xem các đường MA ngắn hạn có cắt lên trên MA dài hạn trong 3–5 phiên tới không. Đó mới là xác nhận xu hướng.

---

### Verify reproducible

Toàn bộ phân tích cohort trên được thực hiện qua [lotusmarket](https://github.com/ducnhd/lotusmarket). Tự chạy lại:

```bash
pip install lotusmarket
```

```python
from lotusmarket import CohortAnalyzer
ca = CohortAnalyzer("SSI")
ca.run(rsi_range=(63, 67), volume_ratio_max=0.25, price_change_min=0.06)
ca.summary()
```

Hoặc dùng CLI:

```bash
lmcli cohort SSI --rsi 65.2 --vol-ratio 0.2 --change 6.96%
```

---

*Bài viết mang tính phân tích dữ liệu, không phải lời khuyên đầu tư. Mọi quyết định giao dịch là trách nhiệm của nhà đầu tư.* 📊
