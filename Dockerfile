# --- BƯỚC 1: Build file thực thi ---
FROM golang:1.22-alpine AS builder

# Thiết lập thư mục làm việc bên trong container
WORKDIR /app

# Copy các file quản lý thư viện trước để tận dụng cache của Docker
COPY go.mod go.sum ./
RUN go mod download

# Copy toàn bộ mã nguồn vào container
COPY . .

# Build code Go thành file thực thi tên là "app"
# CGO_ENABLED=0 giúp chạy tốt trên các hệ điều hành Linux nhẹ như Alpine
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/main .

# --- BƯỚC 2: Tạo Image chạy cuối cùng ---
FROM alpine:latest

# Cài đặt ca-certificates để có thể gọi các API HTTPS (nếu cần)
RUN apk --no-cache add ca-certificates

WORKDIR /root/

# Chỉ lấy duy nhất file thực thi từ bước builder qua đây
COPY --from=builder /app/main .

# Nếu bạn có file .env, hãy bỏ comment dòng dưới đây
# COPY .env .

# Mở port (thay 8080 bằng port mà code Go của bạn đang lắng nghe)
EXPOSE 8080

# Lệnh chạy ứng dụng
CMD ["./main"]