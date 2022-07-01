package provider

import (
	"fmt"

	ookla "github.com/showwin/speedtest-go/speedtest"
	"golang.org/x/sync/errgroup"
)

type Ookla struct {
	target *ookla.Server
}

func NewOokla() (*Ookla, error) {
	target, err := getTargetServer()
	if err != nil {
		return nil, fmt.Errorf("ookla target server: %s", err)
	}

	return &Ookla{target: target}, nil
}

func getTargetServer() (*ookla.Server, error) {
	user, err := ookla.FetchUserInfo()
	if err != nil {
		return nil, fmt.Errorf("fetch user: %s", err)
	}

	serverList, err := ookla.FetchServers(user)
	if err != nil {
		return nil, fmt.Errorf("fetch servers: %s", err)
	}

	return serverList[0], nil
}

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
