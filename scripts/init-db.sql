-- Gateway Worker Database Initialization Script
-- Based on WORKER-PURPOSE.md specifications

-- Enable necessary extensions
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_stat_statements";

-- Create schemas
CREATE SCHEMA IF NOT EXISTS certificates;
CREATE SCHEMA IF NOT EXISTS configurations;
CREATE SCHEMA IF NOT EXISTS logs;
CREATE SCHEMA IF NOT EXISTS analytics;
CREATE SCHEMA IF NOT EXISTS jobs;
CREATE SCHEMA IF NOT EXISTS integrations;

-- Certificates schema tables
CREATE TABLE IF NOT EXISTS certificates.certificates (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    domain VARCHAR(255) NOT NULL UNIQUE,
    cert_data BYTEA NOT NULL,
    private_key BYTEA NOT NULL,
    cert_type VARCHAR(50) NOT NULL DEFAULT 'production',
    expiry_date TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true
);

CREATE TABLE IF NOT EXISTS certificates.certificate_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    certificate_id UUID REFERENCES certificates.certificates(id),
    action VARCHAR(50) NOT NULL,
    details JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Configurations schema tables
CREATE TABLE IF NOT EXISTS configurations.configurations (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    config_type VARCHAR(50) NOT NULL,
    version VARCHAR(50) NOT NULL,
    config_data JSONB NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    is_active BOOLEAN DEFAULT true,
    UNIQUE(config_type, version)
);

CREATE TABLE IF NOT EXISTS configurations.waf_rules (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    rule_id VARCHAR(100) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    pattern TEXT NOT NULL,
    action VARCHAR(50) NOT NULL,
    priority INTEGER NOT NULL DEFAULT 1000,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Jobs schema tables
CREATE TABLE IF NOT EXISTS jobs.jobs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_type VARCHAR(100) NOT NULL,
    status VARCHAR(50) NOT NULL DEFAULT 'queued',
    payload JSONB NOT NULL,
    priority INTEGER NOT NULL DEFAULT 1000,
    retry_count INTEGER DEFAULT 0,
    max_retries INTEGER DEFAULT 3,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    started_at TIMESTAMP WITH TIME ZONE,
    completed_at TIMESTAMP WITH TIME ZONE,
    error_message TEXT
);

CREATE TABLE IF NOT EXISTS jobs.job_history (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    job_id UUID REFERENCES jobs.jobs(id),
    status VARCHAR(50) NOT NULL,
    message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Logs schema tables
CREATE TABLE IF NOT EXISTS logs.access_logs (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    source VARCHAR(100) NOT NULL,
    request_id VARCHAR(100),
    method VARCHAR(10),
    path TEXT,
    status_code INTEGER,
    response_time_ms INTEGER,
    client_ip INET,
    user_agent TEXT,
    headers JSONB,
    metadata JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS logs.security_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(100) NOT NULL,
    severity VARCHAR(20) NOT NULL,
    source VARCHAR(100) NOT NULL,
    client_ip INET,
    description TEXT,
    details JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Analytics schema tables
CREATE TABLE IF NOT EXISTS analytics.metrics (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    metric_name VARCHAR(100) NOT NULL,
    metric_value DECIMAL NOT NULL,
    metric_unit VARCHAR(20),
    source VARCHAR(100) NOT NULL,
    tags JSONB,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS analytics.traffic_summaries (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    period_start TIMESTAMP WITH TIME ZONE NOT NULL,
    period_end TIMESTAMP WITH TIME ZONE NOT NULL,
    total_requests BIGINT NOT NULL,
    unique_ips INTEGER NOT NULL,
    avg_response_time_ms DECIMAL,
    error_rate DECIMAL,
    top_paths JSONB,
    status_codes JSONB,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Integrations schema tables
CREATE TABLE IF NOT EXISTS integrations.external_apis (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    name VARCHAR(100) NOT NULL UNIQUE,
    url TEXT NOT NULL,
    config JSONB NOT NULL,
    enabled BOOLEAN DEFAULT true,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

CREATE TABLE IF NOT EXISTS integrations.api_calls (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    api_id UUID REFERENCES integrations.external_apis(id),
    request_data JSONB,
    response_data JSONB,
    status_code INTEGER,
    duration_ms INTEGER,
    success BOOLEAN,
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create indexes for performance
CREATE INDEX IF NOT EXISTS idx_certificates_domain ON certificates.certificates(domain);
CREATE INDEX IF NOT EXISTS idx_certificates_expiry ON certificates.certificates(expiry_date);
CREATE INDEX IF NOT EXISTS idx_certificates_active ON certificates.certificates(is_active);

CREATE INDEX IF NOT EXISTS idx_configurations_type_version ON configurations.configurations(config_type, version);
CREATE INDEX IF NOT EXISTS idx_configurations_active ON configurations.configurations(is_active);

CREATE INDEX IF NOT EXISTS idx_waf_rules_enabled ON configurations.waf_rules(enabled);
CREATE INDEX IF NOT EXISTS idx_waf_rules_priority ON configurations.waf_rules(priority);

CREATE INDEX IF NOT EXISTS idx_jobs_status ON jobs.jobs(status);
CREATE INDEX IF NOT EXISTS idx_jobs_type ON jobs.jobs(job_type);
CREATE INDEX IF NOT EXISTS idx_jobs_created ON jobs.jobs(created_at);

CREATE INDEX IF NOT EXISTS idx_access_logs_timestamp ON logs.access_logs(timestamp);
CREATE INDEX IF NOT EXISTS idx_access_logs_source ON logs.access_logs(source);
CREATE INDEX IF NOT EXISTS idx_access_logs_client_ip ON logs.access_logs(client_ip);

CREATE INDEX IF NOT EXISTS idx_security_events_timestamp ON logs.security_events(timestamp);
CREATE INDEX IF NOT EXISTS idx_security_events_type ON logs.security_events(event_type);
CREATE INDEX IF NOT EXISTS idx_security_events_severity ON logs.security_events(severity);

CREATE INDEX IF NOT EXISTS idx_metrics_name_timestamp ON analytics.metrics(metric_name, timestamp);
CREATE INDEX IF NOT EXISTS idx_metrics_source ON analytics.metrics(source);

CREATE INDEX IF NOT EXISTS idx_traffic_summaries_period ON analytics.traffic_summaries(period_start, period_end);

CREATE INDEX IF NOT EXISTS idx_api_calls_api_id ON integrations.api_calls(api_id);
CREATE INDEX IF NOT EXISTS idx_api_calls_created ON integrations.api_calls(created_at);

-- Create triggers for updated_at timestamps
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ language 'plpgsql';

CREATE TRIGGER update_certificates_updated_at BEFORE UPDATE ON certificates.certificates FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_configurations_updated_at BEFORE UPDATE ON configurations.configurations FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_waf_rules_updated_at BEFORE UPDATE ON configurations.waf_rules FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_external_apis_updated_at BEFORE UPDATE ON integrations.external_apis FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- Insert default data
INSERT INTO configurations.configurations (config_type, version, config_data, is_active) 
VALUES 
    ('waf', '1.0.0', '{"rules": [], "enabled": true}', true),
    ('routing', '1.0.0', '{"routes": [], "default_backend": ""}', true),
    ('security', '1.0.0', '{"rate_limiting": {"enabled": false}, "cors": {"enabled": false}}', true)
ON CONFLICT (config_type, version) DO NOTHING;

-- Grant permissions (adjust as needed for your security requirements)
GRANT USAGE ON SCHEMA certificates TO worker;
GRANT USAGE ON SCHEMA configurations TO worker;
GRANT USAGE ON SCHEMA logs TO worker;
GRANT USAGE ON SCHEMA analytics TO worker;
GRANT USAGE ON SCHEMA jobs TO worker;
GRANT USAGE ON SCHEMA integrations TO worker;

GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA certificates TO worker;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA configurations TO worker;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA logs TO worker;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA analytics TO worker;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA jobs TO worker;
GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA integrations TO worker;

GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA certificates TO worker;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA configurations TO worker;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA logs TO worker;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA analytics TO worker;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA jobs TO worker;
GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA integrations TO worker;