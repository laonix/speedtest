package provider

import (
	"fmt"
)

const (
	testDownload = "download"
	testUpload   = "upload"
)

type Provider interface {
	RunSpeedTest() (*Result, error)
}

type Result struct {
	Download Speed `json:"download"`
	Upload   Speed `json:"upload"`
}

type Speed float64

func (r Result) String() string {
	return fmt.Sprintf("%8s: %s\n%8s: %s",
		"download", r.Download,
		"upload", r.Upload)
}

const defaultSpeedUnit = "Mbps"

func (s Speed) String() string {
	return fmt.Sprintf("%0.6f %s", s, defaultSpeedUnit)
}

const (
	ProviderNetxlix = "netflix"
	ProviderOokla   = "ookla"
)

func GetProvider(name string) (Provider, error) {
	switch name {
	case ProviderNetxlix:
		return NewNetflix()
	case ProviderOokla:
		return NewOokla()
	default:
		return NewNoopProvider()
	}
}

type NoopProvider struct{}

func NewNoopProvider() (*NoopProvider, error) {
	return &NoopProvider{}, nil
}

func (n *NoopProvider) RunSpeedTest() (*Result, error) {
	return nil, nil
}
