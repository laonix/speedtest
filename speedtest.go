package speedtest

import (
	"fmt"
	prov "github.com/laonix/speedtest/provider"
)

func Run(provider string) (*prov.Result, error) {
	p, err := prov.GetProvider(provider)
	if err != nil {
		return nil, fmt.Errorf("get provider: %s", err)
	}

	return p.RunSpeedTest()
}
