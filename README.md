# Go Function Utilities

There was a time where devs wanted to make things easy (Go), but there are other ways to make things easy.

### Map Tree

Use MapTree to parse an unknown dimentional `goutils.DefaultMap{...}` without performing a crash for access to a key/sub-key.

Use Cases

- External Data
- Unlimited Data to Typed
- Runtimes like JS

```go
tree := goutils.NewTreeMap(map[string]interface{}{ "initial": "my-value" })

tree.Get("initial").AsString()          // "my-value", nil

tree.Set("sub.key", goutils.DefaultMap{ "slice": []interface{}{1, 2, 3} })

tree.IsDefined("sub.key")               // true

tree.Get("sub.key").AsMap()             // goutils.DefaultMap{ "slice": [...] }, nil

tree.Get("sub.key.slice.0").AsInt()     // 1, nil

tree.Get("sub.key.slice").AsSlice()     // []*TreeMap{1, 2, 3}, nil

tree.ToJsonString(true)
// {
//   "initial": "my-value",
//   "sub": {
//     "key": {
//       "slice": [
//         1,
//         2,
//         3
//       ]
//     }
//   }
// }

// AsSliceOf
type MyStruct struct {
	Name string
}
tree.Set("users", []interface{}{
	goutils.DefaultMap{"Name": "Agus"},
})
var users []MyStruct
tree.Get("users").AsSliceOf(&users) // []MyStruct{{"Agus"}}, nil

// AsStruct
tree.Set("config", goutils.DefaultMap{"Name": "TreeMap"})
var cfg MyStruct
tree.Get("config").AsStruct(&cfg) // MyStruct{"TreeMap"}, nil

// AsInt / AsFloat / AsBool
tree.Set("num", 42)
tree.Get("num").AsInt()   // 42, nil
tree.Get("num").AsFloat() // 42.0, nil

tree.Set("flag", "true")
tree.Get("flag").AsBool() // true, nil

mapTree.Delete("sub.key.slice") // *TreeMap{...}, nil
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