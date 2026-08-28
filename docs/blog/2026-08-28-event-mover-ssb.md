---
title: "SSB tăng 6.9% hôm nay 28/08 — cohort lịch sử nói gì?"
date: 2026-08-28
topic: mover
---

## SSB bùng phá +6.88% — nhưng volume nói điều khác

SeABank (SSB) đóng cửa hôm nay ở mức 17.100 VND, tăng gần 6.88% chỉ trong một phiên. Con số đó đủ để thu hút ánh mắt của bất kỳ ai theo dõi bảng điện. Câu hỏi mà data đặt ra ngay lập tức: liệu đây là một cú phá vỡ thật sự, hay chỉ là một nhịp bắn giá thiếu nền tảng?

---

### ① Chuyện gì đã xảy ra hôm nay

Mức tăng 6.88% trong một phiên là một sự kiện không bình thường với cổ phiếu ngân hàng vốn thường di chuyển chậm và ổn định. Giá đóng cửa 17.100 VND đưa SSB về vùng giá đáng chú ý sau một thời gian sideway.

Tuy nhiên, có một chi tiết không thể bỏ qua: volume giao dịch hôm nay chỉ đạt 490.560 cổ phiếu — bằng **0.2× so với MA20 volume**. Nói thẳng ra, khối lượng hôm nay chỉ bằng một phần năm mức trung bình 20 phiên. Một phiên tăng mạnh mà lượng cầu mỏng như vậy là tín hiệu cần đọc kỹ, không phải lý do để hứng khởi.

RSI(14) hiện ở mức **77.5** — đã vượt qua ngưỡng overbought 70 một cách rõ ràng. Trong khi đó, xu hướng MA được ghi nhận là **mixed** — tức là các đường trung bình ngắn và dài hạn chưa xếp hàng đồng thuận theo một hướng nhất định.

---

### ② Vì sao có thể có nhịp tăng này

Nhìn vào bối cảnh rộng hơn, nhóm cổ phiếu ngân hàng vừa và nhỏ trong nước thỉnh thoảng xuất hiện những phiên bùng giá đột biến kiểu này, thường đi kèm với một trong các yếu tố: tin tức nội bộ doanh nghiệp, dòng tiền quay vòng từ nhóm cổ phiếu vừa tăng sang nhóm chưa tăng, hoặc lực cầu từ một số tài khoản tập trung đẩy giá trong bối cảnh thanh khoản thấp.

Đây chính là điểm mấu chốt: **khi volume chỉ bằng 0.2× MA20**, một lực mua tương đối nhỏ cũng có thể tạo ra biến động giá lớn hơn bình thường. Điều đó giải thích tại sao giá có thể nhảy gần 7% mà không cần đến một dòng tiền thực sự mạnh và bền vững đứng phía sau.

MA trend mixed cũng xác nhận: cổ phiếu chưa có xu hướng rõ ràng. Không có đường MA nào đang dẫn dắt giá theo nghĩa kỹ thuật cổ điển.

---

### ③ Cohort lịch sử tương tự đã đi ra sao

Đây là phần mà dữ liệu nói thẳng nhất.

Khi chạy cohort lịch sử với pattern tương tự — giá tăng mạnh trong một phiên, RSI vượt 70, volume thấp hơn đáng kể so với trung bình, MA trend mixed — kết quả cho thấy: **cohort không match pattern mạnh**. Baseline cohort edge xấp xỉ **~0%**.

Điều đó có nghĩa là: trong lịch sử, những lần SSB hoặc các cổ phiếu tương tự rơi vào trạng thái kỹ thuật này, không có edge thống kê rõ ràng nào cho nhịp tiếp theo — cả theo chiều tăng lẫn chiều giảm. Đây là **vùng neutral thực sự**, không phải vùng momentum mạnh.

Lịch sử không cho thấy đây là điểm bắt đầu của một xu hướng tăng bền. Cũng không khẳng định đây là đỉnh để reversal. Cohort edge ~0% nghĩa là dữ liệu quá khứ không cung cấp xác suất thiên lệch đủ mạnh theo hướng nào.

---

### ④ Cần theo dõi gì trong những phiên tới

Ba yếu tố cần quan sát kỹ:

**Volume xác nhận.** Nếu SSB tiếp tục tăng hoặc giữ giá trong các phiên tới mà volume hồi về gần hoặc vượt MA20, đó là dấu hiệu dòng tiền thực sự đang vào. Nếu volume tiếp tục mỏng, nhịp tăng hôm nay khó có nền tảng để duy trì.

**RSI hạ nhiệt hay tiếp tục leo.** RSI 77.5 là mức cao. Lịch sử cho thấy ở vùng này, nếu không có momentum volume hỗ trợ, giá thường consolidate hoặc pullback nhẹ để RSI về vùng 60-65 trước khi có nhịp tiếp theo.

**MA trend thoát khỏi trạng thái mixed.** Khi các đường MA ngắn hạn (5, 10 phiên) bứt lên trên các đường dài hơn một cách rõ ràng và ổn định, đó là lúc cohort lịch sử bắt đầu có edge dương rõ hơn. Hiện tại, mixed trend là tín hiệu chờ đợi.

---

### Verify reproducible 🔍

Toàn bộ phân tích cohort trên có thể tái tạo với:

```bash
pip install lotusmarket
```

```python
from lotusmarket import cohort_scan
result = cohort_scan(ticker="SSB", rsi_min=70, volume_ratio_max=0.3, ma_trend="mixed")
print(result.summary())
```

Hoặc qua CLI:

```bash
lmcli cohort --ticker SSB --rsi 77.5 --vol-ratio 0.2 --ma mixed
```

Tham khảo thêm tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

---

*Bài viết chỉ mang tính phân tích dữ liệu thống kê, không phải lời khuyên đầu tư. Mọi quyết định giao dịch thuộc trách nhiệm cá nhân.*
