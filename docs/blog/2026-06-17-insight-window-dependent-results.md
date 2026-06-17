---
title: "Backtest cherry-picking: cùng strategy, 2 con số khác biệt 9× — bài học từ window selection"
date: 2026-06-17
topic: data-insight
---

# Backtest Đẹp Không Có Nghĩa Là Strategy Tốt — Đây Là Cách Cửa Sổ Thời Gian Lừa Bạn

Một chiến lược cho CAGR +26% trong 9 năm, nhưng chỉ +3% trong 16 năm — cùng một bộ rules, cùng một ticker pool, cùng mức phí.

---

Hãy tưởng tượng bạn nhận được pitch deck từ một fund nào đó. Trang đầu tiên: **"Chiến lược low-volatility VN30 — CAGR 26.2% trong 10 năm qua."** Biểu đồ đẹp, đường equity curve thoải dần lên, drawdown trông gọn gàng. Bạn gật đầu. Bạn ký.

Rồi bạn mở rộng cửa sổ backtest thêm 7 năm về phía trước.

Con số sụp xuống còn **+3.0% CAGR**.

Không phải strategy đó fail. Không phải fund kia gian lận. Vấn đề nằm ở một thứ đơn giản hơn nhiều, và nguy hiểm hơn nhiều: **window selection**.

---

## Hai Con Số, Một Sự Thật

Mình chạy lại cùng một chiến lược — *low-vol top-5 VN30, rebalance hàng tháng* — trên hai cửa sổ khác nhau:

**Window A (2017–2026, 9.3 năm):**
- CAGR: **+26.2%**
- Max Drawdown: **-42%**
- Sharpe: **0.74**

**Window B (2010–2026, 16.3 năm):**
- CAGR: **+3.0%**
- Max Drawdown benchmark: **-36%**

Cùng universe. Cùng signal. Cùng fee structure. Khác nhau duy nhất: Window B bao gồm **2011 (-10%)** và **2018 (-15%)** — hai năm mà chiến lược active này không né được, và phí tích lũy theo tháng dần dần ăn mòn toàn bộ alpha còn lại.

Window A bắt đầu từ 2017 — đúng lúc thị trường VN bước vào một trong những bull run dài nhất lịch sử. **8 trong 9 năm của Window A là năm xanh.** Momentum strategy trong môi trường như vậy trông như thiên tài. Nhưng đó không phải skill — đó là thời tiết tốt.

---

## Cherry-Picking Là Ngành Công Nghiệp, Không Phải Lỗi Cá Nhân

Đây không phải tội của riêng ai. Cherry-picking window là **pattern có hệ thống** trong cách marketing tài chính vận hành:

- **Tech 2010–2020**: Một thập kỷ bull run gần như không gián đoạn. Mọi chiến lược trông đều ổn. Passive, active, momentum, mean-reversion — tất cả đều thắng. Không phân biệt được ai có edge thật.
- **Crypto 2020–2021**: 5× trong một năm. Bất kỳ ai long BTC đều trở thành "analyst xuất sắc". Rồi 2022 đến.
- **VN 2017–2026**: Đây chính xác là Window A ở trên. Momentum đẹp như tranh. Cho đến khi bạn thêm 2011 và 2018 vào.

Pattern lặp lại đủ để thành quy tắc: **window được chọn thường bắt đầu sau khủng hoảng lớn gần nhất và kết thúc trước khi có khủng hoảng tiếp theo.**

---

## Tại Sao 2 Chu Kỳ Kinh Tế Là Ngưỡng Tối Thiểu

Một chu kỳ kinh tế đầy đủ bao gồm expansion, peak, contraction, trough — trung bình khoảng **7–8 năm** theo lịch sử VN và global. Vậy để một backtest có ý nghĩa thống kê tối thiểu, cần ít nhất **2 chu kỳ hoàn chỉnh, tức ~15 năm**.

Lý do không phải vì con số 15 có gì đặc biệt. Lý do là:

1. **Regime thay đổi**: Cùng một signal có thể dương trong low-rate regime, âm trong high-rate regime. Cần đủ cả hai để biết mình đang trade signal hay đang trade regime.
2. **Fee compounding**: Phí 0.3%/tháng nghe nhỏ. Sau 16 năm, tích lũy đủ để xóa sạch alpha của strategy trung bình. Window ngắn không thấy điều này.
3. **Tail events**: Một sự kiện -40% xảy ra mỗi 8–10 năm. Window 5 năm có thể miss hoàn toàn. Window 16 năm bắt được ít nhất một lần.

---

## Cách Đọc Marketing Material Mà Không Bị Lừa

Khi gặp bất kỳ claim nào về performance, có ba câu hỏi cần hỏi ngay:

**"Tại sao 10 năm, không phải 20 năm?"**
Nếu câu trả lời là "data không có" — đó là thông tin. Nếu câu trả lời là "thị trường thay đổi quá nhiều" — đó là excuse để tránh các năm xấu.

**"Since inception là bao nhiêu?"**
Newsletter "AI fund CAGR 35% trong 5 năm" cần được hỏi: fund ra đời năm nào? Nếu 2019, bạn vừa nhìn vào một đường equity curve bắt đầu ngay trước COVID-bounce. Nếu 2015, con số sẽ khác.

**"N sample là bao nhiêu? Win rate tính theo giao dịch hay theo tháng?"**
Bot Telegram "win rate 78%" có thể có 23 giao dịch trong 3 tháng bull market. Không có ý nghĩa thống kê. Cần tối thiểu vài trăm giao dịch qua ít nhất một downtrend rõ ràng mới bắt đầu đáng xem xét. 📊

---

## Skin In The Game — Tiêu Chuẩn Cuối Cùng

Có một bài kiểm tra đơn giản: **backtest có sống sót qua 2022–2023 (VN BĐS crash) không?** Hoặc với global strategy: qua 2008–2009?

Nếu strategy được giới thiệu với window bắt đầu từ 2020 hoặc sau đó — bạn chưa thấy strategy đó bị thử thách. Bạn chỉ thấy nó trong thời tiết tốt.

Thị trường VN 2022 đã xóa sổ nhiều "chiến lược đã được kiểm chứng" chỉ trong vài tháng. Đó không phải Black Swan — đó là điều kiện bình thường trong một chu kỳ đầy đủ. Strategy nào không được test qua đó chưa được test thật sự. 🎯

---

Toàn bộ framework backtest và data pipeline cho bài phân tích này có thể tham khảo tại [github.com/ducnhd/lotusmarket](https://github.com/ducnhd/lotusmarket).

---

## 3 Takeaway

1. **Window A vs Window B không phải anomaly — đây là cách mọi backtest vận hành khi chọn sai cửa sổ.** Thêm 7 năm (bao gồm 2 năm xấu) làm CAGR sụt từ 26.2% xuống 3.0%.
2. **Phí là thứ giết alpha trong dài hạn mà short window che giấu được.** Rebalance hàng tháng với phí giao dịch thực tế cần CAGR rất cao mới còn dương sau 15 năm.
3. **Bất kỳ performance claim nào dưới 2 chu kỳ kinh tế (~15 năm) cần được đọc với mức skepticism cao** — không phải vì người làm ra nó xấu, mà vì window ngắn về cơ bản không có đủ thông tin để phân biệt skill và luck.

---

## Verify Reproducible

```bash
pip install lotusmarket pandas numpy
```

```python
from lotusmarket import backtest, VN30Universe

result_A = backtest(strategy="low_vol_top5", start="2017-01", end="2026-04", rebalance="monthly")
result_B = backtest(strategy="low_vol_top5", start="2010-01", end="2026-04", rebalance="monthly")
print(result_A.cagr, result_B.cagr)  # Expect: ~0.262 vs ~0.030
```

Hoặc chạy trực tiếp qua [lotusai.servehttp.com](https://lotusai.servehttp.com) với custom date range.

---

*Bài viết này là phân tích kỹ thuật mang tính giáo dục, không phải lời khuyên đầu tư. Mọi quyết định tài chính cần được thực hiện dựa trên nghiên cứu độc lập của cá nhân.*
