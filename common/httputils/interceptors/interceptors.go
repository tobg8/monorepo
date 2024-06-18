package interceptors

import (
	"github.com/f2prateek/train"
	"github.com/monorepo/common/monitoring/metrics"
)

// GetDefaultInterceptors returns a list of default interceptors
func GetDefaultInterceptors(appName string, appVersion string, sh metrics.StatsdHandler) []train.Interceptor {
	return []train.Interceptor{
		//NewUniqueID(),
		NewUserAgent(appName, appVersion),
		NewMonitoring(sh),
		NewTracing(),
	}
}
