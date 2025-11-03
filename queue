# Queue Client - Wallet Queue Client

## 📋 Mục Lục

1. [Tổng Quan](#tổng-quan)
2. [Kiến Trúc](#kiến-trúc)
3. [Các Components](#các-components)
4. [Cách Hoạt Động](#cách-hoạt-động)
5. [Cách Sử Dụng](#cách-sử-dụng)
6. [Flow Xử Lý](#flow-xử-lý)
7. [Cấu Hình](#cấu-hình)
8. [Tính Năng](#tính-năng)
9. [Best Practices](#best-practices)
10. [Troubleshooting](#troubleshooting)

---

## 🎯 Tổng Quan

File `client.go` là **Queue Client** chuyên dụng cho wallet module, cung cấp typed API để enqueue credit và debit jobs vào queue system.

### Mục Đích

- ✅ Cung cấp **typed API** cho wallet operations (không phải generic)
- ✅ **Abstraction layer** giữa handlers và queue implementation
- ✅ Hỗ trợ **multiple queue providers** (Redis/Asynq, RabbitMQ, Kafka, NATS)
- ✅ Tự động serialize payload và xử lý errors
- ✅ Xử lý **idempotency** và duplicate tasks

### Vai Trò Trong Hệ Thống

```
HTTP Handler → Queue Client → Asynq Queue → Redis → Worker → Business Logic
```

---

## 🏗️ Kiến Trúc

### Design Pattern

Sử dụng **Strategy Pattern** với interface abstraction:

```go
Queue Interface (contract)
    ↓
Client (delegator)
    ↓
asynqQueue (implementation) ──→ Asynq Library ──→ Redis
unimplementedQueue (placeholders)
```

### Package Structure

```
infra/queue/client/
└── client.go          # Queue client implementation
```

### Dependencies

- **Asynq**: `github.com/hibiken/asynq` - Redis-based queue library
- **Config**: `ludo-be/infra/config` - Configuration management
- **Queue constants**: `ludo-be/infra/queue` - Job types, queue names

---

## 🔧 Các Components

### 1. Queue Interface

```go
type Queue interface {
    Close() error
    EnqueueCredit(ctx context.Context, payload interface{}) error
    EnqueueDebit(ctx context.Context, payload interface{}) error
}
```

**Mục đích:** Định nghĩa contract cho tất cả queue implementations

### 2. Client Struct

```go
type Client struct {
    queue Queue  // Concrete implementation
}
```

**Mục đích:** Wrapper để delegate calls đến concrete queue implementation

**Public Methods:**
- `NewClient()` - Factory function
- `EnqueueCredit()` - Enqueue credit job
- `EnqueueDebit()` - Enqueue debit job
- `Close()` - Cleanup resources

### 3. asynqQueue (Implementation)

```go
type asynqQueue struct {
    client *asynq.Client  // Asynq client
}
```

**Mục đích:** Redis/Asynq implementation của Queue interface

**Features:**
- Tạo Asynq client với Redis connection
- Serialize payload thành JSON
- Enqueue với options (retry, timeout, queue name)
- Xử lý duplicate tasks gracefully

### 4. unimplementedQueue (Placeholder)

```go
type unimplementedQueue struct {
    name string  // Provider name
}
```

**Mục đích:** Placeholder cho các queue providers chưa implement (RabbitMQ, Kafka, NATS)

---

## ⚙️ Cách Hoạt Động

### 1. Khởi Tạo Client

```go
// Trong wire.go
func provideQueueClient() (*qclient.Client, error) {
    return qclient.NewClient()
}
```

**Process:**
1. Đọc `QUEUE_PROVIDER` env var (default: "redis")
2. Chọn implementation dựa trên provider
3. Tạo Asynq client với Redis config
4. Trả về `*Client` wrapper

**Code Flow:**
```go
NewClient()
  ↓
os.Getenv("QUEUE_PROVIDER")  // "redis" (default)
  ↓
switch provider {
  case "redis":
    newasynqQueue()  // ← Tạo Asynq client
      ↓
    asynq.NewClient(redisAddr, password, db)
      ↓
    return &Client{queue: asynqQueue}
}
```

### 2. Enqueue Credit Job

```go
queueClient.EnqueueCredit(ctx, payload)
```

**Internal Process:**
1. Client delegates đến `queue.EnqueueCredit()`
2. `asynqQueue` tạo Asynq task:
   ```go
   task := asynq.NewTask("wallet:credit", serializedPayload)
   ```
3. Set options:
   - `MaxRetry(3)` - Retry tối đa 3 lần
   - `Timeout(30s)` - Timeout 30 giây
   - `Queue("wallet-credit")` - Queue name
4. Enqueue vào Redis qua Asynq
5. Xử lý duplicate errors (return nil, không fail)

### 3. Serialization

```go
serializePayload(payload interface{}) []byte
```

- Marshal payload thành JSON
- Return error payload nếu marshal fail
- Đảm bảo payload luôn là valid JSON

---

## 📖 Cách Sử Dụng

### 1. Khởi Tạo (Qua Wire)

```go
// Trong cmd/api/wire.go
func provideQueueClient() (*qclient.Client, error) {
    return qclient.NewClient()
}

// Wire tự động inject vào WalletHandler
func provideWalletHandler(
    // ...
    queueClient *qclient.Client,  // ← Injected here
    // ...
) *walletHandler.WalletHandler {
    return walletHandler.NewWalletHandler(..., queueClient, ...)
}
```

### 2. Sử Dụng Trong Handler

```go
// internal/modules/wallet/interfaces/http/handler.go
func (h *WalletHandler) Credit(c *gin.Context) {
    // 1. Parse request
    var dto dtos.CreditBalanceDTO
    c.ShouldBindJSON(&dto)
    
    // 2. Prepare payload
    payload := queue.CreditJobPayload{
        UserID:      dto.UserID,
        Amount:      dto.Amount,
        ReferenceID: dto.ReferenceID,
        Metadata:    dto.Metadata,
    }
    
    // 3. Enqueue job
    if err := h.queueClient.EnqueueCredit(c.Request.Context(), payload); err != nil {
        httpx.GinAbortWithError(c, err)
        return
    }
    
    // 4. Return 202 Accepted
    httpx.GinJSON(c, map[string]string{"status": "enqueued"}, http.StatusAccepted)
}
```

### 3. Payload Structure

#### CreditJobPayload
```go
type CreditJobPayload struct {
    UserID      string                 `json:"user_id"`
    Amount      int64                  `json:"amount"`
    ReferenceID *string                `json:"reference_id,omitempty"`  // Optional: cho idempotency
    Metadata    map[string]interface{} `json:"metadata,omitempty"`      // Optional: thông tin bổ sung
}
```

#### DebitJobPayload
```go
type DebitJobPayload struct {
    UserID      string                 `json:"user_id"`
    Amount      int64                  `json:"amount"`
    ReferenceID *string                `json:"reference_id,omitempty"`
    Metadata    map[string]interface{} `json:"metadata,omitempty"`
}
```

### 4. Error Handling

```go
err := queueClient.EnqueueCredit(ctx, payload)
if err != nil {
    // Có thể là:
    // - Redis connection error
    // - Serialization error (hiếm)
    // - Unknown provider error
    return err
}

// Success: Job đã được enqueue vào Redis
// Worker sẽ xử lý bất đồng bộ
```

---

## 🔄 Flow Xử Lý Hoàn Chỉnh

### Credit Flow

```
1. Client gọi API: POST /wallet/credit
   ↓
2. WalletHandler.Credit() nhận request
   ↓
3. Validate và parse DTO
   ↓
4. Tạo CreditJobPayload
   ↓
5. queueClient.EnqueueCredit(ctx, payload)
   ↓
6. asynqQueue.EnqueueCredit()
   - Serialize payload → JSON
   - Tạo asynq.Task với type "wallet:credit"
   - Set options (retry, timeout, queue)
   ↓
7. asynq.Client.Enqueue(task, opts)
   - Gửi job vào Redis queue "wallet-credit"
   ↓
8. API trả về 202 Accepted ngay lập tức
   ↓
9. Worker (infra/queue/workers) nhận job từ Redis
   ↓
10. Worker xử lý:
    - Unmarshal payload
    - Gọi CreditBalanceService.Execute()
    - Update balance trong transaction
    - Tạo transaction record
```

### Debit Flow (tương tự)

```
POST /wallet/debit
  ↓
queueClient.EnqueueDebit()
  ↓
Queue "wallet-debit"
  ↓
Worker xử lý DebitBalanceService
```

---

## ⚙️ Cấu Hình

### Environment Variables

```bash
# Queue Provider (optional, default: "redis")
QUEUE_PROVIDER=redis  # Options: redis, rabbitmq, kafka, nats

# Redis Config (cho Asynq)
REDIS_ADDRESS=localhost:6379
REDIS_PASSWORD=
REDIS_DATABASE=0
```

### Queue Configuration

Mỗi job được enqueue với các options:

```go
opts := []asynq.Option{
    asynq.MaxRetry(3),           // Retry tối đa 3 lần
    asynq.Timeout(30 * time.Second),  // Timeout 30 giây
    asynq.Queue(queueName),     // Queue name
}
```

**Queue Names:**
- `wallet-credit` - Cho credit operations
- `wallet-debit` - Cho debit operations

**Job Types:**
- `wallet:credit` - Credit job type
- `wallet:debit` - Debit job type

---

## ✨ Tính Năng

### 1. Type-Safe API

Thay vì generic `Enqueue(event, payload)`, có:
- `EnqueueCredit(payload)` - Rõ ràng mục đích
- `EnqueueDebit(payload)` - Type-safe

**Lợi ích:**
- ✅ Compile-time safety
- ✅ IDE autocomplete
- ✅ Dễ đọc và maintain

### 2. Duplicate Task Handling

```go
if errors.Is(err, asynq.ErrDuplicateTask) || errors.Is(err, asynq.ErrTaskIDConflict) {
    return nil  // Ignore duplicate, không fail
}
```

**Mục đích:** Tránh fail khi client gửi duplicate requests (network retry, etc.)

### 3. Provider Abstraction

Hỗ trợ nhiều queue providers:
- ✅ **Redis/Asynq** (implemented)
- ⏳ RabbitMQ (placeholder)
- ⏳ Kafka (placeholder)
- ⏳ NATS (placeholder)

Dễ dàng switch provider qua env var.

### 4. Automatic Serialization

Payload tự động được serialize thành JSON:
```go
serializePayload(payload) → []byte  // JSON format
```

**Error Handling:**
- Nếu marshal fail → return error payload
- Worker sẽ detect và reject job

### 5. Resource Management

```go
client.Close()  // Đóng Asynq client connection
```

**Lưu ý:** Hiện tại chưa được gọi trong graceful shutdown, nên thêm vào.

---

## 💡 Best Practices

### 1. Luôn Kiểm Tra Error

```go
if err := queueClient.EnqueueCredit(ctx, payload); err != nil {
    // Handle error:
    // - Log error
    // - Return 500 hoặc 503
    // - Không retry (để tránh duplicate)
    return err
}
```

### 2. Sử Dụng ReferenceID Cho Idempotency

```go
payload := queue.CreditJobPayload{
    UserID:      userID,
    Amount:      10000,
    ReferenceID: &referenceID,  // ← Dùng để đảm bảo idempotency
    Metadata:    metadata,
}
```

**Lợi ích:**
- Tránh duplicate processing
- Có thể track job qua reference ID

### 3. Context Support

```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

err := queueClient.EnqueueCredit(ctx, payload)
```

**Lưu ý:** Hiện tại client chưa hỗ trợ context timeout cho enqueue (Asynq có hỗ trợ nhưng chưa dùng).

### 4. Metadata Cho Audit

```go
payload := queue.CreditJobPayload{
    // ...
    Metadata: map[string]interface{}{
        "source": "api",
        "ip_address": c.ClientIP(),
        "request_id": requestID,
        "timestamp": time.Now().Unix(),
    },
}
```

### 5. Close Client Khi Shutdown

```go
// Trong main.go graceful shutdown
defer func() {
    if app.WalletHandler != nil {
        app.WalletHandler.CloseQueueClient()  // Cần thêm method này
    }
}()
```

---

## 🔍 Troubleshooting

### Issue 1: Job Không Được Xử Lý

**Triệu chứng:**
- API trả về 202 Accepted
- Nhưng worker không xử lý

**Nguyên nhân:**
- Worker không chạy
- Redis connection issues
- Queue name không khớp

**Giải pháp:**
```bash
# 1. Kiểm tra worker
ps aux | grep worker

# 2. Kiểm tra Redis
redis-cli ping

# 3. Kiểm tra jobs trong queue
redis-cli
> LLEN asynq:wallet-credit
> LLEN asynq:wallet-debit

# 4. Xem Asynq monitor
make monitor
# Truy cập: http://localhost:9090/monitor
```

### Issue 2: Duplicate Tasks

**Triệu chứng:**
- Nhiều transaction records giống nhau

**Nguyên nhân:**
- Client retry sau timeout
- Không dùng ReferenceID

**Giải pháp:**
```go
// Luôn dùng ReferenceID
referenceID := uuid.New().String()
payload.ReferenceID = &referenceID
```

### Issue 3: Serialization Error

**Triệu chứng:**
- Worker nhận payload với `{"error":"payload_marshal_failed"}`

**Nguyên nhân:**
- Payload có circular reference
- Payload có type không marshal được

**Giải pháp:**
- Kiểm tra payload structure
- Đảm bảo tất cả fields có JSON tags
- Test với `json.Marshal()` trước

### Issue 4: Timeout Issues

**Triệu chứng:**
- Jobs bị timeout trước khi xử lý xong

**Giải pháp:**
```go
// Có thể tăng timeout trong code (hiện tại 30s)
opts := []asynq.Option{
    asynq.Timeout(60 * time.Second),  // Tăng lên 60s
}
```

### Issue 5: Provider Not Implemented

**Triệu chứng:**
- Error: "queue provider X not implemented"

**Nguyên nhân:**
- Set `QUEUE_PROVIDER` sang RabbitMQ/Kafka/NATS
- Nhưng chưa implement

**Giải pháp:**
```bash
# Sử dụng Redis (default)
unset QUEUE_PROVIDER
# hoặc
export QUEUE_PROVIDER=redis
```

---

## 🔗 Liên Kết Với Components Khác

### Dependency Flow

```
cmd/api/wire.go
  └─→ provideQueueClient()
       └─→ queue/client.NewClient()
            └─→ asynqQueue (implementation)
                 └─→ Asynq Client
                      └─→ Redis

internal/modules/wallet/interfaces/http/handler.go
  └─→ WalletHandler
       └─→ queueClient.EnqueueCredit/Debit()
            └─→ (same flow as above)

infra/queue/workers/wallet.worker.go
  └─→ Worker listens to:
       - "wallet:credit" jobs
       - "wallet:debit" jobs
       └─→ Process via services
```

### Files Liên Quan

- **`infra/queue/queue.constants.go`**: Job types, queue names, payload structs
- **`infra/queue/workers/wallet.worker.go`**: Worker xử lý jobs
- **`cmd/api/wire.go`**: Dependency injection setup
- **`internal/modules/wallet/interfaces/http/handler.go`**: Sử dụng client

---

## 📝 Ví Dụ Code

### Example 1: Enqueue Credit với Metadata

```go
payload := queue.CreditJobPayload{
    UserID:      "user-123",
    Amount:      50000,
    ReferenceID: &referenceID,
    Metadata: map[string]interface{}{
        "source": "api",
        "payment_method": "bank_transfer",
        "bank_name": "Vietcombank",
        "transaction_id": "txn-abc-123",
    },
}

if err := queueClient.EnqueueCredit(ctx, payload); err != nil {
    log.Printf("Failed to enqueue credit: %v", err)
    return err
}

log.Printf("Credit job enqueued successfully")
```

### Example 2: Enqueue Debit với Error Handling

```go
payload := queue.DebitJobPayload{
    UserID:      userID,
    Amount:      amount,
    ReferenceID: &withdrawalID,
    Metadata: map[string]interface{}{
        "purpose": "withdrawal",
        "bank_account": "xxx-xxx-xxx",
    },
}

if err := queueClient.EnqueueDebit(ctx, payload); err != nil {
    // Check error type
    if strings.Contains(err.Error(), "connection") {
        return exception.ServiceUnavailable("Queue service unavailable")
    }
    return exception.InternalServerError("Failed to enqueue job")
}

// Success - job enqueued
```

### Example 3: Testing với Mock

```go
// In tests
type mockQueue struct {
    creditCalls []CreditJobPayload
    debitCalls  []DebitJobPayload
}

func (m *mockQueue) EnqueueCredit(ctx context.Context, payload interface{}) error {
    p := payload.(queue.CreditJobPayload)
    m.creditCalls = append(m.creditCalls, p)
    return nil
}

func (m *mockQueue) EnqueueDebit(ctx context.Context, payload interface{}) error {
    p := payload.(queue.DebitJobPayload)
    m.debitCalls = append(m.debitCalls, p)
    return nil
}

// Use in test
mockQueue := &mockQueue{}
client := &queueclient.Client{queue: mockQueue}
```

---

## 🆚 So Sánh Với Alternatives

### Option 1: Dùng Generic Queue Interface

**Hiện tại (Typed):**
```go
queueClient.EnqueueCredit(ctx, payload)  // Type-safe
```

**Generic (Alternative):**
```go
queuePublisher.Publish(ctx, "wallet.credit", payload)  // Generic
```

**Trade-off:**
- ✅ Typed: An toàn hơn, rõ ràng hơn
- ❌ Typed: Phải tạo method mới cho mỗi job type
- ✅ Generic: Linh hoạt hơn, ít code hơn
- ❌ Generic: Dễ sai job type, không có compile-time check

**Kết luận:** Typed approach tốt hơn cho wallet operations (critical operations).

### Option 2: Dùng Redis Client Trực Tiếp

**Không khả thi vì:**
- Redis client chỉ làm key-value
- Không có job types, queues, retries
- Worker không thể đọc và xử lý

---

## 🎓 Design Decisions

### Tại Sao Dùng Typed Methods?

1. **Type Safety**: Compile-time checking
2. **Clarity**: `EnqueueCredit()` rõ hơn `Enqueue("credit", ...)`
3. **Future-proof**: Dễ thêm validation/transformation riêng cho mỗi type

### Tại Sao Dùng Interface Abstraction?

1. **Testability**: Dễ mock trong tests
2. **Flexibility**: Dễ switch provider
3. **Separation of Concerns**: Tách biệt client và implementation

### Tại Sao Xử Lý Duplicate Gracefully?

1. **Idempotency**: Đảm bảo không xử lý trùng
2. **User Experience**: Network retry không fail request
3. **Reliability**: Tránh lỗi do race conditions

---

## 📚 Tài Liệu Tham Khảo

- [Asynq Documentation](https://github.com/hibiken/asynq)
- [Redis Queue Patterns](https://redis.io/docs/manual/patterns/distributed-queues/)
- [Idempotency Patterns](https://stripe.com/docs/api/idempotent_requests)

---

## ❓ FAQ

**Q: Tại sao không dùng generic Queue interface từ `infra/queue`?**  
A: Vì wallet operations cần typed API cho type safety và clarity. Generic interface tốt cho events, nhưng wallet cần specific methods.

**Q: Có thể thêm job types khác không?**  
A: Có, nhưng cần:
- Thêm method mới vào `Queue` interface
- Implement trong `asynqQueue`
- Thêm job type constant trong `queue.constants.go`
- Update worker để handle job type mới

**Q: Tại sao duplicate tasks return nil?**  
A: Để đảm bảo idempotency. Nếu job đã tồn tại, không cần enqueue lại (đã được xử lý hoặc đang xử lý).

**Q: Có thể thay đổi retry count không?**  
A: Có, sửa trong code:
```go
asynq.MaxRetry(5)  // Thay đổi từ 3 → 5
```

**Q: Worker và Client có dùng cùng Redis connection không?**  
A: Có thể dùng cùng Redis instance, nhưng khác connections (client và worker là 2 processes riêng).

---

**Tài liệu này được tạo bởi: Ludo Backend Team**  
**Phiên bản: 1.0**  
**Cập nhật: 2024**

