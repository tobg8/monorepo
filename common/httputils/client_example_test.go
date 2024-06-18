package httputils_test

import (
	"bytes"
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/monorepo/common/httputils"
)

func ExampleClient_Do() {
	client := httputils.NewClient(5*time.Second, 0)
	r, _ := http.NewRequest("GET", "http://www.leboncoin.fr", nil)
	var b bytes.Buffer
	code, err := client.Do(context.Background(), &b, r)
	fmt.Println("Error:", err)
	fmt.Println("HTTP Code:", code)
	fmt.Println("Response body:", b.String())
}
