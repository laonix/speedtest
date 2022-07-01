package provider

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	netflix "github.com/gesquive/fast-cli/fast"
	"golang.org/x/sync/errgroup"
)

type Netflix struct {
	url    string
	client *http.Client
	ds     float64
	us     float64
}

func NewNetflix() (*Netflix, error) {
	targetUrl, err := getTargetUrl()
	if err != nil {
		return nil, fmt.Errorf("netflix target url: %s", err)
	}

	return &Netflix{
		url:    targetUrl,
		client: &http.Client{},
	}, nil
}

func getTargetUrl() (string, error) {
	urls := netflix.GetDlUrls(1)
	if len(urls) == 0 {
		return "", errors.New("cannot get a target url")
	}

	return urls[0], nil
}

func (n *Netflix) RunSpeedTest() (*Result, error) {
	eg := errgroup.Group{}

	eg.Go(func() error {
		s, err := n.test(testDownload, n.download)
		n.ds = s
		return err
	})

	eg.Go(func() error {
		s, err := n.test(testUpload, n.upload)
		n.us = s
		return err
	})

	if err := eg.Wait(); err != nil {
		return nil, err
	}

	return n.result(), nil
}

func (n *Netflix) result() *Result {
	return &Result{
		Download: Speed(n.ds),
		Upload:   Speed(n.us),
	}
}

const (
	workload  = 10
	payloadMB = 25.0
)

func (n *Netflix) test(testType string, op func() error) (float64, error) {
	eg := errgroup.Group{}

	start := time.Now()
	for i := 0; i < workload; i++ {
		eg.Go(func() error {
			return op()
		})
	}
	if err := eg.Wait(); err != nil {
		return 0, fmt.Errorf("%s: %s", testType, err)
	}
	finish := time.Now()

	return speed(start, finish), nil
}

func speed(start, finish time.Time) float64 {
	return payloadMB * 8.0 * float64(workload) / finish.Sub(start).Seconds()
}

func (n *Netflix) download() error {
	resp, err := n.client.Get(n.url)
	if err != nil {
		return fmt.Errorf("download: %s", err)
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	_, _ = io.ReadAll(resp.Body)

	return nil
}

func (n *Netflix) upload() error {
	values := url.Values{}
	values.Add("content", payload())

	resp, err := n.client.PostForm(n.url, values)
	if err != nil {
		return fmt.Errorf("upload: %s", err)
	}
	defer func(body io.ReadCloser) {
		_ = body.Close()
	}(resp.Body)

	_, _ = io.ReadAll(resp.Body)

	return nil
}

func payload() string {
	return strings.Repeat("0", payloadMB*1024*1024)
}
