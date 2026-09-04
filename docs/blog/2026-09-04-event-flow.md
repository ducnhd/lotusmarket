---
title: "Dòng tiền VN30 nghiêng bán hôm nay 04/09 — regime nào đang dẫn dắt?"
date: 2026-09-04
topic: flow
---

## Áp lực bán đang chi phối VN30+HNX30 — lịch sử nói gì tiếp theo?

Phiên hôm nay, dòng tiền nội tổng của hai rổ cổ phiếu lớn nhất thị trường — VN30 và HNX30 — ghi nhận buy-pressure chỉ còn **26.5%**, một ngưỡng thuộc diện "bán mạnh" theo phân loại định lượng. Câu hỏi mà data 9 năm lịch sử có thể trả lời: khi thị trường rơi vào trạng thái này, điều gì thường xảy ra trong 60 ngày tới?

---

### ① Chuyện gì đang xảy ra

Buy-pressure 26.5% có nghĩa là cứ 100 đồng khớp lệnh trên hai rổ VN30 và HNX30, chỉ có khoảng 26-27 đồng đến từ lực mua chủ động — phần còn lại là lực bán đang áp đảo. Đây không phải trạng thái "dao động nhẹ" hay "điều chỉnh kỹ thuật thông thường". Về mặt định lượng, 26.5% nằm sâu trong vùng bán mạnh, cho thấy tâm lý phân phối đang chiếm ưu thế rõ rệt so với tích lũy.

Điều đáng chú ý là tín hiệu này xuất hiện đồng thời trên cả hai sàn (HOSE và HNX) — không phải hiện tượng cục bộ ở một nhóm cổ phiếu hay một sàn riêng lẻ. Tính đồng thuận của áp lực bán trên diện rộng là yếu tố mà mô hình regime detection của Lotus AI đánh dấu là đáng theo dõi.

---

### ② Bối cảnh regime: CRISIS hay STABLE?

Mô hình phân loại regime 9 năm chia lịch sử giao dịch thành hai trạng thái chính: **CRISIS** (khủng hoảng/biến động cao) và **STABLE** (ổn định/xu hướng bình thường).

Phiên hôm nay, buy-pressure 26.5% khớp với đặc trưng của các giai đoạn **CRISIS** trong cohort lịch sử — những thời điểm mà dòng tiền co cụm, lực mua thụ động, và sàn thường xuyên thấy thanh khoản phân kỳ bất thường. Đây không nhất thiết là khủng hoảng theo nghĩa vĩ mô, mà là một *regime trạng thái* mà mô hình nhận diện dựa trên cấu trúc dòng tiền nội tổng.

Bối cảnh vĩ mô toàn cầu hiện tại — lãi suất ở các thị trường phát triển vẫn ở mức cao, dòng vốn ngoại có xu hướng phân tán, và tâm lý risk-off xuất hiện rải rác — có thể là yếu tố amplifying cho trạng thái này, dù mô hình không định giá thẳng vào các biến vĩ mô mà chỉ đọc từ cấu trúc dòng tiền thực.

---

### ③ Cohort lịch sử 9 năm đã đi tiếp ra sao

Đây là phần data nói thẳng nhất.

Khi nhìn vào toàn bộ 9 năm dữ liệu và phân theo regime:

| Regime | Forward 60d (median) | Win Rate |
|--------|----------------------|----------|
| **CRISIS** | **+8.53%** | **67%** |
| STABLE | +3.20% | 48% |
| Baseline (toàn thị trường) | +4.75% | — |

Kết quả có phần counterintuitive: các phiên thuộc regime **CRISIS** — tức là những lúc buy-pressure thấp, thị trường đang bán mạnh như hôm nay — lại ghi nhận forward return 60 ngày **cao hơn hẳn** so với giai đoạn STABLE. Median +8.53% và win rate 67% so với STABLE chỉ đạt +3.20% với win rate dưới baseline (48% so với benchmark +4.75%).

Lý giải không phải là "mua đáy kiếm lời" — mà là khi lực bán đã áp đảo đến mức này, phần lớn áp lực phân phối đã *đã được hấp thụ*. Các phiên CRISIS trong lịch sử thường đánh dấu vùng mà smart money bắt đầu tích lũy lại ở nền giá thấp hơn, kéo forward return lên cao trong 60 ngày sau đó.

Quan trọng: con số 67% win rate không có nghĩa là chắc chắn. 33% còn lại vẫn là các trường hợp thị trường tiếp tục đi xuống hoặc đi ngang sau regime CRISIS.

---

### ④ Cần theo dõi gì từ đây

Ba điểm mà data cho thấy cần quan sát trong những phiên tới:

**Buy-pressure recovery**: Nếu buy-pressure tăng dần trở lại từ vùng 26-30% lên trên 40%, đó là tín hiệu regime đang chuyển dịch. Ngược lại, nếu tiếp tục dưới 30%, cohort CRISIS kéo dài có thể đẩy biến động tăng thêm trước khi hồi phục.

**Phân kỳ giữa VN30 và HNX30**: Tín hiệu hôm nay là đồng thuận trên cả hai rổ — nếu một trong hai bắt đầu tách ra với buy-pressure cao hơn, điều đó có thể báo hiệu dòng tiền đang tập trung chọn lọc vào nhóm vốn hóa cụ thể.

**Thanh khoản tổng**: Regime CRISIS đi kèm thanh khoản co lại thường ít nguy hiểm hơn CRISIS kèm thanh khoản bùng nổ (panic sell). Cần theo dõi khối lượng khớp lệnh tổng so với trung bình 20 phiên.

---

### Verify reproducible

Toàn bộ phân tích trên có thể tái hiện qua [lotusmarket](https://github.com/ducnhd/lotusmarket):

```python
pip install lotusmarket

from lotusmarket import FlowAnalyzer
fa = FlowAnalyzer(tickers="VN30+HNX30", lookback=9*252)
fa.regime_cohort(forward_days=60).summary()
```

---

⚠️ *Bài viết chỉ mang tính phân tích định lượng từ dữ liệu lịch sử — không phải lời khuyên đầu tư. Mọi quyết định giao dịch thuộc về trách nhiệm cá nhân.*
