# 🔨 Buổi 2: Basic API Development

## Mục Tiêu

- ✅ Implement Clean Architecture với 4 layers
- ✅ CRUD operations: Create, Read (List, Get by ID)
- ✅ In-memory storage với thread-safe
- ✅ Hiểu dependency injection và interfaces

## So Sánh với Buổi 1

| Buổi 1          | Buổi 2              |
| --------------- | ------------------- |
| 1 file main.go  | 4 layers riêng biệt |
| Hello World     | Full CRUD API       |
| No data storage | In-memory storage   |
| Monolithic      | Clean Architecture  |

## Architecture Overview

```
Request Flow:
HTTP Request
  → Handler (parse JSON, HTTP concerns)
    → Service (business logic, validation)
      → Storage (data persistence)
        → Model (domain entities)
```

## Key Changes

### 1. Entity Layer (`internal/model/`)

**New Files:**

- `asset.go` - Asset struct và constants
- `errors.go` - Custom error types

**Teaching Points:**

- Pure domain logic, no dependencies
- Struct tags cho JSON marshalling
- Constants for type safety

### 2. Storage Layer (`internal/storage/`)

**New Files:**

- `storage.go` - Interface definition
- `memory/memory.go` - In-memory implementation

**Teaching Points:**

- Interface cho flexibility (swap implementations)
- Thread-safety với sync.RWMutex
- Why interface? → Buổi 3 sẽ swap sang database!

### 3. Service Layer (`internal/service/`)

**New Files:**

- `asset_service.go` - Business logic
- `service.go` - Service interface

**Teaching Points:**

- Validation logic
- UUID generation
- Business rules (default status = active)
- Dependency injection (nhận Storage interface)

### 4. Handler Layer (`internal/handler/`)

**New Files:**

- `asset_handler.go` - HTTP handlers cho assets
- `health_handler.go` - Health check (refactored từ main)
- `handler.go` - Handler registry

**Teaching Points:**

- HTTP-specific code only
- JSON parsing và encoding
- Status codes (201, 400, 404, 500)
- Dependency injection (nhận Service)

### 5. Main (`cmd/server/main.go`)

**Changes:**

- Wire up all dependencies
- Register routes
- Remove business logic (moved to layers)

## API Endpoints

| Method | Path         | Description      | Status      |
| ------ | ------------ | ---------------- | ----------- |
| GET    | /health      | Health check     | ✅          |
| POST   | /assets      | Create asset     | ✅          |
| GET    | /assets      | List all assets  | ✅          |
| GET    | /assets/{id} | Get single asset | ✅          |
| PUT    | /assets/{id} | Update asset     | 🔜 Homework |
| DELETE | /assets/{id} | Delete asset     | 🔜 Homework |

## Testing

### 1. Health Check

```bash
curl http://localhost:8080/health
```

### 2. Create Asset

```bash
curl -X POST http://localhost:8080/assets \
  -H "Content-Type: application/json" \
  -d '{
    "name": "example.com",
    "type": "domain"
  }'
```

### 3. List Assets

```bash
curl http://localhost:8080/assets
```

### 4. Get Single Asset

```bash
# Replace <id> with actual UUID from create response
curl http://localhost:8080/assets/<id>
```

## Teaching Flow

### Step 1: Review Buổi 1 (10 phút)

- Show main.go từ session1
- Highlight: everything in one place
- Problem: hard to test, hard to maintain

### Step 2: Explain Clean Architecture (15 phút)

- Draw dependency diagram on board
- Show CLEAN_ARCHITECTURE.MD Section 2
- Emphasize: Dependency Rule (inner layers don't know outer)

### Step 3: Bottom-Up Implementation Walkthrough (90 phút)

#### 3.1 Entity Layer (15 phút)

Open `internal/model/asset.go`:

```go
type Asset struct {
    ID        string    `json:"id"`
    Name      string    `json:"name"`
    Type      string    `json:"type"`
    Status    string    `json:"status"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Q&A:**

- Why struct tags? → JSON field names
- Why constants? → Type safety, avoid typos
- Why time.Time? → Proper timestamp handling

#### 3.2 Storage Layer (25 phút)

**Interface first:**

```go
// storage/storage.go
type Storage interface {
    Create(asset *model.Asset) error
    GetAll() ([]*model.Asset, error)
    GetByID(id string) (*model.Asset, error)
}
```

**Q: Why interface?**
→ Buổi 3 sẽ implement PostgresStorage, chỉ cần implement interface này!

**Implementation:**

```go
// storage/memory/memory.go
type MemoryStorage struct {
    data map[string]*model.Asset
    mu   sync.RWMutex
}
```

**Key Points:**

- map[string]\*model.Asset → fast lookup by ID
- sync.RWMutex → thread-safe
  - RLock() for reads (multiple readers OK)
  - Lock() for writes (exclusive)

**Demo race condition** (if time permits):

```bash
# Without mutex → data race
go run -race cmd/server/main.go
# Concurrent requests will show warnings
```

#### 3.3 Service Layer (25 phút)

```go
// service/asset_service.go
func (s *AssetService) CreateAsset(name, assetType string) (*model.Asset, error) {
    // 1. Validation
    if name == "" {
        return nil, errors.New("name is required")
    }

    // 2. Business logic
    asset := &model.Asset{
        ID:        uuid.New().String(),
        Name:      name,
        Type:      assetType,
        Status:    model.StatusActive, // Default
        CreatedAt: time.Now(),
        UpdatedAt: time.Now(),
    }

    // 3. Delegate to storage
    if err := s.storage.Create(asset); err != nil {
        return nil, err
    }

    return asset, nil
}
```

**Teaching Points:**

- Validation BEFORE business logic
- UUID auto-generation
- Default values (status = active)
- Timestamps auto-set
- Service doesn't know HOW data is stored (memory? DB?)

#### 3.4 Handler Layer (25 phút)

```go
// handler/asset_handler.go
func (h *AssetHandler) CreateAsset(w http.ResponseWriter, r *http.Request) {
    // 1. Parse request
    var req CreateAssetRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        RespondError(w, http.StatusBadRequest, "Invalid request body")
        return
    }

    // 2. Call service
    asset, err := h.service.CreateAsset(req.Name, req.Type)
    if err != nil {
        RespondError(w, http.StatusBadRequest, err.Error())
        return
    }

    // 3. Return response
    RespondJSON(w, http.StatusCreated, asset)
}
```

**Teaching Points:**

- Handler only handles HTTP concerns
- No business logic here!
- Status codes: 201 for created, 400 for bad request
- Helper functions: RespondJSON, RespondError

### Step 4: Wire Up in Main (10 phút)

```go
// cmd/server/main.go
func main() {
    // 1. Initialize storage
    store := memory.NewMemoryStorage()

    // 2. Initialize service (inject storage)
    assetService := service.NewAssetService(store)

    // 3. Initialize handlers (inject service)
    assetHandler := handler.NewAssetHandler(assetService)
    healthHandler := handler.NewHealthHandler()

    // 4. Register routes
    mux := http.NewServeMux()
    mux.HandleFunc("GET /health", healthHandler.Check)
    mux.HandleFunc("POST /assets", assetHandler.CreateAsset)
    mux.HandleFunc("GET /assets", assetHandler.ListAssets)
    mux.HandleFunc("GET /assets/{id}", assetHandler.GetAsset)

    // 5. Start server
    log.Println("Server starting on :8080")
    http.ListenAndServe(":8080", mux)
}
```

**Key Point:** Dependency injection từ ngoài vào trong!

### Step 5: Live Demo (20 phút)

1. Run server: `go run cmd/server/main.go`
2. Create 3 assets (domain, IP, service)
3. List all assets
4. Get single asset
5. Try invalid input (empty name) → see validation

### Step 6: Compare with Session 1 (10 phút)

**Session 1 approach (if we continued):**

```go
func createAssetHandler(w http.ResponseWriter, r *http.Request) {
    // Parse JSON
    // Validate
    // Generate UUID
    // Store in global map
    // Return response
    // Everything mixed together!
}
```

**Session 2 approach (Clean Architecture):**

- Each layer has single responsibility
- Easy to test each layer independently
- Easy to swap storage implementation
- Clear separation of concerns

## Homework

### Task 1: Implement Update & Delete

**Update endpoint:**

```go
// PUT /assets/{id}
func (s *AssetService) UpdateAsset(id string, updates map[string]interface{}) (*model.Asset, error) {
    // 1. Get existing asset
    // 2. Validate updates
    // 3. Apply updates
    // 4. Update timestamp
    // 5. Save to storage
}
```

**Delete endpoint:**

```go
// DELETE /assets/{id}
func (s *AssetService) DeleteAsset(id string) error {
    // 1. Check if exists
    // 2. Delete from storage
}
```

### Task 2: Write Tests

```go
// service/asset_service_test.go
func TestCreateAsset(t *testing.T) {
    // Mock storage
    // Test validation
    // Test success case
}
```

### Task 3: Add More Validation

- Name length: 1-255 characters
- Type must be: domain, ip, or service
- Name format validation (domain = must have dot)

## Common Mistakes to Watch For

### ❌ Mistake 1: Business Logic in Handler

```go
// BAD
func (h *Handler) CreateAsset(w http.ResponseWriter, r *http.Request) {
    if req.Name == "" {  // Validation in handler - WRONG!
        // ...
    }
}
```

**Fix:** Move validation to service layer

### ❌ Mistake 2: Handler Calls Storage Directly

```go
// BAD
func (h *Handler) CreateAsset(w http.ResponseWriter, r *http.Request) {
    h.storage.Create(asset)  // Skip service layer - WRONG!
}
```

**Fix:** Always go through service layer

### ❌ Mistake 3: Not Thread-Safe Storage

```go
// BAD
type MemoryStorage struct {
    data map[string]*Asset  // No mutex!
}
```

**Fix:** Add sync.RWMutex

## Resources

- Review: CLEAN_ARCHITECTURE.MD sections 2-4
- [Go Interfaces](https://go.dev/tour/methods/9)
- [Dependency Injection in Go](https://blog.drewolson.org/dependency-injection-in-go)
- [UUID package](https://github.com/google/uuid)
