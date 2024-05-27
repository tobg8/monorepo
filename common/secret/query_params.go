package secret

import (
	"errors"
	"maps"
	"net/url"
)

// urlReplacement is replacement string for secret values in urls.
// Don't use asterisks to avoid percent-encoding. Same string as URL.Redacted.
const urlReplacement = "xxxxx"

// QueryParams provides methods to hide values of secret query string params,
// e.g. api keys. To avoid leak to client, log and monitoring backend.
//
// Also see httputils.WithSecretQueryParams.
//
// Note: When possible, prefer http header over url query string for sensitive
// data, as header values are generally less likely to find their way to logs
// and monitoring metadata than urls.
type QueryParams []string

// HideFromErr hides secret query params values from the url.Error present in
// given error.
//
// These errors are returned by http.Client.Do, http.NewRequest and url.Parse
// among others. They contain full url with query string in their messages e.g.:
//
//	Get "https://domain.example/path?apikey=a1b2c3": net/http:
//	request canceled (Client.Timeout exceeded while awaiting headers)
//
// Note: If an url.Error is found wrapped, returns it unwrapped (due to fmt
// wrapped error message being eagerly evaluated). Wrapping should be done after
// hiding.
func (p QueryParams) HideFromErr(e error) error {
	var urlErr *url.Error
	if len(p) == 0 || !errors.As(e, &urlErr) || urlErr.URL == "" {
		return e
	}

	urlErr.URL = p.HideFromURL(urlErr.URL)
	return urlErr
}

// HideFromURL hides secret query params values from given url.
func (p QueryParams) HideFromURL(s string) string {
	u, err := url.Parse(s)
	if err != nil {
		// Replace whole URL when unparsable.
		return urlReplacement
	}

	u.RawQuery = p.HideFromValues(u.Query()).Encode()
	return u.String()
}

// HideFromValues hides secret query params values from given url.Values.
func (p QueryParams) HideFromValues(original url.Values) url.Values {
	if len(p) == 0 {
		return original
	}

	hidden := maps.Clone(original)

	for _, k := range p {
		vs, present := original[k]
		if !present {
			continue
		}
		hiddenVals := make([]string, len(vs))
		for i, v := range vs {
			// Leave empty values empty to ease debugging.
			if v != "" {
				hiddenVals[i] = urlReplacement
			}
		}
		hidden[k] = hiddenVals
	}

	return hidden
}
