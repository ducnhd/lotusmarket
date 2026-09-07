---
title: "\"Sell in May\" trên VN30 — câu thần chú từ Mỹ có áp dụng được?"
date: 2026-09-07
topic: myth-buster
---

# "Sell in May and Go Away" — Câu thần chú Wall Street có chạy được trên VN-Index không?

Câu trả lời ngắn: trên dữ liệu VN 2010–2025, edge này gần như không tồn tại sau khi trừ phí — nhưng lý do *tại sao* nó không hoạt động mới là phần thú vị.

Cứ mỗi đầu tháng 5, một vòng lặp quen thuộc bắt đầu trên các hội nhóm chứng khoán: ai đó post ảnh câu quote tiếng Anh, vài người gật đầu "đúng rồi anh ơi tháng 5 xuống là chắc", rồi cả thread chìm vào confirmation bias. Không ai hỏi: dữ liệu VN nói gì?

---

## Nguồn gốc: một câu ngạn ngữ từ đất nước khác

"Sell in May and go away" sinh ra từ thị trường Anh và Mỹ — quan sát thực nghiệm rằng giai đoạn tháng 5 đến tháng 10 lịch sử yếu hơn giai đoạn tháng 11 đến tháng 4. Ở Mỹ, seasonality này có một phần cơ sở: dòng tiền tổ chức dịch chuyển theo chu kỳ ngân sách, mùa nghỉ hè giảm thanh khoản, v.v.

Vấn đề: VN-Index không phải S&P 500. Cấu trúc thị trường khác, chu kỳ dòng tiền khác, và quan trọng nhất — **regime thị trường VN không map 1:1 vào lịch dương**.

---

## Thực tế trên VN data: regime quan trọng hơn tháng

Khi nhìn vào forward return 60 ngày theo regime thị trường trên dữ liệu VN, một pattern rõ hơn nhiều so với "tháng mấy" hiện ra:

**STABLE regime** (chiếm 52% số cohort quan sát):
- Forward return 60d trung bình: **+3.20%**
- Win rate: **48%**
- Edge: **-1.55%** (âm — tức là sau điểm vào ở STABLE, return kỳ vọng thực ra yếu)

**VOLATILE regime**:
- Forward return 60d trung bình: **+6.34%**
- Win rate: **56%**

Đây là điểm mấu chốt. Tháng 5–6 trên VN thường rơi vào VOLATILE hoặc STABLE. Tháng 11–12 thường là đuôi của uptrend hoặc EUPHORIA. Nếu bán tháng 5 và mua lại tháng 11, bạn đang **chủ động tránh VOLATILE** (nơi return 60d đạt +6.34%) và **chui vào giai đoạn cuối chu kỳ** (nơi định giá đã cao hơn).

Logic của câu thần chú này, áp dụng nguyên xi vào VN, đang đi ngược chiều dữ liệu.

---

## Bài toán phí: kẻ thù thầm lặng của mọi chiến lược rotate

Giả sử edge có tồn tại — giả sử thôi — thì vẫn còn một bộ lọc nữa: chi phí thực tế.

Mỗi lần thực hiện chiến lược "bán tháng 5, mua tháng 11":
- **Phí giao dịch 2 lần switch/năm**: ~0.7%
- **Thuế bán 0.1%** mỗi lần thực hiện lệnh bán

Tổng chi phí mỗi năm: xấp xỉ **0.8%+**. Với một chiến lược mà edge đo được (nếu có) chỉ ở mức vài phần trăm, 0.8% là con số đủ để xóa toàn bộ alpha.

Và đây là phần dữ liệu hiện tại chưa khẳng định được edge dương rõ ràng — tức là bạn đang trả phí cho một thứ mà ngay cả trước phí cũng chưa chắc đã có lời.

---

## Rủi ro bỏ lỡ: EUPHORIA không báo trước 📅

Một trong những kịch bản tệ nhất của chiến lược rotate theo lịch: bán xong thì thị trường vào EUPHORIA giữa năm.

VN-Index có những đợt tăng mạnh nhất *không* theo lịch cố định. Các giai đoạn bùng nổ 2017, 2021 đều có những chuỗi tháng tăng phi tuyến mà một nhà đầu tư "bán tháng 5" sẽ hoàn toàn đứng ngoài nhìn. Cơ hội cost (opportunity cost) của việc ngồi tiền mặt trong EUPHORIA mid-year có thể lớn hơn bất kỳ downside nào mà chiến lược này cố tránh.

---

## Vậy tháng 5 VN-Index thực sự như thế nào?

Tổng hợp quan sát return tháng 5 trên VN-Index 2010–2025 cho thấy: **không có pattern tháng 5 xuống đủ nhất quán để gọi là seasonality**. Một số năm tháng 5 giảm, một số năm tháng 5 tăng — phân phối này trông giống noise hơn là signal có thể khai thác.

Điều này phù hợp với finding regime ở trên: thị trường VN vận động theo **chu kỳ thanh khoản và tâm lý nội địa**, không theo lịch dương phương Tây.

*Caveat rõ ràng*: đây là finding sơ bộ từ quan sát dữ liệu 16 năm. Cần backtest đầy đủ với walk-forward validation để kết luận chắc chắn hơn. Dữ liệu hiện tại đủ để đặt dấu hỏi lớn, chưa đủ để đóng cửa hoàn toàn vấn đề.

---

## Tại sao myth này vẫn sống?

Vì nó *đôi khi đúng*. Năm nào thị trường tháng 5 giảm, những người đã bán sẽ kể chuyện. Năm nào tháng 5 tăng, không ai nhắc lại. Đây là survivorship bias trong narrative, không phải edge trong dữ liệu.

Thị trường chứng khoán đặc biệt nguy hiểm ở điểm này: bất kỳ chiến lược nào cũng sẽ có những năm đúng. Pattern detector trong não người rất giỏi tìm signal trong noise — đó là tính năng sinh tồn, nhưng là bug trong trading. 🎯

---

## 3 Takeaway

**1. Regime > Calendar.** Trên VN data, VOLATILE regime cho forward return 60d đạt +6.34% với win rate 56% — đây là signal có ý nghĩa hơn nhiều so với "tháng mấy trong năm".

**2. Phí giết edge nhỏ.** Chiến lược rotate 2 lần/năm tốn ~0.8%+ phí và thuế. Với edge chưa được xác nhận rõ ràng, con số này đủ để làm chiến lược âm kỳ vọng.

**3. Myth sống nhờ selective memory.** "Sell in May" đúng một số năm, sai một số năm — đây là đặc điểm của noise, không phải edge. Backtest đầy đủ là bộ lọc duy nhất đáng tin.

---

## Verify Reproducible

Muốn tự kiểm tra với data VN-Index? Lotus Market cung cấp pipeline phân tích mở tại [github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

```bash
pip install lotusmarket pandas numpy
```

```python
from lotusmarket import VNIndex
import pandas as pd

# Load VN-Index data 2010-2025
df = VNIndex.load(start="2010-01-01", end="2025-12-31")

# Tính monthly return và isolate tháng 5
df['month'] = df.index.month
df['monthly_return'] = df['close'].pct_change(21)  # ~21 trading days
may_returns = df[df['month'] == 5]['monthly_return'].dropna()

print(f"Tháng 5 mean return: {may_returns.mean():.2%}")
print(f"Win rate tháng 5: {(may_returns > 0).mean():.2%}")
```

Hoặc khám phá thêm regime-based analysis tại [lotusai.servehttp.com](https://lotusai.servehttp.com).

---

*Bài viết mang tính phân tích quan sát, không phải lời khuyên đầu tư. Mọi quyết định giao dịch cần dựa trên nghiên cứu độc lập và khẩu vị rủi ro cá nhân.*
