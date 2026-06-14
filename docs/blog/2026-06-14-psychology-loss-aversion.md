---
title: "Tại sao bạn bán đáy? Loss aversion + cohort data VN"
date: 2026-06-14
topic: psychology
---

# Tại Sao Não Bạn Là Kẻ Thù Tệ Nhất Trong Thị Trường Chứng Khoán

Retail VN thua không phải vì thiếu thông tin — mà vì não bộ được lập trình sai cho việc đầu tư. Và có một con số cụ thể chứng minh điều đó.

---

Hãy thử một bài toán nhỏ: Bạn mất 1 triệu đồng hôm nay. Cần thắng bao nhiêu để cảm thấy "hòa"?

Nếu câu trả lời trong đầu bạn là "2 triệu", chúc mừng — não bạn đang hoạt động hoàn toàn bình thường. Và đó chính xác là vấn đề.

---

## Kahneman Đã Giải Thích Điều Này Từ Năm 1979

Kahneman & Tversky đặt tên cho hiện tượng này là **loss aversion**: nỗi đau tâm lý khi mất 1 đồng tương đương với niềm vui khi được **2 đồng**. Không phải 1 đồng. Không phải 1.5 đồng. Hai đồng.

Nghe có vẻ như kiến thức tâm lý học thuần túy, xa vời với màn hình xanh đỏ của VN-Index. Nhưng khi map behavior này vào dữ liệu giao dịch thực tế, pattern hiện ra rõ đến mức đáng sợ:

> Mua cổ phiếu A giá 30. A xuống 24, tức -20% → **bán panic**.
> Mua cổ phiếu B giá 30. B lên 36, tức +20% → **bán chốt lời**.

Đây không phải quan sát ngẫu nhiên. Đây là hành vi có hệ thống của phần lớn retail Việt Nam. Và hành vi đó tạo ra một chiến lược hoàn toàn **ngược chiều** với những gì backtest cho thấy là đúng.

---

## "Hold Thua, Run Thắng" — Một Paradox Nghe Quen Mà Không Ai Làm Được

Câu nói kinh điển trong giới đầu tư là: *"Let winners run, cut losses early."*

Retail VN làm chính xác điều ngược lại: **giữ thua, chạy thắng**.

Giữ cổ phiếu lỗ vì "lỡ bán thì lỗ thật rồi" — một loại ảo tưởng kế toán não tự tạo ra để tránh đối mặt với thực tế. Bán cổ phiếu lời sớm vì sợ "lãi bay mất" — cũng chính nỗi sợ loss aversion nhưng đội lốt "bảo toàn lợi nhuận".

Kết quả cuối kỳ? Danh mục chỉ còn toàn cổ phiếu xấu.

---

## Data Nói Gì: Panic-Sell Trong Khủng Hoảng Là Sai Như Thế Nào?

Từ backtest cohort nội bộ trên dữ liệu thị trường VN, regime **CRISIS** — vùng mà panic-selling diễn ra mạnh nhất — cho thấy:

- **Forward return 60 ngày: +8.53%**
- **Win rate: 67%**
- **Edge: +3.78%** — cao nhất trong tất cả các regime được test

Nói thẳng ra: đây là vùng mà retail bán tháo, và cũng là vùng mà data cho thấy xác suất hồi phục **2 trong 3 lần**. Người bán đáy chấp nhận thua chắc chắn 100% để tránh rủi ro thua thêm — trong khi lịch sử cho thấy chỉ có 33% khả năng thị trường tiếp tục giảm sau đó.

Đó là một đánh đổi cực kỳ tệ, và loss aversion là thứ khiến người ta chấp nhận nó.

---

## Kiểm Tra Thực Tế: 30 Ngày Audit

Một audit 30 ngày trên các khuyến nghị được log lại cho thấy kết quả cụ thể:

| Loại khuyến nghị | Đúng / Tổng | Tỷ lệ chính xác |
|---|---|---|
| BUY / HOLD | 14 / 15 | **93%** |
| SELL / CUT_LOSS | 0 / 3 | **0%** |

Ba lần khuyến nghị SELL/CUT_LOSS trong giai đoạn này đều rơi vào **vùng oversold** — tức là tín hiệu panic-sell xuất hiện đúng lúc giá đang ở đáy cục bộ. Cả ba lần đó đều sai.

Đây là finding sơ bộ từ 30 ngày dữ liệu, cần thêm thời gian để confirm. Nhưng pattern đủ rõ để đặt câu hỏi: bạn đang bán vì thesis thay đổi, hay vì màn hình chuyển đỏ?

---

## Defeat Loss Aversion: Không Phải Triết Lý, Là Quy Trình

Cảm xúc không cần bị *triệt tiêu* — cần bị *cô lập khỏi quyết định*. Ba kỹ thuật cụ thể:

### 1. Rule-based stop — viết trước, không chỉnh sau
Ví dụ: *"Không bán dưới -25% trừ khi MA200 đảo chiều bearish."*

Khi rule đã có sẵn, lúc panic bạn không cần *quyết định* nữa — chỉ cần *kiểm tra điều kiện*. Não không có chỗ để sợ hãi xen vào.

### 2. Pre-commit: viết thesis khi bình tĩnh, đọc lại khi hoảng loạn
Trước khi mua, viết ra: *"Mình mua cổ phiếu này vì X, Y, Z. Mình sẽ bán khi X hoặc Y không còn đúng."*

Lúc giá giảm 15%, đọc lại tờ giấy đó. Nếu X và Y vẫn đúng — thesis chưa vỡ. Cảm giác "phải bán" chỉ là tiếng ồn.

### 3. Quy tắc 24 giờ 🕐
Nếu bạn cảm thấy **phải bán ngay hôm nay**, đợi 24 tiếng. Không xem giá. 90% trường hợp cảm giác cấp bách đó tan biến. Nếu sau 24h thesis vẫn vỡ — lúc đó mới là quyết định, không phải phản xạ.

---

## Sarcasm Nhẹ Cho Cộng Đồng: "Cắt Lỗ Kỷ Luật" Không Có Nghĩa Là Bán Mỗi Khi Sợ

Một trong những câu thần chú phổ biến nhất trong các group chứng khoán VN là *"cắt lỗ kỷ luật"*. Câu đó không sai. Nhưng khi áp dụng bừa bãi vào mọi đợt giảm — đặc biệt trong vùng CRISIS/oversold — nó trở thành một công cụ để **hợp lý hóa panic**, không phải kỷ luật thực sự.

Kỷ luật thực sự là: có rule rõ ràng, và follow rule đó nhất quán — bao gồm cả rule *không bán* khi điều kiện chưa kích hoạt. 😤

---

## Takeaway

**1. Loss aversion là bug cấu trúc của não, không phải tính cách cá nhân.** Ai cũng có nó. Vấn đề là có hệ thống để bypass nó hay không.

**2. Vùng panic-sell trong lịch sử là vùng edge dương nhất.** Data từ CRISIS regime cho thấy forward 60d trung bình +8.53%, win rate 67%. Bán đáy là đánh đổi xác suất không có lợi.

**3. Quy trình > cảm xúc.** Rule trước khi mua, thesis viết ra, quy tắc 24 giờ — ba công cụ đơn giản nhưng đủ để cô lập quyết định khỏi noise cảm xúc.

---

## Verify Reproducible

Toàn bộ framework regime detection và backtest cohort có thể tái tạo qua [Lotus AI](https://lotusai.servehttp.com) hoặc clone repo tại [github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

```bash
pip install lotusmarket pandas numpy
```

```python
from lotusmarket import RegimeDetector, Backtest

detector = RegimeDetector(ticker="VN30", lookback=252)
regime = detector.classify_current()  # Returns: CRISIS / BULL / BEAR / SIDEWAYS

bt = Backtest(regime_filter="CRISIS", rsi_threshold=30)
bt.run_cohort(fwd_days=60).summary()
# Expected output: fwd_return=+8.53%, win_rate=0.67, edge=+3.78%
```

---

*Bài viết mang tính phân tích dữ liệu lịch sử, không phải lời khuyên đầu tư. Mọi quyết định mua bán là trách nhiệm cá nhân của nhà đầu tư.*
