# Go Function Utilities

There was a time where devs wanted to make things easy (Go), but there are other ways to make things easy.

```bash
go get github.com/nitsugaro/go-utils@v1.4.1
```

### Map Tree

Use MapTree to parse an unknown dimentional `goutils.DefaultMap{...}` without performing a crash for access to a key/sub-key.

Use Cases

- External Data
- Unlimited Data to Typed
- Runtimes like JS

```go
data := goutils.DefaultMap{
		"user": goutils.DefaultMap{
			"name":  "Alice",
			"age":   28,
			"email": "alice@example.com",
			"roles": []interface{}{"admin", "editor", "viewer"},
		},
	}

	// Create a TreeMap
	tree := goutils.NewTreeMap(data)

	// --- Basic Get & Value Access ---
	fmt.Println("Name:", tree.Get("user.name").AsStringOr("Unknown"))
	fmt.Println("Age:", tree.Get("user.age").AsIntOr(-1))
	fmt.Println("Email:", tree.Get("user.email").AsStringOr("no email"))

	// --- Check Existence ---
	if tree.IsDefined("user.name") {
		fmt.Println("Name is defined.")
	}
	if !tree.Get("user.password").Exists() {
		fmt.Println("Password not found.")
	}
	if tree.Get("user.undefined").IsEmpty() {
		fmt.Println("Field is nil.")
	}

	// --- Set Value ---
	tree.Set("user.country", "Argentina")
	fmt.Println("Country:", tree.Get("user.country").AsStringOr("No country"))

	// --- Delete Field ---
	tree.Delete("user.email")
	fmt.Println("Email after delete:", tree.Get("user.email").AsStringOr("Deleted"))
```

### Http Client

Use to instance Http Client for requests.

```go
/*
ClientConfig {
	BaseURL         string
	Timeout         time.Duration
	SkipVerifySSL   bool
	TrustedCertPEM  []byte
	ClientCertPEM   []byte mTLS
	ClientKeyPEM    []byte mTLS
	DefaultHeaders  map[string]string
	FollowRedirects bool
}
*/

// Create an HTTP client with custom configuration
client, _ := goutils.NewHttpClient(&goutils.ClientConfig{
	BaseURL:         "https://httpbin.org",
	Timeout:         5 * time.Second,
	FollowRedirects: true,
	DefaultHeaders: map[string]string{
		"my-default-header": "1234",
	},
})

// Add an interceptor (e.g., to inject headers dynamically)
client.AddInterceptor(func(req *http.Request) error {
	req.Header.Add("x-transaction-id", "1234")
	return nil
})

// Make a GET request to /redirect/1
res, _ := client.Request("GET", "/redirect/1", nil, nil)

// Access response data
res.Status         // HTTP status code, e.g., 200
res.Headers        // http.Header (response headers)
res.Body           // []byte (raw body)
res.Duration       // Request duration (time.Duration)
res.Raw            // *http.Response (raw response)

// Print response body as string
fmt.Println(string(res.Body)) // (HTML or JSON, depending on the endpoint)
```