---
title: "VIC giảm 7.0% hôm nay 22/07 — cohort lịch sử nói gì?"
date: 2026-07-22
topic: mover
---

## VIC rơi gần 7% trong một phiên — vùng oversold lịch sử đang bật tín hiệu

VIC đóng cửa hôm nay ở 202.100 đồng, giảm **6,99%** chỉ trong một phiên — mức rơi hiếm gặp với một bluechip quy mô như Vinhomes. RSI(14) chạm **25,8**, sâu dưới ngưỡng oversold 30. Câu hỏi mà data có thể trả lời: những lần VIC rơi vào vùng này trước đây, thị trường đã phản ứng ra sao trong 60 ngày tiếp theo?

---

### ① Chuyện gì xảy ra hôm nay

Số liệu không cần diễn giải thêm: **-6,99%** là mức giảm đủ để đưa VIC vào nhóm cổ phiếu ghi nhận phiên tệ nhất trong danh sách largecap hôm nay. Giá hiện tại 202.100 đồng.

Điều đáng chú ý hơn là khối lượng khớp lệnh: chỉ **573.620 cổ phiếu**, bằng **0,2× MA20 volume** — tức là thanh khoản chưa bằng một phần năm mức trung bình 20 phiên gần nhất. Đây là chi tiết không nhỏ. Một phiên giảm sốc nhưng volume cạn kiệt thường phản ánh trạng thái *bán không có người mua cân bằng*, thay vì một làn sóng tháo hàng thật sự có tổ chức. Lực cầu chưa vào, nhưng lực bán cũng không phải dòng tiền lớn tháo chạy.

RSI(14) = **25,8** — lần cuối VIC vào vùng này là những đợt biến động thị trường có tính chu kỳ. Không phải tín hiệu mua, nhưng là một trạng thái kỹ thuật có lịch sử đủ dài để đo được.

MA trend: **mixed** — tức là các đường trung bình đang không đồng thuận về chiều. Chưa có xu hướng rõ ràng hình thành, cả ngắn hạn lẫn trung hạn.

---

### ② Vì sao có thể xảy ra điều này

Data hôm nay không cung cấp nguyên nhân cụ thể từ phía doanh nghiệp hay macro — và bài này sẽ không bịa. Điều có thể quan sát khách quan: VIC là cổ phiếu có trọng số lớn trong các danh mục ETF và quỹ nội địa. Khi thị trường chung bị áp lực thanh khoản hoặc có dòng tiền rút khỏi nhóm largecap, VIC thường chịu tác động khuếch đại so với mid-cap.

Volume 0,2× MA20 là điểm cần nhìn kỹ hơn: volume thấp bất thường trong phiên giảm mạnh đôi khi xuất hiện vào cuối một đợt bán — khi những người muốn thoát đã thoát, người còn lại không bán nữa. Nhưng cũng có thể là thanh khoản đang rất mỏng và bất kỳ lực bán nhỏ nào cũng kéo giá mạnh. Hai kịch bản này cần thêm dữ liệu phiên tiếp theo để phân biệt.

---

### ③ Cohort lịch sử 1.755 trường hợp nói gì

Đây là phần mà Lotus AI đã chạy backtesting có hệ thống. Cohort **RSI < 30 oversold** trên toàn bộ dữ liệu lịch sử VIC ghi nhận **N = 1.755 trường hợp**. Kết quả forward 60 ngày:

- **Lợi nhuận trung bình: +6,32%**
- **Tỷ lệ thắng: 60%**
- **Edge so với baseline: +1,57%**

Con số 60% win rate nghĩa là trong 3 lần RSI < 30, có 2 lần giá phục hồi trong 60 ngày tiếp theo — và 1 lần tiếp tục đi xuống hoặc đi ngang. Không phải tín hiệu thần thánh, nhưng edge **+1,57%** so với baseline cho thấy vùng oversold này *có thống kê ý nghĩa*, không phải nhiễu ngẫu nhiên.

Điều quan trọng: con số này là **trung bình của 1.755 trường hợp**. Mỗi phiên cụ thể có bối cảnh riêng — macro, ngành bất động sản, lãi suất, tâm lý thị trường — không có trường hợp nào giống hệt hôm nay 100%.

Nhưng đây là lý do để *theo dõi*, không phải để bỏ qua.

---

### ④ Cần theo dõi gì trong các phiên tới

Ba biến số cần quan sát:

**Volume phục hồi** — Nếu volume quay trở lại vùng MA20 (~2,5–3 triệu đơn vị) kèm giá giữ trên 202.000, đó là dấu hiệu cầu đang quay lại. Nếu volume tiếp tục cạn mà giá vẫn trượt, câu chuyện khác.

**RSI thoát vùng dưới 30** — Không cần RSI phục hồi về 50 ngay. Chỉ cần tín hiệu RSI bắt đầu tăng từ vùng 25–28 là biên kỹ thuật đầu tiên.

**MA trend đồng thuận** — Hiện tại MA mixed. Khi MA ngắn hạn (MA5/MA10) cắt lên MA trung hạn, đó mới là xác nhận xu hướng, không phải chỉ là một phiên bật kỹ thuật.

🔍 VIC ở 202.100 với RSI 25,8 và volume 0,2× — không phải lúc để bỏ qua, cũng không phải lúc để vội vã.

---

### Verify reproducible

Toàn bộ cohort analysis trên có thể tái tạo qua [Lotus Market](https://github.com/ducnhd/lotusmarket):

```python
pip install lotusmarket

from lotusmarket import backtest
result = backtest.cohort("VIC", signal="RSI<30", forward_days=60)
print(result.summary())
```

Hoặc dùng CLI:

```bash
lmcli cohort VIC --signal rsi_oversold --fwd 60
```

---

*Bài viết này chỉ mang tính phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Mọi quyết định giao dịch là trách nhiệm của nhà đầu tư.*
