package goutils

import (
	"bytes"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"io"
	"net/http"
	"time"
)

type Interceptor func(req *http.Request) error

type Response struct {
	RequestUri string
	Status     int
	Headers    http.Header
	Body       []byte
	Raw        *http.Response
	Duration   time.Duration
}

type ClientConfig struct {
	Timeout         time.Duration
	SkipVerifySSL   bool
	TrustedCertPEM  []byte
	ClientCertPEM   []byte
	ClientKeyPEM    []byte
	DefaultHeaders  map[string]string
	FollowRedirects bool
}

type HttpClient struct {
	client         *http.Client
	baseUrl        string
	defaultHeaders map[string]string
	interceptors   []Interceptor
}

func NewHttpClient(cfg *ClientConfig) (*HttpClient, error) {
	tlsCfg := &tls.Config{InsecureSkipVerify: cfg.SkipVerifySSL}

	// Trusted Cert
	if len(cfg.TrustedCertPEM) > 0 {
		certPool := x509.NewCertPool()
		if !certPool.AppendCertsFromPEM(cfg.TrustedCertPEM) {
			return nil, errors.New("failed to parse trusted cert PEM")
		}
		tlsCfg.InsecureSkipVerify = false
		tlsCfg.RootCAs = certPool
	}

	// mTLS
	if len(cfg.ClientCertPEM) > 0 && len(cfg.ClientKeyPEM) > 0 {
		cert, err := tls.X509KeyPair(cfg.ClientCertPEM, cfg.ClientKeyPEM)
		if err != nil {
			return nil, err
		}
		tlsCfg.Certificates = []tls.Certificate{cert}
	}

	transport := &http.Transport{TLSClientConfig: tlsCfg}

	httpClient := &http.Client{
		Timeout:   cfg.Timeout,
		Transport: transport,
	}

	if !cfg.FollowRedirects {
		httpClient.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}
	}

	return &HttpClient{
		client:         httpClient,
		defaultHeaders: cfg.DefaultHeaders,
	}, nil
}

func (hc *HttpClient) SetBaseUrl(baseUrl string) {
	hc.baseUrl = baseUrl
}

func (hc *HttpClient) AddInterceptor(i Interceptor) {
	hc.interceptors = append(hc.interceptors, i)
}

func (hc *HttpClient) Request(method, uri string, headers map[string]string, body []byte) (*Response, error) {
	fullUri := uri

	if hc.baseUrl != "" {
		fullUri = hc.baseUrl + uri
	}

	req, err := http.NewRequest(method, fullUri, bytes.NewReader(body))
	if err != nil {
		return nil, err
	}

	for k, v := range hc.defaultHeaders {
		req.Header.Set(k, v)
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	for _, i := range hc.interceptors {
		if err := i(req); err != nil {
			return nil, err
		}
	}

	start := time.Now()
	resp, err := hc.client.Do(req)
	duration := time.Since(start)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return &Response{
		RequestUri: fullUri,
		Status:     resp.StatusCode,
		Headers:    resp.Header,
		Body:       data,
		Raw:        resp,
		Duration:   duration,
	}, nil
}
