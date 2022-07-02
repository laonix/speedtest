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

// Netflix is an implementation of Provider running speed tests via https://fast.com/.
type Netflix struct {
	url string
	ds  float64
	us  float64

	client *http.Client
}

// NewNetflix returns a new instance of Netflix.
func NewNetflix() (*Netflix, error) {
	targetUrl, err := getTargetUrl()
	if err != nil {
		return nil, fmt.Errorf("netflix target url: %s", err)
	}

	return &Netflix{
		url:    targetUrl,
		client: http.DefaultClient,
	}, nil
}

// getTargetUrl returns only the first speed test server url from a fetched list of available servers.
func getTargetUrl() (string, error) {
	urls := netflix.GetDlUrls(1)
	if len(urls) == 0 {
		return "", errors.New("cannot get a target url")
	}

	// we are using only the first provided url
	return urls[0], nil
}

// RunSpeedTest performs a speed test via https://fast.com/
// and returns download and upload speed within a result.
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
	workload  = 10   // number of single operations performed within a test
	payloadMB = 25.0 // as long as downloaded payload is 25 MB let's use the same payload size for an upload test
)

// test performs a number of a given operations, tracks their execution time and returns a calculates speed.
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

// speed calculates a speed by the given test duration (start and finish time) assuming than the payload size is defined.
func speed(start, finish time.Time) float64 {
	return payloadMB * 8.0 * float64(workload) / finish.Sub(start).Seconds() // Mbps
}

func (n *Netflix) download() error {
	resp, err := n.client.Get(n.url)
	if err != nil {
		return fmt.Errorf("download: %s", err)
	}
	defer resp.Body.Close()

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
	defer resp.Body.Close()

	_, _ = io.ReadAll(resp.Body)

	return nil
}

// payload returns a string of payloadMB MB size.
func payload() string {
	return strings.Repeat("0", payloadMB*1024*1024)
}
