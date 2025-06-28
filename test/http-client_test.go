package test

import (
	"fmt"
	"net/http"
	"testing"
	"time"

	goutils "github.com/nitsugaro/go-utils"
)

func Test_http_client(t *testing.T) {
	client, err := goutils.NewHttpClient(&goutils.ClientConfig{
		BaseURL:         "https://httpbin.org",
		Timeout:         5 * time.Second,
		FollowRedirects: true,
		DefaultHeaders: map[string]string{
			"my-default-header": "1234",
		},
	})

	client.AddInterceptor(func(req *http.Request) error {
		req.Header.Add("x-transaction-id", "1234")

		return nil
	})

	if err != nil {
		t.Errorf("expected configuration not cause errors and got: %s", err.Error())
	}

	res1, err := client.Request("GET", "/redirect/1", nil, nil)
	if err != nil {
		t.Errorf("request with method 'get' got: %s", err.Error())
	}

	if res1.Status != 200 {
		fmt.Println(res1.Status)
		t.Errorf("expected status 200 and got: %v", res1.Status)
	}

	if res1.Headers.Get("my-default-header") == "1234" {
		fmt.Println(res1.Status)
		t.Errorf("expected header 'my-default-header' be '1234' and got: %v", res1.Headers.Get("my-default-header"))
	}

	if res1.Headers.Get("x-transaction-id") == "1234" {
		fmt.Println(res1.Status)
		t.Errorf("expected header 'x-transaction-id' be '1234' and got: %v", res1.Headers.Get("x-transaction-id"))
	}
}
