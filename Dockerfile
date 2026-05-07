# --- BƯỚC 1: Build file thực thi (Builder) ---
FROM golang:1.25-alpine AS builder

# Thiết lập thư mục làm việc bên trong container
WORKDIR /app

# Chỉ copy file quản lý thư viện trước
COPY go.mod ./

# Nếu bạn có file go.sum thì bỏ comment dòng dưới, 
# nhưng nếu đang lỗi thì cứ để yên thế này cho chắc.
# COPY go.sum ./

# Copy toàn bộ mã nguồn vào container
COPY . .

# Tải các thư viện cần thiết và build code Go thành file thực thi tên là "main"
# Lệnh build sẽ tự động tải các package còn thiếu
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# --- BƯỚC 2: Tạo Image chạy cuối cùng (Runtime) ---
FROM alpine:latest

# Cài đặt ca-certificates để ứng dụng có thể gọi API HTTPS (nếu cần)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Chỉ lấy duy nhất file thực thi từ bước builder sang đây
COPY --from=builder /app/main .

# Nếu code của bạn dùng file .env, hãy bỏ dấu # ở dòng dưới
# COPY .env .

# Mở port (8080 là cổng mặc định của bạn)
EXPOSE 8080

# Lệnh khởi chạy ứng dụng
CMD ["./main"]