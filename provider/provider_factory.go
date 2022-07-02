package provider

import (
	"fmt"
)

// speed test types
const (
	testDownload = "download"
	testUpload   = "upload"
)

// Provider is a contract for running speed test.
type Provider interface {

	// RunSpeedTest runs a speed test and returns its result (download and upload speed in Mbps).
	RunSpeedTest() (*Result, error)
}

// Result holds download and upload speed measured by a speed test.
type Result struct {
	Download Speed `json:"download"`
	Upload   Speed `json:"upload"`
}

// Speed holds speed value.
type Speed float64

// String returns a formated string representation of Result.
func (r Result) String() string {
	return fmt.Sprintf("%8s: %s\n%8s: %s",
		"download", r.Download,
		"upload", r.Upload)
}

const defaultSpeedUnit = "Mbps"

// String returns a formated string representation of Speed using the default speed unit (Mbps).
//
// format: 10.123456 Mbps
func (s Speed) String() string {
	return fmt.Sprintf("%0.6f %s", s, defaultSpeedUnit)
}

// implemented providers
const (
	ProviderNetxlix = "netflix"
	ProviderOokla   = "ookla"
)

// ProviderFactory produces an exact Provider implementation.
type ProviderFactory struct{}

// GetProvider returns an exact provider for a given provider name.
//
// In case of an absence of the exact implementation NoopProvider returns.
func (f ProviderFactory) GetProvider(name string) (Provider, error) {
	switch name {
	case ProviderNetxlix:
		return NewNetflix()
	case ProviderOokla:
		return NewOokla()
	default:
		return NewNoopProvider()
	}
}

// NoopProvider is a no-operation implementation of Provider.
type NoopProvider struct{}

// NewNoopProvider returns a new instance of NoopProvider.
func NewNoopProvider() (*NoopProvider, error) {
	return &NoopProvider{}, nil
}

// RunSpeedTest doesn't perform any test and returns a nil result (without errors).
func (n *NoopProvider) RunSpeedTest() (*Result, error) {
	return nil, nil
}
