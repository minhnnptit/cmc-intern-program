# 🔍 Buổi 4: Advanced Features

## Mục Tiêu

- ✅ Complete CRUD operations (Update, Delete đã có từ Buổi 2)
- ✅ Advanced filtering và searching
- ✅ Input validation comprehensive
- ✅ Pagination
- ✅ Sorting
- ✅ Error handling patterns

## Nội Dung Mới

### 1. Pagination

```go
// Add to Storage interface
type PaginationParams struct {
    Page     int
    PageSize int
    Offset   int
    Limit    int
}

func (s *Storage) GetAllPaginated(params PaginationParams) (*PaginatedResult, error)

type PaginatedResult struct {
    Data       []*model.Asset `json:"data"`
    Total      int            `json:"total"`
    Page       int            `json:"page"`
    PageSize   int            `json:"page_size"`
    TotalPages int            `json:"total_pages"`
}
```

**API Usage:**

```bash
GET /assets?page=1&page_size=10
```

### 2. Sorting

```go
// Query params
GET /assets?sort_by=created_at&sort_order=desc
GET /assets?sort_by=name&sort_order=asc
```

**Implementation:**

```go
func (p *PostgresStorage) GetAll(sortBy, sortOrder string) ([]*model.Asset, error) {
    // Validate sort fields (prevent SQL injection)
    validSortFields := map[string]bool{
        "name": true,
        "created_at": true,
        "updated_at": true,
    }

    if !validSortFields[sortBy] {
        sortBy = "created_at"
    }

    if sortOrder != "asc" && sortOrder != "desc" {
        sortOrder = "desc"
    }

    query := fmt.Sprintf("SELECT * FROM assets ORDER BY %s %s", sortBy, sortOrder)
    // ...
}
```

### 3. Advanced Validation

```go
// internal/validator/asset_validator.go
package validator

type AssetValidator struct{}

func (v *AssetValidator) ValidateCreate(req CreateAssetRequest) error {
    // Name validation
    if len(req.Name) == 0 {
        return errors.New("name is required")
    }
    if len(req.Name) > 255 {
        return errors.New("name too long (max 255 characters)")
    }

    // Type-specific validation
    switch req.Type {
    case "domain":
        if !isDomainValid(req.Name) {
            return errors.New("invalid domain format")
        }
    case "ip":
        if !isIPValid(req.Name) {
            return errors.New("invalid IP address")
        }
    case "service":
        if !isServiceValid(req.Name) {
            return errors.New("invalid service format")
        }
    }

    return nil
}

func isDomainValid(domain string) bool {
    // Must have at least one dot
    // No spaces, special characters
    pattern := `^[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?(\.[a-zA-Z0-9]([a-zA-Z0-9\-]{0,61}[a-zA-Z0-9])?)*$`
    matched, _ := regexp.MatchString(pattern, domain)
    return matched
}

func isIPValid(ip string) bool {
    return net.ParseIP(ip) != nil
}
```

### 4. Complex Filtering

**Multiple filters combined:**

```bash
GET /assets?type=domain&status=active&search=example&page=1&page_size=10&sort_by=name
```

**Implementation:**

```go
type FilterParams struct {
    Type      string
    Status    string
    Search    string
    Page      int
    PageSize  int
    SortBy    string
    SortOrder string
}

func (h *AssetHandler) ListAssets(w http.ResponseWriter, r *http.Request) {
    params := FilterParams{
        Type:      r.URL.Query().Get("type"),
        Status:    r.URL.Query().Get("status"),
        Search:    r.URL.Query().Get("search"),
        Page:      getIntParam(r, "page", 1),
        PageSize:  getIntParam(r, "page_size", 20),
        SortBy:    r.URL.Query().Get("sort_by"),
        SortOrder: r.URL.Query().Get("sort_order"),
    }

    result, err := h.service.FilterAssets(params)
    // ...
}
```

### 5. Bulk Operations

```go
// Bulk create
POST /assets/bulk
{
  "assets": [
    {"name": "example1.com", "type": "domain"},
    {"name": "example2.com", "type": "domain"}
  ]
}

// Bulk delete
DELETE /assets/bulk
{
  "ids": ["uuid1", "uuid2", "uuid3"]
}
```

### 6. Error Response Enhancement

```go
type ErrorResponse struct {
    Error   string                 `json:"error"`
    Code    string                 `json:"code"`
    Details map[string]interface{} `json:"details,omitempty"`
}

// Example response
{
  "error": "Validation failed",
  "code": "VALIDATION_ERROR",
  "details": {
    "name": "name is required",
    "type": "invalid type"
  }
}
```

## Teaching Flow

### 1. Review Query Parameters (15 phút)

- Show basic filtering từ Buổi 2
- Introduce advanced use cases
- API design principles

### 2. Implement Pagination (30 phút)

- Calculate offset từ page/page_size
- SQL LIMIT/OFFSET
- Count total records
- Return metadata

### 3. Implement Sorting (20 phút)

- Dynamic ORDER BY
- SQL injection prevention
- Whitelist sort fields

### 4. Advanced Validation (30 phút)

- Regex patterns
- Type-specific validation
- Validator package structure
- Unit testing validators

### 5. Combine Everything (30 phút)

- Filter + Search + Pagination + Sort
- Query builder pattern
- Performance considerations

### 6. Demo & Testing (15 phút)

- Complex query examples
- Postman collection
- Edge cases

## Performance Considerations

### Database Indexes

```sql
-- Composite indexes for common queries
CREATE INDEX idx_assets_type_status ON assets(type, status);
CREATE INDEX idx_assets_name_search ON assets(name text_pattern_ops);

-- For pagination
CREATE INDEX idx_assets_created_at_id ON assets(created_at DESC, id);
```

### Query Optimization

```go
// ❌ BAD - N+1 problem
for _, asset := range assets {
    details := getAssetDetails(asset.ID) // Multiple queries
}

// ✅ GOOD - Single query with JOIN
assets := getAssetsWithDetails() // One query
```

## Testing Scenarios

```bash
# Test pagination
curl "http://localhost:8080/assets?page=1&page_size=5"
curl "http://localhost:8080/assets?page=2&page_size=5"

# Test sorting
curl "http://localhost:8080/assets?sort_by=name&sort_order=asc"
curl "http://localhost:8080/assets?sort_by=created_at&sort_order=desc"

# Test combined filters
curl "http://localhost:8080/assets?type=domain&status=active&search=example&page=1&sort_by=name"

# Test validation
curl -X POST http://localhost:8080/assets \
  -d '{"name":"", "type":"domain"}' # Should fail

curl -X POST http://localhost:8080/assets \
  -d '{"name":"invalid..domain", "type":"domain"}' # Should fail
```

## Homework

1. **Add date range filtering**

   ```
   GET /assets?created_after=2026-01-01&created_before=2026-12-31
   ```

2. **Add export functionality**

   ```
   GET /assets/export?format=csv
   GET /assets/export?format=json
   ```

3. **Add field selection**

   ```
   GET /assets?fields=id,name,type
   ```

4. **Add aggregation endpoint**
   ```
   GET /assets/stats
   {
     "total": 100,
     "by_type": {"domain": 60, "ip": 30, "service": 10},
     "by_status": {"active": 80, "inactive": 20}
   }
   ```

## Resources

- [SQL Performance Explained](https://use-the-index-luke.com/)
- [Go Validator Package](https://github.com/go-playground/validator)
- [REST API Best Practices](https://stackoverflow.blog/2020/03/02/best-practices-for-rest-api-design/)
