package speedtest

import (
	"fmt"
	prov "github.com/laonix/speedtest/provider"
)

// SpeedtestRunner is a contract for performing network speed test.
type SpeedtestRunner interface {

	// Run performs a speed test via designated provider.
	Run(provider string) (*prov.Result, error)
}

// DefaultRunner is a default implementation of SpeedtestRunner.
type DefaultRunner struct {
	f prov.ProviderFactory
}

// NewRunner returns a new instance of DefaultRunner.
func NewRunner() *DefaultRunner {
	return &DefaultRunner{f: prov.ProviderFactory{}}
}

// Run performs a speed test via designated provider.
//
// In this default implementation provider should be:
//
// - 'ookla' for https://www.speedtest.net/
//
// - 'netflix' for https://fast.com/.
func (r *DefaultRunner) Run(provider string) (*prov.Result, error) {
	p, err := r.f.GetProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("get provider: %s", err)
	}

	return p.RunSpeedTest()
}
