---
title: "BCM giảm 5.5% hôm nay 13/07 — cohort lịch sử nói gì?"
date: 2026-07-13
topic: mover
---

## BCM rơi 5.5% với RSI chạm 16 — vùng oversold hiếm gặp trong lịch sử

Phiên hôm nay, BCM đóng cửa tại 46.000 VND, giảm 5,54% so với phiên trước. Điều đáng chú ý hơn mức giảm là chỉ số RSI(14) vừa chạm 16,2 — một mức hiếm gặp mà dữ liệu lịch sử có câu trả lời khá rõ ràng.

---

### ① Chuyện gì xảy ra hôm nay

Mức giảm 5,54% trong một phiên không phải là con số nhỏ với một cổ phiếu bất động sản khu công nghiệp như BCM. Giá đóng cửa tại 46.000 VND đưa cổ phiếu này về vùng giá thấp đáng để nhìn lại.

Nhưng phần gây chú ý nhất không phải mức giảm giá — mà là khối lượng giao dịch. Volume hôm nay chỉ đạt 68.720 cổ phiếu, tương đương **0,2 lần MA20 volume**. Nói thẳng: lực bán hôm nay không phải đến từ một đợt tháo hàng ồ ạt của tổ chức hay nhà đầu tư lớn. Volume thấp đến mức này cho thấy thị trường gần như không có người mua, nhưng cũng không có áp lực bán mạnh theo nghĩa truyền thống — giá rơi trong trạng thái thanh khoản gần như đóng băng.

Đi kèm với đó, RSI(14) đang ở mức **16,2**. Để so sánh: RSI dưới 30 đã được gọi là oversold, RSI dưới 20 là vùng cực kỳ hiếm trong phần lớn cổ phiếu niêm yết. Chỉ số này phản ánh lực bán đã chiếm ưu thế hoàn toàn trong 14 phiên gần nhất theo tỷ lệ biên độ tích lũy.

MA trend hiện ở trạng thái **mixed** — tức các đường MA ngắn và dài chưa đồng thuận hướng đi, không xác nhận xu hướng rõ ràng ở thời điểm hiện tại.

---

### ② Vì sao điều này xảy ra

Data không cung cấp thông tin về catalyst cụ thể hôm nay, nhưng bức tranh kỹ thuật tự nó kể một câu chuyện. BCM là cổ phiếu thuộc nhóm bất động sản khu công nghiệp — nhóm nhạy cảm với lãi suất, dòng vốn FDI và kỳ vọng vĩ mô. Khi dòng tiền toàn thị trường co hẹp hoặc rủi ro vĩ mô tăng lên, nhóm này thường bị rút vốn trước.

Volume chỉ bằng 0,2× MA20 là dấu hiệu của một trong hai tình huống: hoặc nhà đầu tư đang chờ đợi thêm thông tin trước khi hành động, hoặc thanh khoản đã rút đi và giá đang "rơi tự do" trong khoảng trống khối lượng. Cả hai tình huống đều không cho thấy đây là đáy có xác nhận — nhưng cũng không mâu thuẫn với việc áp lực bán đang cạn dần.

---

### ③ Lịch sử cohort tương tự đã đi ra sao

Đây là phần mà dữ liệu định lượng có thể nói thay cảm tính.

Từ cơ sở dữ liệu lịch sử, cohort **RSI < 30 oversold** gồm **1.755 trường hợp** trên thị trường Việt Nam cho kết quả như sau sau 60 phiên:

- **Lợi nhuận trung bình (forward 60d): +6,32%**
- **Tỷ lệ thắng: 60%**
- **Edge so với baseline: +1,57%**

Nói cụ thể hơn: trong lịch sử, cứ 10 lần một cổ phiếu rơi vào vùng RSI oversold tương tự, có 6 lần giá cao hơn sau 60 phiên so với điểm vào. Edge +1,57% có nghĩa cohort này vượt trội so với mức baseline thị trường một cách có thống kê.

BCM hôm nay không chỉ chạm RSI < 30 — RSI đang ở **16,2**, tức nằm sâu hơn trong vùng oversold so với ngưỡng 30 mà cohort trên đã sàng lọc. Điều đó không tự động có nghĩa kết quả sẽ tốt hơn, nhưng BCM đang nằm trong phần đuôi cực trị của phân phối oversold — loại tín hiệu ít xuất hiện.

---

### ④ Cần theo dõi gì trong các phiên tới

Có ba biến số cần quan sát để xác nhận hoặc bác bỏ kịch bản phục hồi mà cohort lịch sử gợi ý:

**Khối lượng giao dịch**: Phục hồi từ vùng oversold thường đi kèm volume tăng trở lại, ít nhất về mức MA20. Nếu volume tiếp tục ở mức 0,2× MA20 hoặc thấp hơn, thanh khoản chưa quay trở lại — tín hiệu chưa được xác nhận.

**RSI vượt ngưỡng 30**: Đây là điều kiện kỹ thuật cơ bản để xác nhận thoát oversold. Mức 16 xuất hiện hôm nay là điểm khởi đầu, không phải điểm kết thúc.

**MA trend chuyển từ mixed sang aligned**: Khi MA ngắn hạn và dài hạn đồng thuận hướng lên, đó là xác nhận momentum đã thay đổi — thứ hiện tại chưa có.

---

### Verify reproducible

Dữ liệu và tính toán cohort trong bài có thể kiểm tra lại tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket):

```bash
pip install lotusmarket
```

```python
from lotusmarket import cohort
result = cohort.run(filters={"rsi_lt": 30}, forward_days=60)
print(result.summary())
```

Hoặc dùng CLI:

```bash
lmcli cohort --rsi-lt 30 --forward 60d --ticker BCM
```

---

⚠️ *Bài viết chỉ mang tính phân tích dữ liệu, không phải lời khuyên đầu tư. Mọi quyết định mua bán là trách nhiệm của nhà đầu tư.*
