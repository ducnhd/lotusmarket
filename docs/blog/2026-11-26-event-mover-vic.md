---
title: "VIC giảm 6.3% hôm nay 26/11 — cohort lịch sử nói gì?"
date: 2026-11-26
topic: mover
---

## VIC rơi 6.28% với thanh khoản chỉ bằng 1/10 trung bình — thị trường đang im lặng hay lảng tránh?

Phiên hôm nay, Vingroup (VIC) đóng cửa ở mức 230,000 VND, mất 6.28% chỉ trong một phiên. Điều bất thường không nằm ở mức giảm — mà nằm ở chỗ cú rơi này xảy ra trong tình trạng gần như vắng lặng về khối lượng. Câu hỏi mà data đặt ra: đây là nhát chém của kẻ bán tháo, hay chỉ là khoảng trống thanh khoản kéo giá xuống mà không có người mua đỡ?

---

### ① Chuyện gì đã xảy ra hôm nay

Số liệu không che giấu gì: VIC giảm 6.28%, đóng cửa tại 230,000 VND. Volume khớp lệnh đạt 210,620 cổ phiếu — chỉ bằng **0.1× MA20 volume**, tức là thanh khoản hôm nay chỉ bằng 10% mức trung bình 20 phiên gần nhất. RSI(14) hiện ở 37.4, đang tiến vào vùng tiệm cận oversold (ngưỡng 30) nhưng chưa chạm đáy kỹ thuật.

MA trend được ghi nhận là **mixed** — tức là các đường MA ngắn hạn và dài hạn đang không đồng thuận, không có tín hiệu xu hướng rõ ràng theo cả hai chiều.

Kết hợp lại: giá giảm mạnh, volume cực thấp, RSI trung lập-yếu, MA không định hướng. Đây là bộ tứ chỉ báo khá hiếm gặp — và không phải lúc nào cũng dễ đọc.

---

### ② Vì sao cú rơi này đáng chú ý hơn bề ngoài

Một cổ phiếu giảm 6% với volume bình thường hoặc cao thường phản ánh áp lực bán thực sự — dòng tiền rút ra có chủ đích. Nhưng VIC hôm nay rơi 6.28% với volume **chỉ 10% MA20**. Điều này có nghĩa: lực bán không lớn, nhưng lực cầu gần như biến mất hoàn toàn.

Trong cấu trúc thị trường, khi thanh khoản cạn kiệt đến mức này, một lượng bán nhỏ cũng có thể kéo giá xuống phi tuyến. Không cần "cá mập" thoát hàng — chỉ cần người mua ngừng đặt lệnh là giá tự trượt.

Với VIC — cổ phiếu dẫn đầu hệ sinh thái Vingroup, vốn hóa lớn nhất sàn HOSE — việc volume co lại đến mức này trong một phiên giảm sâu cho thấy dòng tiền đang **né tránh chủ động**, không phải bán tháo hoảng loạn. Đây là hai trạng thái hoàn toàn khác nhau về hàm ý.

RSI 37.4 xác nhận áp lực bán tích lũy trong ngắn hạn, nhưng chưa đủ để kích hoạt tín hiệu kỹ thuật phục hồi từ vùng oversold. MA mixed có nghĩa là momentum chưa chọn phe.

---

### ③ Cohort lịch sử tương tự đã đi về đâu

Đây là phần data nói thẳng nhất: **cohort lịch sử không match pattern mạnh** cho tổ hợp hôm nay. Baseline cohort edge xấp xỉ **~0%** — tức là trong các phiên lịch sử có cùng đặc điểm (giảm mạnh + volume cực thấp + RSI vùng 35-40 + MA mixed), kết quả tiếp theo không nghiêng rõ về phía nào.

Nói cụ thể hơn: lịch sử không cho thấy đây là đáy để bật mạnh, cũng không xác nhận đây là điểm bắt đầu của xu hướng giảm kéo dài. Vùng này là **neutral về mặt thống kê**.

Điều đó có nghĩa là: những ai muốn dựa vào precedent lịch sử để ra quyết định sẽ không tìm được cơ sở thống kê đủ mạnh từ phiên hôm nay. Data không ủng hộ bull, cũng không ủng hộ bear — đây là vùng mà việc chờ đợi thêm tín hiệu có độ chắc chắn cao hơn hành động ngay.

---

### ④ Cần theo dõi gì trong các phiên tới

Ba biến số quan trọng nhất cần quan sát:

**Volume phục hồi hay tiếp tục co cụm.** Nếu phiên kế tiếp VIC có volume quay về 0.5× MA20 trở lên kèm giá ổn định hoặc tăng nhẹ, đó là tín hiệu dòng tiền bắt đầu quay lại. Ngược lại, volume tiếp tục dưới 0.2× MA20 kèm giá giảm thêm là dấu hiệu khoảng trống thanh khoản chưa được lấp.

**RSI có xuống dưới 30 không.** RSI 37.4 còn cách ngưỡng oversold kỹ thuật một khoảng. Nếu RSI thủng 30 trong 1-2 phiên tới, lịch sử VIC thường tạo ra phản ứng kỹ thuật ngắn hạn — nhưng độ bền của phản ứng đó phụ thuộc vào volume đi kèm.

**MA trend có hội tụ không.** Trạng thái mixed hiện tại cần được giải quyết. MA ngắn hạn cắt xuống dưới MA dài hạn sẽ xác nhận downtrend; ngược lại là setup phục hồi. Phiên hôm nay chưa cho câu trả lời.

---

### Verify reproducible 🔍

```bash
pip install lotusmarket
```

```python
from lotusmarket import stock
vic = stock("VIC")
vic.show(["rsi", "volume_ratio", "ma_trend"])
```

Hoặc dùng CLI:

```bash
lmcli scan --ticker VIC --indicators rsi,volume_ma_ratio,ma_trend
```

Toàn bộ methodology và cohort backtesting có thể kiểm tra tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

---

*Bài viết phân tích dữ liệu thị trường, không phải lời khuyên đầu tư. Mọi quyết định giao dịch thuộc trách nhiệm cá nhân nhà đầu tư.*
