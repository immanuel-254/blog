package auth

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/immanuel-254/blog/auth/models"
)

// Config defines the configuration for the middleware.
type Config struct {
	Next                      func(r *http.Request) bool
	XSSProtection             string
	ContentTypeNosniff        string
	XFrameOptions             string
	ContentSecurityPolicy     string
	CSPReportOnly             bool
	ReferrerPolicy            string
	PermissionPolicy          string
	CrossOriginEmbedderPolicy string
	CrossOriginOpenerPolicy   string
	CrossOriginResourcePolicy string
	OriginAgentCluster        string
	XDNSPrefetchControl       string
	XDownloadOptions          string
	XPermittedCrossDomain     string
	HSTSMaxAge                int
	HSTSExcludeSubdomains     bool
	HSTSPreloadEnabled        bool
}

// ConfigDefault provides default configuration values.
var ConfigDefault = Config{
	XSSProtection:             "0",
	ContentTypeNosniff:        "nosniff",
	XFrameOptions:             "SAMEORIGIN",
	ReferrerPolicy:            "no-referrer",
	CrossOriginEmbedderPolicy: "require-corp",
	CrossOriginOpenerPolicy:   "same-origin",
	CrossOriginResourcePolicy: "same-origin",
	OriginAgentCluster:        "?1",
	XDNSPrefetchControl:       "off",
	XDownloadOptions:          "noopen",
	XPermittedCrossDomain:     "none",
}

// New creates the middleware handler for `net/http`.
func New(config ...Config) func(next http.Handler) http.Handler {
	// Initialize the configuration.
	cfg := configDefault(config...)

	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Check if the middleware should skip the request.
			if cfg.Next != nil && cfg.Next(r) {
				next.ServeHTTP(w, r)
				return
			}

			// Set headers based on the configuration.
			if cfg.XSSProtection != "" {
				w.Header().Set("X-XSS-Protection", cfg.XSSProtection)
			}
			if cfg.ContentTypeNosniff != "" {
				w.Header().Set("X-Content-Type-Options", cfg.ContentTypeNosniff)
			}
			if cfg.XFrameOptions != "" {
				w.Header().Set("X-Frame-Options", cfg.XFrameOptions)
			}
			if cfg.CrossOriginEmbedderPolicy != "" {
				w.Header().Set("Cross-Origin-Embedder-Policy", cfg.CrossOriginEmbedderPolicy)
			}
			if cfg.CrossOriginOpenerPolicy != "" {
				w.Header().Set("Cross-Origin-Opener-Policy", cfg.CrossOriginOpenerPolicy)
			}
			if cfg.CrossOriginResourcePolicy != "" {
				w.Header().Set("Cross-Origin-Resource-Policy", cfg.CrossOriginResourcePolicy)
			}
			if cfg.OriginAgentCluster != "" {
				w.Header().Set("Origin-Agent-Cluster", cfg.OriginAgentCluster)
			}
			if cfg.ReferrerPolicy != "" {
				w.Header().Set("Referrer-Policy", cfg.ReferrerPolicy)
			}
			if cfg.XDNSPrefetchControl != "" {
				w.Header().Set("X-DNS-Prefetch-Control", cfg.XDNSPrefetchControl)
			}
			if cfg.XDownloadOptions != "" {
				w.Header().Set("X-Download-Options", cfg.XDownloadOptions)
			}
			if cfg.XPermittedCrossDomain != "" {
				w.Header().Set("X-Permitted-Cross-Domain-Policies", cfg.XPermittedCrossDomain)
			}

			// Handle HSTS headers.
			if r.TLS != nil && cfg.HSTSMaxAge > 0 {
				subdomains := ""
				if !cfg.HSTSExcludeSubdomains {
					subdomains = "; includeSubDomains"
				}
				if cfg.HSTSPreloadEnabled {
					subdomains += "; preload"
				}
				w.Header().Set("Strict-Transport-Security", fmt.Sprintf("max-age=%d%s", cfg.HSTSMaxAge, subdomains))
			}

			// Handle Content-Security-Policy headers.
			if cfg.ContentSecurityPolicy != "" {
				if cfg.CSPReportOnly {
					w.Header().Set("Content-Security-Policy-Report-Only", cfg.ContentSecurityPolicy)
				} else {
					w.Header().Set("Content-Security-Policy", cfg.ContentSecurityPolicy)
				}
			}

			// Handle Permissions-Policy headers.
			if cfg.PermissionPolicy != "" {
				w.Header().Set("Permissions-Policy", cfg.PermissionPolicy)
			}

			// Proceed to the next handler.
			next.ServeHTTP(w, r)
		})
	}
}

// Helper function to apply default configuration values.
func configDefault(config ...Config) Config {
	if len(config) < 1 {
		return ConfigDefault
	}

	cfg := config[0]

	// Apply default values where not set.
	if cfg.XSSProtection == "" {
		cfg.XSSProtection = ConfigDefault.XSSProtection
	}
	if cfg.ContentTypeNosniff == "" {
		cfg.ContentTypeNosniff = ConfigDefault.ContentTypeNosniff
	}
	if cfg.XFrameOptions == "" {
		cfg.XFrameOptions = ConfigDefault.XFrameOptions
	}
	if cfg.ReferrerPolicy == "" {
		cfg.ReferrerPolicy = ConfigDefault.ReferrerPolicy
	}
	if cfg.CrossOriginEmbedderPolicy == "" {
		cfg.CrossOriginEmbedderPolicy = ConfigDefault.CrossOriginEmbedderPolicy
	}
	if cfg.CrossOriginOpenerPolicy == "" {
		cfg.CrossOriginOpenerPolicy = ConfigDefault.CrossOriginOpenerPolicy
	}
	if cfg.CrossOriginResourcePolicy == "" {
		cfg.CrossOriginResourcePolicy = ConfigDefault.CrossOriginResourcePolicy
	}
	if cfg.OriginAgentCluster == "" {
		cfg.OriginAgentCluster = ConfigDefault.OriginAgentCluster
	}
	if cfg.XDNSPrefetchControl == "" {
		cfg.XDNSPrefetchControl = ConfigDefault.XDNSPrefetchControl
	}
	if cfg.XDownloadOptions == "" {
		cfg.XDownloadOptions = ConfigDefault.XDownloadOptions
	}
	if cfg.XPermittedCrossDomain == "" {
		cfg.XPermittedCrossDomain = ConfigDefault.XPermittedCrossDomain
	}

	return cfg
}

func Cors(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		allowedOrigin := os.Getenv("DOMAIN") // Replace with your origin
		origin := r.Header.Get("Origin")
		if origin != "" && origin != allowedOrigin {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}
		next.ServeHTTP(w, r)
	})
}

type currentUser string

const current_user currentUser = "current_user"

func RequireAuth(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries := models.New(DB)
		ctx := r.Context()

		var token string

		// 1. Check for token in Authorization header
		token = r.Header.Get("auth")

		// 2. If no token in header, check for the session_token cookie
		if token == "" {
			cookie, err := r.Cookie("session_token")
			if err == nil {
				token = cookie.Value // Use token from cookie if available
			}
		}

		// 3. If no token found in either place, return error
		if token == "" {
			http.Error(w, "missing auth token", http.StatusForbidden)
			return
		}

		session, err := queries.SessionRead(ctx, token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if session.CreatedAt.Time.AddDate(0, 0, 30).Unix() < time.Now().Unix() {
			http.Error(w, "session has expired", http.StatusBadRequest)
			return
		}

		user, err := queries.AuthUserRead(ctx, session.UserID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !user.Isactive.Bool {
			http.Error(w, "inactive user", http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, current_user, user) // Store user in context
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

/*func RequireStaff(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries := models.New(DB)
		ctx := context.Background()

		if w.Header().Get("auth") == "" {
			http.Error(w, "missing auth token", http.StatusForbidden)
			return
		}

		session, err := queries.SessionRead(ctx, w.Header().Get("auth"))

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if session.CreatedAt.Time.AddDate(0, 0, 30).Unix() < time.Now().Unix() {
			http.Error(w, "session has expired", http.StatusBadRequest)
			return
		}

		user, err := queries.AuthUserRead(ctx, session.UserID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !user.Isactive.Bool {
			http.Error(w, "invalid user", http.StatusForbidden)
			return
		}

		if !user.Isstaff.Bool {
			http.Error(w, "invalid user", http.StatusForbidden)
			return
		}

		_ = context.WithValue(ctx, current_user, user)

		next.ServeHTTP(w, r)
	})
}*/

func RequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries := models.New(DB)
		ctx := r.Context()

		var token string

		// 1. Check for token in Authorization header
		token = r.Header.Get("auth")

		// 2. If no token in header, check for the session_token cookie
		if token == "" {
			cookie, err := r.Cookie("session_token")
			if err == nil {
				token = cookie.Value // Use token from cookie if available
			}
		}

		// 3. If no token found in either place, return error
		if token == "" {
			http.Error(w, "Missing auth token", http.StatusForbidden)
			return
		}

		session, err := queries.SessionRead(ctx, token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if session.CreatedAt.Time.AddDate(0, 0, 30).Unix() < time.Now().Unix() {
			http.Error(w, "session has expired", http.StatusBadRequest)
			return
		}

		user, err := queries.AuthUserRead(ctx, session.UserID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !user.Isactive.Bool {
			http.Error(w, "invalid user", http.StatusForbidden)
			return
		}

		if !user.Isadmin.Bool {
			http.Error(w, "invalid user", http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, current_user, user) // Store user in context
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func DashRequireAdmin(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		queries := models.New(DB)
		ctx := r.Context()

		var token string

		// 1. Check for token in Authorization header
		token = r.Header.Get("auth")

		// 2. If no token in header, check for the session_token cookie
		if token == "" {
			cookie, err := r.Cookie("session_token")
			if err == nil {
				token = cookie.Value // Use token from cookie if available
			}
		}

		// 3. If no token found in either place, return error
		if token == "" {
			http.Redirect(w, r, "/dash-login", http.StatusSeeOther)
			return
		}

		session, err := queries.SessionRead(ctx, token)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if session.CreatedAt.Time.AddDate(0, 0, 30).Unix() < time.Now().Unix() {
			http.Redirect(w, r, "/dash-login", http.StatusSeeOther)
			return
		}

		user, err := queries.AuthUserRead(ctx, session.UserID)

		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		if !user.Isactive.Bool {
			http.Error(w, "invalid user", http.StatusForbidden)
			return
		}

		if !user.Isadmin.Bool {
			http.Error(w, "invalid user", http.StatusForbidden)
			return
		}

		ctx = context.WithValue(ctx, current_user, user) // Store user in context
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Custom response writer to capture status code
		lrw := &loggingResponseWriter{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(lrw, r)

		// Log details
		log.Printf("%s %s %d %s", r.Method, r.URL.Path, lrw.statusCode, time.Since(start))
	})
}

// loggingResponseWriter captures the response status code
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
}

// Override WriteHeader to capture status code
func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}
