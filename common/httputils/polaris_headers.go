package httputils

import (
	"net/http"
)

// PolarisHeaders contains all required headers
type PolarisHeaders http.Header

// GetPolarisHeaders return header struct with all internal polaris headers
// This will copy all authentication headers + Request Unique ID
//func GetPolarisHeaders(req *http.Request) PolarisHeaders {
//	headers := PolarisHeaders{}
//	for _, h := range []string{polarisheaders.UniqueID, polarisheaders.Authorization} {
//		v := req.Header[h]
//		if len(v) > 0 {
//			headers[h] = v
//		}
//	}
//	return headers
//}

// SetPolarisHeaders add PolarisHeaders to request headers
func SetPolarisHeaders(req *http.Request, headers PolarisHeaders) {
	for k, v := range headers {
		req.Header[k] = v
	}
}
