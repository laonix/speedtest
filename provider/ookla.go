package provider

import (
	"fmt"

	ookla "github.com/showwin/speedtest-go/speedtest"
	"golang.org/x/sync/errgroup"
)

// Ookla is an implementation of Provider running speed tests via https://www.speedtest.net/.
type Ookla struct {
	target *ookla.Server
}

// NewOokla returns a new instance of Ookla.
func NewOokla() (*Ookla, error) {
	target, err := getTargetServer()
	if err != nil {
		return nil, fmt.Errorf("ookla target server: %s", err)
	}

	return &Ookla{target: target}, nil
}

// getTargetServer returns only the first server from a fetched available speed test servers list.
func getTargetServer() (*ookla.Server, error) {
	user, err := ookla.FetchUserInfo()
	if err != nil {
		return nil, fmt.Errorf("fetch user: %s", err)
	}

	serverList, err := ookla.FetchServers(user)
	if err != nil {
		return nil, fmt.Errorf("fetch servers: %s", err)
	}

	// we are using only the first server
	return serverList[0], nil
}

// RunSpeedTest performs a speed test via https://www.speedtest.net/
// and returns download and upload speed within a result.
func (o *Ookla) RunSpeedTest() (*Result, error) {
	eg := errgroup.Group{}

	eg.Go(func() error {
		return o.test(testDownload, o.target.DownloadTest)
	})

	eg.Go(func() error {
		return o.test(testUpload, o.target.UploadTest)
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return o.result(), nil
}

func (o *Ookla) test(testType string, testFunc func(bool) error) error {
	if err := testFunc(false); err != nil {
		return fmt.Errorf("%s test: %s", testType, err)
	}
	return nil
}

func (o *Ookla) result() *Result {
	return &Result{
		Download: Speed(o.target.DLSpeed),
		Upload:   Speed(o.target.ULSpeed),
	}
}
