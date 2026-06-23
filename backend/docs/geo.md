---
topic: geo
last_verified: 2026-06-23
sources:
  - backend/internal/domain/geolocation.go
  - backend/internal/usecase/geolocation.go
  - backend/internal/infrastructure/ipgeo/ipapi_client.go
  - backend/internal/transport/middleware/geo.go
  - backend/internal/bootstrap/bootstrap.go
---

# IP Geolocation

## What it does
The geolocation feature resolves a client IP address to geographic metadata (country, region, city, timezone, currency, EU membership) using the [ipapi.co](https://ipapi.co) JSON API. Results are cached in Redis under the key `geo:<ip>` for 24 hours when a `CacheService` is available.

The feature is always initialised — ipapi.co's free tier requires no API key. Supplying `IPAPI_KEY` enables the paid tier with higher rate limits.

Private and loopback IPs (`ErrPrivateIP`) are rejected immediately without any HTTP call or cache lookup.

## domain.GeoLocation

```go
// backend/internal/domain/geolocation.go
type GeoLocation struct {
    IP          string
    CountryCode string
    CountryName string
    Region      string
    City        string
    Timezone    string
    Currency    string
    IsEU        bool
}
```

## usecase.GeoLocator interface

```go
// backend/internal/usecase/geolocation.go
type GeoLocator interface {
    Lookup(ctx context.Context, ip string) (*domain.GeoLocation, error)
}
```

## ipgeo.Client constructors

```go
// Production — cache may be nil (disables caching), apiKey may be empty (free tier).
func New(cache usecase.CacheService, apiKey string) *Client

// Tests — points the HTTP client at a custom base URL (e.g. an httptest.Server).
func NewWithBaseURL(cache usecase.CacheService, apiKey, baseURL string) *Client
```

`Client` satisfies `usecase.GeoLocator` via a compile-time assertion:
```go
var _ usecase.GeoLocator = (*Client)(nil)
```

### Lookup behaviour

1. Returns `ErrPrivateIP` immediately for loopback and RFC-1918/RFC-4193 addresses.
2. Checks the cache (`geo:<ip>`) — returns the deserialised value on hit.
3. On cache miss, calls `GET {baseURL}/{ip}/json/` (appends `?key=<apiKey>` when a key is set).
4. On success, writes the serialised result to the cache (TTL 24 h, best-effort — write errors are silently ignored).
5. Returns the populated `*domain.GeoLocation`.

Sentinel error: `ipgeo.ErrPrivateIP`

## GeoFromRequest middleware

```go
// backend/internal/transport/middleware/geo.go
const GeoLocationKey = "geo_location"

func GeoFromRequest(locator usecase.GeoLocator) gin.HandlerFunc
```

Best-effort: if `locator.Lookup` returns any error (private IP, rate-limit, network error), the middleware calls `c.Next()` without storing anything. Handlers must nil-check before reading the key.

## RealIP helper

```go
func RealIP(r *http.Request) string
```

Precedence order:
1. `X-Forwarded-For` header — takes the first (leftmost) comma-separated address.
2. `X-Real-IP` header.
3. `r.RemoteAddr` — `net.SplitHostPort` strips the port; falls back to the raw string on parse error.

`RealIP` is exported so tests and other packages can call it directly.

## Reading geo data in a handler

```go
val, exists := c.Get(middleware.GeoLocationKey)
if !exists {
    // geo unavailable — private IP, rate-limited, or middleware not registered
    return
}
geo, ok := val.(*domain.GeoLocation)
if !ok || geo == nil {
    return
}
// use geo.CountryCode, geo.City, etc.
```

## Bootstrap wiring

`GeoLocator` is always initialised inside `bootstrap.Run()`, unconditionally:

```go
// bootstrap.go — always runs, not guarded by an env var check
geoLocator := ipgeo.New(cache, cfg.IPAPIKey)
log.Info("bootstrap: ipapi geolocation client initialised", "cached", cache != nil)
```

`cache` is the same `usecase.CacheService` used elsewhere — it is `nil` when `REDIS_URL` is not set, which disables caching for geolocation as well.

`App.GeoLocator` is always non-nil after a successful `Run()`.

## Environment variable

| Variable | Required | Default | Description |
|---|---|---|---|
| `IPAPI_KEY` | No | — | ipapi.co API key. Omit or leave empty for the free tier. |

## Testing patterns

### Unit-testing ipgeo.Client
Use `NewWithBaseURL` pointing at an `httptest.NewServer`:

```go
srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
    json.NewEncoder(w).Encode(map[string]any{
        "ip": "1.2.3.4", "country_code": "US", ...
    })
}))
defer srv.Close()

client := ipgeo.NewWithBaseURL(nil, "", srv.URL)
geo, err := client.Lookup(context.Background(), "1.2.3.4")
```

Inject a cache with an inline mock struct:

```go
type mockCacheService struct {
    data map[string]string
}

func (m *mockCacheService) Get(_ context.Context, key string) (string, bool, error) {
    v, ok := m.data[key]
    return v, ok, nil
}
func (m *mockCacheService) Set(_ context.Context, key, value string, _ time.Duration) error {
    m.data[key] = value
    return nil
}
// implement remaining usecase.CacheService methods as no-ops
```

### Unit-testing GeoFromRequest middleware
Use an inline `mockGeoLocator`:

```go
type mockGeoLocator struct {
    geo *domain.GeoLocation
    err error
}

func (m *mockGeoLocator) Lookup(_ context.Context, _ string) (*domain.GeoLocation, error) {
    return m.geo, m.err
}
```

Record the handler under test with `httptest.NewRecorder` and assert `c.Get(middleware.GeoLocationKey)`.
