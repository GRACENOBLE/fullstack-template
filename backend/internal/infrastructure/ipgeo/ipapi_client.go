package ipgeo

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"backend/internal/domain"
	"backend/internal/usecase"
)

// ErrPrivateIP is returned when the caller passes a loopback or RFC-1918
// address. Callers should silently skip the location update in this case.
var ErrPrivateIP = errors.New("ipgeo: private or loopback IP address")

const defaultBaseURL = "https://ipapi.co"

// Client implements usecase.GeoLocator by querying https://ipapi.co/{ip}/json/.
type Client struct {
	httpClient *http.Client
	cache      usecase.CacheService // nil = no caching
	apiKey     string               // empty = free tier
	baseURL    string               // empty = defaultBaseURL
}

// New returns a Client configured for production use.
// cache may be nil (disables caching). apiKey may be empty (free tier).
func New(cache usecase.CacheService, apiKey string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		cache:      cache,
		apiKey:     apiKey,
	}
}

// NewWithBaseURL returns a Client with a custom base URL, intended for tests
// pointing at an httptest server.
func NewWithBaseURL(cache usecase.CacheService, apiKey, baseURL string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 5 * time.Second},
		cache:      cache,
		apiKey:     apiKey,
		baseURL:    baseURL,
	}
}

// Compile-time check: Client satisfies usecase.GeoLocator.
var _ usecase.GeoLocator = (*Client)(nil)

// ipapiResponse matches the JSON shape returned by ipapi.co.
type ipapiResponse struct {
	IP          string `json:"ip"`
	CountryCode string `json:"country_code"`
	CountryName string `json:"country_name"`
	Region      string `json:"region"`
	City        string `json:"city"`
	Timezone    string `json:"timezone"`
	Currency    string `json:"currency"`
	InEU        bool   `json:"in_eu"`
	Error       bool   `json:"error"`
	Reason      string `json:"reason"`
}

// cacheKey returns the Redis key for a given IP.
func cacheKey(ip string) string {
	return "geo:" + ip
}

// Lookup resolves ip to geographic metadata.
// Returns ErrPrivateIP for loopback and RFC-1918 addresses without making
// any outbound HTTP call or cache lookup.
func (c *Client) Lookup(ctx context.Context, ip string) (*domain.GeoLocation, error) {
	if isPrivateIP(ip) {
		return nil, ErrPrivateIP
	}

	// Check cache first.
	if c.cache != nil {
		if val, ok, err := c.cache.Get(ctx, cacheKey(ip)); err == nil && ok {
			var geo domain.GeoLocation
			if jsonErr := json.Unmarshal([]byte(val), &geo); jsonErr == nil {
				return &geo, nil
			}
		}
	}

	geo, err := c.fetch(ctx, ip)
	if err != nil {
		return nil, err
	}

	// Populate cache.
	if c.cache != nil {
		if data, jsonErr := json.Marshal(geo); jsonErr == nil {
			// Best-effort; ignore cache write errors.
			_ = c.cache.Set(ctx, cacheKey(ip), string(data), 24*time.Hour)
		}
	}

	return geo, nil
}

// fetch performs the HTTP GET against ipapi.co and returns a GeoLocation.
func (c *Client) fetch(ctx context.Context, ip string) (*domain.GeoLocation, error) {
	base := c.baseURL
	if base == "" {
		base = defaultBaseURL
	}

	var rawURL string
	if c.apiKey != "" {
		rawURL = fmt.Sprintf("%s/%s/json/?key=%s", base, ip, c.apiKey)
	} else {
		rawURL = fmt.Sprintf("%s/%s/json/", base, ip)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, rawURL, nil)
	if err != nil {
		return nil, fmt.Errorf("ipgeo: build request: %w", err)
	}
	req.Header.Set("User-Agent", "template-backend/1.0")

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("ipgeo: http request: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("ipgeo: unexpected status %d", resp.StatusCode)
	}

	var body ipapiResponse
	if decErr := json.NewDecoder(resp.Body).Decode(&body); decErr != nil {
		return nil, fmt.Errorf("ipgeo: decode response: %w", decErr)
	}

	if body.Error {
		return nil, fmt.Errorf("ipgeo: api error: %s", body.Reason)
	}

	return &domain.GeoLocation{
		IP:          body.IP,
		CountryCode: body.CountryCode,
		CountryName: body.CountryName,
		Region:      body.Region,
		City:        body.City,
		Timezone:    body.Timezone,
		Currency:    body.Currency,
		IsEU:        body.InEU,
	}, nil
}

// isPrivateIP reports whether the given IP string is a loopback or
// RFC-1918 / RFC-4193 private address.
func isPrivateIP(ipStr string) bool {
	ip := net.ParseIP(ipStr)
	if ip == nil {
		return false
	}
	return ip.IsLoopback() || ip.IsPrivate()
}
