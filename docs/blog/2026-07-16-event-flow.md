---
title: "Dòng tiền VN30 nghiêng bán hôm nay 16/07 — regime nào đang dẫn dắt?"
date: 2026-07-16
topic: flow
---

## Áp lực bán bao trùm VN30+HNX30: lịch sử nói gì?

Phiên hôm nay, dòng tiền nội tổng của rổ VN30+HNX30 ghi nhận buy-pressure chỉ ở mức **31.0%** — ngưỡng mà thị trường thường gọi thẳng là "bán mạnh". Câu hỏi không phải là *liệu* áp lực đang nặng hay không, mà là: **9 năm dữ liệu cohort đã cho thấy điều gì xảy ra sau những lần như vậy?**

---

### ① Chuyện gì đang xảy ra

Con số 31.0% buy-pressure không phải một cảnh báo mơ hồ — nó là một phép đo cụ thể trên toàn bộ dòng tiền nội của hai rổ cổ phiếu lớn nhất sàn. Khi tỷ lệ này chìm sâu dưới ngưỡng 50%, bên bán đang kiểm soát giao dịch theo nghĩa chữ số học thuần túy: cứ 10 đồng khớp lệnh, chưa đến 4 đồng đến từ phía mua chủ động.

Đây là tín hiệu regime **CRISIS** theo phân loại của hệ thống Lotus — không phải từ nhãn cảm xúc, mà từ ngưỡng định lượng được xác định trước bằng rolling volatility và breadth.

---

### ② Vì sao có thể hiểu được

Áp lực bán lan rộng trên cả VN30 lẫn HNX30 cùng một lúc thường không phải chuyện của một vài cổ phiếu bị tin xấu riêng lẻ. Khi breadth suy yếu đồng loạt ở cả hai sàn, dòng tiền nội đang phản ánh một tâm lý phòng thủ mang tính hệ thống — nhà đầu tư thu hẹp vị thế trước rủi ro vĩ mô hoặc thanh khoản, thay vì xoay vòng từ nhóm ngành này sang nhóm ngành khác.

Buy-pressure 31.0% đặt phiên hôm nay vào nhóm quan sát hiếm — đủ thấp để kích hoạt phân loại CRISIS trong cơ sở dữ liệu lịch sử 9 năm của Lotus.

---

### ③ Cohort lịch sử đã đi tiếp ra sao

Đây là phần dữ liệu đáng đọc kỹ. 📊

Khi thị trường rơi vào regime **CRISIS** — tức buy-pressure ở vùng tương tự hôm nay — **forward return trung bình sau 60 phiên là +8.53%**, với tỷ lệ thắng **67%**. Nói cụ thể hơn: trong 9 năm, cứ 3 lần thị trường rơi vào cohort này, có 2 lần chỉ số đứng cao hơn sau 60 ngày làm việc.

Để có ngữ cảnh so sánh: regime **STABLE** — trạng thái bình thường khi áp lực mua cân bằng — chỉ cho forward return **+3.20%**, thấp hơn cả baseline lịch sử tổng thể **+4.75%**, và tỷ lệ thắng chỉ **48%** (thực chất là tung đồng xu).

Điều này nghe có vẻ phản trực giác: thị trường bán mạnh lại đi kèm forward return kỳ vọng cao hơn. Nhưng đây là hiện tượng **mean-reversion asymmetry** được ghi nhận trong nhiều nghiên cứu thị trường mới nổi — khi áp lực bán đã bị đẩy đến ngưỡng cực đoan, định giá bắt đầu phản ánh rủi ro quá mức, tạo ra vùng đệm cho phục hồi.

Lưu ý quan trọng: tỷ lệ thắng 67% không có nghĩa là mỗi lần đều phục hồi ngay lập tức, hay mức độ giảm thêm trước khi phục hồi là nhỏ. 33% còn lại là các trường hợp thị trường vẫn thấp hơn sau 60 ngày — và đó là con số không nên bỏ qua.

---

### ④ Cần theo dõi gì từ đây

Ba biến số mình sẽ quan sát trong những phiên tới:

**Buy-pressure có hồi phục về vùng 40–50% không?** Đây là điều kiện cần để xác nhận áp lực bán đã cạn — nếu chỉ số này tiếp tục ở dưới 35%, cohort CRISIS kéo dài và kịch bản phục hồi bị trì hoãn.

**Breadth giữa VN30 và HNX30 có phân kỳ không?** Nếu một sàn bắt đầu cải thiện buy-pressure trong khi sàn kia vẫn ở vùng bán mạnh, đó là dấu hiệu dòng tiền đang chọn lọc lại — thường xuất hiện trước khi toàn thị trường ổn định.

**Forward return 60 ngày theo cohort CRISIS (+8.53%)** là kết quả trung bình — phân phối thực tế rộng hơn nhiều. Các phiên tiếp theo sẽ cho thấy biến động nằm ở đuôi nào của phân phối đó.

---

### Verify reproducible

Toàn bộ tính toán cohort trên có thể tái tạo qua [lotusmarket](https://github.com/ducnhd/lotusmarket):

```bash
pip install lotusmarket
```

```python
from lotusmarket import FlowRegime
regime = FlowRegime(universe=["VN30", "HNX30"], window=9*252)
regime.cohort_analysis(buy_pressure=0.31, forward_days=60)
```

Hoặc dùng CLI:

```bash
lmcli flow-regime --universe VN30,HNX30 --bp 0.31 --fwd 60
```

---

*Bài viết chỉ phản ánh phân tích định lượng từ dữ liệu lịch sử — không phải lời khuyên đầu tư. Mọi quyết định giao dịch thuộc trách nhiệm của nhà đầu tư.*
