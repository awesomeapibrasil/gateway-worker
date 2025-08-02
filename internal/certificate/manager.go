// Package certificate provides certificate management functionality
// as described in WORKER-PURPOSE.md
package certificate

import (
	"context"
	"time"
)

// Manager handles certificate lifecycle management
type Manager struct {
	acmeClient   ACMEClient
	storage      Storage
	distributor  Distributor
	validator    Validator
}

// ACMEClient interface for ACME protocol operations (Let's Encrypt, etc.)
type ACMEClient interface {
	RenewCertificate(ctx context.Context, domain string) (*Certificate, error)
	ValidateCertificate(ctx context.Context, cert *Certificate) error
}

// Storage interface for certificate persistence
type Storage interface {
	Store(ctx context.Context, cert *Certificate) error
	Retrieve(ctx context.Context, domain string) (*Certificate, error)
	List(ctx context.Context) ([]*Certificate, error)
	Delete(ctx context.Context, domain string) error
}

// Distributor interface for certificate distribution to Gateway instances
type Distributor interface {
	Deploy(ctx context.Context, cert *Certificate, instances []string) error
	DeployTemporary(ctx context.Context, cert *TemporaryCertificate) error
}

// Validator interface for certificate validation
type Validator interface {
	ValidateCertificate(ctx context.Context, cert *Certificate) error
	CheckExpiration(ctx context.Context, cert *Certificate) (time.Duration, error)
}

// Certificate represents a TLS certificate
type Certificate struct {
	Domain      string
	CertData    []byte
	PrivateKey  []byte
	Expiry      time.Time
	Type        CertificateType
	Created     time.Time
	LastUpdated time.Time
}

// TemporaryCertificate represents a temporary self-signed certificate
type TemporaryCertificate struct {
	*Certificate
	Reason       string
	ValidityDays int
	IsTemporary  bool
}

// CertificateType represents the type of certificate
type CertificateType string

const (
	CertificateTypeProduction CertificateType = "production"
	CertificateTypeTemporary  CertificateType = "temporary"
	CertificateTypeStaging    CertificateType = "staging"
)

// New creates a new certificate manager
func New(acmeClient ACMEClient, storage Storage, distributor Distributor, validator Validator) *Manager {
	return &Manager{
		acmeClient:  acmeClient,
		storage:     storage,
		distributor: distributor,
		validator:   validator,
	}
}

// MonitorCertificates starts certificate monitoring for expiration
func (m *Manager) MonitorCertificates(ctx context.Context) error {
	ticker := time.NewTicker(24 * time.Hour) // Check daily
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-ticker.C:
			if err := m.checkAndRenewCertificates(ctx); err != nil {
				// Log error but continue monitoring
				continue
			}
		}
	}
}

// checkAndRenewCertificates checks all certificates for expiration and renews if needed
func (m *Manager) checkAndRenewCertificates(ctx context.Context) error {
	certificates, err := m.storage.List(ctx)
	if err != nil {
		return err
	}

	for _, cert := range certificates {
		// Check if certificate expires within 30 days
		if time.Until(cert.Expiry) <= 30*24*time.Hour {
			if err := m.renewCertificate(ctx, cert.Domain); err != nil {
				// If renewal fails, deploy temporary certificate
				if err := m.deployTemporaryCertificate(ctx, cert.Domain, "renewal_failed"); err != nil {
					// Log error for temporary certificate deployment failure
				}
			}
		}
	}

	return nil
}

// renewCertificate renews a certificate using ACME protocol
func (m *Manager) renewCertificate(ctx context.Context, domain string) error {
	// Renew certificate
	newCert, err := m.acmeClient.RenewCertificate(ctx, domain)
	if err != nil {
		return err
	}

	// Validate new certificate
	if err := m.validator.ValidateCertificate(ctx, newCert); err != nil {
		return err
	}

	// Store new certificate
	if err := m.storage.Store(ctx, newCert); err != nil {
		return err
	}

	// Deploy to Gateway instances
	// TODO: Get Gateway instances from configuration
	instances := []string{} // Placeholder
	return m.distributor.Deploy(ctx, newCert, instances)
}

// deployTemporaryCertificate creates and deploys a temporary certificate
func (m *Manager) deployTemporaryCertificate(ctx context.Context, domain, reason string) error {
	tempCert := &TemporaryCertificate{
		Certificate: &Certificate{
			Domain:      domain,
			Type:        CertificateTypeTemporary,
			Created:     time.Now(),
			LastUpdated: time.Now(),
			Expiry:      time.Now().Add(14 * 24 * time.Hour), // 14 days validity
		},
		Reason:       reason,
		ValidityDays: 14,
		IsTemporary:  true,
	}

	// Generate self-signed certificate
	if err := m.generateSelfSignedCertificate(tempCert); err != nil {
		return err
	}

	// Deploy temporary certificate
	return m.distributor.DeployTemporary(ctx, tempCert)
}

// generateSelfSignedCertificate generates a self-signed certificate for temporary use
func (m *Manager) generateSelfSignedCertificate(tempCert *TemporaryCertificate) error {
	// TODO: Implement self-signed certificate generation
	// This should create a certificate with:
	// - Same subject name as original certificate
	// - Extended validity period (14 days maximum)
	// - Clear marking as temporary certificate
	return nil
}

// GetCertificateStatus returns the status of a certificate
func (m *Manager) GetCertificateStatus(ctx context.Context, domain string) (*CertificateStatus, error) {
	cert, err := m.storage.Retrieve(ctx, domain)
	if err != nil {
		return nil, err
	}

	timeUntilExpiry, err := m.validator.CheckExpiration(ctx, cert)
	if err != nil {
		return nil, err
	}

	status := &CertificateStatus{
		Domain:          domain,
		Type:            cert.Type,
		Expiry:          cert.Expiry,
		TimeUntilExpiry: timeUntilExpiry,
		Status:          m.determineStatus(timeUntilExpiry),
	}

	return status, nil
}

// CertificateStatus represents the current status of a certificate
type CertificateStatus struct {
	Domain          string
	Type            CertificateType
	Expiry          time.Time
	TimeUntilExpiry time.Duration
	Status          string
}

// determineStatus determines the certificate status based on expiry time
func (m *Manager) determineStatus(timeUntilExpiry time.Duration) string {
	if timeUntilExpiry <= 0 {
		return "expired"
	}
	if timeUntilExpiry <= 7*24*time.Hour {
		return "expiring"
	}
	if timeUntilExpiry <= 30*24*time.Hour {
		return "renewal_due"
	}
	return "valid"
}