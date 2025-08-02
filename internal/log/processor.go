// Package log provides log processing functionality
// as described in WORKER-PURPOSE.md
package log

import (
	"context"
	"time"
)

// Processor handles log aggregation and processing
type Processor struct {
	aggregator Aggregator
	parser     Parser
	archiver   Archiver
	analyzer   Analyzer
}

// Aggregator interface for log aggregation from Gateway instances
type Aggregator interface {
	Collect(ctx context.Context, source string) ([]LogEntry, error)
	Stream(ctx context.Context, source string) (<-chan LogEntry, error)
}

// Parser interface for log parsing and enrichment
type Parser interface {
	Parse(ctx context.Context, raw string) (*LogEntry, error)
	Enrich(ctx context.Context, entry *LogEntry) error
}

// Archiver interface for log archival and cleanup
type Archiver interface {
	Archive(ctx context.Context, entries []LogEntry) error
	Cleanup(ctx context.Context, before time.Time) error
}

// Analyzer interface for log analysis and pattern detection
type Analyzer interface {
	AnalyzeTraffic(ctx context.Context, entries []LogEntry) (*TrafficAnalysis, error)
	DetectThreats(ctx context.Context, entries []LogEntry) ([]ThreatAlert, error)
}

// LogEntry represents a structured log entry
type LogEntry struct {
	Timestamp   time.Time         `json:"timestamp"`
	Level       string            `json:"level"`
	Source      string            `json:"source"`
	Message     string            `json:"message"`
	RequestID   string            `json:"request_id"`
	Method      string            `json:"method"`
	Path        string            `json:"path"`
	StatusCode  int               `json:"status_code"`
	ResponseTime time.Duration    `json:"response_time"`
	ClientIP    string            `json:"client_ip"`
	UserAgent   string            `json:"user_agent"`
	Headers     map[string]string `json:"headers"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// TrafficAnalysis represents traffic pattern analysis results
type TrafficAnalysis struct {
	Period       time.Duration `json:"period"`
	TotalRequests int          `json:"total_requests"`
	UniqueIPs    int          `json:"unique_ips"`
	TopPaths     []PathStat   `json:"top_paths"`
	StatusCodes  []StatusStat `json:"status_codes"`
	ResponseTimes ResponseTimeStats `json:"response_times"`
}

// PathStat represents statistics for a specific path
type PathStat struct {
	Path   string `json:"path"`
	Count  int    `json:"count"`
	AvgTime time.Duration `json:"avg_time"`
}

// StatusStat represents statistics for HTTP status codes
type StatusStat struct {
	Code  int `json:"code"`
	Count int `json:"count"`
}

// ResponseTimeStats represents response time statistics
type ResponseTimeStats struct {
	Mean time.Duration `json:"mean"`
	P50  time.Duration `json:"p50"`
	P95  time.Duration `json:"p95"`
	P99  time.Duration `json:"p99"`
}

// ThreatAlert represents a detected security threat
type ThreatAlert struct {
	Type        string    `json:"type"`
	Severity    string    `json:"severity"`
	Description string    `json:"description"`
	Source      string    `json:"source"`
	Timestamp   time.Time `json:"timestamp"`
	Details     map[string]interface{} `json:"details"`
}

// New creates a new log processor
func New(aggregator Aggregator, parser Parser, archiver Archiver, analyzer Analyzer) *Processor {
	return &Processor{
		aggregator: aggregator,
		parser:     parser,
		archiver:   archiver,
		analyzer:   analyzer,
	}
}

// ProcessLogs processes logs from Gateway instances
func (p *Processor) ProcessLogs(ctx context.Context, source string) error {
	// TODO: Implement log processing pipeline
	// 1. Collect logs from Gateway instances
	// 2. Parse and enrich log entries
	// 3. Analyze for patterns and threats
	// 4. Archive processed logs
	return nil
}

// AnalyzeTraffic analyzes traffic patterns from logs
func (p *Processor) AnalyzeTraffic(ctx context.Context, period time.Duration) (*TrafficAnalysis, error) {
	// TODO: Implement traffic analysis
	return nil, nil
}

// DetectThreats analyzes logs for security threats
func (p *Processor) DetectThreats(ctx context.Context) ([]ThreatAlert, error) {
	// TODO: Implement threat detection
	return nil, nil
}