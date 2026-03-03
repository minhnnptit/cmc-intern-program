# 🚀 Buổi 1: Foundation & Project Setup

## Mục Tiêu

- ✅ Setup project structure theo Clean Architecture
- ✅ Hiểu Go basics và HTTP server
- ✅ Chạy được Hello World API

## Nội Dung Code

### 1. Project Structure

```
session1-foundation/
├── cmd/
│   └── server/
│       └── main.go          # Entry point - Hello World server
├── internal/
│   ├── model/               # (Empty - chuẩn bị buổi 2)
│   ├── handler/             # (Empty - chuẩn bị buổi 2)
│   ├── service/             # (Empty - chuẩn bị buổi 2)
│   └── storage/             # (Empty - chuẩn bị buổi 2)
├── go.mod
└── README.md
```

### 2. Key Concepts

#### HTTP Server trong Go

- `http.HandleFunc` - register handlers cho routes
- `http.ListenAndServe` - start server
- `ResponseWriter` và `Request` - handle HTTP

#### JSON Response

- `json.NewEncoder(w).Encode()` - convert Go struct sang JSON
- Content-Type header

### 3. Chạy Code

```bash
# Khởi tạo Go module
go mod init mini-asm

# Chạy server
go run cmd/server/main.go

# Test endpoint (terminal khác)
curl http://localhost:8080/health
```

**Expected Output:**

```json
{
  "status": "ok",
  "message": "Mini ASM service is running"
}
```

## Điểm Chú Ý

### 🎯 Teaching Points

1. **Package Structure**
   - `package main` - entry point của application
   - `package internal` - private code, không export ra ngoài

2. **HTTP Basics**
   - Status codes: 200 (OK), 404 (Not Found), 500 (Error)
   - Headers: Content-Type, etc.
   - Request methods: GET, POST, PUT, DELETE

3. **Why Empty Folders?**
   - Chuẩn bị structure trước
   - Sẽ implement từng layer trong buổi 2
   - Giúp hiểu Clean Architecture progression

### ❓ Questions to Ask Students

1. Tại sao cần folder `internal`? → Private code
2. Tại sao main.go ở trong `cmd/server`? → Multiple entry points trong tương lai
3. Response code 200 nghĩa là gì? → Success
4. Làm sao add thêm endpoint `/hello`? → Practice task

### 🏠 Homework

**Task 1:** Thêm endpoint `/hello/{name}`

- Route: `GET /hello/John`
- Response: `{"message": "Hello, John!"}`
- Hint: Use `r.PathValue("name")` trong Go 1.22+

**Task 2:** Thêm endpoint `/info`

- Response: version, start time của server

**Task 3:** Đọc CLEAN_ARCHITECTURE.MD

- Focus: Section 2 (Clean Architecture Layers)
- Chuẩn bị hiểu tại sao cần 4 layers

## So Sánh

### ❌ Bad Practice (Monolithic)

```go
// Tất cả code trong 1 file
func main() {
    http.HandleFunc("/assets", func(w http.ResponseWriter, r *http.Request) {
        // Parse request
        // Validate
        // Business logic
        // Database query
        // Response
    })
}
```

### ✅ Clean Architecture (Preview buổi 2)

```go
// Separation of concerns
handler → service → storage → model
```

**→ Buổi 1 setup structure, Buổi 2 implement layers!**

## Resources

- [Go HTTP Server Tutorial](https://gobyexample.com/http-servers)
- [Go Modules](https://go.dev/blog/using-go-modules)
- [REST API Best Practices](https://restfulapi.net/)
