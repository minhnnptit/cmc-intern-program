package main

import (
	"log"
	"net/http"

	"mini-asm/internal/handler"
	"mini-asm/internal/service"
	"mini-asm/internal/storage/memory"
)

func main() {
	log.Println("🚀 Starting Mini ASM Server...")

	// ============================================
	// DEPENDENCY INJECTION - Wire up all layers
	// ============================================

	// 1. Initialize Storage Layer (Infrastructure)
	//    Using in-memory storage for now
	//    Buổi 3 sẽ swap sang PostgreSQL - chỉ thay đổi dòng này!
	store := memory.NewMemoryStorage()
	log.Println("✅ Storage initialized: In-Memory")

	// 2. Initialize Service Layer (Use Case / Business Logic)
	//    Inject storage dependency
	assetService := service.NewAssetService(store)
	log.Println("✅ Service initialized: AssetService")

	// 3. Initialize Handler Layer (Presentation / HTTP)
	//    Inject service dependency
	assetHandler := handler.NewAssetHandler(assetService)
	healthHandler := handler.NewHealthHandler()
	log.Println("✅ Handlers initialized")

	// ============================================
	// ROUTING - Register HTTP endpoints
	// ============================================

	mux := http.NewServeMux()

	// Health check
	mux.HandleFunc("GET /health", healthHandler.Check)

	// Asset CRUD operations
	mux.HandleFunc("POST /assets", assetHandler.CreateAsset)        // Create
	mux.HandleFunc("GET /assets", assetHandler.ListAssets)          // Read (list with filters)
	mux.HandleFunc("GET /assets/{id}", assetHandler.GetAsset)       // Read (single)
	mux.HandleFunc("PUT /assets/{id}", assetHandler.UpdateAsset)    // Update
	mux.HandleFunc("DELETE /assets/{id}", assetHandler.DeleteAsset) // Delete

	log.Println("✅ Routes registered:")
	log.Println("   GET    /health")
	log.Println("   POST   /assets")
	log.Println("   GET    /assets")
	log.Println("   GET    /assets/{id}")
	log.Println("   PUT    /assets/{id}")
	log.Println("   DELETE /assets/{id}")

	// ============================================
	// START SERVER
	// ============================================

	addr := ":8080"
	log.Printf("🌐 Server listening on http://localhost%s\n", addr)
	log.Println("📖 API Documentation: see docs/api.yml")
	log.Println("Press Ctrl+C to stop")
	log.Println()

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Fatal("❌ Server failed to start:", err)
	}
}

/*
🎓 TEACHING NOTES:

=== SO SÁNH VỚI BUỔI 1 ===

Buổi 1 (main.go):
```
func main() {
    http.HandleFunc("GET /health", func(w, r) {
        // All logic here - mixed responsibilities
        json.NewEncoder(w).Encode(...)
    })
    http.ListenAndServe(":8080", nil)
}
```
- Monolithic
- Hard to test
- Hard to maintain

Buổi 2 (main.go):
```
func main() {
    store := memory.NewMemoryStorage()      // Layer 1: Storage
    service := service.NewAssetService(store) // Layer 2: Service
    handler := handler.NewAssetHandler(service) // Layer 3: Handler

    mux.HandleFunc("POST /assets", handler.CreateAsset)
    http.ListenAndServe(":8080", mux)
}
```
- Clean separation
- Easy to test (inject mocks)
- Easy to maintain (change one layer)

=== DEPENDENCY INJECTION FLOW ===

Outside → Inside (Dependency Rule):
```
main.go
  creates Storage
    ↓ inject to
  creates Service
    ↓ inject to
  creates Handler
```

Inner layers DON'T know about outer layers:
- Service doesn't know about Handler
- Storage doesn't know about Service
- Model doesn't know about anything

=== BUỔI 3 PREVIEW ===

Only 1 line changes:
```go
// Buổi 2
store := memory.NewMemoryStorage()

// Buổi 3
db := connectToPostgres() // New
store := postgres.NewPostgresStorage(db) // Changed
```

Service, Handler, Model: NO CHANGES!
→ This is the power of Clean Architecture!

=== ROUTING WITH Go 1.22+ ===

Method in pattern:
"GET /health"    - only handles GET requests
"POST /assets"   - only handles POST requests
"GET /assets/{id}" - path parameter

Trước Go 1.22:
http.HandleFunc("/assets", func(w, r) {
    if r.Method == "GET" { ... }
    else if r.Method == "POST" { ... }
})

=== KEY CONCEPTS TO EMPHASIZE ===

1. Dependency Injection:
   - Dependencies flow FROM main TO inner layers
   - Not the opposite!
   - Benefits: testability, flexibility

2. Single Responsibility:
   - main.go: wire up dependencies, start server
   - handler: HTTP concerns
   - service: business logic
   - storage: data access

3. Easy to Change:
   - Swap storage implementation
   - Add middleware
   - Change routing
   - All without touching business logic!

=== COMMON QUESTIONS ===

Q: Tại sao không dùng global variables?
A: Hard to test, hard to swap implementations

Q: Tại sao cần separate handlers?
A: Single Responsibility, reusable, testable

Q: Order của dependency injection có quan trọng không?
A: CÓ! Must create inner layers first:
   storage → service → handler

Q: Có thể có multiple services không?
A: CÓ! Example: AssetService, UserService, AuthService
   Each with their own handlers

=== DEMO SCRIPT ===

1. Show code structure:
   cmd/server/main.go  → Entry point
   internal/handler/   → HTTP layer
   internal/service/   → Business logic
   internal/storage/   → Data access
   internal/model/     → Domain entities

2. Explain dependency flow (draw diagram)

3. Run server:
   go run cmd/server/main.go

4. Test endpoints:
   curl http://localhost:8080/health
   curl -X POST http://localhost:8080/assets -d '{"name":"example.com","type":"domain"}'
   curl http://localhost:8080/assets

5. Highlight:
   - Clean logs
   - Clear structure
   - Easy to understand

6. Compare với Buổi 1:
   - Show Buổi 1 code
   - "What if cần thêm database?"
   - "What if cần thêm authentication?"
   - → Buổi 1 approach: mess
   - → Buổi 2 approach: add layer or middleware

=== HOMEWORK HINTS ===

Students có thể:
1. Add logging middleware
2. Add request ID tracking
3. Add CORS headers
4. Add rate limiting
5. Add metrics endpoint

Tất cả là middleware, không touch business logic!
*/
