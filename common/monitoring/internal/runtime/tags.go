package runtime

import (
	"runtime"
	"strings"
)

// Tags returns all tags related to the Go runtime.
//
//   - go_arch:<arch>:       The running program's architecture target (e.g. amd64, arm64).
//   - go_version:<version>: The Go release tag (e.g. 1.19.9, 1.20.4).
//   - go_experiment_<X>:    Any GOEXPERIMENTs (Go experimental feature) are set to non-default values
//     (e.g. boringcrypto, or coverageredesign since Go 1.20).
func Tags() []string {
	fields := strings.Fields(runtime.Version())

	tags := []string{
		"go_arch:" + runtime.GOARCH,
		"go_version:" + strings.TrimPrefix(fields[0], "go"),
	}

	if len(fields) > 1 {
		for _, x := range strings.Split(strings.TrimPrefix(fields[1], "X:"), ",") {
			tags = append(tags, "go_experiment_"+x)
		}
	}

	return tags
}
