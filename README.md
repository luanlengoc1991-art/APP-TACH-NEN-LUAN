# LUAN GROUP | AI Layer Studio — Portable Windows

**Tình trạng:** Bản ứng dụng Windows portable V1.0 để chỉnh sửa/tách layer thủ công. Không có AI auto segmentation và không có kết nối lấy hạn mức Plus làm API.

## Mở trên Windows

1. Tải `LUAN-GROUP-AI-Layer-Studio-Windows-x64.exe` về máy.
2. Nhấp đúp mở. Không cần cài đặt, không yêu cầu quyền Administrator.
3. Ứng dụng tự mở cửa sổ Microsoft Edge (hoặc Google Chrome) dạng app trên `http://127.0.0.1:49753` (hoặc cổng dự phòng nếu đã được sử dụng).
4. Chọn **Dùng poster mẫu** hoặc tải ảnh JPG/PNG/WEBP, khoanh các vùng cần tách, hoàn tất từng vùng. Lặp lại tới khi đủ 9 layer.
5. Chọn background, tỷ lệ xuất rồi tải ZIP (9 PNG trong suốt + 1 PNG nền + manifest).
6. Tab **Kết nối ChatGPT** mở trang `chatgpt.com` CHÍNH THỨC trong cửa sổ Edge/Chrome riêng (chế độ app); có prompt sẵn để copy và dán thủ công. Không đọc cookie, mật khẩu, token hoặc nội dung ChatGPT.

## Cảnh báo minh bạch

- App không có API ChatGPT; đăng nhập ChatGPT trong cửa sổ browser chính thức. Không dùng credit ChatGPT Plus làm API ngoài.
- Tách bằng polygon là thao tác thủ công (đường viền thẳng giữa các đỉnh), không tách từng sợi tóc như model AI.
- Background mode 'lấp vùng' dùng phép nội suy màu từ lân cận; không phải semantic AI inpainting; vùng bị xóa lớn có thể mờ.
- 'Nâng nét 2x/4x' là resampling + sharpen kernel, không thêm chi tiết theo AI.
- Hạn chế 24 megapixel/ảnh xuất. Các PNG layer dùng cùng canvas gốc để dễ ghép trở lại.
- App chỉ mở cổng trên 127.0.0.1. Dữ liệu ảnh không upload đi đâu trong Studio; không lưu cloud.
- Code được viết mới, không sao chép source PixVIP.

## Build từ mã nguồn

Yêu cầu Go 1.23 trở lên, không cần Node/npm.

```powershell
$env:GOOS="windows"
$env:GOARCH="amd64"
go build -ldflags="-s -w -H=windowsgui" -o LUAN-GROUP-AI-Layer-Studio-Windows-x64.exe .
```

Để test server trên máy Linux/macOS:

```bash
go run . -no-browser
```

Trong tab được mở, ảnh và mask lưu IndexedDB trên chính trình duyệt profile đang dùng. Chuyển profile, xoá dữ liệu trang, dùng ẩn danh hoặc app phải chạy bằng cổng dự phòng có thể khiến phiên cũ không hiển thị trong lần mở sau. Dữ liệu không nằm trong gói `.exe`.

## Bảo mật đăng nhập

Website/app không thể tùy tiện dùng ChatGPT credit làm API. Một số ứng dụng đối tác có cơ chế xác thực chính thức; bản này không có quyền được cấp. Không nhập mật khẩu ChatGPT vào form lạ; luôn kiểm tra domain `chatgpt.com` và đăng nhập trên trang OpenAI.

## GitHub là dự án mã nguồn chính

Mọi thay đổi được quản lý trên nhánh `main` của repository này.

- Mã nguồn chính: `main.go`, `go.mod`, `web/index.html`, `web/app.js`.
- Ảnh mẫu: `web/poster-mau.jpg` (là ảnh demo, không phải ảnh mặc định lưu trên server).
- GitHub Actions: `.github/workflows/windows-build.yml` tự kiểm tra và tạo file Windows `.exe` khi có commit lên `main`.
- Vào tab **Actions** → chọn lượt chạy **Build LUAN GROUP Windows App** → tải **Artifacts** để lấy bản Windows mới nhất (artifact được lưu tạm 30 ngày).
- Trong máy tính: cài Go 1.23+, clone repo, sửa `web/app.js`/`web/index.html`, chạy `go run .` rồi push lên `main`.
- Không commit mật khẩu, token ChatGPT, API key, cookie, dữ liệu cá nhân hoặc binary `.exe` vào repository công khai.

**Lưu ý:** thay đổi GitHub **không tự cập nhật file `.exe` đã tải về máy**. Cần tải bản build mới từ GitHub Actions để sử dụng cập nhật.

### Cách đưa source lên GitHub nếu công cụ kết nối không có quyền ghi

Repository: https://github.com/luanlengoc1991-art/APP-TACH-NEN-LUAN

1. Tải ZIP mã nguồn và giải nén trên Windows.
2. Cài [Git for Windows](https://git-scm.com/downloads/win) nếu chưa có.
3. Chạy `UPLOAD-GITHUB.cmd` trong thư mục giải nén, đăng nhập GitHub trong trình duyệt nếu Git yêu cầu.
4. Mở repository và kiểm tra thư mục `web/` và file `main.go` đã xuất hiện.
5. Vào **Actions** để xem bản Windows `.exe` vừa được GitHub tự động biên dịch.

Trong trường hợp Git yêu cầu thông tin người tạo commit, cấu hình Git bằng `git config --global user.name "Ten cua ban"` và `git config --global user.email "email-duoc-xac-minh-tren-GitHub"`, rồi chạy lại lệnh upload. Không bao giờ nhập mật khẩu hay mã token vào file nguồn.
