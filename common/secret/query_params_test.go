package secret

import (
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func Test_QueryParams_HideFromErr(t *testing.T) {
	t.Run("hides from url.Error", func(t *testing.T) {
		e := &url.Error{
			Op:  "foo",
			URL: "https://domain.test/path?bar=banana&baz=lemon&foo=apple",
			Err: errors.New("err"),
		}

		out := QueryParams{"bar"}.HideFromErr(e)

		var ue *url.Error
		assert.ErrorAs(t, out, &ue)
		assert.Equal(t, "https://domain.test/path?bar=xxxxx&baz=lemon&foo=apple", ue.URL)
		assert.NotContains(t, out.Error(), "banana")
	})

	t.Run("hides from wrapped url.Error", func(t *testing.T) {
		e := fmt.Errorf("wrapped: %w", &url.Error{
			Op:  "foo",
			URL: "https://domain.test/path?bar=banana&baz=lemon&foo=apple",
			Err: errors.New("err"),
		})

		out := QueryParams{"bar"}.HideFromErr(e)

		var ue *url.Error
		assert.ErrorAs(t, out, &ue)
		assert.Equal(t, "https://domain.test/path?bar=xxxxx&baz=lemon&foo=apple", ue.URL)
		assert.NotContains(t, out.Error(), "banana")
	})

	t.Run("unwraps", func(t *testing.T) {
		e := fmt.Errorf("wrapped: %w", &url.Error{
			Op:  "foo",
			URL: "https://domain.test/path?bar=banana&baz=lemon&foo=apple",
			Err: errors.New("err"),
		})

		out := QueryParams{"bar"}.HideFromErr(e)
		assert.Error(t, out)
		assert.NotErrorIs(t, out, e)
	})

	t.Run("hides whole url when unparsable", func(t *testing.T) {
		e := &url.Error{
			Op:  "foo",
			URL: "https://domain.test/path?bar=banana&baz=lemon\n&foo=apple",
			Err: errors.New("err"),
		}

		out := QueryParams{"bar"}.HideFromErr(e)

		var ue *url.Error
		assert.ErrorAs(t, out, &ue)
		assert.Equal(t, "xxxxx", ue.URL)
		assert.NotContains(t, out.Error(), "banana")
	})

	t.Run("hides from http.NewRequest error", func(t *testing.T) {
		_, e := http.NewRequest(http.MethodGet,
			":invalid://domain.test/path?bar=banana&baz=lemon&foo=apple", nil)
		require.Error(t, e)

		out := QueryParams{"bar"}.HideFromErr(e)
		assert.Error(t, out)
		assert.NotContains(t, out.Error(), "banana")
	})

	t.Run("hides from http.Client.Do error", func(t *testing.T) {
		req, err := http.NewRequest(http.MethodGet,
			"example://domain.test/path?bar=banana&baz=lemon&foo=apple", nil)
		require.NoError(t, err)

		_, e := (&http.Client{}).Do(req)
		require.Error(t, e)

		out := QueryParams{"bar"}.HideFromErr(e)

		var ue *url.Error
		assert.ErrorAs(t, out, &ue)
		assert.Equal(t, "example://domain.test/path?bar=xxxxx&baz=lemon&foo=apple", ue.URL)
		assert.NotContains(t, out.Error(), "banana")
	})

	t.Run("hides percent-encoded", func(t *testing.T) {
		const key = `🗝️A*UWFh&fnn*xCw"N!Ny QAT?Rs@6@_.L2q!c`

		u, err := url.Parse("https://domain.test/path")
		require.NoError(t, err)

		u.RawQuery = url.Values{
			"foo": []string{"apple"},
			"bar": []string{key},
			"baz": []string{"lemon"},
		}.Encode()

		e := &url.Error{
			Op:  "foo",
			URL: u.String(),
			Err: errors.New("err"),
		}

		out := QueryParams{"bar"}.HideFromErr(e)

		var ue *url.Error
		assert.ErrorAs(t, out, &ue)
		assert.Equal(t, "https://domain.test/path?bar=xxxxx&baz=lemon&foo=apple", ue.URL)
		assert.NotContains(t, out.Error(), key)
		assert.NotContains(t, out.Error(), url.QueryEscape(key))
		assert.NotContains(t, out.Error(), url.PathEscape(key))
	})

	t.Run("does nothing with empty secret params", func(t *testing.T) {
		e := fmt.Errorf("wrapped: %w", &url.Error{
			Op:  "foo",
			URL: "https://domain.test/path?bar=banana&baz=lemon&foo=apple",
			Err: errors.New("err"),
		})

		out := QueryParams{}.HideFromErr(e)
		assert.ErrorIs(t, out, e)
	})

	t.Run("does nothing with nil error", func(t *testing.T) {
		out := QueryParams{}.HideFromErr(nil)
		assert.Nil(t, out)
	})

	t.Run("does nothing with non url.Error", func(t *testing.T) {
		e := errors.New("err")
		out := QueryParams{}.HideFromErr(e)
		assert.ErrorIs(t, out, e)
	})

	t.Run("does nothing with empty URL", func(t *testing.T) {
		e := fmt.Errorf("wrapped: %w", &url.Error{
			Op:  "foo",
			Err: errors.New("err"),
		})

		out := QueryParams{}.HideFromErr(e)
		assert.ErrorIs(t, out, e)
	})
}

func Test_QueryParams_HideFromURL(t *testing.T) {
	t.Run("hides from url", func(t *testing.T) {
		u := "https://domain.test/path?bar=banana&baz=lemon&foo=apple"
		out := QueryParams{"bar"}.HideFromURL(u)
		assert.Equal(t, "https://domain.test/path?bar=xxxxx&baz=lemon&foo=apple", out)
	})

	t.Run("hides percent-encoded", func(t *testing.T) {
		const key = `🗝️A*UWFh&fnn*xCw"N!Ny QAT?Rs@6@_.L2q!c`

		u, err := url.Parse("https://domain.test/path")
		require.NoError(t, err)

		u.RawQuery = url.Values{
			"foo": []string{"apple"},
			"bar": []string{key},
			"baz": []string{"lemon"},
		}.Encode()

		out := QueryParams{"bar"}.HideFromURL(u.String())
		assert.Equal(t, "https://domain.test/path?bar=xxxxx&baz=lemon&foo=apple", out)
	})

	t.Run("hides whole url when unparsable", func(t *testing.T) {
		u := "https://domain.test/path?bar=banana&baz=lemon\n&foo=apple"
		out := QueryParams{"bar"}.HideFromURL(u)
		assert.Equal(t, "xxxxx", out)
	})
}

func Test_QueryParams_HideFromValues(t *testing.T) {
	t.Run("hides secret params", func(t *testing.T) {
		params := url.Values{}
		params.Add("foo", "apple")
		params.Add("bar", "banana")
		params.Add("bar", "cherry")
		params.Add("baz", "lemon")
		params.Add("qux", "orange")

		out := QueryParams{"bar", "qux"}.HideFromValues(params)
		assert.Equal(t, url.Values{
			"foo": []string{"apple"},
			"bar": []string{"xxxxx", "xxxxx"},
			"baz": []string{"lemon"},
			"qux": []string{"xxxxx"},
		}, out)
	})

	t.Run("leaves absent params absent", func(t *testing.T) {
		params := url.Values{}
		params.Add("foo", "apple")
		params.Add("qux", "orange")

		out := QueryParams{"bar", "qux"}.HideFromValues(params)
		assert.Equal(t, url.Values{
			"foo": []string{"apple"},
			"qux": []string{"xxxxx"},
		}, out)
	})

	t.Run("leaves empty values empty", func(t *testing.T) {
		params := url.Values{}
		params.Add("foo", "apple")
		params.Add("bar", "")
		params.Add("bar", "banana")
		params.Add("bar", "cherry")
		params.Add("baz", "lemon")
		params.Add("qux", "")

		out := QueryParams{"bar", "qux"}.HideFromValues(params)
		assert.Equal(t, url.Values{
			"foo": []string{"apple"},
			"bar": []string{"", "xxxxx", "xxxxx"},
			"baz": []string{"lemon"},
			"qux": []string{""},
		}, out)
	})

	t.Run("doesn't modify original", func(t *testing.T) {
		params := url.Values{}
		params.Add("foo", "apple")
		params.Add("bar", "banana")

		out := QueryParams{"bar"}.HideFromValues(params)
		assert.NotEqual(t, params, out)
		assert.Equal(t, url.Values{
			"foo": []string{"apple"},
			"bar": []string{"banana"},
		}, params)
	})

	t.Run("does nothing with empty secret params", func(t *testing.T) {
		params := url.Values{}
		params.Add("foo", "apple")
		params.Add("bar", "banana")

		out := QueryParams{}.HideFromValues(params)
		assert.Equal(t, params, out)
	})
}
