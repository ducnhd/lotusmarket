---
title: "Dòng tiền VN30 nghiêng bán hôm nay 13/08 — regime nào đang dẫn dắt?"
date: 2026-08-13
topic: flow
---

## Áp lực bán phủ khắp VN30+HNX30 — lịch sử nói gì tiếp theo?

Phiên hôm nay, dòng tiền nội tổng trên rổ VN30+HNX30 ghi nhận buy-pressure chỉ còn 28.7% — tức gần 3/4 lực giao dịch đang nghiêng về phía bán. Con số này không phải dao động thông thường; nó phản ánh một trạng thái mà thị trường đang ở vùng căng thẳng thực sự. Câu hỏi đặt ra: trong những lần tương tự trước đây, thị trường đã đi về đâu sau 60 ngày?

---

### ① Chuyện gì đang xảy ra

Buy-pressure 28.7% là tín hiệu "bán mạnh" theo phân loại dòng tiền nội. Để hình dung rõ hơn: ngưỡng trung tính thường nằm quanh 50%, tức bên mua và bên bán cân bằng nhau. Khi con số rơi xuống dưới 30%, lực bán đang áp đảo gần gấp đôi lực mua.

Điều này có nghĩa là không chỉ một vài cổ phiếu bị xả — cả hai rổ blue-chip lớn nhất sàn đều chịu áp lực cùng lúc. Đây là kiểu dòng tiền mà giới quant thường gọi là "coordinated selling": không phải hoảng loạn lẻ tẻ, mà là một lực rút khá đồng thuận.

---

### ② Vì sao áp lực bán đang lấn át

Không cần phỏng đoán nguyên nhân vĩ mô cụ thể — data tự phân loại theo **regime** (trạng thái thị trường) thay vì theo tin tức. Và regime hiện tại được nhận diện là **CRISIS**.

Regime CRISIS không nhất thiết có nghĩa là "thị trường sụp đổ". Nó có nghĩa là các chỉ số nội tại — biến động, dòng tiền, breadth — đang ở trạng thái bất ổn đủ để hệ thống phân loại tách nó ra khỏi STABLE. Sự phân loại này có giá trị: nó giúp so sánh đúng cohort lịch sử, thay vì gộp chung mọi phiên bán mạnh vào cùng một nhóm.

Điều đáng chú ý là bán mạnh trong regime CRISIS và bán mạnh trong regime STABLE là hai thứ hoàn toàn khác nhau về mặt forward return — và data 9 năm xác nhận điều này rõ ràng.

---

### ③ Cohort lịch sử đã đi tiếp ra sao

Đây là phần quan trọng nhất.

Hệ thống nhìn lại 9 năm dữ liệu, phân loại từng phiên theo regime, rồi tính **forward return 60 ngày** (tương đương khoảng 3 tháng giao dịch):

| Regime | Fwd Return 60d | Win Rate |
|--------|---------------|----------|
| **CRISIS** | **+8.53%** | **67%** |
| STABLE | +3.20% | 48% |
| Baseline (toàn bộ) | +4.75% | — |

Kết quả này đảo ngược hoàn toàn trực giác thông thường. Những phiên bán mạnh trong regime CRISIS — thay vì dẫn đến suy giảm kéo dài — lại cho forward return trung bình **+8.53%** sau 60 ngày, với win rate **67%**. Tức là trong 2/3 trường hợp lịch sử, thị trường đã hồi phục và tạo lợi nhuận dương trong vòng 3 tháng.

So sánh với regime STABLE: win rate chỉ 48% — thấp hơn cả baseline — và return trung bình +3.20%, dưới mức baseline +4.75%. Điều này cho thấy khi thị trường "bình yên" mà có áp lực bán, tín hiệu đó ít ý nghĩa hơn; còn khi áp lực bán xuất hiện giữa vùng CRISIS, lịch sử cho thấy đây thường là vùng giá mà lực cầu cuối cùng sẽ quay lại mạnh hơn.

Cơ chế đằng sau: trong regime CRISIS, áp lực bán thường đến từ margin call hoặc tâm lý sợ hãi lan rộng — những lực lượng có chu kỳ kết thúc rõ ràng. Một khi lực bán kỹ thuật này cạn kiệt, thị trường thường bật lại nhanh hơn trung bình. 📊

---

### ④ Cần theo dõi gì trong các phiên tới

Ba tín hiệu đáng quan sát để xác nhận hay phủ nhận kịch bản lịch sử:

**Buy-pressure có hồi phục về trên 40% không?** Đây là dấu hiệu đầu tiên cho thấy lực bán đang giảm bớt. Nếu buy-pressure tiếp tục dưới 30% nhiều phiên liên tiếp, cohort lịch sử có thể chưa kích hoạt.

**Regime có chuyển từ CRISIS sang STABLE không?** Sự chuyển đổi regime thường đi kèm với sự cải thiện breadth — tức số lượng cổ phiếu tăng giá trong rổ bắt đầu vượt số giảm giá.

**Độ sâu của đợt bán**: 60 ngày trong cohort CRISIS là khoảng thời gian đủ dài để thị trường hoàn tất một chu kỳ phân phối — tích lũy. Điểm đáng chú ý là win rate 67% không có nghĩa là thị trường đi thẳng lên; nhiều trường hợp trong cohort vẫn có thêm nhịp giảm trước khi hồi phục.

---

### Verify reproducible

Toàn bộ số liệu trên có thể kiểm tra lại tại [https://github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket):

```bash
pip install lotusmarket
```

```python
from lotusmarket import FlowRegime
result = FlowRegime().cohort_analysis(
    universe="VN30+HNX30",
    forward_days=60,
    lookback_years=9
)
print(result)
```

Hoặc dùng CLI:

```bash
lmcli flow-regime --universe VN30+HNX30 --fwd 60 --years 9
```

---

⚠️ *Bài viết chỉ mang tính phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Kết quả quá khứ không đảm bảo kết quả tương lai.*
