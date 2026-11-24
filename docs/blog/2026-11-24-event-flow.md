---
title: "Dòng tiền VN30 nghiêng bán hôm nay 24/11 — regime nào đang dẫn dắt?"
date: 2026-11-24
topic: flow
---

## Áp lực bán phủ kín VN30+HNX30 — lịch sử nói gì về những lần như thế này?

Hôm nay dòng tiền nội tổng của rổ VN30+HNX30 chỉ ghi nhận **buy-pressure 28.0%** — tức là chưa đến 3 phần 10 thanh khoản thuần đang nghiêng về phía mua. Con số này đủ để xếp phiên hôm nay vào nhóm "bán mạnh" theo cách phân loại của Lotus AI. Câu hỏi không phải là thị trường có đang xấu hay không — mà là **những lần buy-pressure xuống thấp như vậy trong lịch sử, thị trường đã đi tiếp ra sao?**

---

### ① Chuyện gì đang xảy ra

Buy-pressure 28% nghĩa là gần **72% dòng tiền trong rổ blue-chip hai sàn đang nghiêng về phía bán**. Đây không phải trạng thái thận trọng bình thường — đây là trạng thái mà người ta thường gọi là "xả hàng có chủ đích." Khi cả VN30 lẫn HNX30 cùng thể hiện áp lực này một lúc, tín hiệu có trọng lượng hơn so với việc chỉ một rổ đơn lẻ bị bán.

Điểm đáng chú ý: buy-pressure là chỉ số **dòng tiền thực**, không phải giá. Giá có thể chưa phản ánh hết — hoặc đã phản ánh xong — tùy thuộc vào regime thị trường tại thời điểm đó. Và đây chính là chỗ dữ liệu lịch sử 9 năm bắt đầu có ích.

---

### ② Hai thế giới khác nhau: CRISIS vs STABLE

Trong 9 năm dữ liệu, Lotus AI phân loại thị trường theo **regime** — tạm hiểu là "trạng thái nền" của thị trường tại mỗi thời điểm. Kết quả cohort khi buy-pressure xuống vùng bán mạnh như hôm nay cho thấy hai kịch bản tách biệt rõ ràng:

**Regime CRISIS:**
- Forward return 60 ngày: **+8.53%**
- Win rate: **67%**

**Regime STABLE:**
- Forward return 60 ngày: **+3.20%**
- Win rate: **48%** (thấp hơn cả baseline +4.75%)

Đọc số này cần cẩn thận. Trong STABLE, khi áp lực bán mạnh xuất hiện, lịch sử cho thấy win rate **thấp hơn cả baseline** — tức là tín hiệu bán mạnh trong môi trường ổn định không phải cơ hội, mà gần hơn với cảnh báo. Ngược lại, trong CRISIS, cùng mức buy-pressure đó lại đi kèm với forward return gần **gấp đôi baseline** và win rate 67%.

Nói ngắn gọn hơn: **bán mạnh trong CRISIS thường là đáy cục bộ — bán mạnh trong STABLE thường là phần đầu của đợt giảm tiếp.**

---

### ③ Tại sao lại có sự phân kỳ này?

Giải thích hợp lý nhất từ góc độ cấu trúc dòng tiền: trong CRISIS, bán mạnh phản ánh **panic** — cảm xúc, không phải thông tin. Người bán không nhất thiết biết nhiều hơn người mua; họ chỉ sợ hơn. Kết quả là giá bị đẩy xuống dưới giá trị, và hồi phục về sau là tự nhiên.

Trong STABLE, bán mạnh lại thường phản ánh **rotation** hoặc **thông tin dẫn trước** — tức là ai đó đang thoát ra trước khi câu chuyện xấu được công bố rộng rãi. Lúc này, áp lực bán là tín hiệu, không phải nhiễu.

Vấn đề là hôm nay thị trường đang ở regime nào? Đó là điều cần xác định — và cũng là điều dữ liệu realtime của Lotus AI theo dõi liên tục.

---

### ④ Cần theo dõi gì trong những phiên tới

Dựa trên cấu trúc cohort, có **ba biến số** đáng quan sát:

1. **Regime tiếp tục CRISIS hay chuyển STABLE?** — Nếu regime xác nhận CRISIS, cohort lịch sử hỗ trợ luận điểm hồi phục trong 60 ngày. Nếu chuyển STABLE trong trạng thái bán mạnh hiện tại, lịch sử cho thấy cẩn trọng hơn.

2. **Buy-pressure có phục hồi lên trên 40% trong 3-5 phiên tiếp theo không?** — Đây là mức mà dòng tiền thường bắt đầu cân bằng lại trong các cohort CRISIS hồi phục.

3. **HNX30 hay VN30 dẫn đầu trong việc phục hồi buy-pressure?** — HNX30 thường phản ứng sớm hơn vì thanh khoản thấp hơn, dễ bị ảnh hưởng bởi dòng tiền tổ chức khi họ bắt đầu giải ngân.

---

### Verify reproducible

Toàn bộ tính toán cohort trên có thể tái hiện qua [lotusmarket](https://github.com/ducnhd/lotusmarket):

```python
pip install lotusmarket

from lotusmarket import FlowAnalyzer
fa = FlowAnalyzer()
fa.cohort_regime(basket=["VN30","HNX30"], signal="buy_pressure_lt_30", fwd_days=60)
```

Hoặc nếu dùng CLI:

```bash
lmcli cohort --basket VN30,HNX30 --signal bp_low --fwd 60 --by regime
```

---

⚠️ *Bài viết phân tích dữ liệu lịch sử và dòng tiền — không phải lời khuyên đầu tư. Mọi quyết định giao dịch là trách nhiệm cá nhân.*
