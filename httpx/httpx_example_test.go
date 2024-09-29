package httpx_test

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/gophero/goal/httpx"
)

// When running tests, the example code will be executed, and the example method must use fmt.println to output.
// In the end, the agreed output must be consistent with the Output to print the results, as inconsistent results
// will cause the test to fail.

func ExampleNewBuilder() {
	// start a test http server
	startServer()

	var ret string
	// build a get request
	httpx.NewBuilder("http://localhost:1234/html").WhenSuccess(func(resp *http.Response) { // request success
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Fatalf("response body should be readable but not: %v", err)
		}
		ret = string(body)
	}).WhenFailed(func(err error) { // request failed
		panic(err)
	}).Get() // request completed, clear resources

	fmt.Println(ret)
	// Output:
	// <html><head>test</head><body><h1>test page!</h1></body></html>
}

func ExampleGetString() {
	startServer()

	s, err := httpx.GetString("http://localhost:1234/json")
	if err != nil {
		log.Fatalf("should has no error but found: %v", err)
	}
	fmt.Println(s)
	// Output:
	// {"name":"lily","age":20,"height":70.5}
}

func ExampleMustGetString() {
	startServer()

	s := httpx.MustGetString("http://localhost:1234/json")
	fmt.Println(s)
	// Output:
	// {"name":"lily","age":20,"height":70.5}
}

func ExampleMustGet() {
	startServer()

	var s string
	httpx.MustGet("http://localhost:1234/json", func(resp *http.Response) {
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		s = string(bs)
	})
	fmt.Println(s)
	// Output:
	// {"name":"lily","age":20,"height":70.5}
}

func ExampleGet() {
	startServer()

	var s string
	httpx.Get("http://localhost:1234/json", func(resp *http.Response) {
		bs, err := io.ReadAll(resp.Body)
		if err != nil {
			panic(err)
		}
		s = string(bs)
	}, func(err error) {
		panic(err)
	})

	fmt.Println(s)
	// Output:
	// {"name":"lily","age":20,"height":70.5}
}

func ExampleGetBytes() {
	startServer()

	bs, err := httpx.GetBytes("http://localhost:1234/json")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(string(bs))
	// Output:
	// {"name":"lily","age":20,"height":70.5}
}

func ExampleMustGetBytes() {
	startServer()

	bs := httpx.MustGetBytes("http://localhost:1234/json")
	fmt.Println(string(bs))
	// Output:
	// {"name":"lily","age":20,"height":70.5}
}

type user struct {
	Name   string  `json:"name"`
	Age    int     `json:"age"`
	Height float32 `json:"height"`
}

func ExampleGetJson() {
	startServer()

	u, err := httpx.GetJson("http://localhost:1234/json", &user{})
	if err != nil {
		fmt.Println(err)
	}
	bs, _ := json.Marshal(u)
	fmt.Println(string(bs))
	// Output:
	// {"name":"lily","age":20,"height":70.5}
}

func ExampleMustGetJson() {
	startServer()

	u := httpx.MustGetJson("http://localhost:1234/json", &user{})

	bs, _ := json.Marshal(u)
	fmt.Println(string(bs))
	// Output:
	// {"name":"lily","age":20,"height":70.5}
}
