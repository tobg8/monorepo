package httptester

import (
	"io"
	"net/http"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func Test_Head_can_call_the_head_method_with_only_an_url_on_a_server(t *testing.T) {
	app := gin.New()
	app.HEAD("/url", func(c *gin.Context) {
		c.Status(http.StatusOK)
	})

	response := Head(app, "/url")

	assert.Equal(t, http.StatusOK, response.Code)
}

func Test_Head_can_call_the_head_method_with_a_body(t *testing.T) {
	app := gin.New()
	app.HEAD("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		c.Status(http.StatusOK)
	})

	response := Head(app, "/url", "request body")

	assert.Equal(t, http.StatusOK, response.Code)
}

func Test_Head_can_call_the_head_method_with_a_string_as_the_body(t *testing.T) {
	app := gin.New()
	app.HEAD("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		c.Status(http.StatusOK)
	})

	response := Head(app, "/url", "request body")

	assert.Equal(t, http.StatusOK, response.Code)
}

func Test_Head_can_call_the_head_method_with_headers_on_a_server(t *testing.T) {
	app := gin.New()
	app.HEAD("/url", func(c *gin.Context) {
		assert.Equal(t, "header value", c.Request.Header.Get("header name"))
		c.Status(http.StatusOK)
	})

	response := Head(app, "/url", map[string]string{"header name": "header value"})

	assert.Equal(t, http.StatusOK, response.Code)
}

func Test_Head_can_call_the_head_method_with_a_body_and_headers_on_a_server(t *testing.T) {
	app := gin.New()
	app.HEAD("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		assert.Equal(t, "header value", c.Request.Header.Get("header name"))
		c.Status(http.StatusOK)
	})

	response := Head(app, "/url", "request body", map[string]string{"header name": "header value"})

	assert.Equal(t, http.StatusOK, response.Code)
}

func Test_Head_panics_when_two_bodies_are_provided(t *testing.T) {
	app := gin.New()
	assert.Panics(t, func() {
		_ = Head(app, "/url", "body", "body")
	})
}

func Test_Head_panics_when_two_headers_are_provided(t *testing.T) {
	app := gin.New()
	assert.Panics(t, func() {
		_ = Head(app, "/url", map[string]string{"header name": "header value"}, map[string]string{})
	})
}

func Test_Get_can_call_the_get_method_with_only_an_url_on_a_server(t *testing.T) {
	app := gin.New()
	app.GET("/url", func(c *gin.Context) {
		c.String(http.StatusOK, "body")
	})

	response := Get(app, "/url")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Get_can_call_the_get_method_with_a_body(t *testing.T) {
	app := gin.New()
	app.GET("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		c.String(http.StatusOK, "body")
	})

	response := Get(app, "/url", "request body")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Get_can_call_the_get_method_with_a_string_as_the_body(t *testing.T) {
	app := gin.New()
	app.GET("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		c.String(http.StatusOK, "body")
	})

	response := Get(app, "/url", "request body")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Get_can_call_the_get_method_with_headers_on_a_server(t *testing.T) {
	app := gin.New()
	app.GET("/url", func(c *gin.Context) {
		assert.Equal(t, "header value", c.Request.Header.Get("header name"))
		c.String(http.StatusOK, "body")
	})

	response := Get(app, "/url", map[string]string{"header name": "header value"})

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Get_can_call_the_get_method_with_a_body_and_headers_on_a_server(t *testing.T) {
	app := gin.New()
	app.GET("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		assert.Equal(t, "header value", c.Request.Header.Get("header name"))
		c.String(http.StatusOK, "body")
	})

	response := Get(app, "/url", "request body", map[string]string{"header name": "header value"})

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Get_panics_when_two_bodies_are_provided(t *testing.T) {
	app := gin.New()
	assert.Panics(t, func() {
		_ = Get(app, "/url", "body", "body")
	})
}

func Test_Get_panics_when_two_headers_are_provided(t *testing.T) {
	app := gin.New()
	assert.Panics(t, func() {
		_ = Get(app, "/url", map[string]string{"header name": "header value"}, map[string]string{})
	})
}

func Test_Post_can_call_the_post_method(t *testing.T) {
	app := gin.New()
	app.POST("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		c.String(http.StatusOK, "body")
	})

	response := Post(app, "/url", "request body")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Post_can_call_the_post_method_with_headers_on_a_server(t *testing.T) {
	app := gin.New()
	app.POST("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		assert.Equal(t, "header value", c.Request.Header.Get("header name"))
		c.String(http.StatusOK, "body")
	})

	response := Post(app, "/url", "request body", map[string]string{"header name": "header value"})

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Post_panics_when_two_headers_are_provided(t *testing.T) {
	app := gin.New()
	assert.Panics(t, func() {
		_ = Post(app, "/url", "request body", map[string]string{"header name": "header value"}, map[string]string{"header name": "header value"})
	})
}

func Test_Delete_can_call_the_delete_method_with_only_an_url_on_a_server(t *testing.T) {
	app := gin.New()
	app.DELETE("/url", func(c *gin.Context) {
		c.String(http.StatusOK, "body")
	})

	response := Delete(app, "/url")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Delete_can_call_the_delete_method_with_a_body(t *testing.T) {
	app := gin.New()
	app.DELETE("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		c.String(http.StatusOK, "body")
	})

	response := Delete(app, "/url", "request body")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Delete_can_call_the_delete_method_with_a_string_as_the_body(t *testing.T) {
	app := gin.New()
	app.DELETE("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		c.String(http.StatusOK, "body")
	})

	response := Delete(app, "/url", "request body")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Delete_can_call_the_delete_method_with_headers_on_a_server(t *testing.T) {
	app := gin.New()
	app.DELETE("/url", func(c *gin.Context) {
		assert.Equal(t, "header value", c.Request.Header.Get("header name"))
		c.String(http.StatusOK, "body")
	})

	response := Delete(app, "/url", map[string]string{"header name": "header value"})

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Delete_can_call_the_delete_method_with_a_body_and_headers_on_a_server(t *testing.T) {
	app := gin.New()
	app.DELETE("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		assert.Equal(t, "header value", c.Request.Header.Get("header name"))
		c.String(http.StatusOK, "body")
	})

	response := Delete(app, "/url", "request body", map[string]string{"header name": "header value"})

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Delete_panics_when_two_bodies_are_provided(t *testing.T) {
	app := gin.New()
	assert.Panics(t, func() {
		_ = Delete(app, "/url", "body", "body")
	})
}

func Test_Delete_panics_when_two_headers_are_provided(t *testing.T) {
	app := gin.New()
	assert.Panics(t, func() {
		_ = Delete(app, "/url", map[string]string{"header name": "header value"}, map[string]string{})
	})
}

func Test_Put_can_call_the_put_method(t *testing.T) {
	app := gin.New()
	app.PUT("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		c.String(http.StatusOK, "body")
	})

	response := Put(app, "/url", "request body")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Put_can_call_the_put_method_with_headers_on_a_server(t *testing.T) {
	app := gin.New()
	app.PUT("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		assert.Equal(t, "header value", c.Request.Header.Get("header name"))
		c.String(http.StatusOK, "body")
	})

	response := Put(app, "/url", "request body", map[string]string{"header name": "header value"})

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Put_panics_when_two_headers_are_provided(t *testing.T) {
	app := gin.New()
	assert.Panics(t, func() {
		_ = Put(app, "/url", "request body", map[string]string{"header name": "header value"}, map[string]string{"header name": "header value"})
	})
}

func Test_Patch_can_call_the_patch_method(t *testing.T) {
	app := gin.New()
	app.PATCH("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		c.String(http.StatusOK, "body")
	})

	response := Patch(app, "/url", "request body")

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Patch_can_call_the_patch_method_with_headers_on_a_server(t *testing.T) {
	app := gin.New()
	app.PATCH("/url", func(c *gin.Context) {
		requestBodyBytes, _ := io.ReadAll(c.Request.Body)
		assert.Equal(t, "request body", string(requestBodyBytes))
		assert.Equal(t, "header value", c.Request.Header.Get("header name"))
		c.String(http.StatusOK, "body")
	})

	response := Patch(app, "/url", "request body", map[string]string{"header name": "header value"})

	assert.Equal(t, http.StatusOK, response.Code)
	assert.Equal(t, "body", response.Body.String())
}

func Test_Patch_panics_when_two_headers_are_provided(t *testing.T) {
	app := gin.New()
	assert.Panics(t, func() {
		_ = Patch(app, "/url", "request body", map[string]string{"header name": "header value"}, map[string]string{"header name": "header value"})
	})
}
