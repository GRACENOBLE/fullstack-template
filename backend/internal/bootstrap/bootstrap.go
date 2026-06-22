package bootstrap

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"math"
	"math/rand/v2"
	"os"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	_ "github.com/joho/godotenv/autoload"

	"backend/internal/infrastructure/cache/redis"
	"backend/internal/infrastructure/database/postgres"
	"backend/internal/infrastructure/email"
	"backend/internal/infrastructure/queue"
	"backend/internal/infrastructure/storage/r2"
	"backend/internal/usecase"
	"backend/pkg/firebase"
	"backend/pkg/logger"
)

const (
	maxAttempts  = 5
	baseDelay    = 500 * time.Millisecond
	maxDelay     = 16 * time.Second
	totalTimeout = 60 * time.Second
	pingTimeout  = 15 * time.Second // accommodates Neon cold starts (~8-15 s)
)

// App holds all initialised, validated shared dependencies.
// Constructed once by Run and passed to the HTTP server.
type App struct {
	DB             *sql.DB
	Cache          usecase.CacheService        // nil when REDIS_URL is not set
	Enqueuer       usecase.Enqueuer            // nil when REDIS_URL is not set
	Firebase       usecase.FirebaseAdminClient // nil when FIREBASE_PROJECT_ID is not set
	FCMSender      usecase.NotificationSender  // nil when FIREBASE_PROJECT_ID is not set
	EmailSender    usecase.EmailSender         // nil when MAILJET_API_KEY is not set
	StorageService usecase.StorageService      // nil when R2_ACCOUNT_ID is not set
	Config         Config
	Log            *slog.Logger
}

// Config holds all validated configuration values read from environment variables.
type Config struct {
	Port                       int
	Env                        string
	DB                         postgres.DBConfig
	RedisURL                   string
	RateLimitRPS               float64
	RateLimitBurst             int
	FirebaseProjectID          string
	FirebaseServiceAccountJSON string
	SentryDSN                  string
	MailjetAPIKey              string
	MailjetSecretKey           string
	FromEmail                  string
	FromName                   string
	R2AccountID                string
	R2AccessKey                string
	R2SecretKey                string
	R2Bucket                   string
	R2PublicURL                string
}

// ConfigError is returned when required configuration is absent or invalid.
type ConfigError struct {
	Issues []string
}

func (e *ConfigError) Error() string {
	return "invalid configuration: " + strings.Join(e.Issues, "; ")
}

// Pinger is satisfied by any dependency that can report its own liveness.
type Pinger interface {
	PingContext(ctx context.Context) error
}

// Run loads configuration, initialises all shared dependencies, validates required
// config, and probes services for readiness before returning. A non-nil error means
// the process should not start; callers should exit with a non-zero status code.
func Run(ctx context.Context) (*App, error) {
	log := logger.New(os.Getenv("ENV"))
	slog.SetDefault(log)

	log.Info("bootstrap: starting")

	cfg := loadConfig()

	if err := validateConfig(cfg, log); err != nil {
		return nil, err
	}

	db, err := postgres.NewPostgresDB(cfg.DB)
	if err != nil {
		return nil, fmt.Errorf("bootstrap: database: %w", err)
	}

	probeCtx, cancel := context.WithTimeout(ctx, totalTimeout)
	defer cancel()

	if err := probeWithRetry(probeCtx, "postgres", db, log); err != nil {
		return nil, err
	}

	var cache usecase.CacheService
	if cfg.RedisURL != "" {
		c, err := redis.New(cfg.RedisURL)
		if err != nil {
			return nil, fmt.Errorf("bootstrap: redis: %w", err)
		}
		if err := probeWithRetry(probeCtx, "redis", c, log); err != nil {
			return nil, err
		}
		cache = c
	}

	var enqueuer usecase.Enqueuer
	if cfg.RedisURL != "" {
		eq, err := queue.NewClient(cfg.RedisURL)
		if err != nil {
			return nil, fmt.Errorf("bootstrap: queue: %w", err)
		}
		enqueuer = eq
	}

	var firebaseClient usecase.FirebaseAdminClient
	var fcmSender usecase.NotificationSender
	if cfg.FirebaseProjectID != "" {
		fbApp, err := firebase.NewApp(ctx, cfg.FirebaseProjectID, cfg.FirebaseServiceAccountJSON)
		if err != nil {
			return nil, fmt.Errorf("bootstrap: firebase: %w", err)
		}
		authClient, err := firebase.NewAuthClient(ctx, fbApp)
		if err != nil {
			return nil, fmt.Errorf("bootstrap: firebase auth: %w", err)
		}
		firebaseClient = authClient
		msgClient, err := firebase.NewMessagingClient(ctx, fbApp)
		if err != nil {
			return nil, fmt.Errorf("bootstrap: firebase messaging: %w", err)
		}
		fcmSender = msgClient
		log.Info("bootstrap: firebase clients initialised", "project_id", cfg.FirebaseProjectID)
	}

	var emailSender usecase.EmailSender
	if cfg.MailjetAPIKey != "" && cfg.MailjetSecretKey != "" {
		emailSender = email.NewMailjetSender(
			cfg.MailjetAPIKey,
			cfg.MailjetSecretKey,
			cfg.FromEmail,
			cfg.FromName,
		)
		log.Info("bootstrap: mailjet email sender initialised", "from_email", cfg.FromEmail)
	}

	var storageService usecase.StorageService
	if cfg.R2AccountID != "" {
		svc, err := r2.New(cfg.R2AccountID, cfg.R2AccessKey, cfg.R2SecretKey, cfg.R2Bucket, cfg.R2PublicURL)
		if err != nil {
			return nil, fmt.Errorf("bootstrap: r2: %w", err)
		}
		storageService = svc
		log.Info("bootstrap: R2 storage client initialised", "bucket", cfg.R2Bucket)
	}

	log.Info("bootstrap: all checks passed — ready to serve")

	return &App{
		DB:             db,
		Cache:          cache,
		Enqueuer:       enqueuer,
		Firebase:       firebaseClient,
		FCMSender:      fcmSender,
		EmailSender:    emailSender,
		StorageService: storageService,
		Config:         cfg,
		Log:            log,
	}, nil
}

func loadConfig() Config {
	port, _ := strconv.Atoi(os.Getenv("PORT"))
	if port == 0 {
		port = 8080
	}

	schema := os.Getenv("BLUEPRINT_DB_SCHEMA")
	if schema == "" {
		schema = "public"
	}

	sslMode := os.Getenv("BLUEPRINT_DB_SSLMODE")
	if sslMode == "" {
		sslMode = "disable"
	}

	rps, _ := strconv.ParseFloat(os.Getenv("RATE_LIMIT_RPS"), 64)
	burst, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_BURST"))
	if burst == 0 && rps > 0 {
		burst = int(rps) * 5
	}

	return Config{
		Port:                       port,
		Env:                        os.Getenv("ENV"),
		RedisURL:                   os.Getenv("REDIS_URL"),
		RateLimitRPS:               rps,
		RateLimitBurst:             burst,
		FirebaseProjectID:          os.Getenv("FIREBASE_PROJECT_ID"),
		FirebaseServiceAccountJSON: os.Getenv("FIREBASE_SERVICE_ACCOUNT_JSON"),
		SentryDSN:                  os.Getenv("SENTRY_DSN"),
		MailjetAPIKey:              os.Getenv("MAILJET_API_KEY"),
		MailjetSecretKey:           os.Getenv("MAILJET_SECRET_KEY"),
		FromEmail:                  os.Getenv("FROM_EMAIL"),
		FromName:                   os.Getenv("FROM_NAME"),
		R2AccountID:                os.Getenv("R2_ACCOUNT_ID"),
		R2AccessKey:                os.Getenv("R2_ACCESS_KEY"),
		R2SecretKey:                os.Getenv("R2_SECRET_KEY"),
		R2Bucket:                   os.Getenv("R2_BUCKET"),
		R2PublicURL:                os.Getenv("R2_PUBLIC_URL"),
		DB: postgres.DBConfig{
			Host:     os.Getenv("BLUEPRINT_DB_HOST"),
			Port:     os.Getenv("BLUEPRINT_DB_PORT"),
			Database: os.Getenv("BLUEPRINT_DB_DATABASE"),
			Username: os.Getenv("BLUEPRINT_DB_USERNAME"),
			Password: os.Getenv("BLUEPRINT_DB_PASSWORD"),
			Schema:   schema,
			SSLMode:  sslMode,
		},
	}
}

func validateConfig(cfg Config, log *slog.Logger) error {
	log.Info("bootstrap: validating configuration")

	var issues []string

	requireNonEmpty := func(name, val string) {
		if strings.TrimSpace(val) == "" {
			issues = append(issues, fmt.Sprintf("%s must not be empty", name))
		}
	}

	requireNonEmpty("BLUEPRINT_DB_HOST", cfg.DB.Host)
	requireNonEmpty("BLUEPRINT_DB_PORT", cfg.DB.Port)
	requireNonEmpty("BLUEPRINT_DB_DATABASE", cfg.DB.Database)
	requireNonEmpty("BLUEPRINT_DB_USERNAME", cfg.DB.Username)
	requireNonEmpty("BLUEPRINT_DB_PASSWORD", cfg.DB.Password)

	// Mailjet: if any credential is provided, the full set is required.
	if cfg.MailjetAPIKey != "" || cfg.MailjetSecretKey != "" || cfg.FromEmail != "" {
		requireNonEmpty("MAILJET_API_KEY", cfg.MailjetAPIKey)
		requireNonEmpty("MAILJET_SECRET_KEY", cfg.MailjetSecretKey)
		requireNonEmpty("FROM_EMAIL", cfg.FromEmail)
	}

	if len(issues) > 0 {
		for _, issue := range issues {
			log.Error("bootstrap: config invalid", "detail", issue)
		}
		return &ConfigError{Issues: issues}
	}

	log.Info("bootstrap: configuration valid")
	return nil
}

func probeWithRetry(ctx context.Context, name string, p Pinger, log *slog.Logger) error {
	var lastErr error
	for attempt := 1; attempt <= maxAttempts; attempt++ {
		if attempt > 1 {
			delay := jitteredBackoff(attempt - 1)
			log.Info("bootstrap: waiting before retry",
				"service", name, "attempt", attempt, "delay", delay.String())
			select {
			case <-ctx.Done():
				return fmt.Errorf("bootstrap: %s: timed out after %d attempt(s): %w", name, attempt-1, lastErr)
			case <-time.After(delay):
			}
		}

		log.Info("bootstrap: probing service",
			"service", name, "attempt", attempt, "max_attempts", maxAttempts)

		attemptCtx, cancel := context.WithTimeout(ctx, pingTimeout)
		pingErr := p.PingContext(attemptCtx)
		cancel()

		if pingErr == nil {
			log.Info("bootstrap: service ready", "service", name, "attempts", attempt)
			return nil
		}
		lastErr = pingErr
		log.Warn("bootstrap: service not ready",
			"service", name, "attempt", attempt, "error", pingErr)
	}
	return fmt.Errorf("bootstrap: %s: not reachable after %d attempts: %w", name, maxAttempts, lastErr)
}

// jitteredBackoff returns a random duration in [0, min(maxDelay, baseDelay*2^attempt)].
// Full jitter avoids thundering-herd on simultaneous restarts.
func jitteredBackoff(attempt int) time.Duration {
	cap := time.Duration(math.Min(float64(maxDelay), float64(baseDelay)*math.Pow(2, float64(attempt))))
	return time.Duration(rand.Int64N(int64(cap) + 1))
}
