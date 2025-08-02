// Package config provides configuration management functionality
// as described in WORKER-PURPOSE.md
package config

import (
	"context"
	"encoding/json"
	"time"
)

// Manager handles configuration lifecycle management
type Manager struct {
	storage    Storage
	validator  Validator
	distributor Distributor
}

// Storage interface for configuration persistence
type Storage interface {
	Store(ctx context.Context, config *Configuration) error
	Retrieve(ctx context.Context, configType, version string) (*Configuration, error)
	List(ctx context.Context, configType string) ([]*Configuration, error)
	GetLatest(ctx context.Context, configType string) (*Configuration, error)
}

// Validator interface for configuration validation
type Validator interface {
	ValidateWAFRules(ctx context.Context, rules []WAFRule) []ValidationError
	ValidateRoutingConfig(ctx context.Context, config RoutingConfig) []ValidationError
	ValidateBackendConfig(ctx context.Context, config BackendConfig) []ValidationError
	ValidateSecurityPolicy(ctx context.Context, policy SecurityPolicy) []ValidationError
}

// Distributor interface for configuration distribution to Gateway instances
type Distributor interface {
	DeployWAFRules(ctx context.Context, rules []WAFRule, instances []string) error
	DeployRoutingConfig(ctx context.Context, config RoutingConfig, instances []string) error
	DeployBackendConfig(ctx context.Context, config BackendConfig, instances []string) error
	DeploySecurityPolicy(ctx context.Context, policy SecurityPolicy, instances []string) error
}

// Configuration represents a configuration item
type Configuration struct {
	Type        ConfigurationType `json:"type"`
	Version     string           `json:"version"`
	Data        json.RawMessage  `json:"data"`
	Created     time.Time        `json:"created"`
	LastUpdated time.Time        `json:"last_updated"`
	Active      bool             `json:"active"`
}

// ConfigurationType represents the type of configuration
type ConfigurationType string

const (
	ConfigurationTypeWAF      ConfigurationType = "waf"
	ConfigurationTypeRouting  ConfigurationType = "routing"
	ConfigurationTypeBackend  ConfigurationType = "backend"
	ConfigurationTypeSecurity ConfigurationType = "security"
)

// WAFRule represents a Web Application Firewall rule
type WAFRule struct {
	ID          string    `json:"id"`
	Name        string    `json:"name"`
	Pattern     string    `json:"pattern"`
	Action      string    `json:"action"` // "block", "log", "rate_limit"
	Priority    int       `json:"priority"`
	Enabled     bool      `json:"enabled"`
	Created     time.Time `json:"created"`
	LastUpdated time.Time `json:"last_updated"`
}

// RoutingConfig represents routing configuration
type RoutingConfig struct {
	Routes      []Route           `json:"routes"`
	DefaultBackend string         `json:"default_backend"`
	Middleware  []MiddlewareConfig `json:"middleware"`
}

// Route represents a routing rule
type Route struct {
	Path        string            `json:"path"`
	Method      string            `json:"method"`
	Backend     string            `json:"backend"`
	Headers     map[string]string `json:"headers"`
	Priority    int               `json:"priority"`
	Enabled     bool              `json:"enabled"`
}

// MiddlewareConfig represents middleware configuration
type MiddlewareConfig struct {
	Name    string                 `json:"name"`
	Config  map[string]interface{} `json:"config"`
	Enabled bool                   `json:"enabled"`
}

// BackendConfig represents backend service configuration
type BackendConfig struct {
	Services     []BackendService `json:"services"`
	HealthChecks HealthCheckConfig `json:"health_checks"`
	LoadBalancing LoadBalancingConfig `json:"load_balancing"`
}

// BackendService represents a backend service
type BackendService struct {
	Name     string   `json:"name"`
	Hosts    []string `json:"hosts"`
	Port     int      `json:"port"`
	Protocol string   `json:"protocol"`
	Weight   int      `json:"weight"`
	Enabled  bool     `json:"enabled"`
}

// HealthCheckConfig represents health check configuration
type HealthCheckConfig struct {
	Enabled  bool          `json:"enabled"`
	Path     string        `json:"path"`
	Interval time.Duration `json:"interval"`
	Timeout  time.Duration `json:"timeout"`
	Retries  int           `json:"retries"`
}

// LoadBalancingConfig represents load balancing configuration
type LoadBalancingConfig struct {
	Algorithm string                 `json:"algorithm"` // "round_robin", "least_connections", "ip_hash"
	Config    map[string]interface{} `json:"config"`
}

// SecurityPolicy represents security policy configuration
type SecurityPolicy struct {
	RateLimiting   RateLimitingConfig `json:"rate_limiting"`
	Authentication AuthConfig         `json:"authentication"`
	CORS           CORSConfig         `json:"cors"`
	Headers        HeadersConfig      `json:"headers"`
}

// RateLimitingConfig represents rate limiting configuration
type RateLimitingConfig struct {
	Enabled bool          `json:"enabled"`
	Rate    int           `json:"rate"`
	Window  time.Duration `json:"window"`
	Burst   int           `json:"burst"`
}

// AuthConfig represents authentication configuration
type AuthConfig struct {
	Enabled  bool     `json:"enabled"`
	Type     string   `json:"type"` // "jwt", "oauth2", "basic"
	Providers []string `json:"providers"`
}

// CORSConfig represents CORS configuration
type CORSConfig struct {
	Enabled         bool     `json:"enabled"`
	AllowedOrigins  []string `json:"allowed_origins"`
	AllowedMethods  []string `json:"allowed_methods"`
	AllowedHeaders  []string `json:"allowed_headers"`
	ExposedHeaders  []string `json:"exposed_headers"`
	AllowCredentials bool    `json:"allow_credentials"`
	MaxAge          int      `json:"max_age"`
}

// HeadersConfig represents security headers configuration
type HeadersConfig struct {
	HSTS            HSTSConfig `json:"hsts"`
	ContentSecurity CSPConfig  `json:"content_security"`
	XFrameOptions   string     `json:"x_frame_options"`
	XContentType    string     `json:"x_content_type"`
}

// HSTSConfig represents HSTS configuration
type HSTSConfig struct {
	Enabled           bool `json:"enabled"`
	MaxAge            int  `json:"max_age"`
	IncludeSubdomains bool `json:"include_subdomains"`
	Preload           bool `json:"preload"`
}

// CSPConfig represents Content Security Policy configuration
type CSPConfig struct {
	Enabled   bool   `json:"enabled"`
	Policy    string `json:"policy"`
	ReportOnly bool  `json:"report_only"`
}

// ValidationError represents a configuration validation error
type ValidationError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
	Code    string `json:"code"`
}

// New creates a new configuration manager
func New(storage Storage, validator Validator, distributor Distributor) *Manager {
	return &Manager{
		storage:    storage,
		validator:  validator,
		distributor: distributor,
	}
}

// UpdateWAFRules updates WAF rules with validation and deployment
func (m *Manager) UpdateWAFRules(ctx context.Context, rules []WAFRule, emergencyDeployment bool) error {
	// Validate rules
	if !emergencyDeployment {
		if errors := m.validator.ValidateWAFRules(ctx, rules); len(errors) > 0 {
			return NewValidationError(errors)
		}
	}

	// Create configuration
	rulesData, err := json.Marshal(rules)
	if err != nil {
		return err
	}

	config := &Configuration{
		Type:        ConfigurationTypeWAF,
		Version:     generateVersion(),
		Data:        rulesData,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		Active:      true,
	}

	// Store configuration
	if err := m.storage.Store(ctx, config); err != nil {
		return err
	}

	// Deploy to Gateway instances
	// TODO: Get Gateway instances from configuration
	instances := []string{} // Placeholder
	return m.distributor.DeployWAFRules(ctx, rules, instances)
}

// UpdateRoutingConfig updates routing configuration
func (m *Manager) UpdateRoutingConfig(ctx context.Context, routingConfig RoutingConfig) error {
	// Validate configuration
	if errors := m.validator.ValidateRoutingConfig(ctx, routingConfig); len(errors) > 0 {
		return NewValidationError(errors)
	}

	// Create configuration
	configData, err := json.Marshal(routingConfig)
	if err != nil {
		return err
	}

	config := &Configuration{
		Type:        ConfigurationTypeRouting,
		Version:     generateVersion(),
		Data:        configData,
		Created:     time.Now(),
		LastUpdated: time.Now(),
		Active:      true,
	}

	// Store configuration
	if err := m.storage.Store(ctx, config); err != nil {
		return err
	}

	// Deploy to Gateway instances
	instances := []string{} // Placeholder
	return m.distributor.DeployRoutingConfig(ctx, routingConfig, instances)
}

// GetConfiguration retrieves a configuration
func (m *Manager) GetConfiguration(ctx context.Context, configType ConfigurationType, version string) (*Configuration, error) {
	if version == "" {
		return m.storage.GetLatest(ctx, string(configType))
	}
	return m.storage.Retrieve(ctx, string(configType), version)
}

// generateVersion generates a version string based on timestamp
func generateVersion() string {
	return time.Now().Format("20060102150405")
}

// NewValidationError creates a validation error
func NewValidationError(errors []ValidationError) error {
	// TODO: Implement proper validation error type
	return nil
}