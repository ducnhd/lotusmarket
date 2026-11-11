---
title: "Dòng tiền VN30 nghiêng bán hôm nay 11/11 — regime nào đang dẫn dắt?"
date: 2026-11-11
topic: flow
---

## Áp lực bán phủ kín VN30+HNX30 — lịch sử 9 năm nói gì về đợt này?

Hôm nay dòng tiền nội tổng của rổ VN30+HNX30 chỉ ghi nhận buy-pressure ở mức **27.4%** — ngưỡng mà thị trường thường gọi thẳng là "bán mạnh". Con số này không phải cảm tính, mà là tỷ lệ khối lượng mua chủ động so với tổng khớp lệnh. Câu hỏi mà data 9 năm sẽ trả lời: khi áp lực bán áp đảo đến mức này, thị trường Việt thường đi về đâu trong 60 phiên tiếp theo?

---

### ① Chuyện gì đang xảy ra

Buy-pressure 27.4% có nghĩa là gần **3/4 khối lượng khớp lệnh** trong ngày hôm nay đến từ phía bán chủ động. Nói cách khác, người bán không chờ — họ đẩy lệnh ra thị trường, chấp nhận giá thấp hơn để thoát hàng. Đây là tín hiệu áp lực, không phải xu hướng dài hạn, nhưng nó đủ mạnh để kích hoạt phân loại **regime CRISIS** trong hệ thống phân tích cohort của Lotus AI.

Cần phân biệt rõ: buy-pressure thấp không tự động nghĩa là giá sẽ tiếp tục giảm. Nó nghĩa là **cung đang áp đảo cầu ở thời điểm hiện tại** — một trạng thái tạm thời mà lịch sử đã lặp đi lặp lại nhiều lần.

---

### ② Vì sao trạng thái này xuất hiện

Dữ liệu không ghi nhận một sự kiện đơn lẻ nào làm nguyên nhân, nhưng bức tranh dòng tiền nội tổng trải đều trên cả hai rổ VN30 và HNX30 cho thấy áp lực bán không chỉ tập trung vào một nhóm ngành. Khi áp lực lan rộng như vậy — không phải chỉ một vài mã blue-chip hay midcap đơn lẻ chịu xả — thị trường thường đang phản ứng với một yếu tố vĩ mô hoặc tâm lý chung, chứ không phải thông tin riêng lẻ của doanh nghiệp.

Điều đáng chú ý là buy-pressure 27.4% đủ thấp để hệ thống phân loại đây là vùng **CRISIS**, tách biệt hoàn toàn khỏi vùng **STABLE** (nơi thị trường giao dịch cân bằng hơn). Hai regime này có hành vi thống kê rất khác nhau — và đó là lý do việc phân loại đúng regime quan trọng hơn việc đoán hướng thị trường ngay lập tức.

---

### ③ Cohort lịch sử 9 năm đã đi ra sao

Đây là phần data thực sự nói chuyện. Trong 9 năm dữ liệu lịch sử, hai regime cho kết quả **60 phiên tới** như sau:

| Regime | Fwd Return 60d | Win Rate |
|--------|---------------|----------|
| **CRISIS** | **+8.53%** | **67%** |
| STABLE | +3.20% | 48% |
| *(Baseline toàn bộ)* | *+4.75%* | *—* |

Số liệu này cần đọc đúng cách. Khi thị trường rơi vào regime CRISIS — tức buy-pressure xuống đáy, bán mạnh áp đảo — **67% các lần trong lịch sử**, chỉ số tiến lên trong 60 phiên tiếp theo với mức bình quân **+8.53%**. Con số này gấp đôi so với return của regime STABLE (+3.20%), và vượt xa baseline +4.75%.

Ngược lại, khi thị trường ở trạng thái STABLE — cân bằng, không có áp lực cực đoan — win rate chỉ đạt **48%**, tức thấp hơn cả tung đồng xu. Return bình quân +3.20% cũng thấp hơn baseline.

Lý giải thống kê ở đây có logic rõ ràng: khi bán tháo đạt đỉnh, phần lớn áp lực xả hàng đã được tiêu hóa. Các phiên STABLE ngược lại thường là giai đoạn thị trường "lơ lửng" — không có catalyst rõ, dòng tiền chờ đợi, và kết quả phân kỳ lớn hơn nhiều.

**Lưu ý quan trọng**: 67% win rate không phải 100%. Trong 33% trường hợp còn lại, thị trường tiếp tục giảm hoặc đi ngang quá 60 phiên. Data cho thấy cohort thắng ở xác suất cao hơn hẳn, nhưng không có gì chắc chắn tuyệt đối.

---

### ④ Cần theo dõi gì trong các phiên tới

Ba điểm cần quan sát:

**Buy-pressure phục hồi về trên 40%** — đây là dấu hiệu cầu bắt đầu quay lại, thường xuất hiện 3-7 phiên sau đáy áp lực. Nếu buy-pressure tiếp tục ở dưới 30% thêm nhiều phiên liên tiếp, regime CRISIS có thể kéo dài và cohort lịch sử sẽ cần cập nhật lại.

**Phân kỳ giữa VN30 và HNX30** — nếu hai rổ bắt đầu tách nhau (một rổ phục hồi buy-pressure nhanh hơn), đó là tín hiệu dòng tiền đang chọn lọc nhóm cổ phiếu cụ thể thay vì bán/mua đồng loạt.

**Khối lượng khớp lệnh tổng** — nếu áp lực bán 27.4% đi kèm khối lượng thấp, đó là bán của thiểu số. Nếu khối lượng cao bất thường, áp lực thực sự lớn hơn nhiều.

---

### Verify reproducible 🔍

```bash
pip install lotusmarket
```

```python
from lotusmarket import flow, cohort

bp = flow.buy_pressure(tickers="VN30+HNX30", date="today")
result = cohort.regime_analysis(buy_pressure=bp, forward_days=60, lookback_years=9)
print(result.summary())
```

Hoặc dùng CLI:

```bash
lmcli flow --ticker VN30+HNX30 --cohort --fwd 60
```

Toàn bộ methodology và data source có tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

---

*Bài viết phân tích dữ liệu thống kê lịch sử, không phải lời khuyên đầu tư. Mọi quyết định giao dịch là trách nhiệm của nhà đầu tư.*
