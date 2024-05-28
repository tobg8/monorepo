package benchmarks

import (
	"flag"
)

type benchFlags struct {
	LoggerName     string
	LoopIterations int
}

var globalBenchFlags benchFlags

func init() {
	flag.StringVar(
		&globalBenchFlags.LoggerName, "test.logger", "",
		"The name of the logger kind to benchmark. "+
			"Empty string for benchmark all available loggers.",
	)

	flag.IntVar(
		&globalBenchFlags.LoopIterations, "test.loop", 10,
		"The number of iteration of the loop in each benchmark.",
	)
}
