package httputils_test

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/monorepo/common/httputils"
	"github.com/stretchr/testify/mock"
)

func ExampleMockClient_Do_reader() {
	r, _ := http.NewRequest("GET", "https://www.leboncoin.test", nil)

	client := new(httputils.MockClient)
	client.On("Do", mock.Anything, r).Return(http.StatusOK, nil, strings.NewReader("Hello, world!"))

	var b bytes.Buffer
	code, err := client.Do(context.Background(), &b, r)

	fmt.Println("Error:", err)
	fmt.Println("HTTP Code:", code)
	fmt.Println("Response body:", b.String())
	// Output:
	// Error: <nil>
	// HTTP Code: 200
	// Response body: Hello, world!
}

func ExampleMockClient_Do_string() {
	r, _ := http.NewRequest("GET", "https://www.leboncoin.test", nil)

	client := new(httputils.MockClient)
	client.On("Do", mock.Anything, r).Return(http.StatusOK, nil, "Hello, world!")

	var b bytes.Buffer
	code, err := client.Do(context.Background(), &b, r)

	fmt.Println("Error:", err)
	fmt.Println("HTTP Code:", code)
	fmt.Println("Response body:", b.String())
	// Output:
	// Error: <nil>
	// HTTP Code: 200
	// Response body: Hello, world!
}

func ExampleMockClient_Do_bytes() {
	r, _ := http.NewRequest("GET", "https://www.leboncoin.test", nil)

	client := new(httputils.MockClient)
	client.On("Do", mock.Anything, r).Return(http.StatusOK, nil, []byte("Hello, world!"))

	var b bytes.Buffer
	code, err := client.Do(context.Background(), &b, r)

	fmt.Println("Error:", err)
	fmt.Println("HTTP Code:", code)
	fmt.Println("Response body:", b.String())
	// Output:
	// Error: <nil>
	// HTTP Code: 200
	// Response body: Hello, world!
}

func ExampleMockClient_Do_error() {
	r, _ := http.NewRequest("GET", "https://www.leboncoin.test", nil)

	client := new(httputils.MockClient)
	client.On("Do", mock.Anything, r).Return(http.StatusBadRequest, errors.New("oops"), nil)

	var b bytes.Buffer
	code, err := client.Do(context.Background(), &b, r)

	fmt.Println("Error:", err)
	fmt.Println("HTTP Code:", code)
	fmt.Println("Response body:", b.String())
	// Output:
	// Error: oops
	// HTTP Code: 400
	// Response body:
}

func ExampleMockClient_DoAndUnmarshalJSON_reader() {
	r, _ := http.NewRequest("GET", "https://www.leboncoin.test", nil)

	client := new(httputils.MockClient)
	client.On("DoAndUnmarshalJSON", mock.Anything, r).
		Return(http.StatusOK, nil, strings.NewReader(`{"message": "Hello, world!"}`))

	var m struct{ Message string }
	code, err := client.DoAndUnmarshalJSON(context.Background(), &m, r)

	fmt.Println("Error:", err)
	fmt.Println("HTTP Code:", code)
	fmt.Println("Response body.message:", m.Message)
	// Output:
	// Error: <nil>
	// HTTP Code: 200
	// Response body.message: Hello, world!
}

func ExampleMockClient_DoAndUnmarshalJSON_string() {
	r, _ := http.NewRequest("GET", "https://www.leboncoin.test", nil)

	client := new(httputils.MockClient)
	client.On("DoAndUnmarshalJSON", mock.Anything, r).
		Return(http.StatusOK, nil, `{"message": "Hello, world!"}`)

	var m struct{ Message string }
	code, err := client.DoAndUnmarshalJSON(context.Background(), &m, r)

	fmt.Println("Error:", err)
	fmt.Println("HTTP Code:", code)
	fmt.Println("Response body.message:", m.Message)
	// Output:
	// Error: <nil>
	// HTTP Code: 200
	// Response body.message: Hello, world!
}

func ExampleMockClient_DoAndUnmarshalJSON_bytes() {
	r, _ := http.NewRequest("GET", "https://www.leboncoin.test", nil)

	client := new(httputils.MockClient)
	client.On("DoAndUnmarshalJSON", mock.Anything, r).
		Return(http.StatusOK, nil, []byte(`{"message": "Hello, world!"}`))

	var m struct{ Message string }
	code, err := client.DoAndUnmarshalJSON(context.Background(), &m, r)

	fmt.Println("Error:", err)
	fmt.Println("HTTP Code:", code)
	fmt.Println("Response body.message:", m.Message)
	// Output:
	// Error: <nil>
	// HTTP Code: 200
	// Response body.message: Hello, world!
}

func ExampleMockClient_DoAndUnmarshalJSON_error() {
	r, _ := http.NewRequest("GET", "https://www.leboncoin.test", nil)

	client := new(httputils.MockClient)
	client.On("DoAndUnmarshalJSON", mock.Anything, r).
		Return(http.StatusBadRequest, errors.New("oops"), nil)

	var m struct{ Message string }
	code, err := client.DoAndUnmarshalJSON(context.Background(), &m, r)

	fmt.Println("Error:", err)
	fmt.Println("HTTP Code:", code)
	fmt.Println("Response body.message:", m.Message)
	// Output:
	// Error: oops
	// HTTP Code: 400
	// Response body.message:
}
