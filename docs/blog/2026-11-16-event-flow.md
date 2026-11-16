---
title: "Dòng tiền VN30 nghiêng bán hôm nay 16/11 — regime nào đang dẫn dắt?"
date: 2026-11-16
topic: flow
---

## Áp lực bán phủ kín VN30+HNX30 — lịch sử nói gì?

Hôm nay dòng tiền nội trong rổ VN30+HNX30 chỉ còn **31.8% buy-pressure** — tức gần 7 phần 10 thanh khoản đang nghiêng về phía bán. Con số này đủ để xếp phiên hôm nay vào nhóm **"bán mạnh"** theo thang phân loại của Lotus AI. Câu hỏi mà data 9 năm sẽ trả lời: sau những lần thị trường rơi vào trạng thái tương tự, giá đi về đâu trong 60 ngày tiếp theo?

---

### ① Chuyện gì đang xảy ra

Buy-pressure 31.8% không phải con số ngẫu nhiên — đây là tín hiệu dòng tiền nội **áp đảo chiều bán** trên cả hai rổ blue-chip lớn nhất sàn. Khi buy-pressure rơi xuống vùng này, thị trường thường đang trải qua một trong hai kịch bản: hoặc lực cầu thực sự yếu, hoặc bên mua đang rút lui chờ giá tốt hơn trước khi tái phân bổ.

Điểm đáng chú ý là tín hiệu này xuất hiện đồng thời trên **cả VN30 lẫn HNX30** — không phải cô lập ở một sàn. Điều đó cho thấy đây là áp lực mang tính hệ thống, không phải câu chuyện riêng của một vài bluechip.

---

### ② Bối cảnh: tại sao áp lực bán lại dồn về đây

Dữ liệu không ghi nhận một sự kiện đơn lẻ nào làm bùng phát đợt bán này, nhưng cấu trúc dòng tiền nội cho thấy bên mua đang **không sẵn sàng đỡ giá** ở vùng hiện tại. Trong môi trường lãi suất và tỷ giá còn nhiều biến số, nhóm tổ chức nội thường chọn cách quan sát thay vì dẫn dắt — và con số 31.8% phản ánh đúng tâm lý đó.

Một cách đọc khác: buy-pressure thấp không tự động có nghĩa là thị trường sẽ tiếp tục giảm. Nó có nghĩa là **bên mua đang vắng mặt**, và khoảng trống đó đang bị lấp bởi lực bán. Hai điều này khác nhau về bản chất — và chính sự phân biệt này mới quan trọng khi nhìn vào dữ liệu lịch sử.

---

### ③ Cohort 9 năm nói gì

Đây là phần data thực sự đáng đọc kỹ.

Lotus AI phân loại thị trường theo **regime** — tức trạng thái nền của thị trường tại thời điểm tín hiệu xuất hiện. Trong 9 năm dữ liệu, các phiên có cấu trúc dòng tiền tương tự hôm nay được chia thành hai nhóm chính:

**CRISIS regime** — những lần thị trường đang trong giai đoạn căng thẳng, bất ổn:
- Forward return 60 ngày: **+8.53%**
- Tỷ lệ thắng: **67%**

**STABLE regime** — những lần thị trường đang ở trạng thái bình thường:
- Forward return 60 ngày: **+3.20%**
- Tỷ lệ thắng: **48%** (thấp hơn baseline +4.75%)

Hai con số này kể một câu chuyện rõ ràng: **áp lực bán mạnh trong giai đoạn khủng hoảng lại đi kèm với kỳ vọng phục hồi cao hơn hẳn giai đoạn ổn định.** Lý do không khó hiểu — khi mọi người đang bán tháo vì sợ hãi, định giá thường bị đẩy xuống thấp hơn giá trị thực, và phục hồi là kết quả tự nhiên.

Ngược lại, trong STABLE regime, buy-pressure thấp lại cho kết quả **kém hơn baseline** (3.20% so với 4.75%). Đây là dấu hiệu của sự phân phối có kiểm soát — bên mua không hoảng loạn, nhưng bên bán đang có lý do chiến lược để thoát hàng. 🔍

Hiện tại chúng ta đang ở regime nào? Data không tự trả lời câu hỏi đó trong một phiên — nhưng sự phân kỳ giữa hai cohort đủ để xác định rằng **việc đọc đúng regime là biến số quan trọng nhất**, không phải bản thân con số 31.8%.

---

### ④ Cần theo dõi gì trong những phiên tới

Ba điểm cần quan sát:

1. **Buy-pressure có hồi phục về trên 40% không?** Vùng 40-45% thường là ngưỡng mà bên mua bắt đầu tái xuất hiện có tổ chức.

2. **Regime hiện tại xác nhận là CRISIS hay STABLE?** Nếu là CRISIS, cohort lịch sử cho thấy 67% phiên tương tự đã có forward return dương sau 60 ngày. Nếu là STABLE, baseline còn cao hơn kết quả cohort — tức tín hiệu này không có lợi thế thống kê.

3. **Sự đồng thuận giữa VN30 và HNX30 có tiếp diễn không?** Hai rổ cùng bán mạnh hôm nay — nếu HNX30 phục hồi buy-pressure trước, đó thường là dấu hiệu dòng tiền đang tìm vùng trú ẩn ở mid-cap trước khi quay lại bluechip. 📊

---

### Verify reproducible

```bash
pip install lotusmarket
```

```python
from lotusmarket import FlowScanner, RegimeCohort

flow = FlowScanner(tickers="VN30+HNX30").buy_pressure()
cohort = RegimeCohort(signal="buy_pressure_low", fwd_days=60).fit(years=9)
print(cohort.summary())
```

Hoặc dùng CLI: `lmcli cohort --signal bp_low --fwd 60 --regime all`

Toàn bộ methodology và data source có tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

---

*Bài viết phân tích dữ liệu thống kê lịch sử, không phải lời khuyên đầu tư. Mọi quyết định giao dịch là trách nhiệm của nhà đầu tư.*
